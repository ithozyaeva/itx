<script setup lang="ts">
import type { ReferalSearchFilters } from '@/services/referals'
import { Input } from '@/components/ui/input'
import { useDictionary } from '@/composables/useDictionary'
import { reactive, watch } from 'vue'

const emit = defineEmits<{
  change: [filters: ReferalSearchFilters]
}>()

const filters = reactive<ReferalSearchFilters>({
  grade: '',
  company: '',
  status: '',
})

const { grades, referalLinkStatuses } = useDictionary(['grades', 'referalLinkStatuses'])

watch(filters, () => {
  emit('change', { ...filters })
}, { deep: true })
</script>

<template>
  <div class="flex flex-col sm:flex-row gap-4 flex-wrap">
    <Input
      v-model="filters.company"
      placeholder="Поиск по компании..."
      class="max-w-xs"
    />
    <select
      v-model="filters.grade"
      class="border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring max-w-xs"
    >
      <option value="">
        Все грейды
      </option>
      <option v-for="g in grades" :key="g.value" :value="g.value">
        {{ g.label }}
      </option>
    </select>
    <select
      v-model="filters.status"
      class="border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring max-w-xs"
    >
      <option value="">
        Все статусы
      </option>
      <option v-for="s in referalLinkStatuses" :key="s.value" :value="s.value">
        {{ s.label }}
      </option>
    </select>
  </div>
</template>
