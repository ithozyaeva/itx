<script setup lang="ts">
import type { Mentor } from '@/models/profile'
import { Typography } from '@/components/ui/typography'

const props = defineProps<{
  mentor: Mentor
}>()

const maxVisibleTags = 5

function getInitials(mentor: Mentor) {
  return `${mentor.firstName.charAt(0)}${mentor.lastName.charAt(0)}`.toUpperCase()
}
</script>

<template>
  <RouterLink
    :to="`/mentors/${props.mentor.id}`"
    data-reveal
    class="bg-card rounded-sm border p-4 hover:border-accent/50 hover:-translate-y-0.5 transition-all duration-200 flex gap-4 cursor-pointer terminal-card"
  >
    <!-- Avatar -->
    <div class="shrink-0">
      <img
        v-if="mentor.avatarUrl"
        :src="mentor.avatarUrl"
        :alt="`${mentor.firstName} ${mentor.lastName}`"
        class="w-12 h-12 rounded-full object-cover"
      >
      <div
        v-else
        class="w-12 h-12 rounded-full bg-accent text-white flex items-center justify-center text-sm font-medium"
      >
        {{ getInitials(mentor) }}
      </div>
    </div>

    <!-- Info -->
    <div class="flex flex-col gap-1.5 min-w-0">
      <Typography
        variant="h4"
        as="h3"
        class="truncate"
      >
        {{ mentor.firstName }} {{ mentor.lastName }}
      </Typography>

      <p v-if="mentor.occupation" class="text-sm text-muted-foreground truncate">
        {{ mentor.occupation }}
      </p>

      <p v-if="mentor.experience" class="text-sm line-clamp-2">
        {{ mentor.experience }}
      </p>

      <div v-if="mentor.profTags?.length" class="flex flex-wrap gap-1 mt-1">
        <span
          v-for="tag in mentor.profTags.slice(0, maxVisibleTags)"
          :key="tag.id"
          class="inline-flex items-center rounded-full border border-accent/30 px-2 py-0.5 text-xs text-accent"
        >
          {{ tag.title }}
        </span>
        <span
          v-if="mentor.profTags.length > maxVisibleTags"
          class="inline-flex items-center rounded-full border border-muted px-2 py-0.5 text-xs text-muted-foreground"
        >
          +{{ mentor.profTags.length - maxVisibleTags }}
        </span>
      </div>

      <a
        v-if="mentor.tg"
        :href="`https://t.me/${mentor.tg}`"
        target="_blank"
        class="text-sm text-primary underline mt-auto"
        @click.stop
      >
        @{{ mentor.tg }}
      </a>
    </div>
  </RouterLink>
</template>
