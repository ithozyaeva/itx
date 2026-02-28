<script setup lang="ts">
import type { EventSearchFilters } from '@/services/events'
import { reactive, watch } from 'vue'
import { Input } from '@/components/ui/input'
import { useDictionary } from '@/composables/useDictionary'

const emit = defineEmits<{
  change: [filters: EventSearchFilters]
}>()

const filters = reactive<EventSearchFilters>({
  title: '',
  placeType: '',
})

const { placeTypes } = useDictionary(['placeTypes'])

watch(filters, () => {
  emit('change', { ...filters })
}, { deep: true })
</script>

<template>
  <div class="flex flex-col sm:flex-row gap-4">
    <Input
      v-model="filters.title"
      placeholder="Поиск по названию..."
      class="max-w-xs"
    />
    <select
      v-model="filters.placeType"
      class="border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring max-w-xs"
    >
      <option value="">
        Все форматы
      </option>
      <option v-for="pt in placeTypes" :key="pt.value" :value="pt.value">
        {{ pt.label }}
      </option>
    </select>
  </div>
</template>
