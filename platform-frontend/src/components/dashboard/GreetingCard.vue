<script setup lang="ts">
import { computed } from 'vue'
import { useUser, useUserLevel } from '@/composables/useUser'
import { SUBSCRIPTION_LEVELS } from '@/models/profile'

const user = useUser()
const { level, levelIndex } = useUserLevel()

function pluralizeDays(n: number): string {
  if (n % 10 === 1 && n % 100 !== 11)
    return 'день'
  if (n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20))
    return 'дня'
  return 'дней'
}

const daysSinceJoined = computed(() => {
  if (!user.value?.createdAt)
    return 1
  const created = new Date(user.value.createdAt)
  const now = new Date()
  const diff = Math.floor((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24))
  return Math.max(diff, 1)
})
</script>

<template>
  <div class="rounded-3xl border bg-card p-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold">
          Привет, {{ user?.firstName }}
        </h1>
        <p class="text-muted-foreground mt-1">
          Ты в IT-Хозяевах уже {{ daysSinceJoined }} {{ pluralizeDays(daysSinceJoined) }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-sm font-medium">{{ level }}</span>
        <div class="flex gap-1">
          <span
            v-for="(lvl, i) in SUBSCRIPTION_LEVELS"
            :key="lvl"
            class="h-2.5 w-2.5 rounded-full transition-colors"
            :class="i <= levelIndex ? 'bg-green-500' : 'bg-muted'"
          />
        </div>
      </div>
    </div>
  </div>
</template>
