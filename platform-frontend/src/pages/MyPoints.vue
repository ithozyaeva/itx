<script setup lang="ts">
import type { PointsSummary, PointTransaction } from '@/models/points'
import type { ChatQuestWithProgress } from '@/services/chatQuestService'
import { Typography } from 'itx-ui-kit'
import { Calendar, CheckCircle, FileText, Folder, Loader2, MessageCircle, MessageSquare, Mic, Share2, Star, User } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast'
import { useSSE } from '@/composables/useSSE'
import { reasonLabels } from '@/lib/reasonLabels'
import { formatShortDate } from '@/lib/utils'
import { chatQuestService } from '@/services/chatQuestService'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'

const PAGE_SIZE = 10
const { toast } = useToast()
const data = ref<PointsSummary | null>(null)
const chatQuests = ref<ChatQuestWithProgress[]>([])
const isLoading = ref(true)
const visibleCount = ref(PAGE_SIZE)

const rewards = [
  { action: 'Участие в событии', points: 10 },
  { action: 'Проведение события', points: 25 },
  { action: 'Отзыв на сообщество', points: 15 },
  { action: 'Отзыв на услугу', points: 10 },
  { action: 'Загрузка резюме', points: 10 },
  { action: 'Создание реферала', points: 5 },
  { action: 'Конверсия реферала', points: 30 },
  { action: 'Заполненный профиль', points: 20 },
  { action: 'Еженедельная активность', points: 5 },
  { action: '3+ события за месяц', points: 30 },
  { action: 'Серия 4 недели подряд', points: 50 },
]

const quests = [
  { label: 'Запишись на событие', to: '/events', points: 10, icon: Calendar, reason: 'event_attend' },
  { label: 'Проведи событие', to: '/events', points: 25, icon: Mic, reason: 'event_host' },
  { label: 'Оставь отзыв на сообщество', to: '/my-reviews', points: 15, icon: MessageSquare, reason: 'review_community' },
  { label: 'Загрузи резюме', to: '/resumes', points: 10, icon: FileText, reason: 'resume_upload' },
  { label: 'Заполни профиль полностью', to: '/me', points: 20, icon: User, reason: 'profile_complete' },
  { label: 'Создай реферальную ссылку', to: '/referals', points: 5, icon: Folder, reason: 'referal_create' },
  { label: 'Получи конверсию по рефералу', to: '/referals', points: 30, icon: Share2, reason: 'referal_conversion' },
]

const weeklyQuests = [
  { label: 'Еженедельная активность', points: 5, reason: 'weekly_activity' },
  { label: '3+ события за месяц', points: 30, reason: 'monthly_active' },
  { label: 'Серия 4 недели подряд', points: 50, reason: 'streak_4weeks' },
]

const completedReasons = computed(() => {
  if (!data.value)
    return new Set<string>()
  return new Set(data.value.transactions.map(tx => tx.reason))
})

const visibleTransactions = computed<PointTransaction[]>(() => {
  if (!data.value)
    return []
  return data.value.transactions.slice(0, visibleCount.value)
})

const hasMoreTransactions = computed(() => {
  if (!data.value)
    return false
  return visibleCount.value < data.value.transactions.length
})

function loadMore() {
  visibleCount.value += PAGE_SIZE
}

function isQuestCompleted(reason: string) {
  return completedReasons.value.has(reason)
}

const SEEN_KEY = 'points_seen_quests'

function getSeenQuests(): Set<string> {
  try {
    const raw = localStorage.getItem(SEEN_KEY)
    return raw ? new Set(JSON.parse(raw)) : new Set()
  }
  catch {
    return new Set()
  }
}

function saveSeenQuests(seen: Set<string>) {
  localStorage.setItem(SEEN_KEY, JSON.stringify([...seen]))
}

function showNewCompletions() {
  const seen = getSeenQuests()
  const newlyCompleted = quests.filter(
    q => isQuestCompleted(q.reason) && !seen.has(q.reason),
  )

  if (newlyCompleted?.length > 0) {
    const labels = newlyCompleted.map(q => q.label).join(', ')
    const totalPoints = newlyCompleted.reduce((sum, q) => sum + q.points, 0)
    toast({
      title: `Задание выполнено! +${totalPoints}`,
      description: labels,
    })
    newlyCompleted.forEach(q => seen.add(q.reason))
    saveSeenQuests(seen)
  }
}

