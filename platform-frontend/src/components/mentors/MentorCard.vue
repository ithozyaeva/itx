<script setup lang="ts">
import type { Mentor } from '@/models/profile'
import { Tag, Typography } from 'itx-ui-kit'

defineProps<{
  mentor: Mentor
}>()
</script>

<template>
  <RouterLink
    :to="`/mentors/${mentor.id}`"
    data-reveal
    class="bg-card rounded-3xl border p-4 hover:shadow-md transition-shadow flex flex-col gap-2 cursor-pointer"
  >
    <Typography variant="h4" as="h3">
      {{ mentor.firstName }} {{ mentor.lastName }}
    </Typography>
    <p v-if="mentor.occupation" class="text-sm text-muted-foreground">
      {{ mentor.occupation }}
    </p>
    <p v-if="mentor.experience" class="text-sm">
      {{ mentor.experience }}
    </p>
    <div v-if="mentor.profTags?.length" class="flex flex-wrap gap-1 mt-1">
      <Tag v-for="tag in mentor.profTags" :key="tag.id">
        {{ tag.title }}
      </Tag>
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
  </RouterLink>
</template>
