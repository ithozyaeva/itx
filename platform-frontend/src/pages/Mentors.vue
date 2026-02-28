<script setup lang="ts">
import type { Mentor } from '@/models/profile'
import MentorCard from '@/components/mentors/MentorCard.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useCardReveal } from '@/composables/useCardReveal'
import { mentorsService } from '@/services/mentors'
import { Typography } from 'itx-ui-kit'
import { Loader2, UserX } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'

const PAGE_SIZE = 12

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const mentors = ref<Mentor[]>([])
const total = ref(0)
const isLoading = ref(false)
const isLoadingMore = ref(false)
const searchQuery = ref('')
const selectedTag = ref('')

async function loadMentors() {
  isLoading.value = true
  try {
    const result = await mentorsService.getAll(PAGE_SIZE, 0)
    mentors.value = result.items
    total.value = result.total
  }
  finally {
    isLoading.value = false
  }
}

async function loadMore() {
  isLoadingMore.value = true
  try {
    const result = await mentorsService.getAll(PAGE_SIZE, mentors.value.length)
    mentors.value.push(...result.items)
    total.value = result.total
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

const filteredMentors = computed(() => {
  return mentors.value.filter((mentor) => {
    const matchesSearch = !searchQuery.value
      || `${mentor.firstName} ${mentor.lastName}`.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesTag = !selectedTag.value
      || mentor.profTags?.some(tag => tag.id === Number(selectedTag.value))
    return matchesSearch && matchesTag
  })
})

onMounted(loadMentors)
</script>

<template>
  <div ref="containerRef" class="container mx-auto px-4 py-8">
    <Typography variant="h2" as="h1" class="mb-8">
      Менторы
    </Typography>

    <div class="flex flex-col sm:flex-row gap-4 mb-6">
      <Input
        v-model="searchQuery"
        placeholder="Поиск по имени..."
        class="max-w-xs"
      />
      <select
        v-model="selectedTag"
        class="border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring max-w-xs"
      >
        <option value="">
          Все теги
        </option>
        <option v-for="tag in allTags" :key="tag.id" :value="tag.id">
          {{ tag.title }}
        </option>
      </select>
    </div>

    <div v-if="isLoading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <div v-if="filteredMentors.length === 0" class="flex flex-col items-center gap-2 py-8 text-muted-foreground">
        <UserX class="h-10 w-10" />
        <p>Менторы не найдены</p>
      </div>
      <template v-else>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <MentorCard
            v-for="mentor in filteredMentors"
            :key="mentor.id"
            :mentor="mentor"
          />
        </div>
        <div v-if="mentors.length < total" class="mt-6 flex justify-center">
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
