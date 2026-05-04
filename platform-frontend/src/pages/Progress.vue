<script setup lang="ts">
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
import { computed, onMounted, ref, watch } from 'vue'
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
import { reasonLabels } from '@/lib/reasonLabels'
import { formatShortDate } from '@/lib/utils'
import { chatQuestService } from '@/services/chatQuestService'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'

type TabKey = 'today' | 'period' | 'history' | 'sources'
type PeriodFilter = ChallengeKind | 'chats'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()

const VALID_TABS: TabKey[] = ['today', 'period', 'history', 'sources']

const initialTab: TabKey = (() => {
  const t = route.query.tab as string | undefined
  return (VALID_TABS as string[]).includes(t ?? '') ? (t as TabKey) : 'today'
})()
const activeTab = ref<TabKey>(initialTab)

watch(activeTab, (val) => {
  if (route.query.tab !== val)
    router.replace({ query: { ...route.query, tab: val } })
})

watch(() => route.query.tab, (val) => {
  if (val && (VALID_TABS as string[]).includes(val as string) && val !== activeTab.value)
    activeTab.value = val as TabKey
})

// ─── Today (dailies) ─────────────────────────────────────────────────
const { today, loading: dailiesLoading, refresh: refreshDailies } = useDailies()
const dailyTasks = computed(() => today.value?.tasks ?? [])
const allBonus = computed(() => today.value?.allBonus ?? { points: 50, awarded: false })
const completedTaskCount = computed(() => dailyTasks.value.filter(t => t.awarded).length)
const totalTaskCount = computed(() => dailyTasks.value.length)
const allCompleted = computed(() => totalTaskCount.value > 0 && completedTaskCount.value === totalTaskCount.value)

const tierBadge: Record<string, { label: string, classes: string }> = {
  engagement: { label: 'Просмотр', classes: 'bg-blue-500/10 text-blue-500' },
  light: { label: 'Действие', classes: 'bg-green-500/10 text-green-500' },
  meaningful: { label: 'Контент', classes: 'bg-orange-500/10 text-orange-500' },
  big: { label: 'Серьёзное', classes: 'bg-purple-500/10 text-purple-500' },
}

// ─── Period (challenges + chat quests) ────────────────────────────────
const { data: challengesData, loading: challengesLoading, fetchAll: fetchChallenges } = useChallenges()
const periodFilter = ref<PeriodFilter>('weekly')
const periodFilters: { key: PeriodFilter, label: string, icon: any }[] = [
  { key: 'weekly', label: 'Неделя', icon: Sword },
  { key: 'monthly', label: 'Месяц', icon: Trophy },
  { key: 'chats', label: 'Чаты', icon: Flame },
]

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

useSSE('points', () => {
  fetchPoints()
  refreshDailies()
})

// ─── Sources (статичный справочник способов) ──────────────────────────
// Источники, у которых есть отдельный экран (Биржа заданий уже отдельный
// раздел; рефералы / события / отзывы / резюме — обычные разделы платформы).
const sources = [
  { label: 'Запишись на событие', to: '/events', points: 10, icon: Calendar, tone: 'blue' as const },
  { label: 'Проведи событие', to: '/events', points: 25, icon: Mic, tone: 'orange' as const },
  { label: 'Оставь отзыв на сообщество', to: '/my-reviews', points: 15, icon: MessageSquare, tone: 'green' as const },
  { label: 'Загрузи резюме', to: '/resumes', points: 10, icon: FileText, tone: 'accent' as const },
  { label: 'Заполни профиль', to: '/me', points: 20, icon: User, tone: 'purple' as const },
  { label: 'Создай реферальную ссылку', to: '/referals', points: 5, icon: Folder, tone: 'accent' as const },
  { label: 'Получи конверсию по рефералу', to: '/referals', points: 30, icon: Share2, tone: 'green' as const },
  { label: 'Биржа заданий', to: '/tasks', points: 0, icon: Target, tone: 'blue' as const },
  { label: 'Откройте AI-материалы', to: '/ai-materials', points: 0, icon: Sparkles, tone: 'purple' as const },
]

// ─── Toast о новых выполненных «способах» ─────────────────────────────
// Раньше жил в MyPoints.vue:80-112; оставлен на месте, чтобы юзер
// не потерял уведомление о завершении старых статичных квестов.
const SEEN_KEY = 'progress_seen_sources'
function showNewCompletions(transactions: PointTransaction[]) {
  let seen: Set<string>
  try {
    const raw = localStorage.getItem(SEEN_KEY)
    seen = new Set(raw ? JSON.parse(raw) as string[] : [])
  }
  catch {
    seen = new Set()
  }
  const reasons = new Set(transactions.map(tx => tx.reason))
  const reasonToSource = new Map<string, { label: string, points: number }>([
    ['event_attend', { label: 'Запишись на событие', points: 10 }],
    ['event_host', { label: 'Проведи событие', points: 25 }],
    ['review_community', { label: 'Оставь отзыв', points: 15 }],
    ['resume_upload', { label: 'Загрузи резюме', points: 10 }],
    ['profile_complete', { label: 'Заполни профиль', points: 20 }],
    ['referal_create', { label: 'Создай реферальную ссылку', points: 5 }],
    ['referal_conversion', { label: 'Конверсия реферала', points: 30 }],
  ])
  const newly: { label: string, points: number, reason: string }[] = []
  for (const [reason, src] of reasonToSource) {
    if (reasons.has(reason) && !seen.has(reason))
      newly.push({ ...src, reason })
  }
  if (newly.length === 0)
    return
  const total = newly.reduce((s, x) => s + x.points, 0)
  toast({
    title: `Задание выполнено! +${total}`,
    description: newly.map(n => n.label).join(', '),
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
  refreshDailies()
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
              :pill-classes="tierBadge[task.tier]?.classes"
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
        <div class="flex gap-2 flex-wrap" role="tablist" aria-label="Тип периода">
          <button
            v-for="f in periodFilters"
            :key="f.key"
            type="button"
            role="tab"
            :aria-selected="periodFilter === f.key"
            class="inline-flex items-center gap-2 px-3 py-2 rounded-full text-sm font-medium transition-colors min-h-[36px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            :class="periodFilter === f.key
              ? 'bg-primary text-primary-foreground'
              : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
            @click="periodFilter = f.key"
          >
            <component :is="f.icon" class="h-4 w-4" aria-hidden="true" />
            {{ f.label }}
          </button>
        </div>

        <!-- Weekly / Monthly -->
        <template v-if="periodFilter !== 'chats'">
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
            :title="periodFilter === 'weekly' ? 'Еженедельные челленджи появятся в понедельник' : 'Ежемесячный челлендж появится 1-го числа'"
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
              :icon="periodFilter === 'weekly' ? Sword : Trophy"
              icon-tone="purple"
              :has-achievement="!!c.achievementCode"
            />
          </div>
        </template>

        <!-- Chat quests -->
        <template v-else>
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
        </template>
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
            v-for="src in sources"
            :key="src.label"
            :to="src.to"
            class="flex items-center gap-3 rounded-sm p-4 bg-card border border-border hover:border-accent/40 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <TintedIcon :icon="src.icon" :tone="src.tone" size="md" />
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm">
                {{ src.label }}
              </div>
            </div>
            <PointsBadge v-if="src.points > 0" :amount="src.points" />
            <span v-else class="text-xs text-muted-foreground">
              <CheckCircle class="h-3 w-3 inline" aria-hidden="true" /> разное
            </span>
          </RouterLink>
        </div>
      </TabsContent>
    </Tabs>
  </div>
</template>
