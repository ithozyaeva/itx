<script setup lang="ts">
import type { EventSearchFilters } from '@/services/events'
import { onBeforeUnmount, reactive, ref, watch } from 'vue'
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

const debounceTimer = ref<ReturnType<typeof setTimeout>>()

watch(filters, () => {
  clearTimeout(debounceTimer.value)
  debounceTimer.value = setTimeout(() => {
    emit('change', { ...filters })
  }, 350)
}, { deep: true })

onBeforeUnmount(() => clearTimeout(debounceTimer.value))
</script>

<template>
  <div class="flex flex-col sm:flex-row gap-3">
    <Input
      v-model="filters.title"
      placeholder="Поиск по названию..."
      class="w-full sm:flex-1"
    />
    <Select v-model="filters.placeType" class="w-full sm:w-auto">
      <SelectTrigger class="w-full sm:w-44">
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
