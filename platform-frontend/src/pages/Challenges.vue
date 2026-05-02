<script setup lang="ts">
import type { ChallengeKind, ChallengeWithProgress } from '@/models/challenges'
import { Calendar, CheckCircle2, Medal, Star, Sword, Trophy } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { Typography } from '@/components/ui/typography'
import { useChallenges } from '@/composables/useChallenges'

const { data, loading, fetchAll } = useChallenges()

const tab = ref<ChallengeKind>('weekly')
const list = computed(() => (tab.value === 'weekly' ? data.value?.weekly : data.value?.monthly) ?? [])

const tabs: { key: ChallengeKind, label: string, icon: any }[] = [
  { key: 'weekly', label: 'Еженедельные', icon: Sword },
  { key: 'monthly', label: 'Ежемесячные', icon: Trophy },
]

function progressPct(c: ChallengeWithProgress) {
  if (c.target === 0)
    return 0
  return Math.min(100, Math.round((c.progress / c.target) * 100))
}

function daysLeft(endsAt: string): number {
  const diffMs = new Date(endsAt).getTime() - Date.now()
  return Math.max(0, Math.ceil(diffMs / 86400000))
}

function deadlineLabel(endsAt: string) {
  const d = daysLeft(endsAt)
  if (d <= 0)
    return 'Истекает'
  if (d === 1)
    return '1 день'
  if (d <= 7)
    return `${d} дн.`
  return `${d} дн.`
}

onMounted(() => {
  fetchAll()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-4xl">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/challenges
    </div>
    <div class="flex items-center justify-between mb-6">
      <Typography variant="h2" as="h1">
        Челленджи
      </Typography>
    </div>

    <!-- Tabs -->
    <div class="flex gap-2 mb-6 flex-wrap">
      <button
        v-for="t in tabs"
        :key="t.key"
        class="flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
        :class="tab === t.key
          ? 'bg-primary text-primary-foreground'
          : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
        @click="tab = t.key"
      >
        <component :is="t.icon" class="h-4 w-4" />
        {{ t.label }}
      </button>
    </div>

    <div
      v-if="loading && !data"
      class="grid grid-cols-1 sm:grid-cols-2 gap-4"
    >
      <div
        v-for="i in 3"
        :key="i"
        class="rounded-sm border bg-card p-5"
      >
        <div class="h-5 w-2/3 bg-muted animate-pulse rounded mb-2" />
        <div class="h-3 w-full bg-muted animate-pulse rounded mb-3" />
        <div class="h-1 w-full bg-muted animate-pulse rounded" />
      </div>
    </div>

    <div
      v-else-if="list.length === 0"
      class="rounded-sm border border-dashed bg-muted/20 p-12 text-center"
    >
      <Sword class="h-10 w-10 text-muted-foreground/50 mx-auto mb-3" />
      <p class="text-sm text-muted-foreground">
        {{ tab === 'weekly' ? 'Еженедельные челленджи появятся в понедельник' : 'Ежемесячный челлендж появится 1-го числа' }}
      </p>
    </div>

    <div
      v-else
      class="grid grid-cols-1 sm:grid-cols-2 gap-4"
    >
      <div
        v-for="c in list"
        :key="c.instanceId"
        class="rounded-sm border bg-card p-5 transition-colors terminal-card"
        :class="c.awarded ? 'border-yellow-500/30 bg-yellow-500/5' : ''"
      >
        <div class="flex items-start justify-between gap-3 mb-3">
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold leading-snug">
              {{ c.title }}
            </p>
            <p
              v-if="c.description"
              class="text-xs text-muted-foreground mt-1"
            >
              {{ c.description }}
            </p>
          </div>
          <span class="text-xs font-bold text-yellow-500 bg-yellow-500/10 px-2 py-0.5 rounded-full shrink-0">
            +{{ c.rewardPoints }}
          </span>
        </div>
        <div class="mt-3">
          <div class="flex items-center justify-between text-xs text-muted-foreground mb-1">
            <span v-if="c.awarded" class="text-yellow-500 font-medium flex items-center gap-1">
              <Medal class="h-3 w-3" />
              Завершено
            </span>
            <span v-else>{{ c.progress }} / {{ c.target }}</span>
            <span class="flex items-center gap-1" :class="daysLeft(c.endsAt) <= 1 ? 'text-red-500' : daysLeft(c.endsAt) <= 3 ? 'text-orange-500' : ''">
              <Calendar class="h-3 w-3" />
              {{ deadlineLabel(c.endsAt) }}
            </span>
          </div>
          <div class="w-full h-1.5 rounded-full bg-muted overflow-hidden">
            <div
              class="h-full rounded-full transition-all"
              :class="c.awarded ? 'bg-yellow-500' : 'bg-accent'"
              :style="{ width: `${progressPct(c)}%` }"
            />
          </div>
        </div>
        <div
          v-if="c.achievementCode"
          class="mt-3 text-[11px] text-purple-500 flex items-center gap-1"
        >
          <Star class="h-3 w-3" />
          + ачивка
        </div>
        <div
          v-if="c.awarded"
          class="mt-3 inline-flex items-center gap-1 text-xs text-green-500 font-medium"
        >
          <CheckCircle2 class="h-3 w-3" />
          Награда зачислена
        </div>
      </div>
    </div>
  </div>
</template>
