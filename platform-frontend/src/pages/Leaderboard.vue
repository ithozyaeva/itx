<script setup lang="ts">
import type { LeaderboardEntry } from '@/models/points'
import { Loader2, Trophy } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import ErrorState from '@/components/common/ErrorState.vue'
import { Button } from '@/components/ui/button'
import { Typography } from '@/components/ui/typography'
import { useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'

const PAGE_SIZE = 50

const allEntries = ref<LeaderboardEntry[]>([])
const displayCount = ref(PAGE_SIZE)
const isLoading = ref(true)
const loadError = ref<string | null>(null)

const user = useUser()

const visibleEntries = computed(() => allEntries.value.slice(0, displayCount.value))
const hasMore = computed(() => displayCount.value < allEntries.value.length)

const currentUserEntry = computed(() => {
  if (!user.value)
    return null
  const userId = user.value.id
  const index = allEntries.value.findIndex(e => e.memberId === userId)
  if (index === -1)
    return null
  return { entry: allEntries.value[index], rank: index + 1 }
})

const isCurrentUserVisible = computed(() => {
  if (!currentUserEntry.value)
    return false
  return currentUserEntry.value.rank <= displayCount.value
})

async function fetchLeaderboard() {
  isLoading.value = true
  loadError.value = null
  try {
    const response = await pointsService.getLeaderboard(10000)
    allEntries.value = response.items || []
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

function showMore() {
  displayCount.value += PAGE_SIZE
}

onMounted(() => {
  fetchLeaderboard()
})

function getAvatarSrc(entry: LeaderboardEntry) {
  if (entry.avatarUrl)
    return entry.avatarUrl
  if (entry.tg)
    return `https://t.me/i/userpic/160/${entry.tg}.jpg`
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(entry.firstName || '?')}&background=random`
}

function isCurrentUser(entry: LeaderboardEntry) {
  return user.value ? entry.memberId === user.value.id : false
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/leaderboard
    </div>
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

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="fetchLeaderboard()"
    />

    <div
      v-else-if="allEntries.length === 0"
      class="text-center py-12 text-muted-foreground"
    >
      Пока нет данных о баллах. Участвуйте в событиях, чтобы зарабатывать баллы!
    </div>

    <div
      v-else
      class="space-y-3"
    >
      <div
        v-if="currentUserEntry && !isCurrentUserVisible"
        class="p-4 bg-primary/10 border border-primary/30 rounded-sm shadow-sm mb-4"
      >
        <RouterLink
          :to="`/members/${currentUserEntry.entry.memberId}`"
          class="flex items-center gap-4"
        >
          <div class="flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold shrink-0 bg-primary/20 text-primary">
            {{ currentUserEntry.rank }}
          </div>

          <div class="w-10 h-10 rounded-full overflow-hidden shrink-0 bg-accent/20">
            <img
              :src="getAvatarSrc(currentUserEntry.entry)"
              :alt="`${currentUserEntry.entry.firstName} ${currentUserEntry.entry.lastName}`"
              :style="{ opacity: 0 }"
              loading="lazy"
              class="w-full h-full object-cover transition-opacity duration-300"
              @load="($event.target as HTMLImageElement).style.opacity = '1'"
              @error="(e: Event) => { const img = e.target as HTMLImageElement; if (img.dataset.fallback) return; img.dataset.fallback = '1'; img.src = `https://ui-avatars.com/api/?name=${encodeURIComponent(currentUserEntry!.entry.firstName || '?')}&background=random`; img.style.opacity = '1' }"
            >
          </div>

          <div class="flex-1 min-w-0">
            <div class="text-sm text-primary font-medium mb-0.5">
              Ваша позиция
            </div>
            <div class="font-medium truncate">
              {{ currentUserEntry.entry.firstName }} {{ currentUserEntry.entry.lastName }}
            </div>
          </div>

          <div class="flex items-center gap-1.5 shrink-0">
            <Trophy class="h-4 w-4 text-primary" />
            <span class="font-bold text-primary">
              {{ currentUserEntry.entry.total }}
            </span>
          </div>
        </RouterLink>
      </div>

      <RouterLink
        v-for="(entry, index) in visibleEntries"
        :key="entry.memberId"
        :to="`/members/${entry.memberId}`"
        class="flex items-center gap-4 p-4 bg-card border border-border rounded-sm shadow-sm hover:border-primary/30 transition-colors"
        :class="{
          'border-yellow-500/50 bg-yellow-500/5': index < 3 && !isCurrentUser(entry),
          'border-primary/50 bg-primary/10 ring-1 ring-primary/20': isCurrentUser(entry),
        }"
      >
        <div
          class="flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold shrink-0"
          :class="isCurrentUser(entry)
            ? 'bg-primary/20 text-primary'
            : index < 3 ? 'bg-yellow-500/20 text-yellow-500' : 'bg-muted text-muted-foreground'"
        >
          {{ index + 1 }}
        </div>

        <div class="w-10 h-10 rounded-full overflow-hidden shrink-0 bg-accent/20">
          <img
            :src="getAvatarSrc(entry)"
            :alt="`${entry.firstName} ${entry.lastName}`"
            :style="{ opacity: 0 }"
            loading="lazy"
            class="w-full h-full object-cover transition-opacity duration-300"
            @load="($event.target as HTMLImageElement).style.opacity = '1'"
            @error="(e: Event) => { const img = e.target as HTMLImageElement; if (img.dataset.fallback) return; img.dataset.fallback = '1'; img.src = `https://ui-avatars.com/api/?name=${encodeURIComponent(entry.firstName || '?')}&background=random`; img.style.opacity = '1' }"
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
            :class="isCurrentUser(entry)
              ? 'text-primary'
              : index < 3 ? 'text-yellow-500' : 'text-muted-foreground'"
          />
          <span
            class="font-bold"
            :class="isCurrentUser(entry)
              ? 'text-primary'
              : index < 3 ? 'text-yellow-500' : 'text-foreground'"
          >
            {{ entry.total }}
          </span>
        </div>
      </RouterLink>

      <div
        v-if="hasMore"
        class="mt-4 flex justify-center"
      >
        <Button
          variant="outline"
          @click="showMore"
        >
          Показать ещё
        </Button>
      </div>
    </div>
  </div>
</template>
