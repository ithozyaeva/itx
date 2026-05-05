<script setup lang="ts">
import type { Component } from 'vue'
import type { IconTone } from '@/components/progress/TintedIcon.vue'
import type { PointSource } from '@/lib/reasonLabels'
import type { ChallengeKind } from '@/models/challenges'
import type { PointsSummary, PointTransaction } from '@/models/points'
import type { ChatQuestWithProgress } from '@/services/chatQuestService'
import {
  Calendar,
  CheckCircle,
  FileText,
  Flame,
  Folder,
  Inbox,
  MessageCircle,
  MessageSquare,
  Mic,
  Share2,
  Sparkles,
  Star,
  Sword,
  Target,
  Trophy,
  User,
} from 'lucide-vue-next'
import { computed, defineAsyncComponent, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import { DailyCheckInWidget, PointsBadge, TaskCard, TaskCardSkeleton, TintedIcon } from '@/components/progress'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { useChallenges } from '@/composables/useChallenges'
import { useDailies } from '@/composables/useDailies'
import { useSSE } from '@/composables/useSSE'
import { tierBadge } from '@/lib/dailyTier'
import { pointSources, reasonLabels } from '@/lib/reasonLabels'
import { formatShortDate } from '@/lib/utils'
import { chatQuestService } from '@/services/chatQuestService'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'

// Панели тяжёлых табов грузим лениво — иначе при открытии /progress?tab=today
// (дефолт) пользователь скачивает JS всех 4 «соседних» табов впустую.
// MyStatsPanel особенно ощутим: SVG-charts + contribution graph.
const LeaderboardPanel = defineAsyncComponent(() => import('@/components/progress/LeaderboardPanel.vue'))
const AchievementsPanel = defineAsyncComponent(() => import('@/components/progress/AchievementsPanel.vue'))
const MyStatsPanel = defineAsyncComponent(() => import('@/components/progress/MyStatsPanel.vue'))
const KudosPanel = defineAsyncComponent(() => import('@/components/progress/KudosPanel.vue'))

type TabKey = 'today' | 'period' | 'history' | 'sources' | 'leaderboard' | 'achievements' | 'stats' | 'kudos'
type PeriodFilter = ChallengeKind | 'chats'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()

const VALID_TABS: TabKey[] = ['today', 'period', 'history', 'sources', 'leaderboard', 'achievements', 'stats', 'kudos']
const VALID_PERIODS: PeriodFilter[] = ['weekly', 'monthly', 'chats']

const initialTab: TabKey = (() => {
  const t = route.query.tab as string | undefined
  return (VALID_TABS as string[]).includes(t ?? '') ? (t as TabKey) : 'today'
})()
const initialPeriod: PeriodFilter = (() => {
  const k = route.query.kind as string | undefined
  return (VALID_PERIODS as string[]).includes(k ?? '') ? (k as PeriodFilter) : 'weekly'
})()

const activeTab = ref<TabKey>(initialTab)
const periodFilter = ref<PeriodFilter>(initialPeriod)

watch(activeTab, (val) => {
  if (route.query.tab !== val)
    router.replace({ query: { ...route.query, tab: val } })
})

watch(periodFilter, (val) => {
  if (route.query.kind !== val)
    router.replace({ query: { ...route.query, kind: val } })
})

watch(() => route.query.tab, (val) => {
  if (val && (VALID_TABS as string[]).includes(val as string) && val !== activeTab.value)
    activeTab.value = val as TabKey
})

watch(() => route.query.kind, (val) => {
  if (val && (VALID_PERIODS as string[]).includes(val as string) && val !== periodFilter.value)
    periodFilter.value = val as PeriodFilter
})

// ─── Today (dailies) ─────────────────────────────────────────────────
// Дейлики: API возвращает только task.awarded (зачислена ли награда).
// Промежуточного состояния «выполнено, но баллы не упали» у дейликов
// нет — поэтому в TaskCard прокидываем и done, и awarded из awarded.
// Карточка получает «жёлтый» (awarded) стейт; если в будущем появится
// разделение, передадим разные значения.
const { today, loading: dailiesLoading } = useDailies()
const dailyTasks = computed(() => today.value?.tasks ?? [])
const allBonus = computed(() => today.value?.allBonus ?? { points: 50, awarded: false })
const completedTaskCount = computed(() => dailyTasks.value.filter(t => t.awarded).length)
const totalTaskCount = computed(() => dailyTasks.value.length)
const allCompleted = computed(() => totalTaskCount.value > 0 && completedTaskCount.value === totalTaskCount.value)

// ─── Period (challenges + chat quests) ────────────────────────────────
const { data: challengesData, loading: challengesLoading, fetchAll: fetchChallenges } = useChallenges()

const challengeList = computed(() => {
  if (periodFilter.value === 'weekly')
    return challengesData.value?.weekly ?? []
  if (periodFilter.value === 'monthly')
    return challengesData.value?.monthly ?? []
  return []
})

const chatQuests = ref<ChatQuestWithProgress[]>([])
const chatQuestsLoading = ref(true)
const chatQuestsError = ref<string | null>(null)

async function fetchChatQuests() {
  chatQuestsLoading.value = true
  chatQuestsError.value = null
  try {
    chatQuests.value = await chatQuestService.getAllQuests() ?? []
  }
  catch (error) {
    chatQuestsError.value = (await handleError(error)).message
  }
  finally {
    chatQuestsLoading.value = false
  }
}

useSSE('quests', () => fetchChatQuests())

// ─── History (points) ─────────────────────────────────────────────────
const PAGE_SIZE = 10
const pointsData = ref<PointsSummary | null>(null)
const pointsLoading = ref(true)
const visibleCount = ref(PAGE_SIZE)

const visibleTransactions = computed<PointTransaction[]>(() =>
  pointsData.value?.transactions.slice(0, visibleCount.value) ?? [],
)
const hasMoreTransactions = computed(() =>
  pointsData.value ? visibleCount.value < pointsData.value.transactions.length : false,
)

async function fetchPoints() {
  pointsLoading.value = true
  try {
    pointsData.value = await pointsService.getMyPoints()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    pointsLoading.value = false
  }
}

// SSE 'points' и без нас обновляет дейлики через подписку внутри useDailies()
// — здесь просто перетягиваем сводку транзакций.
useSSE('points', fetchPoints)

// ─── Sources (статичный справочник способов) ──────────────────────────
// pointSources живёт в lib/reasonLabels.ts — общий словарь для табы и тостов.
// Тут добавляем UI-метаданные (иконка/тон), которые в чисто-data-словаре
// были бы лишним балластом.
const sourceVisuals: Record<string, { icon: Component, tone: IconTone }> = {
  event_attend: { icon: Calendar, tone: 'blue' },
  event_host: { icon: Mic, tone: 'orange' },
  review_community: { icon: MessageSquare, tone: 'green' },
  resume_upload: { icon: FileText, tone: 'accent' },
  profile_complete: { icon: User, tone: 'purple' },
  referal_create: { icon: Folder, tone: 'accent' },
  referal_conversion: { icon: Share2, tone: 'green' },
}

const extraSources: { label: string, to: string, icon: Component, tone: IconTone }[] = [
  { label: 'Биржа заданий', to: '/tasks', icon: Target, tone: 'blue' },
  { label: 'AI-материалы', to: '/ai-materials', icon: Sparkles, tone: 'purple' },
]

// ─── Toast о новых выполненных «способах» ─────────────────────────────
// Раньше жил в MyPoints.vue:80-112; перенесли сюда, чтобы поведение
// не пропало. Активация: первое появление транзакции с reason из pointSources.
const SEEN_KEY = 'progress_seen_sources'
function showNewCompletions(transactions: PointTransaction[]) {
  let seen: Set<string>
  try {
    const raw = localStorage.getItem(SEEN_KEY)
    const parsed = raw ? JSON.parse(raw) : []
    seen = new Set(Array.isArray(parsed) ? parsed.filter((s): s is string => typeof s === 'string') : [])
  }
  catch {
    seen = new Set()
  }
  const reasons = new Set(transactions.map(tx => tx.reason))
  const newly: PointSource[] = []
  for (const src of pointSources) {
    if (reasons.has(src.reason) && !seen.has(src.reason))
      newly.push(src)
  }
  if (newly.length === 0)
    return
  const total = newly.reduce((s, x) => s + x.points, 0)
  toast({
    title: `Задание выполнено! +${total}`,
    description: newly.map(n => n.shortLabel ?? n.label).join(', '),
  })
  newly.forEach(n => seen.add(n.reason))
  localStorage.setItem(SEEN_KEY, JSON.stringify([...seen]))
}

watch(pointsData, (val) => {
  if (val?.transactions)
    showNewCompletions(val.transactions)
})

// ─── Init ─────────────────────────────────────────────────────────────
onMounted(() => {
  fetchChallenges()
  fetchChatQuests()
  fetchPoints()
})

function loadMoreTransactions() {
  visibleCount.value += PAGE_SIZE
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-4xl">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/progress
    </div>
    <div class="flex items-end justify-between mb-6 gap-4 flex-wrap">
      <Typography variant="h2" as="h1">
        Прогресс
      </Typography>
      <div
        v-if="pointsData"
        class="flex items-center gap-2.5 rounded-sm border border-border bg-card px-3 py-2"
      >
        <TintedIcon :icon="Star" tone="yellow" size="sm" />
        <div>
          <p class="text-lg font-bold leading-none">
            {{ pointsData.balance }}
          </p>
          <p class="text-[11px] text-muted-foreground mt-0.5">
            баллов
          </p>
        </div>
      </div>
    </div>

    <Tabs v-model="activeTab" class="w-full">
      <TabsList class="mb-6 w-full justify-start overflow-x-auto">
        <TabsTrigger value="today">
          Сегодня
        </TabsTrigger>
        <TabsTrigger value="period">
          Период
        </TabsTrigger>
        <TabsTrigger value="history">
          История
        </TabsTrigger>
        <TabsTrigger value="sources">
          Способы
        </TabsTrigger>
        <TabsTrigger value="leaderboard">
          Рейтинг
        </TabsTrigger>
        <TabsTrigger value="achievements">
          Достижения
        </TabsTrigger>
        <TabsTrigger value="stats">
          Моя статистика
        </TabsTrigger>
        <TabsTrigger value="kudos">
          Благодарности
        </TabsTrigger>
      </TabsList>

      <!-- ───────── ТАБ «Сегодня» ───────── -->
      <TabsContent value="today" class="space-y-6">
        <DailyCheckInWidget variant="hero" />

        <div v-if="dailyTasks.length > 0">
          <div class="flex items-center justify-between mb-3">
            <Typography variant="h4" as="h2">
              Задания на сегодня
            </Typography>
            <div class="text-sm text-muted-foreground flex items-center gap-2">
              <Target class="h-4 w-4" aria-hidden="true" />
              {{ completedTaskCount }} / {{ totalTaskCount }}
            </div>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <TaskCard
              v-for="task in dailyTasks"
              :key="task.id"
              :title="task.title"
              :description="task.description"
              :points="task.points"
              :progress="task.progress"
              :target="task.target"
              :done="task.awarded"
              :awarded="task.awarded"
              :pill-label="tierBadge[task.tier]?.label ?? task.tier"
              :pill-tone="tierBadge[task.tier]?.tone"
              size="compact"
            />
          </div>

          <div
            class="mt-4 rounded-sm border p-4 flex items-center gap-3 transition-colors"
            :class="allBonus.awarded ? 'border-yellow-500/40 bg-yellow-500/10' : 'border-dashed border-border bg-muted/30'"
          >
            <Star
              class="h-6 w-6 shrink-0"
              :class="allBonus.awarded ? 'text-yellow-500' : 'text-muted-foreground'"
              aria-hidden="true"
            />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold">
                Бонус +{{ allBonus.points }} за все {{ totalTaskCount }} заданий
              </p>
              <p class="text-xs text-muted-foreground">
                {{ allBonus.awarded
                  ? 'Бонус уже зачислен'
                  : allCompleted
                    ? 'Зачисление…'
                    : `Выполни оставшиеся ${totalTaskCount - completedTaskCount}, чтобы получить бонус` }}
              </p>
            </div>
          </div>
        </div>

        <EmptyState
          v-else-if="!dailiesLoading"
          :icon="Inbox"
          variant="dashed"
          size="sm"
          title="Задания формируются"
          description="Загляни через минуту."
        />
      </TabsContent>

      <!-- ───────── ТАБ «Период» ───────── -->
      <TabsContent value="period" class="space-y-6">
        <Tabs v-model="periodFilter">
          <TabsList class="w-full justify-start overflow-x-auto">
            <TabsTrigger value="weekly" class="gap-2">
              <Sword class="h-4 w-4" aria-hidden="true" /> Неделя
            </TabsTrigger>
            <TabsTrigger value="monthly" class="gap-2">
              <Trophy class="h-4 w-4" aria-hidden="true" /> Месяц
            </TabsTrigger>
            <TabsTrigger value="chats" class="gap-2">
              <Flame class="h-4 w-4" aria-hidden="true" /> Чаты
            </TabsTrigger>
          </TabsList>

          <TabsContent value="weekly" class="mt-6">
            <div
              v-if="challengesLoading && !challengesData"
              class="grid grid-cols-1 sm:grid-cols-2 gap-4"
            >
              <TaskCardSkeleton v-for="i in 3" :key="i" />
            </div>
            <EmptyState
              v-else-if="challengeList.length === 0"
              :icon="Sword"
              variant="dashed"
              title="Еженедельные челленджи появятся в понедельник"
            />
            <div
              v-else
              class="grid grid-cols-1 sm:grid-cols-2 gap-4"
            >
              <TaskCard
                v-for="c in challengeList"
                :key="c.instanceId"
                :title="c.title"
                :description="c.description"
                :points="c.rewardPoints"
                :progress="c.progress"
                :target="c.target"
                :done="!!c.completedAt"
                :awarded="c.awarded"
                :ends-at="c.endsAt"
                :icon="Sword"
                icon-tone="purple"
                :has-achievement="!!c.achievementCode"
              />
            </div>
          </TabsContent>

          <TabsContent value="monthly" class="mt-6">
            <div
              v-if="challengesLoading && !challengesData"
              class="grid grid-cols-1 sm:grid-cols-2 gap-4"
            >
              <TaskCardSkeleton v-for="i in 3" :key="i" />
            </div>
            <EmptyState
              v-else-if="challengeList.length === 0"
              :icon="Trophy"
              variant="dashed"
              title="Ежемесячный челлендж появится 1-го числа"
            />
            <div
              v-else
              class="grid grid-cols-1 sm:grid-cols-2 gap-4"
            >
              <TaskCard
                v-for="c in challengeList"
                :key="c.instanceId"
                :title="c.title"
                :description="c.description"
                :points="c.rewardPoints"
                :progress="c.progress"
                :target="c.target"
                :done="!!c.completedAt"
                :awarded="c.awarded"
                :ends-at="c.endsAt"
                :icon="Trophy"
                icon-tone="purple"
                :has-achievement="!!c.achievementCode"
              />
            </div>
          </TabsContent>

          <TabsContent value="chats" class="mt-6">
            <div
              v-if="chatQuestsLoading"
              class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
            >
              <TaskCardSkeleton v-for="i in 3" :key="i" />
            </div>
            <ErrorState
              v-else-if="chatQuestsError"
              :message="chatQuestsError"
              @retry="fetchChatQuests()"
            />
            <EmptyState
              v-else-if="chatQuests.length === 0"
              :icon="Flame"
              variant="dashed"
              title="Заданий в чатах пока нет"
              description="Пиши в Telegram-чатах сообщества — задания появятся автоматически."
            />
            <div
              v-else
              class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
            >
              <TaskCard
                v-for="quest in chatQuests"
                :key="quest.id"
                :title="quest.title"
                :description="quest.description"
                :points="quest.pointsReward"
                :progress="quest.currentCount"
                :target="quest.targetCount"
                :done="quest.completed"
                :awarded="quest.completed"
                :ends-at="quest.endsAt"
                :icon="quest.questType === 'daily_streak' ? Flame : MessageCircle"
                icon-tone="orange"
                :progress-label="quest.questType === 'daily_streak'
                  ? `${quest.currentCount} / ${quest.targetCount} дней подряд`
                  : `${quest.currentCount} / ${quest.targetCount} сообщений`"
                size="compact"
              />
            </div>
          </TabsContent>
        </Tabs>
      </TabsContent>

      <!-- ───────── ТАБ «История» ───────── -->
      <TabsContent value="history" class="space-y-4">
        <div
          v-if="pointsData"
          class="rounded-sm border bg-card terminal-card p-5 flex items-center gap-4"
        >
          <TintedIcon :icon="Star" tone="yellow" size="lg" />
          <div>
            <div class="text-sm text-muted-foreground">
              Текущий баланс
            </div>
            <div class="text-3xl font-bold">
              {{ pointsData.balance }}
            </div>
          </div>
        </div>

        <Typography variant="h4" as="h2">
          История транзакций
        </Typography>

        <EmptyState
          v-if="pointsData && pointsData.transactions.length === 0 && !pointsLoading"
          :icon="Inbox"
          variant="dashed"
          title="Пока нет транзакций"
          description="Выполняй задания, чтобы зарабатывать баллы."
        />

        <div v-else-if="pointsData" class="space-y-3">
          <div
            v-for="tx in visibleTransactions"
            :key="tx.id"
            class="flex items-center gap-4 p-4 bg-card border border-border rounded-sm"
          >
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm truncate">
                {{ tx.description || reasonLabels[tx.reason] || tx.reason }}
              </div>
              <div class="text-xs text-muted-foreground mt-0.5">
                {{ reasonLabels[tx.reason] || reasonLabels[tx.reason?.toLowerCase()] || tx.reason }}
              </div>
            </div>
            <div class="shrink-0 text-right">
              <div
                class="font-bold"
                :class="tx.amount > 0 ? 'text-green-500' : 'text-red-500'"
              >
                {{ tx.amount > 0 ? '+' : '' }}{{ tx.amount }}
              </div>
              <div class="text-xs text-muted-foreground">
                {{ formatShortDate(tx.createdAt) }}
              </div>
            </div>
          </div>
          <div
            v-if="hasMoreTransactions"
            class="mt-4 flex justify-center"
          >
            <Button variant="outline" @click="loadMoreTransactions">
              Показать ещё
            </Button>
          </div>
        </div>
      </TabsContent>

      <!-- ───────── ТАБ «Способы» ───────── -->
      <TabsContent value="sources" class="space-y-4">
        <p class="text-sm text-muted-foreground">
          Регулярные дейлики, недельные челленджи и активность в чатах — на соседних табах. Здесь — крупные действия и интеграции платформы, за которые тоже начисляются баллы.
        </p>
        <div class="grid gap-3 sm:grid-cols-2">
          <RouterLink
            v-for="src in pointSources"
            :key="src.reason"
            :to="src.to"
            class="flex items-center gap-3 rounded-sm p-4 bg-card border border-border hover:border-accent/40 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring min-h-[44px]"
          >
            <TintedIcon
              :icon="sourceVisuals[src.reason]?.icon ?? Star"
              :tone="sourceVisuals[src.reason]?.tone ?? 'accent'"
              size="md"
            />
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm">
                {{ src.label }}
              </div>
            </div>
            <PointsBadge :amount="src.points" />
          </RouterLink>

          <RouterLink
            v-for="extra in extraSources"
            :key="extra.label"
            :to="extra.to"
            class="flex items-center gap-3 rounded-sm p-4 bg-card border border-border hover:border-accent/40 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring min-h-[44px]"
          >
            <TintedIcon :icon="extra.icon" :tone="extra.tone" size="md" />
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm">
                {{ extra.label }}
              </div>
            </div>
            <span class="text-xs text-muted-foreground inline-flex items-center gap-1">
              <CheckCircle class="h-3 w-3" aria-hidden="true" />
              разное
            </span>
          </RouterLink>
        </div>
      </TabsContent>

      <!-- ───────── ТАБ «Рейтинг» ───────── -->
      <TabsContent value="leaderboard">
        <LeaderboardPanel />
      </TabsContent>

      <!-- ───────── ТАБ «Достижения» ───────── -->
      <TabsContent value="achievements">
        <AchievementsPanel />
      </TabsContent>

      <!-- ───────── ТАБ «Моя статистика» ───────── -->
      <TabsContent value="stats">
        <MyStatsPanel />
      </TabsContent>

      <!-- ───────── ТАБ «Благодарности» ───────── -->
      <TabsContent value="kudos">
        <KudosPanel />
      </TabsContent>
    </Tabs>
  </div>
</template>
