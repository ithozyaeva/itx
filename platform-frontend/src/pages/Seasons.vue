<script setup lang="ts">
import type { Season, SeasonWithLeaderboard } from '@/models/season'
import { Calendar, Crown, Loader2, Medal, Trophy } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { Typography } from '@/components/ui/typography'
import { displayName, formatShortDate } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { seasonService } from '@/services/seasons'

const activeSeason = ref<SeasonWithLeaderboard | null>(null)
const allSeasons = ref<Season[]>([])
const selectedSeason = ref<SeasonWithLeaderboard | null>(null)
const isLoading = ref(true)
const isLoadingSeason = ref(false)

const rankIcons = [Crown, Trophy, Medal]
const rankColors = ['text-yellow-500', 'text-zinc-400', 'text-amber-700']

async function fetchActive() {
  isLoading.value = true
  try {
    activeSeason.value = await seasonService.getActive()
    selectedSeason.value = activeSeason.value
  }
  catch {
    activeSeason.value = null
  }
  try {
    allSeasons.value = await seasonService.getAll()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

async function selectSeason(id: number) {
  isLoadingSeason.value = true
  try {
    selectedSeason.value = await seasonService.getLeaderboard(id, 50)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoadingSeason.value = false
  }
}

function formatDate(d: string) {
  return formatShortDate(d)
}

function getAvatarSrc(entry: SeasonWithLeaderboard['leaderboard'][number]) {
  if (entry.avatarUrl)
    return entry.avatarUrl
  if (entry.tg)
    return `https://t.me/i/userpic/160/${entry.tg}.jpg`
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(entry.firstName || '?')}&background=random`
}

onMounted(() => {
  fetchActive()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <Typography
      variant="h2"
      as="h1"
      class="mb-6"
    >
      Сезоны
    </Typography>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <!-- Season tabs -->
      <div
        v-if="allSeasons.length > 0"
        class="flex gap-2 mb-6 flex-wrap"
      >
        <button
          v-for="season in allSeasons"
          :key="season.id"
          class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
          :class="selectedSeason?.season.id === season.id
            ? 'bg-primary text-primary-foreground'
            : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
          @click="selectSeason(season.id)"
        >
          {{ season.title }}
          <span
            v-if="season.status === 'ACTIVE'"
            class="ml-1 inline-block h-1.5 w-1.5 rounded-full bg-green-500"
          />
        </button>
      </div>

      <div
        v-if="!selectedSeason"
        class="text-center py-12 text-muted-foreground"
      >
        <Calendar class="h-12 w-12 mx-auto mb-3 opacity-50" />
        <p>Нет активных сезонов</p>
      </div>

      <template v-else>
        <!-- Season info -->
        <div class="rounded-2xl border bg-card border-border p-4 mb-6">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="font-semibold text-lg">
                {{ selectedSeason.season.title }}
              </h3>
              <p class="text-sm text-muted-foreground">
                {{ formatDate(selectedSeason.season.startDate) }} — {{ formatDate(selectedSeason.season.endDate) }}
              </p>
            </div>
            <span
              class="px-2.5 py-1 rounded-full text-xs font-medium"
              :class="selectedSeason.season.status === 'ACTIVE'
                ? 'bg-green-500/10 text-green-500'
                : 'bg-muted text-muted-foreground'"
            >
              {{ selectedSeason.season.status === 'ACTIVE' ? 'Активный' : 'Завершён' }}
            </span>
          </div>
        </div>

        <div
          v-if="isLoadingSeason"
          class="flex justify-center py-8"
        >
          <Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
        </div>

        <!-- Leaderboard -->
        <div
          v-else
          class="space-y-2"
        >
          <div
            v-for="entry in selectedSeason.leaderboard"
            :key="entry.memberId"
            class="flex items-center gap-3 rounded-xl border bg-card border-border p-3"
          >
            <div class="w-8 text-center shrink-0">
              <component
                :is="rankIcons[entry.rank - 1]"
                v-if="entry.rank <= 3"
                class="h-5 w-5 mx-auto"
                :class="rankColors[entry.rank - 1]"
              />
              <span
                v-else
                class="text-sm text-muted-foreground font-medium"
              >{{ entry.rank }}</span>
            </div>
            <div class="w-8 h-8 rounded-full overflow-hidden shrink-0 bg-accent/20">
              <img
                :src="getAvatarSrc(entry)"
                :alt="displayName(entry.firstName, entry.lastName)"
                :style="{ opacity: 0 }"
                loading="lazy"
                class="w-full h-full object-cover transition-opacity duration-300"
                @load="($event.target as HTMLImageElement).style.opacity = '1'"
                @error="(e: Event) => { const img = e.target as HTMLImageElement; if (img.dataset.fallback) return; img.dataset.fallback = '1'; img.src = `https://ui-avatars.com/api/?name=${encodeURIComponent(entry.firstName || '?')}&background=random`; img.style.opacity = '1' }"
              >
            </div>
            <router-link
              :to="`/members/${entry.memberId}`"
              class="flex-1 text-sm font-medium hover:underline truncate"
            >
              {{ displayName(entry.firstName, entry.lastName) }}
            </router-link>
            <span class="text-sm font-bold tabular-nums">{{ entry.total }}</span>
          </div>

          <div
            v-if="selectedSeason.leaderboard.length === 0"
            class="text-center py-8 text-muted-foreground"
          >
            <p>Пока нет данных</p>
          </div>
        </div>
      </template>
    </template>
  </div>
</template>
