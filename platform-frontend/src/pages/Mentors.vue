<script setup lang="ts">
import type { Mentor } from '@/models/profile'
import MentorCard from '@/components/mentors/MentorCard.vue'
import { Input } from '@/components/ui/input'
import { useCardReveal } from '@/composables/useCardReveal'
import { mentorsService } from '@/services/mentors'
import { Typography } from 'itx-ui-kit'
import { computed, onMounted, ref } from 'vue'

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const mentors = ref<Mentor[]>([])
const searchQuery = ref('')
const selectedTag = ref('')

async function loadMentors() {
  const result = await mentorsService.getAll()
  mentors.value = result.items
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

    <div v-if="filteredMentors.length === 0" class="text-muted-foreground">
      Менторы не найдены
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <MentorCard
        v-for="mentor in filteredMentors"
        :key="mentor.id"
        :mentor="mentor"
      />
    </div>
  </div>
</template>
