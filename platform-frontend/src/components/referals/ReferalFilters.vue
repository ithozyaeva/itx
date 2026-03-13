<script setup lang="ts">
import type { ReferalSearchFilters } from '@/services/referals'
import { onBeforeUnmount, reactive, ref, watch } from 'vue'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useDictionary } from '@/composables/useDictionary'

const emit = defineEmits<{
  change: [filters: ReferalSearchFilters]
}>()

const filters = reactive<ReferalSearchFilters>({
  grade: '',
  company: '',
  status: '',
})

const { grades, referalLinkStatuses } = useDictionary(['grades', 'referalLinkStatuses'])

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
  <div class="flex flex-col gap-4">
    <div class="flex flex-col sm:flex-row gap-4 flex-wrap">
      <Input
        v-model="filters.company"
        placeholder="Поиск по компании..."
        class="max-w-xs"
      />
      <Select v-model="filters.grade" class="max-w-xs">
        <SelectTrigger class="max-w-xs">
          <SelectValue placeholder="Все грейды" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="">
            Все грейды
          </SelectItem>
          <SelectItem
            v-for="g in grades"
            :key="g.value"
            :value="g.value"
          >
            {{ g.label }}
          </SelectItem>
        </SelectContent>
      </Select>
      <Select v-model="filters.status" class="max-w-xs">
        <SelectTrigger class="max-w-xs">
          <SelectValue placeholder="Все статусы" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="">
            Все статусы
          </SelectItem>
          <SelectItem
            v-for="s in referalLinkStatuses"
            :key="s.value"
            :value="s.value"
          >
            {{ s.label }}
          </SelectItem>
        </SelectContent>
      </Select>
    </div>
  </div>
</template>
