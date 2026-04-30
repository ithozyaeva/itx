<script setup lang="ts">
import type { AcceptableInputValue } from 'reka-ui'
import { computed, onMounted, ref, watch } from 'vue'
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxList } from '@/components/ui/combobox'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import { AI_MATERIAL_LIMITS } from '@/models/aiMaterial'
import { aiMaterialsService } from '@/services/aiMaterials'

const props = defineProps<{
  modelValue: string[]
  readonly?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const tags = ref<string[]>([...props.modelValue])
watch(() => props.modelValue, val => (tags.value = [...val]))

const search = ref('')
const open = ref(false)
const allTags = ref<string[]>([])

const suggestions = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q)
    return []
  return allTags.value
    .filter(t => t.toLowerCase().includes(q) && !tags.value.includes(t))
    .slice(0, 8)
})

function normalize(raw: string): string | null {
  const v = raw.trim().toLowerCase()
  if (!v)
    return null
  if (v.length > AI_MATERIAL_LIMITS.tagLenMax)
    return null
  return v
}

function commit(next: string[]) {
  // Дедуп + лимит — единственный источник истины при любых добавлениях.
  const out: string[] = []
  for (const raw of next) {
    const v = normalize(raw)
    if (!v)
      continue
    if (out.includes(v))
      continue
    out.push(v)
    if (out.length >= AI_MATERIAL_LIMITS.tagsMax)
      break
  }
  tags.value = out
  emit('update:modelValue', out)
}

function onTagsUpdate(next: AcceptableInputValue[]) {
  commit(next as string[])
  search.value = ''
}

function selectSuggestion(tag: string) {
  commit([...tags.value, tag])
  search.value = ''
  open.value = false
}

onMounted(async () => {
  try {
    const res = await aiMaterialsService.topTags()
    allTags.value = res.tags
  }
  catch {
    // suggestions не критичны — autocomplete просто не заполнится
  }
})
</script>

<template>
  <div class="w-full">
    <Combobox v-model:model-value="tags" v-model:open="open" :ignore-filter="true" :multiple="true">
      <ComboboxAnchor class="w-full" as-child>
        <TagsInput
          :model-value="tags"
          :readonly="readonly"
          @update:model-value="onTagsUpdate"
        >
          <TagsInputItem v-for="t in tags" :key="t" :value="t">
            <TagsInputItemText>{{ t }}</TagsInputItemText>
            <TagsInputItemDelete v-if="!readonly" />
          </TagsInputItem>
          <ComboboxInput v-model="search" :readonly="readonly" class="w-full" as-child>
            <TagsInputInput
              placeholder="Теги (до 5)"
              :readonly="readonly"
              :maxlength="AI_MATERIAL_LIMITS.tagLenMax"
              class="w-full box-shadow-0 p-0 shadow-none border-none focus-visible:ring-0 h-auto"
            />
          </ComboboxInput>
        </TagsInput>
      </ComboboxAnchor>
      <ComboboxList v-if="suggestions.length > 0" class="w-[--reka-popper-anchor-width]">
        <ComboboxEmpty />
        <ComboboxGroup>
          <ComboboxItem
            v-for="t in suggestions"
            :key="t"
            :value="t"
            @select.prevent="() => selectSuggestion(t)"
          >
            {{ t }}
          </ComboboxItem>
        </ComboboxGroup>
      </ComboboxList>
    </Combobox>
    <p class="mt-1 text-xs text-muted-foreground">
      Свободные теги, по 5 максимум. Помогают другим найти материал.
    </p>
  </div>
</template>
