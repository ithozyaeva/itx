<script setup lang="ts">
import {
  CheckCircle,
  Flame,
  Snowflake,
  Star,
  Target,
  Trophy,
} from 'lucide-vue-next'
import { computed, onMounted } from 'vue'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { useDailies } from '@/composables/useDailies'

const { today, streak, loading, checkingIn, refresh, checkIn } = useDailies()
const { toast } = useToast()

const checkInDone = computed(() => today.value?.checkIn.done ?? false)
const tasks = computed(() => today.value?.tasks ?? [])
const allBonus = computed(() => today.value?.allBonus ?? { points: 50, awarded: false })

const completedTaskCount = computed(() => tasks.value.filter(t => t.awarded).length)
const totalTaskCount = computed(() => tasks.value.length)
const allCompleted = computed(() => totalTaskCount.value > 0 && completedTaskCount.value === totalTaskCount.value)

function progressPct(t: { progress: number, target: number }) {
  if (t.target === 0)
    return 0
  return Math.min(100, Math.round((t.progress / t.target) * 100))
}

const tierBadge: Record<string, { label: string, classes: string }> = {
  engagement: { label: 'Просмотр', classes: 'bg-blue-500/10 text-blue-500' },
  light: { label: 'Действие', classes: 'bg-green-500/10 text-green-500' },
  meaningful: { label: 'Контент', classes: 'bg-orange-500/10 text-orange-500' },
  big: { label: 'Серьёзное', classes: 'bg-purple-500/10 text-purple-500' },
}

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
  refresh()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-4xl">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/dailies
    </div>
    <div class="flex items-center justify-between mb-6">
      <Typography variant="h2" as="h1">
        Ежедневные задания
      </Typography>
    </div>

    <!-- Hero: check-in + streak -->
    <div
      class="rounded-sm border bg-card p-6 mb-6 transition-colors terminal-card"
      :class="checkInDone ? 'border-green-500/30 bg-green-500/5' : ''"
    >
      <div class="flex items-start gap-5 flex-col sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-4">
          <div
            class="flex items-center justify-center w-14 h-14 rounded-sm shrink-0"
            :class="checkInDone ? 'bg-green-500/20' : 'bg-orange-500/10'"
          >
            <Flame
              v-if="!checkInDone"
              class="h-7 w-7 text-orange-500"
            />
            <CheckCircle
              v-else
              class="h-7 w-7 text-green-500"
            />
          </div>
          <div>
            <p class="text-sm text-muted-foreground">
              Текущий стрик
            </p>
            <p class="text-2xl font-bold leading-tight">
              {{ streak?.current ?? 0 }} <span class="text-base font-normal text-muted-foreground">дн.</span>
            </p>
            <p
              v-if="streak && streak.longest > (streak.current ?? 0)"
              class="text-xs text-muted-foreground mt-0.5"
            >
              Личный рекорд: {{ streak.longest }} дн.
            </p>
          </div>
        </div>

        <button
          v-if="!checkInDone"
          class="px-6 py-3 rounded-sm bg-primary text-primary-foreground font-medium text-sm hover:bg-primary/90 transition-colors disabled:opacity-50 w-full sm:w-auto"
          :disabled="checkingIn"
          @click="handleCheckIn"
        >
          {{ checkingIn ? '...' : 'Сделать check-in (+5 баллов)' }}
        </button>
        <div
          v-else
          class="text-xs text-green-500 flex items-center gap-1.5 font-medium"
        >
          <CheckCircle class="h-4 w-4" />
          Сегодня отмечен
        </div>
      </div>

      <!-- Milestones progress bar -->
      <div
        v-if="streak && streak.milestones.length"
        class="mt-5 pt-5 border-t border-border/50"
      >
        <div class="flex items-center justify-between text-xs text-muted-foreground mb-2">
          <span>До следующей награды</span>
          <span v-if="streak.nextThreshold">
            {{ streak.daysToNext }} дн. → +{{ streak.milestones.find(m => m.days === streak!.nextThreshold)?.reward }} баллов
          </span>
          <span v-else class="text-yellow-500 font-medium flex items-center gap-1">
            <Trophy class="h-3 w-3" />
            Все пороги пройдены
          </span>
        </div>
        <div class="grid grid-cols-4 gap-2">
          <div
            v-for="m in streak.milestones"
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
          v-if="streak && streak.freezesAvailable > 0"
          class="mt-3 text-[11px] text-muted-foreground flex items-center gap-1.5"
        >
          <Snowflake class="h-3 w-3" />
          Доступна 1 заморозка на этой неделе — пропуск 1 дня не сбросит стрик.
        </div>
      </div>
    </div>

    <!-- Tasks -->
    <div v-if="tasks.length > 0">
      <div class="flex items-center justify-between mb-3">
        <Typography variant="h4" as="h2">
          Задания на сегодня
        </Typography>
        <div class="text-sm text-muted-foreground flex items-center gap-2">
          <Target class="h-4 w-4" />
          {{ completedTaskCount }} / {{ totalTaskCount }}
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="rounded-sm border bg-card p-4 transition-colors"
          :class="task.awarded ? 'border-green-500/30 bg-green-500/5' : ''"
        >
          <div class="flex items-start justify-between gap-3 mb-2">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold leading-snug">
                {{ task.title }}
              </p>
              <p
                v-if="task.description"
                class="text-xs text-muted-foreground mt-1"
              >
                {{ task.description }}
              </p>
            </div>
            <span class="text-xs font-bold text-yellow-500 bg-yellow-500/10 px-2 py-0.5 rounded-full shrink-0">
              +{{ task.points }}
            </span>
          </div>
          <div class="mt-3">
            <div class="flex items-center justify-between text-xs text-muted-foreground mb-1">
              <span :class="tierBadge[task.tier]?.classes" class="px-1.5 py-0.5 rounded text-[10px] font-medium">
                {{ tierBadge[task.tier]?.label ?? task.tier }}
              </span>
              <span v-if="task.awarded" class="text-green-500 font-medium flex items-center gap-1">
                <CheckCircle class="h-3 w-3" />
                Готово
              </span>
              <span v-else>{{ task.progress }} / {{ task.target }}</span>
            </div>
            <div class="w-full h-1 rounded-full bg-muted overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="task.awarded ? 'bg-green-500' : 'bg-accent'"
                :style="{ width: `${progressPct(task)}%` }"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- All-five bonus -->
      <div
        class="mt-4 rounded-sm border p-4 flex items-center gap-3 transition-colors"
        :class="allBonus.awarded ? 'border-yellow-500/40 bg-yellow-500/10' : 'border-dashed border-border bg-muted/30'"
      >
        <Star
          class="h-6 w-6 shrink-0"
          :class="allBonus.awarded ? 'text-yellow-500' : 'text-muted-foreground'"
        />
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold">
            Бонус +{{ allBonus.points }} за все {{ totalTaskCount }} заданий
          </p>
          <p class="text-xs text-muted-foreground">
            {{ allBonus.awarded
              ? `Бонус уже зачислен`
              : allCompleted
                ? `Зачисление…`
                : `Выполни оставшиеся ${totalTaskCount - completedTaskCount}, чтобы получить бонус` }}
          </p>
        </div>
      </div>
    </div>

    <div
      v-else-if="!loading"
      class="rounded-sm border border-dashed bg-muted/20 p-6 text-center text-sm text-muted-foreground"
    >
      Задания на сегодня формируются. Загляни через минуту.
    </div>
  </div>
</template>