async function fetchPoints() {
  isLoading.value = true
  try {
    const [points, quests] = await Promise.all([
      pointsService.getMyPoints(),
      chatQuestService.getActiveQuests().catch(() => [] as ChatQuestWithProgress[]),
    ])
    data.value = points
    chatQuests.value = quests
    showNewCompletions()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

function questProgress(quest: ChatQuestWithProgress) {
  if (!quest.targetCount)
    return 0
  return Math.min(100, Math.round((quest.currentCount / quest.targetCount) * 100))
}

function formatQuestDeadline(dateStr: string) {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = date.getTime() - now.getTime()
  const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
  if (diffDays <= 0)
    return 'Истекло'
  if (diffDays === 1)
    return 'Остался 1 день'
  if (diffDays <= 7)
    return `Осталось ${diffDays} дн.`
  return date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
}

useSSE('points', () => fetchPoints())

onMounted(() => {
  fetchPoints()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <Typography
      variant="h2"
      as="h1"
      class="mb-6"
    >
      Мои баллы
    </Typography>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else-if="data">
      <!-- Баланс -->
      <div class="bg-card border border-border rounded-2xl p-6 mb-8 flex items-center gap-4">
        <div class="flex items-center justify-center w-14 h-14 rounded-full bg-yellow-500/20">
          <Star class="h-7 w-7 text-yellow-500" />
        </div>
        <div>
          <div class="text-sm text-muted-foreground">
            Текущий баланс
          </div>
          <div class="text-3xl font-bold">
            {{ data.balance }}
          </div>
        </div>
      </div>

      <!-- Задания -->
      <Typography
        variant="h3"
        as="h2"
        class="mb-4"
      >
        Задания
      </Typography>
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 mb-8">
        <RouterLink
          v-for="quest in quests"
          :key="quest.to"
          :to="quest.to"
          class="flex items-center gap-3 rounded-2xl p-4 transition-colors"
          :class="isQuestCompleted(quest.reason)
            ? 'bg-green-500/5 border border-green-500/30'
            : 'bg-card border border-border hover:border-primary/50'"
        >
          <div
            class="flex items-center justify-center w-10 h-10 rounded-full shrink-0"
            :class="isQuestCompleted(quest.reason) ? 'bg-green-500/20' : 'bg-primary/10'"
          >
            <CheckCircle
              v-if="isQuestCompleted(quest.reason)"
              class="h-5 w-5 text-green-500"
            />
            <component
              :is="quest.icon"
              v-else
              class="h-5 w-5 text-primary"
            />
          </div>
          <div class="flex-1 min-w-0">
            <div
              class="font-medium text-sm"
              :class="isQuestCompleted(quest.reason) ? 'line-through text-muted-foreground' : ''"
            >
              {{ quest.label }}
            </div>
          </div>
          <div
            class="shrink-0 text-xs font-bold px-2 py-1 rounded-full"
            :class="isQuestCompleted(quest.reason)
              ? 'text-green-500 bg-green-500/10'
              : 'text-yellow-500 bg-yellow-500/10'"
          >
            {{ isQuestCompleted(quest.reason) ? '' : '+' }}{{ quest.points }}
          </div>
        </RouterLink>
      </div>

      <!-- Еженедельные задания -->
      <Typography
        variant="h3"
        as="h2"
        class="mb-4"
      >
        Бонусы за активность
      </Typography>
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 mb-8">
        <div
          v-for="quest in weeklyQuests"
          :key="quest.reason"
          class="flex items-center gap-3 rounded-2xl p-4 transition-colors"
          :class="isQuestCompleted(quest.reason)
            ? 'bg-green-500/5 border border-green-500/30'
            : 'bg-card border border-border'"
        >
          <div
            class="flex items-center justify-center w-10 h-10 rounded-full shrink-0"
            :class="isQuestCompleted(quest.reason) ? 'bg-green-500/20' : 'bg-primary/10'"
          >
            <CheckCircle
              v-if="isQuestCompleted(quest.reason)"
              class="h-5 w-5 text-green-500"
            />
            <Star
              v-else
              class="h-5 w-5 text-primary"
            />
          </div>
          <div class="flex-1 min-w-0">
            <div
              class="font-medium text-sm"
              :class="isQuestCompleted(quest.reason) ? 'line-through text-muted-foreground' : ''"
            >
              {{ quest.label }}
            </div>
          </div>
          <div
            class="shrink-0 text-xs font-bold px-2 py-1 rounded-full"
            :class="isQuestCompleted(quest.reason)
              ? 'text-green-500 bg-green-500/10'
              : 'text-yellow-500 bg-yellow-500/10'"
          >
            {{ isQuestCompleted(quest.reason) ? '' : '+' }}{{ quest.points }}
          </div>
        </div>
      </div>

      <!-- Задания в чатах -->
      <template v-if="chatQuests?.length > 0">
        <Typography
          variant="h3"
          as="h2"
          class="mb-4"
        >
          Задания в чатах
        </Typography>
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 mb-8">
          <div
            v-for="quest in chatQuests"
            :key="quest.id"
            class="rounded-2xl p-4 transition-colors"
            :class="quest.completed
              ? 'bg-green-500/5 border border-green-500/30'
              : 'bg-card border border-border'"
          >
            <div class="flex items-start gap-3">
              <div
                class="flex items-center justify-center w-10 h-10 rounded-full shrink-0"
                :class="quest.completed ? 'bg-green-500/20' : 'bg-primary/10'"
              >
                <CheckCircle
                  v-if="quest.completed"
                  class="h-5 w-5 text-green-500"
                />
                <MessageCircle
                  v-else
                  class="h-5 w-5 text-primary"
                />
              </div>
              <div class="flex-1 min-w-0">
                <div
                  class="font-medium text-sm"
                  :class="quest.completed ? 'line-through text-muted-foreground' : ''"
                >
                  {{ quest.title }}
                </div>
                <div
                  v-if="quest.description"
                  class="text-xs text-muted-foreground mt-0.5"
                >
                  {{ quest.description }}
                </div>
                <!-- Прогресс-бар -->
                <div v-if="!quest.completed" class="mt-2">
                  <div class="flex items-center justify-between text-xs text-muted-foreground mb-1">
                    <span>{{ quest.currentCount }} / {{ quest.targetCount }}</span>
                    <span>{{ formatQuestDeadline(quest.endsAt) }}</span>
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
                class="shrink-0 text-xs font-bold px-2 py-1 rounded-full"
                :class="quest.completed
                  ? 'text-green-500 bg-green-500/10'
                  : 'text-yellow-500 bg-yellow-500/10'"
              >
                {{ quest.completed ? '' : '+' }}{{ quest.pointsReward }}
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- За что начисляются баллы -->
      <Typography
        variant="h3"
        as="h2"
        class="mb-4"
      >
        За что начисляются баллы
      </Typography>
      <div class="bg-card border border-border rounded-2xl overflow-x-auto mb-8">
        <table class="w-full text-sm min-w-[320px]">
          <thead>
            <tr class="border-b border-border">
              <th class="text-left px-4 py-3 font-medium text-muted-foreground">
                Действие
              </th>
              <th class="text-right px-4 py-3 font-medium text-muted-foreground">
                Баллы
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="reward in rewards"
              :key="reward.action"
              class="border-b border-border last:border-0"
            >
              <td class="px-4 py-3">
                {{ reward.action }}
              </td>
              <td class="px-4 py-3 text-right font-bold text-yellow-500">
                +{{ reward.points }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- История транзакций -->
      <Typography
        variant="h3"
        as="h2"
        class="mb-4"
      >
        История транзакций
      </Typography>
      <div
        v-if="(data.transactions?.length ?? 0) === 0"
        class="text-center py-12 text-muted-foreground"
      >
        Пока нет транзакций. Выполняйте задания, чтобы зарабатывать баллы!
      </div>
      <template v-else>
        <div class="space-y-3">
          <div
            v-for="tx in visibleTransactions"
            :key="tx.id"
            class="flex items-center gap-4 p-4 bg-card border border-border rounded-2xl"
          >
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm">
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
        </div>

        <div
          v-if="hasMoreTransactions"
          class="mt-4 flex justify-center"
        >
          <Button
            variant="outline"
            @click="loadMore"
          >
            Показать ещё
          </Button>
        </div>
      </template>
    </template>
  </div>
</template>
