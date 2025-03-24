<script setup lang="ts">
import type { ProfTag, ProfTagFormData } from '@/models/mentors'
import type { AcceptableInputValue } from 'reka-ui'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxList } from '@/components/ui/combobox'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import { profTagService } from '@/services/profTagService'
import { computed, onMounted, ref } from 'vue'

const props = defineProps<{
  profTags: ProfTagFormData[]
}>()

const emit = defineEmits(['update:profTags'])
const localProfTags = ref<ProfTagFormData[]>(props.profTags)
const allProfTags = ref<ProfTag[]>([])
const searchProfTag = ref('')
const openCombobox = ref(false)

// Фильтрованные профессиональные теги
const filteredProfTags = computed(() => {
  if (!searchProfTag.value)
    return []

  return allProfTags.value.filter(tag =>
    tag.title.toLowerCase().includes(searchProfTag.value.toLowerCase())
    && !props.profTags.some(t => t.id === tag.id || t.title === tag.title),
  )
})

function addProfTag(newTags: AcceptableInputValue[]) {
  const uniqueTags: ProfTagFormData[] = []
  const titleSet = new Set<string>()

  for (const tag of newTags as ProfTagFormData[]) {
    if (!titleSet.has(tag.title)) {
      titleSet.add(tag.title)
      uniqueTags.push(tag)
    }
  }
  emit('update:profTags', uniqueTags)
  searchProfTag.value = ''
}

function pushTag(newTag: ProfTag) {
  localProfTags.value.push(newTag)
  emit('update:profTags', localProfTags.value)
}

// Удаление профессионального тега
function removeProfTag(tag: AcceptableInputValue) {
  tag = tag as ProfTag
  const index = props.profTags.findIndex(t => t.id === tag.id || t.title === tag.title)
  removeTag(index)
}
// Удаление тега
function removeTag(index: number) {
  const newTags = [...props.profTags]
  newTags.splice(index, 1)
  emit('update:profTags', newTags)
}

// Загрузка всех профессиональных тегов
async function loadProfTags() {
  try {
    const response = await profTagService.getAll()
    allProfTags.value = response.items
  }
  catch (error) {
    console.error('Ошибка при загрузке профессиональных тегов:', error)
  }
}

function render(profTag: AcceptableInputValue) {
  return (profTag as ProfTag).title
}
function convertValue(text: string) {
  return {
    title: text,
  }
}

onMounted(loadProfTags)
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Профессиональные теги</CardTitle>
    </CardHeader>
    <CardContent>
      <Combobox v-model="localProfTags" v-model:open="openCombobox" :ignore-filter="true">
        <ComboboxAnchor class="w-full" as-child>
          <TagsInput v-model:model-value="localProfTags" :display-value="render" :convert-value="convertValue" @update:model-value="addProfTag" @remove-tag="removeProfTag">
            <TagsInputItem v-for="item in props.profTags" :key="item.title" :value="item">
              <TagsInputItemText />
              <TagsInputItemDelete />
            </TagsInputItem>
            <ComboboxInput v-model="searchProfTag" class="w-full" as-child>
              <TagsInputInput placeholder="Теги..." class="w-full p-0 border-none focus-visible:ring-0 h-auto" />
            </ComboboxInput>
          </TagsInput>
          <ComboboxList v-if="filteredProfTags.length > 0" class="w-[--reka-popper-anchor-width]">
            <ComboboxEmpty />
            <ComboboxGroup>
              <ComboboxItem
                v-for="profTag in filteredProfTags" :key="profTag.title" :value="profTag"
                @select.prevent="(ev: any) => {
                  searchProfTag = ''
                  console.log(ev.detail)
                  pushTag(ev.detail.value)

                  if (filteredProfTags.length === 0) {
                    openCombobox = false
                  }
                }"
              >
                {{ profTag.title }}
              </ComboboxItem>
            </ComboboxGroup>
          </ComboboxList>
        </ComboboxAnchor>
      </Combobox>
    </CardContent>
  </Card>
</template>
