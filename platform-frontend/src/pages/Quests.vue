<script setup lang="ts">
import type { ChatQuestWithProgress } from '@/services/chatQuestService'
import {
  Calendar,
  CheckCircle,
  Flame,
  MessageCircle,
  Star,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import QuestCardSkeleton from '@/components/quests/QuestCardSkeleton.vue'
import { Typography } from '@/components/ui/typography'
import { useSSE } from '@/composables/useSSE'
import { formatShortDate } from '@/lib/utils'
import { chatQuestService } from '@/services/chatQuestService'
import { handleError } from '@/services/errorService'

const quests = ref<ChatQuestWithProgress[]>([])
const isLoading = ref(true)
const activeFilter = ref<'all' | 'active' | 'completed'>('all')

const filters: { key: 'all' | 'active' | 'completed', label: string }[] = [
  { key: 'all', label: 'Все' },
  { key: 'active', label: 'Активные' },
  { key: 'completed', label: 'Выполненные' },
]

const filteredQuests = computed(() => {
  if (activeFilter.value === 'all')
    return quests.value
  if (activeFilter.value === 'active')
    return quests.value.filter(q => !q.completed)
  return quests.value.filter(q => q.completed)
})

function questProgress(quest: ChatQuestWithProgress) {
  if (quest.targetCount === 0)
    return 0
  return Math.min(100, Math.round((quest.currentCount / quest.targetCount) * 100))
}

function questDeadlineDays(dateStr: string) {
  const date = new Date(dateStr)
  const now = new Date()
  return Math.ceil((date.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
}

function formatQuestDeadline(dateStr: string) {
  const days = questDeadlineDays(dateStr)
  if (days <= 0)
    return 'Истекает'
  if (days === 1)
    return '1 день'
  if (days <= 7)
    return `${days} дн.`
  return formatShortDate(dateStr)
}

function questTypeLabel(quest: ChatQuestWithProgress) {
  if (quest.questType === 'daily_streak')
    return `${quest.currentCount} / ${quest.targetCount} дней подряд`
  return `${quest.currentCount} / ${quest.targetCount} сообщений`
}

async function fetchQuests() {
  isLoading.value = true
  try {
    quests.value = await chatQuestService.getAllQuests() ?? []
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

useSSE('quests', () => fetchQuests())

onMounted(() => {
  fetchQuests()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/quests
    </div>
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Задания в чатах
      </Typography>
      <div
        v-if="!isLoading"
        class="flex items-center gap-2 text-sm"
      >
        <Flame class="h-5 w-5 text-orange-500" />
        <span class="font-bold">{{ quests.filter(q => q.completed).length }}</span>
        <span class="text-muted-foreground">/ {{ quests.length }}</span>
      </div>
    </div>

    <div
      v-if="isLoading"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <QuestCardSkeleton v-for="i in 3" :key="i" />
    </div>

    <template v-else>
      <div class="flex gap-2 mb-6 flex-wrap">
        <button
          v-for="f in filters"
          :key="f.key"
          class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
          :class="activeFilter === f.key
            ? 'bg-primary text-primary-foreground'
            : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
          @click="activeFilter = f.key"
        >
          {{ f.label }}
        </button>
      </div>

      <div
        v-if="filteredQuests.length > 0"
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
      >
        <div
          v-for="quest in filteredQuests"
          :key="quest.id"
          class="rounded-sm border bg-card p-5 transition-colors terminal-card"
          :class="quest.completed ? 'border-green-500/30 bg-green-500/5' : ''"
        >
          <div class="flex items-start gap-3">
            <div
              class="flex items-center justify-center w-10 h-10 rounded-xl shrink-0"
              :class="quest.completed ? 'bg-green-500/20' : 'bg-orange-500/10'"
            >
              <component
                :is="quest.questType === 'daily_streak' ? Flame : MessageCircle"
                class="h-5 w-5"
                :class="quest.completed ? 'text-green-500' : 'text-orange-500'"
              />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between gap-2">
                <p class="text-sm font-semibold">
                  {{ quest.title }}
                </p>
                <span class="text-xs font-bold text-yellow-500 bg-yellow-500/10 px-2 py-0.5 rounded-full shrink-0">
                  +{{ quest.pointsReward }}
                </span>
              </div>
              <p
                v-if="quest.description"
                class="text-xs text-muted-foreground mt-1"
              >
                {{ quest.description }}
              </p>
            </div>
          </div>

          <div class="mt-4">
            <div class="flex items-center justify-between text-xs text-muted-foreground mb-1.5">
              <span>{{ questTypeLabel(quest) }}</span>
              <span v-if="quest.completed" class="text-green-500 font-medium flex items-center gap-1">
                <CheckCircle class="h-3 w-3" />
                Выполнено
              </span>
              <span v-else>{{ questProgress(quest) }}%</span>
            </div>
            <div class="w-full h-1.5 rounded-full bg-muted overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="quest.completed ? 'bg-green-500' : 'bg-accent'"
                :style="{ width: `${questProgress(quest)}%` }"
              />
            </div>
          </div>

          <div
            v-if="!quest.completed"
            class="flex items-center gap-3 mt-3 text-xs"
          >
            <span
              class="flex items-center gap-1"
              :class="questDeadlineDays(quest.endsAt) <= 1 ? 'text-red-500' : questDeadlineDays(quest.endsAt) <= 3 ? 'text-orange-500' : 'text-muted-foreground'"
            >
              <Calendar class="h-3 w-3" />
              {{ formatQuestDeadline(quest.endsAt) }}
            </span>
            <span class="flex items-center gap-1 text-yellow-500">
              <Star class="h-3 w-3" />
              {{ quest.pointsReward }} баллов
            </span>
          </div>
        </div>
      </div>

      <div
        v-else
        class="text-center py-12"
      >
        <Flame class="h-12 w-12 text-muted-foreground/30 mx-auto mb-3" />
        <p class="text-muted-foreground">
          {{ activeFilter === 'completed' ? 'Нет выполненных заданий' : activeFilter === 'active' ? 'Нет активных заданий' : 'Заданий пока нет' }}
        </p>
      </div>
    </template>
  </div>
</template>
