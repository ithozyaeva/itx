<script setup lang="ts">
import type { AIMaterialFilters, AIMaterialKind, AIMaterialSort } from '@/models/aiMaterial'
import { Search } from 'lucide-vue-next'
import { computed } from 'vue'
import { AI_MATERIAL_KIND_OPTIONS } from '@/models/aiMaterial'

const props = defineProps<{ modelValue: AIMaterialFilters }>()
const emit = defineEmits<{ 'update:modelValue': [v: AIMaterialFilters] }>()

const filters = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

function setKind(kind: AIMaterialKind | '') {
  filters.value = { ...filters.value, kind, offset: 0 }
}

function setSort(sort: AIMaterialSort) {
  filters.value = { ...filters.value, sort, offset: 0 }
}

function setQuery(q: string) {
  filters.value = { ...filters.value, q, offset: 0 }
}

const kindTabs: { key: AIMaterialKind | '', label: string }[] = [
  { key: '', label: 'Все' },
  ...AI_MATERIAL_KIND_OPTIONS.map(o => ({ key: o.value, label: o.label })),
]

const sortOptions: { key: AIMaterialSort, label: string }[] = [
  { key: 'new', label: 'Сначала новые' },
  { key: 'popular', label: 'Популярные' },
]
</script>

<template>
  <div class="flex flex-col gap-3 mb-6">
    <div class="flex flex-col sm:flex-row sm:items-center gap-3">
      <div class="flex gap-2 flex-wrap">
        <button
          v-for="tab in kindTabs"
          :key="tab.key"
          class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
          :class="(filters.kind || '') === tab.key
            ? 'bg-primary text-primary-foreground'
            : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
          @click="setKind(tab.key)"
        >
          {{ tab.label }}
        </button>
      </div>
      <div class="relative sm:ml-auto">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <input
          :value="filters.q ?? ''"
          type="search"
          placeholder="Поиск по названию или описанию..."
          class="w-full sm:w-72 rounded-sm border border-border bg-background pl-9 pr-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
          @input="setQuery(($event.target as HTMLInputElement).value)"
        >
      </div>
    </div>

    <div class="flex items-center gap-2 flex-wrap">
      <span class="text-sm text-muted-foreground">Сортировка:</span>
      <button
        v-for="opt in sortOptions"
        :key="opt.key"
        class="px-3 py-1 rounded-full text-xs font-medium transition-colors"
        :class="(filters.sort || 'new') === opt.key
          ? 'bg-primary text-primary-foreground'
          : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
        @click="setSort(opt.key)"
      >
        {{ opt.label }}
      </button>
    </div>
  </div>
</template>
