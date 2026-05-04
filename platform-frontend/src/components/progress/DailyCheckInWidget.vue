<script setup lang="ts">
import { ArrowRight, CheckCircle, Flame, Snowflake, Trophy } from 'lucide-vue-next'
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useToast } from '@/components/ui/toast'
import { useDailies } from '@/composables/useDailies'
import TintedIcon from './TintedIcon.vue'

const props = withDefaults(defineProps<{
  variant?: 'hero' | 'compact'
  /** В compact-режиме кликабельная карточка ведёт сюда (по умолчанию /progress?tab=today). */
  linkTo?: string
}>(), {
  variant: 'hero',
  linkTo: '/progress?tab=today',
})

const { today, streak, checkingIn, refresh, checkIn } = useDailies()
const { toast } = useToast()

const checkInDone = computed(() => today.value?.checkIn.done ?? false)
const streakValue = computed(() => streak.value?.current ?? 0)
const longest = computed(() => streak.value?.longest ?? 0)
const milestones = computed(() => streak.value?.milestones ?? [])
const nextThreshold = computed(() => streak.value?.nextThreshold ?? null)
const daysToNext = computed(() => streak.value?.daysToNext ?? null)
const freezesAvailable = computed(() => streak.value?.freezesAvailable ?? 0)
const tasksTotal = computed(() => today.value?.tasks.length ?? 5)
const tasksAwarded = computed(() => today.value?.tasks.filter(t => t.awarded).length ?? 0)

async function handleCheckIn() {
  const resp = await checkIn()
  if (!resp)
    return
  if (resp.alreadyToday) {
    toast({ title: 'Сегодня уже отмечались', description: 'Возвращайтесь завтра — стрик растёт.' })
    return
  }
  toast({ title: '+5 баллов', description: `Стрик: ${resp.streak.current} дней подряд 🔥` })
}

onMounted(() => {
  if (!today.value)
    refresh()
})
</script>

<template>
  <!-- HERO variant — большая карточка для страницы Progress -->
  <div
    v-if="props.variant === 'hero'"
    class="rounded-sm border bg-card p-6 transition-colors terminal-card"
    :class="checkInDone ? 'border-green-500/30 bg-green-500/5' : ''"
  >
    <div class="flex items-start gap-5 flex-col sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-4">
        <TintedIcon
          :icon="Flame"
          tone="orange"
          :done="checkInDone"
          size="lg"
        />
        <div>
          <p class="text-sm text-muted-foreground">
            Текущий стрик
          </p>
          <p class="text-2xl font-bold leading-tight">
            {{ streakValue }} <span class="text-base font-normal text-muted-foreground">дн.</span>
          </p>
          <p
            v-if="longest > streakValue"
            class="text-xs text-muted-foreground mt-0.5"
          >
            Личный рекорд: {{ longest }} дн.
          </p>
        </div>
      </div>

      <button
        v-if="!checkInDone"
        type="button"
        class="px-6 py-3 rounded-sm bg-primary text-primary-foreground font-medium text-sm hover:bg-primary/90 transition-colors disabled:opacity-50 w-full sm:w-auto min-h-[44px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        :disabled="checkingIn"
        @click="handleCheckIn"
      >
        {{ checkingIn ? '...' : 'Сделать check-in (+5 баллов)' }}
      </button>
      <div
        v-else
        class="text-xs text-green-500 flex items-center gap-1.5 font-medium"
      >
        <CheckCircle class="h-4 w-4" aria-hidden="true" />
        Сегодня отмечен
      </div>
    </div>

    <div
      v-if="milestones.length"
      class="mt-5 pt-5 border-t border-border/50"
    >
      <div class="flex items-center justify-between text-xs text-muted-foreground mb-2">
        <span>До следующей награды</span>
        <span v-if="nextThreshold">
          {{ daysToNext }} дн. → +{{ milestones.find(m => m.days === nextThreshold)?.reward }} баллов
        </span>
        <span v-else class="text-yellow-500 font-medium flex items-center gap-1">
          <Trophy class="h-3 w-3" aria-hidden="true" />
          Все пороги пройдены
        </span>
      </div>
      <div class="grid grid-cols-4 gap-2">
        <div
          v-for="m in milestones"
          :key="m.days"
          class="rounded-sm border px-2 py-2 text-center transition-colors"
          :class="m.reached
            ? 'border-yellow-500/40 bg-yellow-500/10 text-yellow-500'
            : 'border-border bg-muted/30 text-muted-foreground'"
        >
          <div class="text-xs font-bold">
            {{ m.days }} дн.
          </div>
          <div class="text-[10px]">
            +{{ m.reward }}
          </div>
        </div>
      </div>
      <div
        v-if="freezesAvailable > 0"
        class="mt-3 text-[11px] text-muted-foreground flex items-center gap-1.5"
      >
        <Snowflake class="h-3 w-3" aria-hidden="true" />
        Доступна 1 заморозка на этой неделе — пропуск 1 дня не сбросит стрик.
      </div>
    </div>
  </div>

  <!-- COMPACT variant — карточка для дашборда / сайдбаров -->
  <RouterLink
    v-else
    :to="linkTo"
    class="block rounded-sm border bg-card p-4 transition-colors terminal-card hover:border-accent/40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
    :class="checkInDone ? 'border-green-500/30 bg-green-500/5' : ''"
  >
    <div class="flex items-center justify-between gap-3 flex-wrap">
      <div class="flex items-center gap-3 min-w-0">
        <TintedIcon
          :icon="Flame"
          tone="orange"
          :done="checkInDone"
          size="md"
        />
        <div class="min-w-0">
          <div class="flex items-baseline gap-2">
            <p class="text-sm font-semibold">
              Стрик: {{ streakValue }} дн.
            </p>
            <p class="text-xs text-muted-foreground">
              · {{ tasksAwarded }} / {{ tasksTotal }} дейликов
            </p>
          </div>
          <p class="text-xs text-muted-foreground mt-0.5">
            {{ checkInDone ? 'Сегодня отмечен — приходи завтра' : 'Сделай check-in и фарми баллы' }}
          </p>
        </div>
      </div>

      <button
        v-if="!checkInDone"
        type="button"
        class="px-4 py-2 rounded-sm bg-primary text-primary-foreground text-xs font-medium hover:bg-primary/90 transition-colors disabled:opacity-50 min-h-[44px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        :disabled="checkingIn"
        @click.stop.prevent="handleCheckIn"
      >
        {{ checkingIn ? '...' : 'Check-in (+5)' }}
      </button>
      <span
        v-else
        class="text-xs text-muted-foreground flex items-center gap-1"
      >
        Все детали <ArrowRight class="h-3 w-3" aria-hidden="true" />
      </span>
    </div>
  </RouterLink>
</template>
