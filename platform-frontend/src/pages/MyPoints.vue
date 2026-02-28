<script setup lang="ts">
import type { PointsSummary } from '@/models/points'
import { Typography } from 'itx-ui-kit'
import { Calendar, FileText, Folder, Loader2, MessageSquare, Star, User } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'

const data = ref<PointsSummary | null>(null)
const isLoading = ref(true)

const reasonLabels: Record<string, string> = {
  event_attend: 'Участие в событии',
  event_host: 'Проведение события',
  review_community: 'Отзыв на сообщество',
  review_service: 'Отзыв на услугу',
  resume_upload: 'Загрузка резюме',
  referal_create: 'Создание реферала',
  referal_conversion: 'Конверсия реферала',
  profile_complete: 'Заполнение профиля',
  weekly_activity: 'Еженедельная активность',
  monthly_active: 'Месячная активность',
  streak_4weeks: 'Серия 4 недели',
  admin_manual: 'Ручное начисление',
}

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
  { label: 'Запишись на событие', to: '/events', points: 10, icon: Calendar },
  { label: 'Оставь отзыв на сообщество', to: '/my-reviews', points: 15, icon: MessageSquare },
  { label: 'Загрузи резюме', to: '/resumes', points: 10, icon: FileText },
  { label: 'Заполни профиль полностью', to: '/me', points: 20, icon: User },
  { label: 'Создай реферальную ссылку', to: '/referals', points: 5, icon: Folder },
]

async function fetchPoints() {
  isLoading.value = true
  try {
    data.value = await pointsService.getMyPoints()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

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
          class="flex items-center gap-3 bg-card border border-border rounded-2xl p-4 hover:border-primary/50 transition-colors"
        >
          <div class="flex items-center justify-center w-10 h-10 rounded-full bg-primary/10 shrink-0">
            <component
              :is="quest.icon"
              class="h-5 w-5 text-primary"
            />
          </div>
          <div class="flex-1 min-w-0">
            <div class="font-medium text-sm">
              {{ quest.label }}
            </div>
          </div>
          <div class="shrink-0 text-xs font-bold text-yellow-500 bg-yellow-500/10 px-2 py-1 rounded-full">
            +{{ quest.points }}
          </div>
        </RouterLink>
      </div>

      <!-- За что начисляются баллы -->
      <Typography
        variant="h3"
        as="h2"
        class="mb-4"
      >
        За что начисляются баллы
      </Typography>
      <div class="bg-card border border-border rounded-2xl overflow-hidden mb-8">
        <table class="w-full text-sm">
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
        v-if="data.transactions.length === 0"
        class="text-center py-12 text-muted-foreground"
      >
        Пока нет транзакций. Выполняйте задания, чтобы зарабатывать баллы!
      </div>
      <div
        v-else
        class="space-y-3"
      >
        <div
          v-for="tx in data.transactions"
          :key="tx.id"
          class="flex items-center gap-4 p-4 bg-card border border-border rounded-2xl"
        >
          <div class="flex-1 min-w-0">
            <div class="font-medium text-sm">
              {{ tx.description || reasonLabels[tx.reason] || tx.reason }}
            </div>
            <div class="text-xs text-muted-foreground mt-0.5">
              {{ reasonLabels[tx.reason] || tx.reason }}
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
              {{ new Date(tx.createdAt).toLocaleDateString() }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
