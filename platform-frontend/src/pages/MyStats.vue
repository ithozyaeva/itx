<script setup lang="ts">
import type { ProfileStats } from '@/models/profileStats'
import { Typography } from 'itx-ui-kit'
import {
  Calendar,
  CheckCircle,
  ClipboardList,
  Heart,
  Loader2,
  MessageSquare,
  Mic,
  Share2,
  Star,
  TrendingUp,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { handleError } from '@/services/errorService'
import { profileStatsService } from '@/services/profileStats'

const stats = ref<ProfileStats | null>(null)
const isLoading = ref(true)

async function fetchStats() {
  isLoading.value = true
  try {
    stats.value = await profileStatsService.getMyStats()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

const statCards = computed(() => {
  if (!stats.value)
    return []
  return [
    { label: 'Баллы', value: stats.value.pointsBalance, icon: Star, color: 'text-yellow-500' },
    { label: 'Посещено событий', value: stats.value.eventsAttended, icon: Calendar, color: 'text-blue-500' },
    { label: 'Проведено событий', value: stats.value.eventsHosted, icon: Mic, color: 'text-purple-500' },
    { label: 'Отзывов', value: stats.value.reviewsCount, icon: MessageSquare, color: 'text-green-500' },
    { label: 'Рефералов', value: stats.value.referralsCount, icon: Share2, color: 'text-orange-500' },
    { label: 'Благодарностей получено', value: stats.value.kudosReceived, icon: Heart, color: 'text-red-500' },
    { label: 'Благодарностей отправлено', value: stats.value.kudosSent, icon: Heart, color: 'text-pink-400' },
    { label: 'Заданий создано', value: stats.value.tasksCreated, icon: ClipboardList, color: 'text-cyan-500' },
    { label: 'Заданий выполнено', value: stats.value.tasksDone, icon: CheckCircle, color: 'text-emerald-500' },
  ]
})

const maxHistoryValue = computed(() => {
  if (!stats.value?.pointsHistory?.length)
    return 1
  return Math.max(...stats.value.pointsHistory.map(h => h.total), 1)
})

function formatMonth(m: string) {
  const [_, month] = m.split('-')
  const months = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']
  return months[Number.parseInt(month) - 1] || m
}

function memberSinceFormatted() {
  if (!stats.value?.memberSince)
    return ''
  return new Date(stats.value.memberSince).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <Typography
      variant="h2"
      as="h1"
      class="mb-6"
    >
      Моя статистика
    </Typography>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else-if="stats">
      <p
        v-if="stats.memberSince"
        class="text-sm text-muted-foreground mb-6"
      >
        Участник с {{ memberSinceFormatted() }}
      </p>

      <!-- Stat grid -->
      <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 mb-8">
        <div
          v-for="card in statCards"
          :key="card.label"
          class="rounded-2xl border bg-card border-border p-4"
        >
          <div class="flex items-center gap-2 mb-2">
            <component
              :is="card.icon"
              class="h-4 w-4"
              :class="card.color"
            />
            <span class="text-xs text-muted-foreground">{{ card.label }}</span>
          </div>
          <p class="text-2xl font-bold tabular-nums">
            {{ card.value }}
          </p>
        </div>
      </div>

      <!-- Points history chart -->
      <div
        v-if="stats.pointsHistory && stats.pointsHistory.length > 0"
        class="rounded-2xl border bg-card border-border p-4"
      >
        <div class="flex items-center gap-2 mb-4">
          <TrendingUp class="h-4 w-4 text-primary" />
          <h3 class="font-semibold text-sm">
            Баллы по месяцам
          </h3>
        </div>
        <div class="flex items-end gap-2 h-32">
          <div
            v-for="month in stats.pointsHistory"
            :key="month.month"
            class="flex-1 flex flex-col items-center gap-1"
          >
            <span class="text-xs font-medium tabular-nums">{{ month.total }}</span>
            <div
              class="w-full rounded-t-md bg-primary/80 transition-all min-h-1"
              :style="{ height: `${(month.total / maxHistoryValue) * 100}%` }"
            />
            <span class="text-xs text-muted-foreground">{{ formatMonth(month.month) }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
