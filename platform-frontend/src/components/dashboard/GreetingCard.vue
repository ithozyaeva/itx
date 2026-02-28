<script setup lang="ts">
import { useUser, useUserLevel } from '@/composables/useUser'
import { SUBSCRIPTION_LEVELS } from '@/models/profile'

const user = useUser()
const { level, levelIndex } = useUserLevel()

const daysSinceJoined = 1
</script>

<template>
  <div class="rounded-3xl border bg-card p-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold">
          Привет, {{ user?.firstName }} 👋
        </h1>
        <p class="text-muted-foreground mt-1">
          Ты в IT-Хозяевах уже {{ daysSinceJoined }} {{ daysSinceJoined === 1 ? 'день' : 'дней' }}
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
