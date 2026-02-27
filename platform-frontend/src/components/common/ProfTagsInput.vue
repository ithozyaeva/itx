<script setup lang="ts">
import type { ProfTag } from '@/models/profile'
import type { AcceptableInputValue } from 'reka-ui'
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxList } from '@/components/ui/combobox'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import { profileSiervice } from '@/services/profile'
import { computed, onMounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array as () => ProfTag[],
    default: () => [],
  },
  readonly: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const localTags = ref<ProfTag[]>(props.modelValue)

watch(() => props.modelValue, (newValue) => {
  localTags.value = newValue
})

const allProfTags = ref<ProfTag[]>([])
const searchProfTag = ref('')
const openCombobox = ref(false)

const filteredProfTags = computed(() => {
  if (!searchProfTag.value)
    return []

  return allProfTags.value.filter(tag =>
    tag.title.toLowerCase().includes(searchProfTag.value.toLowerCase())
    && !localTags.value.some(t => t.id === tag.id || t.title === tag.title),
  )
})

function addProfTag(newTags: AcceptableInputValue[]) {
  const uniqueTags: ProfTag[] = []
  const titleSet = new Set<string>(localTags.value.map(tag => tag.title))

  for (const tag of newTags as ProfTag[]) {
    const normalizedTitle = tag.title.trim().toLowerCase()
    if (normalizedTitle && !titleSet.has(normalizedTitle)) {
      titleSet.add(normalizedTitle)
      uniqueTags.push({ id: tag.id || 0, title: tag.title.trim() } as ProfTag)
    }
  }
  const updatedTags = [...localTags.value, ...uniqueTags]
  localTags.value = updatedTags
  emit('update:modelValue', updatedTags)
  searchProfTag.value = ''
}

function removeProfTag(tag: AcceptableInputValue) {
  tag = tag as ProfTag
  const index = localTags.value.findIndex(t => t.id === tag.id || t.title === tag.title)
  if (index !== -1) {
    removeTagByIndex(index)
  }
}

function removeTagByIndex(index: number) {
  const newTags = [...localTags.value]
  newTags.splice(index, 1)
  localTags.value = newTags
  emit('update:modelValue', newTags)
}

async function loadProfTags() {
  try {
    const response = await profileSiervice.getAllProfTags()
    allProfTags.value = response
  }
  catch (error) {
    console.error('Ошибка при загрузке профессиональных тегов:', error)
    // TODO: Показать сообщение об ошибке пользователю
  }
}

function render(profTag: AcceptableInputValue) {
  return (profTag as ProfTag).title
}

function convertValue(text: string): ProfTag {
  return { id: 0, title: text.trim() } as ProfTag
}

function handleComboboxItemSelect(profTag: ProfTag) {
  if (!localTags.value.some(t => t.id === profTag.id || t.title === profTag.title)) {
    const updatedTags = [...localTags.value, profTag]
    localTags.value = updatedTags
    emit('update:modelValue', updatedTags)
  }
  searchProfTag.value = ''
  openCombobox.value = false
}

onMounted(loadProfTags)
</script>

<template>
  <div class="w-full">
    <Combobox v-model:model-value="localTags" v-model:open="openCombobox" :ignore-filter="true" :multiple="true">
      <ComboboxAnchor class="w-full" as-child>
        <TagsInput
          v-model:model-value="localTags"
          :display-value="render"
          :convert-value="convertValue"
          :readonly="readonly"
          @update:model-value="addProfTag"
          @remove-tag="removeProfTag"
        >
          <TagsInputItem v-for="item in localTags" :key="item.id || item.title" :value="item">
            <TagsInputItemText>{{ item.title }}</TagsInputItemText>
            <TagsInputItemDelete v-if="!readonly" />
          </TagsInputItem>
          <ComboboxInput v-model="searchProfTag" :readonly="readonly" class="w-full" as-child>
            <TagsInputInput
              placeholder="Теги..."
              :readonly="readonly"
              class="w-full box-shadow-0 p-0 shadow-none border-none focus-visible:ring-0 h-auto"
            />
          </ComboboxInput>
        </TagsInput>
      </ComboboxAnchor>
      <ComboboxList v-if="filteredProfTags.length > 0" class="w-[--reka-popper-anchor-width]">
        <ComboboxEmpty />
        <ComboboxGroup>
          <ComboboxItem
            v-for="profTag in filteredProfTags" :key="profTag.id"
            :value="profTag"
            @select.prevent="() => handleComboboxItemSelect(profTag)"
          >
            {{ profTag.title }}
          </ComboboxItem>
        </ComboboxGroup>
      </ComboboxList>
    </Combobox>
  </div>
</template>

<style scoped>
/* Add any specific styles here if needed */
</style>
