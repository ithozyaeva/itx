<script setup lang="ts">
import type { UserAchievement } from '@/models/achievement'
import type { CommunityEvent } from '@/models/event'
import type { PointsSummary } from '@/models/points'
import type { TaskExchange } from '@/models/taskExchange'
import type { ChatQuestWithProgress } from '@/services/chatQuestService'
import { Typography } from 'itx-ui-kit'
import {
  ArrowRight,
  Award,
  Calendar,
  CheckCircle,
  ClipboardList,
  Flame,
  Loader2,
  MessageCircle,
  Radio,
  Star,
  Trophy,
  Users,
  Zap,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import EventCard from '@/components/events/EventCard.vue'
import { useCardReveal } from '@/composables/useCardReveal'
import { useUser, useUserLevel } from '@/composables/useUser'
import { dateFormatter } from '@/lib/utils'
import { SUBSCRIPTION_LEVELS } from '@/models/profile'
import { achievementsService } from '@/services/achievements'
import { chatQuestService } from '@/services/chatQuestService'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'
import { pointsService } from '@/services/points'
import { taskExchangeService } from '@/services/taskExchange'

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const user = useUser()
const { level, levelIndex } = useUserLevel()

const isLoading = ref(true)
const nearestEvent = ref<CommunityEvent | null>(null)
const upcomingEvents = ref<CommunityEvent[]>([])
const pointsSummary = ref<PointsSummary | null>(null)
const chatQuests = ref<ChatQuestWithProgress[]>([])
const openTasks = ref<TaskExchange[]>([])
const achievements = ref<{ total: number, unlocked: number, recent: UserAchievement[] }>({ total: 0, unlocked: 0, recent: [] })

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
  created.setHours(0, 0, 0, 0)
  now.setHours(0, 0, 0, 0)
  const diff = Math.floor((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24))
  return diff + 1
})

const timeGreeting = computed(() => {
  const hour = new Date().getHours()
  if (hour >= 5 && hour < 12)
    return 'Доброе утро'
  if (hour >= 12 && hour < 18)
    return 'Добрый день'
  return 'Добрый вечер'
})

const activeQuests = computed(() => chatQuests.value.filter(q => !q.completed))
const completedQuests = computed(() => chatQuests.value.filter(q => q.completed))

function questProgress(quest: ChatQuestWithProgress) {
  return Math.min(100, Math.round((quest.currentCount / quest.targetCount) * 100))
}

function formatQuestDeadline(dateStr: string) {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = date.getTime() - now.getTime()
  const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
  if (diffDays <= 0)
    return 'Истекает'
  if (diffDays === 1)
    return '1 день'
  if (diffDays <= 7)
    return `${diffDays} дн.`
  return date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
}

const isNearestEventLive = computed(() => {
  if (!nearestEvent.value)
    return false
  const now = new Date()
  const eventDate = new Date(nearestEvent.value.date)
  const diffMs = now.getTime() - eventDate.getTime()
  return diffMs >= 0 && diffMs < 2 * 60 * 60 * 1000
})

const nearestEventDate = computed(() =>
  nearestEvent.value ? dateFormatter.format(new Date(nearestEvent.value.date)) : '',
)

const isMemberOfNearest = computed(() =>
  user.value && nearestEvent.value
    ? nearestEvent.value.members.some(m => m.id === user.value!.id)
    : false,
)

const isHostOfNearest = computed(() =>
  user.value && nearestEvent.value
    ? nearestEvent.value.hosts.some(h => h.id === user.value!.id)
    : false,
)

