<script setup lang="ts">
import type { ProfileStats } from '@/models/profileStats'
import {
  ArrowDown,
  ArrowUp,
  Calendar,
  CheckCircle,
  ClipboardList,
  Heart,
  MessageSquare,
  Mic,
  Minus,
  Share2,
  Star,
  TrendingUp,
  Trophy,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import ErrorState from '@/components/common/ErrorState.vue'
import { Skeleton } from '@/components/ui/skeleton'
import { Typography } from '@/components/ui/typography'
import { useUser } from '@/composables/useUser'
import { reasonLabels } from '@/lib/reasonLabels'
import { formatShortDate } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'
import { profileStatsService } from '@/services/profileStats'

const user = useUser()

const stats = ref<ProfileStats | null>(null)
const isLoading = ref(true)
const loadError = ref<string | null>(null)
const leaderboardPosition = ref<number | null>(null)

async function fetchStats() {
  isLoading.value = true
  loadError.value = null
  try {
    const [statsData, leaderboard] = await Promise.allSettled([
      profileStatsService.getMyStats(),
      pointsService.getLeaderboard(),
    ])

    if (statsData.status === 'fulfilled')
      stats.value = statsData.value
    else
      throw statsData.reason

    if (leaderboard.status === 'fulfilled' && user.value) {
      const entries = leaderboard.value.items ?? []
      const userId = user.value.id
      const idx = entries.findIndex(e => e.memberId === userId)
      leaderboardPosition.value = idx >= 0 ? idx + 1 : null
    }
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

const statCards = computed(() => {
  if (!stats.value)
    return []
  return [
    { label: 'Баллы', value: stats.value.pointsBalance, icon: Star, color: 'text-yellow-500', bg: 'bg-yellow-500/10' },
    { label: 'Посещено событий', value: stats.value.eventsAttended, icon: Calendar, color: 'text-blue-500', bg: 'bg-blue-500/10' },
    { label: 'Проведено событий', value: stats.value.eventsHosted, icon: Mic, color: 'text-purple-500', bg: 'bg-purple-500/10' },
    { label: 'Отзывов', value: stats.value.reviewsCount, icon: MessageSquare, color: 'text-green-500', bg: 'bg-green-500/10' },
    { label: 'Рефералов', value: stats.value.referralsCount, icon: Share2, color: 'text-orange-500', bg: 'bg-orange-500/10' },
    { label: 'Благодарностей получено', value: stats.value.kudosReceived, icon: Heart, color: 'text-red-500', bg: 'bg-red-500/10' },
    { label: 'Благодарностей отправлено', value: stats.value.kudosSent, icon: Heart, color: 'text-pink-400', bg: 'bg-pink-400/10' },
    { label: 'Заданий создано', value: stats.value.tasksCreated, icon: ClipboardList, color: 'text-cyan-500', bg: 'bg-cyan-500/10' },
    { label: 'Заданий выполнено', value: stats.value.tasksDone, icon: CheckCircle, color: 'text-emerald-500', bg: 'bg-emerald-500/10' },
  ]
})

const pointsTrend = computed(() => {
  const history = stats.value?.pointsHistory
  if (!history || history.length < 2)
    return null

  const current = history.at(-1)?.total ?? 0
  const previous = history[history.length - 2]?.total ?? 0

  if (previous === 0)
    return null

  const change = ((current - previous) / previous) * 100
  return {
    value: Math.abs(Math.round(change)),
    direction: change > 0 ? 'up' : change < 0 ? 'down' : 'same',
    current,
    previous,
  }
})

const maxHistoryValue = computed(() => {
  if (!stats.value?.pointsHistory?.length)
    return 1
  return Math.max(...stats.value.pointsHistory.map(h => h.total), 1)
})

const activitySummary = computed(() => {
  if (!stats.value)
    return null

  const total = stats.value.eventsAttended
    + stats.value.eventsHosted
    + stats.value.reviewsCount
    + stats.value.kudosSent
    + stats.value.tasksCreated
    + stats.value.tasksDone

  return {
    totalActions: total,
    mostActive: getMostActiveArea(),
  }
})

function getMostActiveArea(): string {
  if (!stats.value)
    return ''

  const areas = [
    { label: 'события', value: stats.value.eventsAttended + stats.value.eventsHosted },
    { label: 'задания', value: stats.value.tasksCreated + stats.value.tasksDone },
    { label: 'благодарности', value: stats.value.kudosSent + stats.value.kudosReceived },
    { label: 'отзывы', value: stats.value.reviewsCount },
    { label: 'рефералы', value: stats.value.referralsCount },
  ]

  const sorted = [...areas].sort((a, b) => b.value - a.value)
  return sorted.length > 0 && sorted[0].value > 0 ? sorted[0].label : 'нет активности'
}

function formatMonth(m: string) {
  const parts = m.split('-')
  const month = parts[1]
  if (!month)
    return m
  const months = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']
  return months[Number.parseInt(month) - 1] || m
}

function memberSinceFormatted() {
  if (!stats.value?.memberSince)
    return ''
  return formatShortDate(stats.value.memberSince)
}

function daysSinceMember(): number {
  if (!stats.value?.memberSince)
    return 0
  const diff = Date.now() - new Date(stats.value.memberSince).getTime()
  return Math.floor(diff / (1000 * 60 * 60 * 24))
}

const topSources = computed(() => {
  if (!stats.value?.pointsBySource?.length)
    return []
  const positive = stats.value.pointsBySource.filter(s => s.total > 0)
  const totalPositive = positive.reduce((sum, s) => sum + s.total, 0)
  return positive.slice(0, 5).map(s => ({
    reason: s.reason,
    label: reasonLabels[s.reason] || s.reason,
    total: s.total,
    percent: totalPositive > 0 ? Math.round((s.total / totalPositive) * 100) : 0,
  }))
})

const maxSourceValue = computed(() => {
  if (!topSources.value.length)
    return 1
  return Math.max(...topSources.value.map(s => s.total), 1)
})

// Contribution graph
const contributionData = computed(() => {
  const weeks: { date: string, count: number, dayOfWeek: number }[][] = []
  const activityMap = new Map<string, number>()

  for (const day of stats.value?.activityHistory ?? []) {
    activityMap.set(day.date, day.count)
  }

  const today = new Date()
  const todayTime = today.getTime()
  const startDate = new Date(today)
  startDate.setDate(startDate.getDate() - 83)
  startDate.setDate(startDate.getDate() - ((startDate.getDay() + 6) % 7))

  let currentWeek: { date: string, count: number, dayOfWeek: number }[] = []
  const d = new Date(startDate)
  let maxCount = 0
  let activeDays = 0

  while (d.getTime() <= todayTime) {
    const dateStr = d.toISOString().slice(0, 10)
    const dayOfWeek = (d.getDay() + 6) % 7
    const count = activityMap.get(dateStr) ?? 0
    if (count > maxCount)
      maxCount = count
    if (count > 0)
      activeDays++
    currentWeek.push({ date: dateStr, count, dayOfWeek })
    if (dayOfWeek === 6) {
      weeks.push([...currentWeek])
      currentWeek = []
    }
    d.setDate(d.getDate() + 1)
  }

  if (currentWeek.length > 0)
    weeks.push(currentWeek)

  return { weeks, maxCount: Math.max(maxCount, 1), activeDays }
})

const contributionWeeks = computed(() => contributionData.value.weeks)
const maxActivityCount = computed(() => contributionData.value.maxCount)

function activityLevel(count: number): number {
  if (count === 0)
    return 0
  const ratio = count / maxActivityCount.value
  if (ratio <= 0.25)
    return 1
  if (ratio <= 0.5)
    return 2
  if (ratio <= 0.75)
    return 3
  return 4
}

const activityLevelClasses: Record<number, string> = {
  0: 'bg-muted',
  1: 'bg-accent/20',
  2: 'bg-accent/40',
  3: 'bg-accent/70',
  4: 'bg-accent',
}

const totalActivityDays = computed(() => contributionData.value.activeDays)

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/my-stats
    </div>
    <Typography
      variant="h2"
      as="h1"
      class="mb-6"
    >
      Моя статистика
    </Typography>

    <!-- Skeleton loading state -->
    <div v-if="isLoading" class="space-y-6">
      <!-- Summary row skeleton -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div v-for="i in 3" :key="i" class="rounded-sm border bg-card border-border terminal-card p-4">
          <Skeleton class="h-3 w-20 rounded mb-2" />
          <Skeleton class="h-7 w-24 rounded-lg mb-1" />
          <Skeleton class="h-3 w-32 rounded" />
        </div>
      </div>
      <!-- Stat grid skeleton -->
      <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
        <div v-for="i in 9" :key="i" class="rounded-sm border bg-card border-border terminal-card p-4">
          <div class="flex items-center gap-2 mb-2">
            <Skeleton class="w-8 h-8 rounded-lg" />
            <Skeleton class="h-3 w-20 rounded" />
          </div>
          <Skeleton class="h-8 w-12 rounded-lg" />
        </div>
      </div>
      <!-- Chart skeleton -->
      <div class="rounded-sm border bg-card border-border terminal-card p-4">
        <Skeleton class="h-5 w-40 rounded mb-4" />
        <Skeleton class="h-40 w-full rounded-lg" />
      </div>
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="fetchStats"
    />

    <template v-else-if="stats">
      <!-- Summary row -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6">
        <div class="rounded-sm border bg-card border-border terminal-card p-4">
          <p class="text-xs text-muted-foreground mb-1">
            Участник с
          </p>
          <p class="text-sm font-semibold">
            {{ memberSinceFormatted() }}
          </p>
          <p class="text-xs text-muted-foreground mt-1">
            {{ daysSinceMember() }} дней в сообществе
          </p>
        </div>

        <div class="rounded-sm border bg-card border-border terminal-card p-4">
          <p class="text-xs text-muted-foreground mb-1">
            Общая активность
          </p>
          <p class="text-2xl font-bold tabular-nums">
            {{ activitySummary?.totalActions ?? 0 }}
          </p>
          <p class="text-xs text-muted-foreground mt-1">
            Больше всего: {{ activitySummary?.mostActive }}
          </p>
        </div>

        <div
          v-if="leaderboardPosition"
          class="rounded-sm border bg-card border-border terminal-card p-4"
        >
          <p class="text-xs text-muted-foreground mb-1">
            Место в рейтинге
          </p>
          <div class="flex items-center gap-2">
            <Trophy class="h-5 w-5 text-yellow-500" />
            <p class="text-2xl font-bold tabular-nums">
              #{{ leaderboardPosition }}
            </p>
          </div>
        </div>
      </div>

      <!-- Stat grid -->
      <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 mb-6">
        <div
          v-for="card in statCards"
          :key="card.label"
          class="rounded-sm border bg-card border-border terminal-card p-4"
        >
          <div class="flex items-center gap-2 mb-2">
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg"
              :class="card.bg"
            >
              <component
                :is="card.icon"
                class="h-4 w-4"
                :class="card.color"
              />
            </div>
            <span class="text-xs text-muted-foreground">{{ card.label }}</span>
          </div>
          <p class="text-2xl font-bold tabular-nums">
            {{ card.value }}
          </p>
        </div>
      </div>

      <!-- Points by source -->
      <div
        v-if="topSources.length > 0"
        class="rounded-sm border bg-card border-border terminal-card p-4 mb-6"
      >
        <div class="flex items-center gap-2 mb-4">
          <Star class="h-4 w-4 text-yellow-500" />
          <h3 class="font-semibold text-sm">
            Баллы по источникам
          </h3>
        </div>
        <div class="space-y-3">
          <div
            v-for="source in topSources"
            :key="source.reason"
            class="space-y-1"
          >
            <div class="flex items-center justify-between text-xs">
              <span class="text-muted-foreground">{{ source.label }}</span>
              <span class="font-medium tabular-nums">{{ source.total }} ({{ source.percent }}%)</span>
            </div>
            <div class="w-full h-2 rounded-full bg-muted overflow-hidden">
              <div
                class="h-full rounded-full bg-accent transition-all"
                :style="{ width: `${(source.total / maxSourceValue) * 100}%` }"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Achievements progress -->
      <div
        v-if="stats.achievementsTotal > 0"
        class="rounded-sm border bg-card border-border terminal-card p-4 mb-6"
      >
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <Trophy class="h-4 w-4 text-yellow-500" />
            <h3 class="font-semibold text-sm">
              Достижения
            </h3>
          </div>
          <RouterLink
            to="/achievements"
            class="text-xs text-primary hover:underline"
          >
            Все достижения
          </RouterLink>
        </div>
        <div class="flex items-center gap-3 mb-2">
          <p class="text-2xl font-bold tabular-nums">
            {{ stats.achievementsEarned }}
          </p>
          <span class="text-sm text-muted-foreground">из {{ stats.achievementsTotal }}</span>
        </div>
        <div class="w-full h-2.5 rounded-full bg-muted overflow-hidden">
          <div
            class="h-full rounded-full bg-yellow-500 transition-all"
            :style="{ width: `${(stats.achievementsEarned / stats.achievementsTotal) * 100}%` }"
          />
        </div>
      </div>

      <!-- Points history chart -->
      <div
        v-if="stats.pointsHistory && stats.pointsHistory.length > 0"
        class="rounded-sm border bg-card border-border terminal-card p-4 mb-6"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <TrendingUp class="h-4 w-4 text-primary" />
            <h3 class="font-semibold text-sm">
              Баллы по месяцам
            </h3>
          </div>
          <!-- Trend indicator -->
          <div
            v-if="pointsTrend"
            class="flex items-center gap-1 text-xs font-medium"
            :class="{
              'text-green-500': pointsTrend.direction === 'up',
              'text-red-500': pointsTrend.direction === 'down',
              'text-muted-foreground': pointsTrend.direction === 'same',
            }"
          >
            <ArrowUp
              v-if="pointsTrend.direction === 'up'"
              class="h-3.5 w-3.5"
            />
            <ArrowDown
              v-else-if="pointsTrend.direction === 'down'"
              class="h-3.5 w-3.5"
            />
            <Minus
              v-else
              class="h-3.5 w-3.5"
            />
            {{ pointsTrend.value }}% за месяц
          </div>
        </div>

        <!-- SVG area chart -->
        <div class="relative h-40">
          <svg
            class="w-full h-full"
            :viewBox="`0 0 ${stats.pointsHistory.length * 100} 160`"
            preserveAspectRatio="xMidYMid meet"
          >
            <!-- Area -->
            <path
              :d="(() => {
                const points = stats.pointsHistory.map((m, i) => {
                  const x = i * 100 + 50
                  const y = 150 - (m.total / maxHistoryValue) * 130
                  return `${x},${y}`
                })
                const first = stats.pointsHistory.length > 0 ? 50 : 0
                const last = (stats.pointsHistory.length - 1) * 100 + 50
                return `M${first},150 L${points.join(' L')} L${last},150 Z`
              })()"
              class="fill-primary/10"
            />
            <!-- Line -->
            <polyline
              :points="stats.pointsHistory.map((m, i) => {
                const x = i * 100 + 50
                const y = 150 - (m.total / maxHistoryValue) * 130
                return `${x},${y}`
              }).join(' ')"
              class="stroke-primary"
              fill="none"
              stroke-width="2.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <!-- Dots -->
            <circle
              v-for="(m, i) in stats.pointsHistory"
              :key="m.month"
              :cx="i * 100 + 50"
              :cy="150 - (m.total / maxHistoryValue) * 130"
              r="4"
              class="fill-primary"
            />
          </svg>
        </div>

        <!-- Labels -->
        <div class="flex justify-between mt-2 px-1">
          <div
            v-for="month in stats.pointsHistory"
            :key="month.month"
            class="flex flex-col items-center text-center flex-1"
          >
            <span class="text-xs font-medium tabular-nums">{{ month.total }}</span>
            <span class="text-[10px] text-muted-foreground">{{ formatMonth(month.month) }}</span>
          </div>
        </div>
      </div>

      <!-- Contribution graph -->
      <div
        v-if="contributionWeeks.length > 0"
        class="rounded-sm border bg-card border-border terminal-card p-4"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <Calendar class="h-4 w-4 text-primary" />
            <h3 class="font-semibold text-sm">
              Активность за 12 недель
            </h3>
          </div>
          <span class="text-xs text-muted-foreground">
            {{ totalActivityDays }} активных дней
          </span>
        </div>

        <div class="flex gap-1 overflow-x-auto pb-1">
          <div
            v-for="(week, wi) in contributionWeeks"
            :key="wi"
            class="flex flex-col gap-1"
          >
            <div
              v-for="day in week"
              :key="day.date"
              class="w-3 h-3 rounded-sm transition-colors"
              :class="activityLevelClasses[activityLevel(day.count)]"
              :title="`${day.date}: ${day.count} действий`"
            />
          </div>
        </div>

        <div class="flex items-center justify-end gap-1.5 mt-2">
          <span class="text-[10px] text-muted-foreground">Меньше</span>
          <div
            v-for="level in [0, 1, 2, 3, 4]"
            :key="level"
            class="w-3 h-3 rounded-sm"
            :class="activityLevelClasses[level]"
          />
          <span class="text-[10px] text-muted-foreground">Больше</span>
        </div>
      </div>
    </template>
  </div>
</template>
