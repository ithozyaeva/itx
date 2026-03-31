<script setup lang="ts">
import type { Mentor } from '@/models/profile'
import { Typography } from 'itx-ui-kit'
import { Loader2, Users } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import MentorCard from '@/components/mentors/MentorCard.vue'
import MentorCardSkeleton from '@/components/mentors/MentorCardSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useCardReveal } from '@/composables/useCardReveal'
import { handleError } from '@/services/errorService'
import { mentorsService } from '@/services/mentors'

const PAGE_SIZE = 12

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const mentors = ref<Mentor[]>([])
const total = ref(0)
const isLoading = ref(false)
const isLoadingMore = ref(false)
const loadError = ref<string | null>(null)
const searchQuery = ref('')
const selectedTag = ref('')
const sortBy = ref<'name' | 'services'>('name')

async function loadMentors() {
  isLoading.value = true
  loadError.value = null
  try {
    const result = await mentorsService.getAll(PAGE_SIZE, 0)
    mentors.value = result.items ?? []
    total.value = result.total
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

async function loadMore() {
  isLoadingMore.value = true
  try {
    const result = await mentorsService.getAll(PAGE_SIZE, mentors.value.length)
    mentors.value.push(...(result.items ?? []))
    total.value = result.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoadingMore.value = false
  }
}

const allTags = computed(() => {
  const tags = new Map<number, string>()
  for (const mentor of mentors.value) {
    for (const tag of mentor.profTags ?? []) {
      tags.set(tag.id, tag.title)
    }
  }
  return Array.from(tags, ([id, title]) => ({ id, title }))
})

const hasActiveFilters = computed(() => !!searchQuery.value || !!selectedTag.value)

const sortOptions: { key: 'name' | 'services', label: string }[] = [
  { key: 'name', label: 'По имени' },
  { key: 'services', label: 'По количеству услуг' },
]

const filteredMentors = computed(() => {
  const filtered = mentors.value.filter((mentor) => {
    const matchesSearch = !searchQuery.value
      || `${mentor.firstName} ${mentor.lastName}`.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesTag = !selectedTag.value
      || mentor.profTags?.some(tag => tag.id === Number(selectedTag.value))
    return matchesSearch && matchesTag
  })
  return filtered.sort((a, b) => {
    if (sortBy.value === 'services')
      return (b.services?.length ?? 0) - (a.services?.length ?? 0)
    const nameA = `${a.firstName} ${a.lastName}`.toLowerCase()
    const nameB = `${b.firstName} ${b.lastName}`.toLowerCase()
    return nameA.localeCompare(nameB, 'ru')
  })
})

onMounted(loadMentors)
</script>

<template>
  <div ref="containerRef" class="container mx-auto px-4 py-8">
    <Typography variant="h2" as="h1" class="mb-8">
      Менторы
    </Typography>

    <div class="flex flex-col gap-4 mb-6">
      <div class="flex flex-col sm:flex-row gap-4">
        <Input
          v-model="searchQuery"
          placeholder="Поиск по имени..."
          class="max-w-xs"
        />
        <Select v-model="selectedTag" class="max-w-xs">
          <SelectTrigger class="max-w-xs">
            <SelectValue placeholder="Все теги" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">
              Все теги
            </SelectItem>
            <SelectItem
              v-for="tag in allTags"
              :key="tag.id"
              :value="String(tag.id)"
            >
              {{ tag.title }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="flex gap-2 flex-wrap items-center">
        <span class="text-sm text-muted-foreground">Сортировка:</span>
        <button
          v-for="option in sortOptions"
          :key="option.key"
          class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
          :class="sortBy === option.key
            ? 'bg-primary text-primary-foreground'
            : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
          @click="sortBy = option.key"
        >
          {{ option.label }}
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <MentorCardSkeleton v-for="i in 6" :key="i" />
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="loadMentors"
    />

    <template v-else>
      <EmptyState
        v-if="filteredMentors.length === 0"
        :icon="Users"
        title="Менторов пока нет"
        description="Скоро здесь появятся менторы сообщества"
      />
      <template v-else>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <MentorCard
            v-for="mentor in filteredMentors"
            :key="mentor.id"
            :mentor="mentor"
          />
        </div>
        <div v-if="mentors.length < total && !hasActiveFilters" class="mt-6 flex justify-center">
          <Button
            variant="outline"
            :disabled="isLoadingMore"
            @click="loadMore"
          >
            <Loader2 v-if="isLoadingMore" class="mr-2 h-4 w-4 animate-spin" />
            Показать ещё
          </Button>
        </div>
      </template>
    </template>
  </div>
</template>
