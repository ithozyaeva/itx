<script setup lang="ts">
import type { AcceptableInputValue } from 'reka-ui'
import type { EventTag } from '@/models/events.ts'
import { computed, onMounted, ref, watch } from 'vue'
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxList } from '@/components/ui/combobox'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import { eventTagService } from '@/services/eventTagService.ts'

const props = defineProps<{
  eventTags: EventTag[]
}>()

const emit = defineEmits(['update:eventTags'])
const localEventTags = ref<EventTag[]>(props.eventTags)
const allEventTags = ref<EventTag[]>([])
const searchEventTag = ref('')
const openCombobox = ref(false)

watch(() => props.eventTags, (newVal) => {
  localEventTags.value = [...newVal]
}, { deep: true })

// Отфильтрованные теги ивентов
const filteredEventTags = computed(() => {
  if (!searchEventTag.value)
    return []

  return allEventTags.value.filter(tag =>
    tag.name.toLowerCase().includes(searchEventTag.value.toLowerCase())
    && !props.eventTags.some(t => t.id === tag.id || t.name === tag.name),
  )
})

function addEventTag(newTags: AcceptableInputValue[]) {
  const uniqueTags: EventTag[] = []
  const nameSet = new Set<string>()

  for (const tag of newTags as EventTag[]) {
    if (!nameSet.has(tag.name)) {
      nameSet.add(tag.name)
      uniqueTags.push(tag)
    }
  }
  emit('update:eventTags', uniqueTags)
  searchEventTag.value = ''
}

function pushTag(newTag: EventTag) {
  localEventTags.value.push(newTag)
  emit('update:eventTags', localEventTags.value)
}

// Удаление тега ивента
function removeEventTag(tag: AcceptableInputValue) {
  tag = tag as EventTag
  const index = props.eventTags.findIndex(t => t.id === tag.id || t.name === tag.name)
  removeTag(index)
}
// Удаление тега ивента
function removeTag(index: number) {
  const newTags = [...props.eventTags]
  newTags.splice(index, 1)
  emit('update:eventTags', newTags)
}

// Загрузка всех тегов ивентов
async function loadEventTags() {
  try {
    const response = await eventTagService.getAll()
    allEventTags.value = response.items
  }
  catch (error) {
    console.error('Ошибка при загрузке профессиональных тегов:', error)
  }
}

function render(eventTag: AcceptableInputValue) {
  return (eventTag as EventTag).name
}
function convertValue(text: string) {
  return {
    name: text,
  }
}

onMounted(loadEventTags)
</script>

<template>
  <Combobox v-model="localEventTags" v-model:open="openCombobox" :ignore-filter="true">
    <ComboboxAnchor class="w-full" as-child>
      <TagsInput v-model:model-value="localEventTags" :display-value="render" :convert-value="convertValue" @update:model-value="addEventTag" @remove-tag="removeEventTag">
        <TagsInputItem v-for="item in localEventTags" :key="item.name" :value="item">
          <TagsInputItemText />
          <TagsInputItemDelete @click="removeEventTag(item)" />
        </TagsInputItem>
        <ComboboxInput v-model="searchEventTag" class="w-full" as-child>
          <TagsInputInput placeholder="Теги..." class="w-full p-0 border-none focus-visible:ring-0 h-auto" />
        </ComboboxInput>
      </TagsInput>
      <ComboboxList v-if="filteredEventTags.length > 0" class="w-[--reka-popper-anchor-width]">
        <ComboboxEmpty />
        <ComboboxGroup>
          <ComboboxItem
            v-for="eventTag in filteredEventTags" :key="eventTag.name" :value="eventTag"
            @select.prevent="(ev: any) => {
              searchEventTag = ''
              pushTag(ev.detail.value)

              if (filteredEventTags.length === 0) {
                openCombobox = false
              }
            }"
          >
            {{ eventTag.name }}
          </ComboboxItem>
        </ComboboxGroup>
      </ComboboxList>
    </ComboboxAnchor>
  </Combobox>
</template>