onMounted(async () => {
  try {
    const results = await Promise.allSettled([
      eventsService.searchNext(1, 0),
      eventsService.searchNext(4, 1),
      pointsService.getMyPoints(),
      chatQuestService.getActiveQuests(),
      taskExchangeService.getAll({ status: 'OPEN', limit: 3 }),
      achievementsService.getMyAchievements(),
    ])

    if (results[0].status === 'fulfilled')
      nearestEvent.value = results[0].value.items[0] ?? null
    if (results[1].status === 'fulfilled')
      upcomingEvents.value = results[1].value.items
    if (results[2].status === 'fulfilled')
      pointsSummary.value = results[2].value
    if (results[3].status === 'fulfilled')
      chatQuests.value = results[3].value
    if (results[4].status === 'fulfilled')
      openTasks.value = results[4].value.items
    if (results[5].status === 'fulfilled') {
      const a = results[5].value
      achievements.value = {
        total: a.totalCount,
        unlocked: a.unlockedCount,
        recent: a.items.filter(i => i.unlocked).slice(0, 3),
      }
    }
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div
    ref="containerRef"
    class="px-4 py-6 md:py-8 max-w-4xl mx-auto"
  >
    <!-- Loading -->
    <div
      v-if="isLoading"
      class="flex justify-center py-24"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <!-- Hero Greeting -->
      <div
        class="relative rounded-3xl border bg-card p-6 md:p-8 overflow-hidden"
        data-reveal
      >
        <!-- Subtle accent gradient -->
        <div class="absolute top-0 right-0 w-48 h-48 rounded-full bg-accent/5 -translate-y-1/2 translate-x-1/4 blur-2xl" />
        <div class="relative">
          <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
            <div>
              <h1 class="text-2xl md:text-3xl font-bold tracking-tight">
                {{ timeGreeting }}, {{ user?.firstName }}
              </h1>
              <p class="text-muted-foreground mt-1.5">
                Ты в IT-Хозяевах уже <span class="font-medium text-foreground">{{ daysSinceJoined }}</span> {{ pluralizeDays(daysSinceJoined) }}
              </p>
            </div>

            <!-- Level badge -->
            <div class="flex items-center gap-3 shrink-0">
              <div class="text-right hidden sm:block">
                <p class="text-sm font-medium">
                  {{ level }}
                </p>
                <p class="text-xs text-muted-foreground">
                  Уровень {{ levelIndex + 1 }}/{{ SUBSCRIPTION_LEVELS.length }}
                </p>
              </div>
              <div class="flex gap-1">
                <span
                  v-for="(lvl, i) in SUBSCRIPTION_LEVELS"
                  :key="lvl"
                  class="h-2 w-5 rounded-full transition-colors"
                  :class="i <= levelIndex ? 'bg-accent' : 'bg-muted'"
                />
              </div>
            </div>
          </div>

          <!-- Quick stats -->
          <div class="grid grid-cols-3 gap-3 mt-6">
            <RouterLink
              to="/points"
              class="flex items-center gap-2.5 rounded-2xl bg-muted/50 hover:bg-muted/80 transition-colors p-3"
            >
              <div class="flex items-center justify-center w-9 h-9 rounded-xl bg-yellow-500/10">
                <Star class="h-4 w-4 text-yellow-500" />
              </div>
              <div class="min-w-0">
                <p class="text-lg font-bold leading-none">
                  {{ pointsSummary?.balance ?? 0 }}
                </p>
                <p class="text-[11px] text-muted-foreground mt-0.5">
                  баллов
                </p>
              </div>
            </RouterLink>

            <RouterLink
              to="/achievements"
              class="flex items-center gap-2.5 rounded-2xl bg-muted/50 hover:bg-muted/80 transition-colors p-3"
            >
              <div class="flex items-center justify-center w-9 h-9 rounded-xl bg-purple-500/10">
                <Trophy class="h-4 w-4 text-purple-500" />
              </div>
              <div class="min-w-0">
                <p class="text-lg font-bold leading-none">
                  {{ achievements.unlocked }}<span class="text-sm font-normal text-muted-foreground">/{{ achievements.total }}</span>
                </p>
                <p class="text-[11px] text-muted-foreground mt-0.5">
                  достижений
                </p>
              </div>
            </RouterLink>

            <RouterLink
              to="/events"
              class="flex items-center gap-2.5 rounded-2xl bg-muted/50 hover:bg-muted/80 transition-colors p-3"
            >
              <div class="flex items-center justify-center w-9 h-9 rounded-xl bg-blue-500/10">
                <Calendar class="h-4 w-4 text-blue-500" />
              </div>
              <div class="min-w-0">
                <p class="text-lg font-bold leading-none">
                  {{ daysSinceJoined }}
                </p>
                <p class="text-[11px] text-muted-foreground mt-0.5">
                  {{ pluralizeDays(daysSinceJoined) }} с нами
                </p>
              </div>
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Nearest Event -->
      <div
        v-if="nearestEvent"
        class="mt-5 rounded-3xl border bg-card overflow-hidden"
        data-reveal
      >
        <div class="p-5 md:p-6">
          <div class="flex items-center gap-2 mb-3">
            <Zap class="h-4 w-4 text-accent" />
            <span class="text-sm font-medium text-accent">Ближайшее событие</span>
            <span
              v-if="isNearestEventLive"
              class="inline-flex items-center gap-1 rounded-full bg-red-500 px-2 py-0.5 text-xs font-medium text-white ml-auto"
            >
              <Radio class="h-3 w-3" />
              LIVE
            </span>
          </div>

          <h2 class="text-xl font-bold">
            {{ nearestEvent.title }}
          </h2>
          <p
            v-if="nearestEvent.description"
            class="text-sm text-muted-foreground mt-1.5 line-clamp-2"
          >
            {{ nearestEvent.description }}
          </p>

          <div class="flex flex-wrap items-center gap-x-4 gap-y-1.5 mt-3 text-sm text-muted-foreground">
            <span class="flex items-center gap-1.5">
              <Calendar class="h-3.5 w-3.5" />
              {{ nearestEventDate }}
            </span>
            <span
              v-if="nearestEvent.hosts.length"
              class="flex items-center gap-1.5"
            >
              <Users class="h-3.5 w-3.5" />
              {{ nearestEvent.hosts.map(h => h.firstName).join(', ') }}
            </span>
            <span class="flex items-center gap-1.5">
              {{ nearestEvent.members.length }}{{ nearestEvent.maxParticipants > 0 ? `/${nearestEvent.maxParticipants}` : '' }} участников
            </span>
          </div>

          <div class="flex items-center gap-3 mt-4">
            <RouterLink
              v-if="isMemberOfNearest || isHostOfNearest"
              to="/events"
              class="inline-flex items-center gap-1.5 rounded-xl bg-accent/10 text-accent px-4 py-2 text-sm font-medium"
            >
              <CheckCircle class="h-4 w-4" />
              {{ isHostOfNearest ? 'Вы ведущий' : 'Вы записаны' }}
            </RouterLink>
            <RouterLink
              v-else
              to="/events"
              class="inline-flex items-center gap-1.5 rounded-xl bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 transition-colors"
            >
              Записаться
              <ArrowRight class="h-3.5 w-3.5" />
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Active Quests & Tasks Grid -->
      <div
        v-if="activeQuests.length > 0 || openTasks.length > 0"
        class="mt-5 grid gap-5 lg:grid-cols-2"
      >
        <!-- Chat Quests -->
        <div
          v-if="activeQuests.length > 0"
          class="rounded-3xl border bg-card p-5"
          data-reveal
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <Flame class="h-4 w-4 text-orange-500" />
              <span class="text-sm font-semibold">Задания в чатах</span>
            </div>
            <RouterLink
              to="/points"
              class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              Все →
            </RouterLink>
          </div>

          <div class="space-y-3">
            <div
              v-for="quest in activeQuests.slice(0, 3)"
              :key="quest.id"
              class="rounded-2xl bg-muted/40 p-3.5"
            >
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-center gap-2.5 min-w-0">
                  <div class="flex items-center justify-center w-8 h-8 rounded-lg bg-orange-500/10 shrink-0">
                    <MessageCircle class="h-4 w-4 text-orange-500" />
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-medium truncate">
                      {{ quest.title }}
                    </p>
                    <p class="text-xs text-muted-foreground">
                      {{ formatQuestDeadline(quest.endsAt) }}
                    </p>
                  </div>
                </div>
                <span class="text-xs font-bold text-yellow-500 bg-yellow-500/10 px-2 py-0.5 rounded-full shrink-0">
                  +{{ quest.pointsReward }}
                </span>
              </div>
              <div class="mt-2.5">
                <div class="flex items-center justify-between text-xs text-muted-foreground mb-1">
                  <span>{{ quest.currentCount }} / {{ quest.targetCount }}</span>
                  <span>{{ questProgress(quest) }}%</span>
                </div>
                <div class="w-full bg-muted rounded-full h-1.5">
                  <div
                    class="bg-accent rounded-full h-1.5 transition-all"
                    :style="{ width: `${questProgress(quest)}%` }"
                  />
                </div>
              </div>
            </div>

            <div
              v-if="completedQuests.length > 0"
              class="text-xs text-muted-foreground text-center pt-1"
            >
              <CheckCircle class="h-3 w-3 inline mr-1" />
              {{ completedQuests.length }} выполнено
            </div>
          </div>
        </div>

        <!-- Open Tasks -->
        <div
          v-if="openTasks.length > 0"
          class="rounded-3xl border bg-card p-5"
          data-reveal
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <ClipboardList class="h-4 w-4 text-blue-500" />
              <span class="text-sm font-semibold">Биржа заданий</span>
            </div>
            <RouterLink
              to="/tasks"
              class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              Все →
            </RouterLink>
          </div>

          <div class="space-y-3">
            <RouterLink
              v-for="task in openTasks"
              :key="task.id"
              to="/tasks"
              class="block rounded-2xl bg-muted/40 hover:bg-muted/60 transition-colors p-3.5"
            >
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0">
                  <p class="text-sm font-medium truncate">
                    {{ task.title }}
                  </p>
                  <p class="text-xs text-muted-foreground mt-0.5 line-clamp-1">
                    {{ task.description }}
                  </p>
                </div>
                <span class="text-xs font-medium text-blue-500 bg-blue-500/10 px-2 py-0.5 rounded-full shrink-0">
                  Открыто
                </span>
              </div>
              <div class="flex items-center gap-2 mt-2 text-xs text-muted-foreground">
                <span>от {{ task.creator.firstName }}</span>
              </div>
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Achievements Preview -->
      <div
        v-if="achievements.recent.length > 0"
        class="mt-5 rounded-3xl border bg-card p-5"
        data-reveal
      >
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <Award class="h-4 w-4 text-purple-500" />
            <span class="text-sm font-semibold">Последние достижения</span>
          </div>
          <RouterLink
            to="/achievements"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
          >
            Все →
          </RouterLink>
        </div>
        <div class="flex gap-3 overflow-x-auto pb-1">
          <div
            v-for="ach in achievements.recent"
            :key="ach.id"
            class="flex items-center gap-2.5 rounded-2xl bg-muted/40 px-4 py-3 shrink-0"
          >
            <span class="text-2xl">{{ ach.icon }}</span>
            <div>
              <p class="text-sm font-medium">
                {{ ach.title }}
              </p>
              <p class="text-xs text-muted-foreground">
                {{ ach.description }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Upcoming Events -->
      <div
        v-if="upcomingEvents.length > 0"
        class="mt-5"
        data-reveal
      >
        <div class="flex items-center justify-between mb-4">
          <Typography
            variant="h4"
            as="h2"
          >
            Предстоящие события
          </Typography>
          <RouterLink
            to="/events"
            class="text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            Все события →
          </RouterLink>
        </div>
        <div class="space-y-4">
          <EventCard
            v-for="event in upcomingEvents"
            :key="event.id"
            :event="event"
          />
        </div>
      </div>

      <!-- Quick Links -->
      <div
        class="mt-5 grid grid-cols-2 sm:grid-cols-4 gap-3"
        data-reveal
      >
        <RouterLink
          to="/leaderboard"
          class="rounded-2xl border bg-card hover:bg-muted/50 transition-colors p-4 text-center"
        >
          <Trophy class="h-5 w-5 mx-auto text-yellow-500 mb-1.5" />
          <p class="text-xs font-medium">
            Лидерборд
          </p>
        </RouterLink>
        <RouterLink
          to="/marketplace"
          class="rounded-2xl border bg-card hover:bg-muted/50 transition-colors p-4 text-center"
        >
          <Star class="h-5 w-5 mx-auto text-accent mb-1.5" />
          <p class="text-xs font-medium">
            Барахолка
          </p>
        </RouterLink>
        <RouterLink
          to="/mentors"
          class="rounded-2xl border bg-card hover:bg-muted/50 transition-colors p-4 text-center"
        >
          <Users class="h-5 w-5 mx-auto text-blue-500 mb-1.5" />
          <p class="text-xs font-medium">
            Менторы
          </p>
        </RouterLink>
        <RouterLink
          to="/referals"
          class="rounded-2xl border bg-card hover:bg-muted/50 transition-colors p-4 text-center"
        >
          <Zap class="h-5 w-5 mx-auto text-orange-500 mb-1.5" />
          <p class="text-xs font-medium">
            Рефералы
          </p>
        </RouterLink>
      </div>
    </template>
  </div>
</template>
