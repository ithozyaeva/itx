<script setup lang="ts">
import type { AcceptableInputValue } from 'reka-ui'
import type { Member } from '@/models/members'
import { asyncComputed, useDebounce } from '@vueuse/core'
import { computed, ref } from 'vue'
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxList } from '@/components/ui/combobox'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import { memberService } from '@/services/memberService'

const props = defineProps<{
  hosts: Member[]
}>()

const emit = defineEmits(['update:hosts'])

const localHosts = computed({
  get: () => props.hosts,
  set: (val: Member[]) => emit('update:hosts', val),
})

const searchMemberString = ref('')
const debouncedSearch = useDebounce(searchMemberString, 400)
const openCombobox = ref(false)

// Фильтрованные профессиональные теги
const findedMembers = asyncComputed(async () => {
  if (!debouncedSearch.value)
    return []

  return (await memberService.search({ username: debouncedSearch.value, offset: 0, limit: 100 })).items!
}, [])

function addHost(newHosts: AcceptableInputValue[]) {
  const uniqueMembers: Member[] = []
  const tgIdSet = new Set<number>()

  for (const member of newHosts as Member[]) {
    if (!member?.id) {
      return
    }

    if (!tgIdSet.has(member.telegramID)) {
      tgIdSet.add(member.telegramID)
      uniqueMembers.push(member)
    }
  }
  emit('update:hosts', uniqueMembers)
  searchMemberString.value = ''
}

function pushHost(newHost: Member) {
  if (!newHost?.id) {
    return
  }

  localHosts.value.push(newHost)
}

// Удаление профессионального тега
function removeHost(host: AcceptableInputValue) {
  host = host as Member
  const index = localHosts.value.findIndex(t => t.id === host.id)

  removeHostByIndex(index)
}
// Удаление тега
function removeHostByIndex(index: number) {
  localHosts.value.splice(index, 1)
}

function render(host: AcceptableInputValue) {
  return `${(host as Member).firstName} ${(host as Member).lastName}`
}
function convertValue(text: string) {
  return {
    title: text,
  }
}
</script>

<template>
  <Combobox v-model="localHosts" v-model:open="openCombobox" :ignore-filter="true">
    <ComboboxAnchor class="w-full" as-child>
      <TagsInput v-model:model-value="localHosts" :display-value="render" :convert-value="convertValue" @update:model-value="addHost" @remove-tag="removeHost">
        <TagsInputItem v-for="item in localHosts" :key="item.id" :value="item">
          <TagsInputItemText />
          <TagsInputItemDelete @click="removeHost(item)" />
        </TagsInputItem>
        <ComboboxInput v-model="searchMemberString" class="w-full" as-child>
          <TagsInputInput placeholder="Спикеры... (искать по tg_username)" class="w-full p-0 border-none focus-visible:ring-0 h-auto" />
        </ComboboxInput>
      </TagsInput>
      <ComboboxList v-if="findedMembers.length > 0" class="w-[--reka-popper-anchor-width]">
        <ComboboxEmpty />
        <ComboboxGroup>
          <ComboboxItem
            v-for="member in findedMembers" :key="member.id" :value="member"
            @select.prevent="(ev: any) => {
              searchMemberString = ''
              pushHost(ev.detail.value)

              if (findedMembers.length === 0) {
                openCombobox = false
              }
            }"
          >
            {{ member.firstName }} {{ member.lastName }}
          </ComboboxItem>
        </ComboboxGroup>
      </ComboboxList>
    </ComboboxAnchor>
  </Combobox>
</template>
