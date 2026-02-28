<script setup lang="ts">
import type { ReferalLink } from '@/models/referals'
import { computed } from 'vue'

const props = defineProps<{
  referral: ReferalLink
}>()

const gradeLabel = computed(() => {
  const map: Record<string, string> = {
    junior: 'Junior',
    middle: 'Middle',
    senior: 'Senior',
  }
  return map[props.referral.grade] || props.referral.grade
})
</script>

<template>
  <div class="rounded-2xl border bg-card p-4 hover:shadow-md transition-shadow">
    <div class="flex items-start justify-between gap-2">
      <div class="min-w-0">
        <h4 class="font-medium truncate">
          {{ referral.company }}
        </h4>
        <p class="text-sm text-muted-foreground mt-0.5">
          {{ gradeLabel }}
        </p>
      </div>
      <span
        class="shrink-0 rounded-full px-2 py-0.5 text-xs font-medium"
        :class="referral.status === 'active'
          ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
          : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300'"
      >
        {{ referral.status === 'active' ? 'Активна' : 'Заморожена' }}
      </span>
    </div>
    <div
      v-if="referral.profTags.length"
      class="flex flex-wrap gap-1 mt-2"
    >
      <span
        v-for="tag in referral.profTags"
        :key="tag.id"
        class="rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground"
      >
        {{ tag.title }}
      </span>
    </div>
    <div class="flex items-center justify-between mt-3 text-xs text-muted-foreground">
      <span>{{ referral.vacationsCount }} вакансий</span>
      <span>{{ referral.conversionsCount }} откликов</span>
    </div>
  </div>
</template>
