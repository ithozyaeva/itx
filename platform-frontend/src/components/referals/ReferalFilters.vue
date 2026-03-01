<script setup lang="ts">
import type { ProfTag } from '@/models/profile'
import type { ReferalSearchFilters } from '@/services/referals'
import { onMounted, reactive, ref, watch } from 'vue'
import { Input } from '@/components/ui/input'
import { useDictionary } from '@/composables/useDictionary'

const emit = defineEmits<{
  change: [filters: ReferalSearchFilters]
}>()

const filters = reactive<ReferalSearchFilters>({
  grade: '',
  company: '',
  status: '',
  profTagIds: [],
})

const { grades, referalLinkStatuses } = useDictionary(['grades', 'referalLinkStatuses'])

const profTags = ref<ProfTag[]>([])

onMounted(async () => {
  try {
    const response = await fetch('/api/profTags')
    const data = await response.json()
    profTags.value = data.items ?? []
  }
  catch {
    profTags.value = []
  }
})

function toggleTag(tagId: number) {
  const idx = filters.profTagIds?.indexOf(tagId) ?? -1
  if (idx === -1) {
    filters.profTagIds = [...(filters.profTagIds ?? []), tagId]
  }
  else {
    filters.profTagIds = filters.profTagIds?.filter(id => id !== tagId) ?? []
  }
}

watch(filters, () => {
  emit('change', { ...filters })
}, { deep: true })
</script>

<template>
  <div class="flex flex-col gap-4">
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
        <option
          v-for="g in grades"
          :key="g.value"
          :value="g.value"
        >
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
        <option
          v-for="s in referalLinkStatuses"
          :key="s.value"
          :value="s.value"
        >
          {{ s.label }}
        </option>
      </select>
    </div>
    <div
      v-if="profTags.length > 0"
      class="flex flex-wrap gap-2"
    >
      <button
        v-for="tag in profTags"
        :key="tag.id"
        class="px-3 py-1 text-xs rounded-full border transition-colors"
        :class="filters.profTagIds?.includes(tag.id)
          ? 'bg-primary text-primary-foreground border-primary'
          : 'bg-transparent border-input hover:border-primary/50'"
        @click="toggleTag(tag.id)"
      >
        {{ tag.title }}
      </button>
    </div>
  </div>
</template>
