<script setup lang="ts">
import type { UserAchievement } from '@/models/achievement'
import type { CommunityEvent } from '@/models/event'
import type { ChatHighlight } from '@/models/highlight'
import type { PointsSummary } from '@/models/points'
import type { TaskExchange } from '@/models/taskExchange'
import type { ChatQuestWithProgress } from '@/services/chatQuestService'
import {
  ArrowRight,
  Award,
  BookOpen,
  Briefcase,
  Calendar,
  CalendarCheck,
  CheckCircle,
  ClipboardList,
  Crown,
  FileText,
  Flame,
  Footprints,
  GraduationCap,
  HardHat,
  History,
  ListChecks,
  Medal,
  MessageCircle,
  MessageSquare,
  MessagesSquare,
  Mic,
  Package,
  Presentation,
  Radio,
  Share2,
  ShoppingCart,
  Star,
  Swords,
  Target,
  Trophy,
  UserCheck,
  UserPlus,
  Users,
  Zap,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { getNextOccurrenceDate } from '@/composables/useEventOccurrence'
import { useSSE } from '@/composables/useSSE'
import { isUserSubscribed, useUser, useUserLevel } from '@/composables/useUser'
import { dateFormatter, formatShortDate } from '@/lib/utils'
import { SUBSCRIPTION_LEVELS } from '@/models/profile'
import { achievementsService } from '@/services/achievements'
import { chatQuestService } from '@/services/chatQuestService'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'
import { highlightsService } from '@/services/highlights'
import { pointsService } from '@/services/points'
import { taskExchangeService } from '@/services/taskExchange'

const achievementIconMap: Record<string, any> = {
  'footprints': Footprints,
  'flame': Flame,
  'calendar-check': CalendarCheck,
  'medal': Medal,
  'mic': Mic,
  'presentation': Presentation,
  'star': Star,
  'trophy': Trophy,
  'crown': Crown,
  'gem': Award,
  'message-square': MessageSquare,
  'messages-square': MessagesSquare,
  'book-open': BookOpen,
  'share-2': Share2,
  'users': Users,
  'user-plus': UserPlus,
  'user-check': UserCheck,
  'zap': Zap,
  'file-text': FileText,
  'package': Package,
  'shopping-cart': ShoppingCart,
  'clipboard-list': ClipboardList,
  'list-checks': ListChecks,
  'check-circle': CheckCircle,
  'hard-hat': HardHat,
  'briefcase': Briefcase,
  'target': Target,
  'swords': Swords,
  'history': History,
  'graduation-cap': GraduationCap,
}

const user = useUser()
const { level, levelIndex } = useUserLevel()
const isSubscribed = isUserSubscribed()

const isLoading = ref(true)
const nearestEvents = ref<CommunityEvent[]>([])
const pointsSummary = ref<PointsSummary | null>(null)
const chatQuests = ref<ChatQuestWithProgress[]>([])
const openTasks = ref<TaskExchange[]>([])
const achievements = ref<{ total: number, unlocked: number, recent: UserAchievement[] }>({ total: 0, unlocked: 0, recent: [] })
const highlights = ref<ChatHighlight[]>([])

async function refreshPoints() {
  try {
    pointsSummary.value = await pointsService.getMyPoints()
  }
  catch {
    // silent refresh
  }
}

useSSE('points', () => refreshPoints())

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
  if (!quest.targetCount)
    return 0
  return Math.min(100, Math.round((quest.currentCount / quest.targetCount) * 100))
}

function questDeadlineDays(dateStr: string) {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = date.getTime() - now.getTime()
  return Math.ceil(diffMs / (1000 * 60 * 60 * 24))
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
  return formatShortDate(date)
}

function isEventLive(event: CommunityEvent) {
  const now = new Date()
  const eventDate = getNextOccurrenceDate(event, now)
  const diffMs = now.getTime() - eventDate.getTime()
  return diffMs >= 0 && diffMs < 2 * 60 * 60 * 1000
}

function isMemberOfEvent(event: CommunityEvent) {
  return event.members.some(m => m.id === user.value?.id)
}

function isHostOfEvent(event: CommunityEvent) {
  return event.hosts.some(h => h.id === user.value?.id)
}

