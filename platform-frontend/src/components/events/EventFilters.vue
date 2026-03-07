<script setup lang="ts">
import type { EventSearchFilters } from '@/services/events'
import { reactive, watch } from 'vue'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
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
  <div class="flex flex-col sm:flex-row gap-3">
    <Input
      v-model="filters.title"
      placeholder="Поиск по названию..."
      class="w-full sm:max-w-xs"
    />
    <Select v-model="filters.placeType" class="w-full sm:max-w-xs">
      <SelectTrigger class="w-full sm:max-w-xs">
        <SelectValue placeholder="Все форматы" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="">
          Все форматы
        </SelectItem>
        <SelectItem v-for="pt in placeTypes" :key="pt.value" :value="pt.value">
          {{ pt.label }}
        </SelectItem>
      </SelectContent>
    </Select>
  </div>
</template>
