<script setup lang="ts">
import type { LeaderboardEntry } from '@/models/points'
import { Typography } from 'itx-ui-kit'
import { Loader2, Trophy } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'

const entries = ref<LeaderboardEntry[]>([])
const isLoading = ref(true)

async function fetchLeaderboard() {
  isLoading.value = true
  try {
    const response = await pointsService.getLeaderboard(50)
    entries.value = response.items || []
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchLeaderboard()
})

function getAvatarSrc(entry: LeaderboardEntry) {
  return entry.avatarUrl || `https://t.me/i/userpic/160/${entry.tg}.jpg`
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <Typography
      variant="h2"
      as="h1"
      class="mb-6"
    >
      Рейтинг участников
    </Typography>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div
      v-else-if="entries.length === 0"
      class="text-center py-12 text-muted-foreground"
    >
      Пока нет данных о баллах. Участвуйте в событиях, чтобы зарабатывать баллы!
    </div>

    <div
      v-else
      class="space-y-3"
    >
      <div
        v-for="(entry, index) in entries"
        :key="entry.memberId"
        class="flex items-center gap-4 p-4 bg-card border border-border rounded-2xl shadow-sm"
        :class="{ 'border-yellow-500/50 bg-yellow-500/5': index < 3 }"
      >
        <div
          class="flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold shrink-0"
          :class="index < 3 ? 'bg-yellow-500/20 text-yellow-500' : 'bg-muted text-muted-foreground'"
        >
          {{ index + 1 }}
        </div>

        <div class="w-10 h-10 rounded-full overflow-hidden shrink-0 bg-accent/20">
          <img
            :src="getAvatarSrc(entry)"
            class="w-full h-full object-cover"
          >
        </div>

        <div class="flex-1 min-w-0">
          <div class="font-medium truncate">
            {{ entry.firstName }} {{ entry.lastName }}
          </div>
          <div
            v-if="entry.tg"
            class="text-sm text-muted-foreground truncate"
          >
            {{ entry.tg }}
          </div>
        </div>

        <div class="flex items-center gap-1.5 shrink-0">
          <Trophy
            class="h-4 w-4"
            :class="index < 3 ? 'text-yellow-500' : 'text-muted-foreground'"
          />
          <span
            class="font-bold"
            :class="index < 3 ? 'text-yellow-500' : 'text-foreground'"
          >
            {{ entry.total }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