onMounted(async () => {
  try {
    const results = await Promise.allSettled([
      eventsService.searchNext(20, 0),
      pointsService.getMyPoints(),
      chatQuestService.getActiveQuests(),
      Promise.all([
        taskExchangeService.getAll({ status: 'OPEN', limit: 10 }),
        taskExchangeService.getAll({ status: 'IN_PROGRESS', limit: 10 }),
      ]),
      achievementsService.getMyAchievements(),
      highlightsService.getRecent(5),
    ])

    if (results[0].status === 'fulfilled') {
      const now = new Date()
      const items = results[0].value?.items ?? []
      nearestEvents.value = [...items]
        .sort((a, b) => getNextOccurrenceDate(a, now).getTime() - getNextOccurrenceDate(b, now).getTime())
        .slice(0, 3)
    }
    if (results[1].status === 'fulfilled')
      pointsSummary.value = results[1].value
    if (results[2].status === 'fulfilled')
      chatQuests.value = results[2].value ?? []
    if (results[3].status === 'fulfilled') {
      const [openRes, inProgressRes] = results[3].value
      const open = openRes?.items ?? []
      const inProgress = (inProgressRes?.items ?? []).filter(t => (t.assignees?.length ?? 0) < t.maxAssignees)
      openTasks.value = [...open, ...inProgress].slice(0, 3)
    }
    if (results[4].status === 'fulfilled') {
      const a = results[4].value
      achievements.value = {
        total: a?.totalCount ?? 0,
        unlocked: a?.unlockedCount ?? 0,
        recent: (a?.items ?? []).filter(i => i.unlocked).slice(0, 3),
      }
    }
    if (results[5].status === 'fulfilled')
      highlights.value = results[5].value ?? []
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
  <div class="px-4 py-6 md:py-8 max-w-4xl mx-auto">
    <!-- Loading Skeleton -->
    <div
      v-if="isLoading"
      class="space-y-5"
    >
      <!-- Greeting card skeleton -->
      <div class="rounded-sm border bg-card p-6 md:p-8">
        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          <div class="space-y-2">
            <div class="h-8 w-64 animate-pulse rounded-lg bg-muted" />
            <div class="h-4 w-48 animate-pulse rounded-lg bg-muted" />
          </div>
          <div class="flex items-center gap-3">
            <div class="flex gap-1">
              <div
                v-for="i in 5"
                :key="i"
                class="h-2 w-5 animate-pulse rounded-full bg-muted"
              />
            </div>
          </div>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 mt-6">
          <div
            v-for="i in 3"
            :key="i"
            class="flex items-center gap-2.5 rounded-sm bg-muted/50 p-3"
          >
            <div class="w-9 h-9 animate-pulse rounded-sm bg-muted" />
            <div class="space-y-1.5">
              <div class="h-5 w-12 animate-pulse rounded bg-muted" />
              <div class="h-3 w-16 animate-pulse rounded bg-muted" />
            </div>
          </div>
        </div>
      </div>

      <!-- Event section skeleton -->
      <div class="rounded-sm border bg-card p-5 md:p-6">
        <div class="flex items-center gap-2 mb-3">
          <div class="h-4 w-4 animate-pulse rounded bg-muted" />
          <div class="h-4 w-36 animate-pulse rounded bg-muted" />
        </div>
        <div class="h-6 w-3/4 animate-pulse rounded-lg bg-muted mb-2" />
        <div class="h-4 w-1/2 animate-pulse rounded bg-muted mb-3" />
        <div class="flex gap-4">
          <div class="h-4 w-28 animate-pulse rounded bg-muted" />
          <div class="h-4 w-24 animate-pulse rounded bg-muted" />
        </div>
      </div>

      <!-- Two-column grid skeleton -->
      <div class="grid gap-5 md:grid-cols-2">
        <div
          v-for="i in 2"
          :key="i"
          class="rounded-sm border bg-card p-5"
        >
          <div class="flex items-center gap-2 mb-4">
            <div class="h-4 w-4 animate-pulse rounded bg-muted" />
            <div class="h-4 w-32 animate-pulse rounded bg-muted" />
          </div>
          <div class="space-y-3">
            <div
              v-for="j in 2"
              :key="j"
              class="rounded-sm bg-muted/30 border border-border/50 p-3.5"
            >
              <div class="h-4 w-3/4 animate-pulse rounded bg-muted mb-2" />
              <div class="h-3 w-1/2 animate-pulse rounded bg-muted" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <template v-else>
      <!-- Subscription teaser — показываем UNSUBSCRIBER'у, чтобы понимал куда платить -->
      <RouterLink
        v-if="!isSubscribed"
        to="/tariffs"
        class="block relative rounded-sm border border-accent/40 bg-accent/[0.04] hover:bg-accent/[0.07] p-5 md:p-6 transition-colors group"
      >
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div class="flex-1 min-w-0">
            <div class="font-mono text-[10px] uppercase tracking-widest text-accent/70 mb-1.5">
              // подписка
            </div>
            <h3 class="font-semibold text-lg mb-1">
              Открой полный доступ к платформе
            </h3>
            <p class="text-sm text-muted-foreground">
              События, биржа заданий, ИИ-чаты и менторские контакты — от 259 ₽/мес.
            </p>
          </div>
          <span class="shrink-0 inline-flex items-center gap-2 px-4 py-2 rounded-sm bg-accent text-accent-foreground font-medium text-sm group-hover:bg-accent/90 transition-colors">
            Выбрать тариф →
          </span>
        </div>
      </RouterLink>

      <!-- Hero Greeting -->
      <div
        data-onboarding="dashboard"
        class="relative rounded-sm border bg-card p-6 md:p-8 overflow-hidden terminal-card"
      >
        <!-- Subtle accent gradient -->
        <div class="absolute top-0 right-0 w-48 h-48 rounded-full bg-accent/5 -translate-y-1/2 translate-x-1/4 blur-2xl" />
        <div class="relative">
          <div class="font-mono text-[11px] text-accent/70 tracking-wider mb-3">
            <span class="text-muted-foreground">{{ new Date().toISOString().slice(0, 10) }}</span> //
            session.active
          </div>
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
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 mt-6">
            <RouterLink
              to="/points"
              class="flex items-center gap-2.5 rounded-sm border border-border bg-muted/30 hover:border-accent/40 transition-colors p-3"
            >
              <div class="flex items-center justify-center w-9 h-9 rounded-sm bg-yellow-500/10">
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
              class="flex items-center gap-2.5 rounded-sm border border-border bg-muted/30 hover:border-accent/40 transition-colors p-3"
            >
              <div class="flex items-center justify-center w-9 h-9 rounded-sm bg-purple-500/10">
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
              class="flex items-center gap-2.5 rounded-sm border border-border bg-muted/30 hover:border-accent/40 transition-colors p-3"
            >
              <div class="flex items-center justify-center w-9 h-9 rounded-sm bg-blue-500/10">
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

      <!-- Nearest Events -->
      <div
        v-if="nearestEvents.length > 0"
        class="mt-5 space-y-4"
      >
        <div
          v-for="event in nearestEvents"
          :key="event.id"
          class="rounded-sm terminal-card overflow-hidden"
          :class="event.exclusiveChatId
            ? 'bg-gradient-to-br from-amber-50/80 to-yellow-50/50 dark:from-amber-950/30 dark:to-yellow-950/20 border-amber-300/60 dark:border-amber-600/40'
            : 'bg-card'"
        >
          <div class="p-5 md:p-6">
            <div v-if="event.exclusiveChatId" class="flex items-center gap-1.5 text-amber-600 dark:text-amber-400 mb-2">
              <Crown class="h-4 w-4" />
              <span class="text-xs font-semibold uppercase tracking-wide">{{ event.exclusiveChatTitle || 'Эксклюзив' }}</span>
            </div>
            <div class="flex items-center gap-2 mb-3">
              <span class="font-mono text-[10px] text-muted-foreground/60 mr-1">[01]</span>
              <Zap class="h-4 w-4 text-accent" />
              <span class="text-sm font-medium text-accent">Ближайшее событие</span>
              <span
                v-if="isEventLive(event)"
                class="inline-flex items-center gap-1 rounded-full bg-red-500 px-2 py-0.5 text-xs font-medium text-white ml-auto"
              >
                <Radio class="h-3 w-3" />
                LIVE
              </span>
            </div>

            <h2 class="text-xl font-bold">
              {{ event.title }}
            </h2>
            <p
              v-if="event.description"
              class="text-sm text-muted-foreground mt-1.5 line-clamp-2"
            >
              {{ event.description }}
            </p>

            <div class="flex flex-wrap items-center gap-x-4 gap-y-1.5 mt-3 text-sm text-muted-foreground">
              <span class="flex items-center gap-1.5">
                <Calendar class="h-3.5 w-3.5" />
                {{ dateFormatter.format(getNextOccurrenceDate(event)) }}
              </span>
              <span
                v-if="event.hosts.length"
                class="flex items-center gap-1.5"
              >
                <Users class="h-3.5 w-3.5" />
                {{ event.hosts.map(h => h.firstName).join(', ') }}
              </span>
              <span class="flex items-center gap-1.5">
                {{ event.members.length }}{{ event.maxParticipants > 0 ? `/${event.maxParticipants}` : '' }} участников
              </span>
            </div>

            <div class="flex items-center gap-3 mt-4">
              <RouterLink
                v-if="isMemberOfEvent(event) || isHostOfEvent(event)"
                to="/events"
                class="inline-flex items-center gap-1.5 rounded-sm bg-accent/10 text-accent px-4 py-2 text-sm font-medium"
              >
                <CheckCircle class="h-4 w-4" />
                {{ isHostOfEvent(event) ? 'Вы ведущий' : 'Вы записаны' }}
              </RouterLink>
              <RouterLink
                v-else
                to="/events"
                class="inline-flex items-center gap-1.5 rounded-sm bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 transition-colors"
              >
                Записаться
                <ArrowRight class="h-3.5 w-3.5" />
              </RouterLink>
            </div>
          </div>
        </div>
      </div>

      <div
        v-else
        class="mt-5 rounded-sm terminal-card bg-card overflow-hidden"
      >
        <div class="p-5 md:p-6">
          <div class="flex items-center gap-2 mb-3">
            <span class="font-mono text-[10px] text-muted-foreground/60 mr-1">[01]</span>
            <Zap class="h-4 w-4 text-accent" />
            <span class="text-sm font-medium text-accent">Ближайшие события</span>
          </div>
          <p class="text-sm text-muted-foreground">
            Пока нет запланированных событий
          </p>
          <RouterLink
            to="/events"
            class="inline-flex items-center gap-1.5 rounded-sm bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 transition-colors mt-3"
          >
            Посмотреть события
            <ArrowRight class="h-3.5 w-3.5" />
          </RouterLink>
        </div>
      </div>

      <!-- Active Quests & Tasks Grid -->
      <div class="mt-5 grid gap-5 md:grid-cols-2">
        <!-- Chat Quests -->
        <div
          class="rounded-sm terminal-card bg-card p-5"
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <span class="font-mono text-[10px] text-muted-foreground/60 mr-1">[02]</span>
              <Flame class="h-4 w-4 text-orange-500" />
              <span class="text-sm font-semibold">Задания в чатах</span>
            </div>
            <RouterLink
              to="/quests"
              class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              Все →
            </RouterLink>
          </div>

          <div
            v-if="activeQuests.length > 0"
            class="space-y-3"
          >
            <div
              v-for="quest in activeQuests.slice(0, 3)"
              :key="quest.id"
              class="rounded-sm bg-muted/30 border border-border/50 p-3.5"
            >
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-center gap-2.5 min-w-0">
                  <div class="flex items-center justify-center w-8 h-8 rounded-lg bg-orange-500/10 shrink-0">
                    <component
                      :is="quest.questType === 'daily_streak' ? Flame : MessageCircle"
                      class="h-4 w-4 text-orange-500"
                    />
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-medium truncate">
                      {{ quest.title }}
                    </p>
                    <p
                      v-if="quest.description"
                      class="text-xs text-muted-foreground line-clamp-1"
                    >
                      {{ quest.description }}
                    </p>
                    <p
                      class="text-xs flex items-center gap-1"
                      :class="questDeadlineDays(quest.endsAt) <= 1 ? 'text-red-500' : questDeadlineDays(quest.endsAt) <= 3 ? 'text-orange-500' : 'text-muted-foreground'"
                    >
                      <Calendar class="h-3 w-3" />
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
                  <span>{{ quest.currentCount }} / {{ quest.targetCount }} {{ quest.questType === 'daily_streak' ? 'дней подряд' : 'сообщений' }}</span>
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
          <p
            v-else
            class="text-sm text-muted-foreground"
          >
            Нет активных заданий. Пиши в чатах, чтобы зарабатывать баллы!
          </p>
        </div>

        <!-- Open Tasks -->
        <div
          class="rounded-sm terminal-card bg-card p-5"
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <span class="font-mono text-[10px] text-muted-foreground/60 mr-1">[03]</span>
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

          <div
            v-if="openTasks.length > 0"
            class="space-y-3"
          >
            <RouterLink
              v-for="task in openTasks"
              :key="task.id"
              to="/tasks"
              class="block rounded-sm bg-muted/30 border border-border/50 hover:bg-muted/60 transition-colors p-3.5"
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
                <span
                  class="text-xs font-medium px-2 py-0.5 rounded-full shrink-0"
                  :class="task.status === 'IN_PROGRESS' ? 'text-orange-500 bg-orange-500/10' : 'text-blue-500 bg-blue-500/10'"
                >
                  {{ task.status === 'IN_PROGRESS' ? `${task.assignees?.length ?? 0}/${task.maxAssignees}` : 'Открыто' }}
                </span>
              </div>
              <div class="flex items-center gap-2 mt-2 text-xs text-muted-foreground">
                <span>от {{ task.creator.firstName }}</span>
              </div>
            </RouterLink>
          </div>
          <p
            v-else
            class="text-sm text-muted-foreground"
          >
            Нет открытых заданий
          </p>
        </div>
      </div>

      <!-- Chat Highlights -->
      <div
        v-if="highlights.length > 0"
        class="mt-5 rounded-sm terminal-card bg-card p-5"
      >
        <div class="flex items-center gap-2 mb-4">
          <span class="font-mono text-[10px] text-muted-foreground/60 mr-1">[04]</span>
          <MessageSquare class="h-4 w-4 text-yellow-500" />
          <span class="text-sm font-semibold">Лучшее из чатов</span>
        </div>
        <div class="space-y-3">
          <div
            v-for="hl in highlights"
            :key="hl.id"
            class="rounded-sm bg-muted/30 border border-border/50 p-3.5"
          >
            <div class="flex items-center gap-2 mb-1.5">
              <div class="flex items-center justify-center w-7 h-7 rounded-full bg-accent/10 shrink-0">
                <span class="text-xs font-bold text-accent">
                  {{ (hl.authorFirstName || hl.authorUsername || '?').charAt(0).toUpperCase() }}
                </span>
              </div>
              <span class="text-sm font-medium truncate">
                {{ hl.authorFirstName || hl.authorUsername || 'Аноним' }}
              </span>
              <span
                v-if="hl.authorUsername"
                class="text-xs text-muted-foreground"
              >
                @{{ hl.authorUsername }}
              </span>
            </div>
            <p class="text-sm text-foreground/90 whitespace-pre-line line-clamp-3">
              {{ hl.messageText }}
            </p>
            <p class="text-[11px] text-muted-foreground mt-1.5">
              {{ formatShortDate(hl.createdAt) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Achievements Preview -->
      <div
        class="mt-5 rounded-sm terminal-card bg-card p-5"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <span class="font-mono text-[10px] text-muted-foreground/60 mr-1">[05]</span>
            <Award class="h-4 w-4 text-purple-500" />
            <span class="text-sm font-semibold">Достижения</span>
          </div>
          <RouterLink
            to="/achievements"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
          >
            Все →
          </RouterLink>
        </div>
        <div
          v-if="achievements.recent.length > 0"
          class="flex gap-3 overflow-x-auto pb-1"
        >
          <div
            v-for="ach in achievements.recent"
            :key="ach.id"
            class="flex items-center gap-2.5 rounded-sm bg-muted/30 border border-border/50 px-4 py-3 shrink-0"
          >
            <div class="flex items-center justify-center w-10 h-10 rounded-full bg-green-500/20 shrink-0">
              <component
                :is="achievementIconMap[ach.icon] || Award"
                class="h-5 w-5 text-green-500"
              />
            </div>
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
        <div v-else>
          <p class="text-sm text-muted-foreground">
            {{ achievements.unlocked }} из {{ achievements.total }} открыто
          </p>
          <RouterLink
            to="/achievements"
            class="text-sm text-accent hover:underline mt-1 inline-block"
          >
            Посмотреть все достижения →
          </RouterLink>
        </div>
      </div>

      <!-- Quick Links -->
      <div
        class="mt-5 grid grid-cols-2 sm:grid-cols-4 gap-3"
      >
        <RouterLink
          to="/leaderboard"
          class="rounded-sm border bg-card hover:bg-muted/50 transition-colors p-4 text-center terminal-card"
        >
          <Trophy class="h-5 w-5 mx-auto text-yellow-500 mb-1.5" />
          <p class="text-xs font-medium">
            Лидерборд
          </p>
        </RouterLink>
        <RouterLink
          to="/marketplace"
          class="rounded-sm border bg-card hover:bg-muted/50 transition-colors p-4 text-center terminal-card"
        >
          <Star class="h-5 w-5 mx-auto text-accent mb-1.5" />
          <p class="text-xs font-medium">
            Барахолка
          </p>
        </RouterLink>
        <RouterLink
          to="/mentors"
          class="rounded-sm border bg-card hover:bg-muted/50 transition-colors p-4 text-center terminal-card"
        >
          <Users class="h-5 w-5 mx-auto text-blue-500 mb-1.5" />
          <p class="text-xs font-medium">
            Менторы
          </p>
        </RouterLink>
        <RouterLink
          to="/referals"
          class="rounded-sm border bg-card hover:bg-muted/50 transition-colors p-4 text-center terminal-card"
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
