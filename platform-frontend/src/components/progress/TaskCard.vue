<script setup lang="ts">
import type { Component } from 'vue'
import type { RouteLocationRaw } from 'vue-router'
import type { IconTone } from './TintedIcon.vue'
import { CheckCircle, Star } from 'lucide-vue-next'
import { computed } from 'vue'
import Deadline from './Deadline.vue'
import PointsBadge from './PointsBadge.vue'
import ProgressBar from './ProgressBar.vue'
import TintedIcon from './TintedIcon.vue'

const props = withDefaults(defineProps<{
  title: string
  description?: string
  points: number
  // Прогресс
  progress?: number
  target?: number
  /** Кастомный текст в строке прогресса вместо «N / M» (напр. «3 / 7 дней подряд»). */
  progressLabel?: string
  // Состояния
  /** Цель достигнута. */
  done?: boolean
  /** Награда зачислена (для челленджей: completed != awarded). Подразумевает done. */
  awarded?: boolean
  // Иконка-аватар
  icon?: Component
  iconTone?: IconTone
  /** Дедлайн — показывается, пока задача не done. */
  endsAt?: string
  /** Бейдж типа задачи (напр. tier дейлика). Не заменяет счётчик прогресса. */
  pillLabel?: string
  pillTone?: IconTone
  /** Маркер «+ ачивка» (для челленджей с achievementCode). */
  hasAchievement?: boolean
  /** Если задан — карточка становится <RouterLink>. */
  to?: RouteLocationRaw
  size?: 'compact' | 'normal'
  /** Зачёркивать заголовок при done (action-карточки в духе todo-листа). */
  showDoneStrike?: boolean
}>(), {
  done: false,
  awarded: false,
  iconTone: 'orange',
  pillTone: 'accent',
  size: 'normal',
  hasAchievement: false,
  showDoneStrike: false,
})

const showProgressBlock = computed(() =>
  typeof props.progress === 'number' && typeof props.target === 'number' && props.target > 0,
)

const completionLabel = computed(() => {
  if (!props.done && !props.awarded)
    return null
  return props.awarded ? 'Награда зачислена' : 'Выполнено'
})

const counterLabel = computed(() => {
  if (props.progressLabel)
    return props.progressLabel
  if (typeof props.progress === 'number' && typeof props.target === 'number')
    return `${props.progress} / ${props.target}`
  return ''
})

const tagClass = computed(() => {
  if (props.awarded)
    return 'border-yellow-500/30 bg-yellow-500/5'
  if (props.done)
    return 'border-green-500/30 bg-green-500/5'
  return 'border-border bg-card'
})

const progressState = computed(() => {
  if (props.awarded)
    return 'awarded' as const
  if (props.done)
    return 'done' as const
  return 'active' as const
})

const pillToneClass = computed(() => {
  switch (props.pillTone) {
    case 'orange': return 'bg-orange-500/10 text-orange-500'
    case 'green': return 'bg-green-500/10 text-green-500'
    case 'blue': return 'bg-blue-500/10 text-blue-500'
    case 'yellow': return 'bg-yellow-500/10 text-yellow-500'
    case 'purple': return 'bg-purple-500/10 text-purple-500'
    default: return 'bg-accent/10 text-accent'
  }
})

const Wrapper = computed(() => props.to ? 'router-link' : 'div')
</script>

<template>
  <component
    :is="Wrapper"
    :to="to"
    class="block rounded-sm border transition-colors terminal-card"
    :class="[
      tagClass,
      to ? 'hover:border-accent/40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring' : '',
      size === 'compact' ? 'p-4' : 'p-5',
    ]"
  >
    <div class="flex items-start gap-3">
      <TintedIcon
        v-if="icon"
        :icon="icon"
        :tone="iconTone"
        :done="done"
        :size="size === 'compact' ? 'sm' : 'md'"
      />
      <div class="flex-1 min-w-0">
        <div class="flex items-start justify-between gap-2">
          <p
            class="text-sm font-semibold leading-snug"
            :class="showDoneStrike && done ? 'line-through text-muted-foreground' : ''"
          >
            {{ title }}
          </p>
          <PointsBadge
            :amount="points"
            :earned="awarded || (done && !showProgressBlock)"
          />
        </div>
        <p
          v-if="description"
          class="text-xs text-muted-foreground mt-1"
        >
          {{ description }}
        </p>
      </div>
    </div>

    <div v-if="showProgressBlock" class="mt-4">
      <div class="flex items-center justify-between gap-2 text-xs text-muted-foreground mb-1.5">
        <div class="flex items-center gap-2 min-w-0">
          <span
            v-if="pillLabel"
            class="px-1.5 py-0.5 rounded text-[10px] font-medium shrink-0"
            :class="pillToneClass"
          >
            {{ pillLabel }}
          </span>
          <span
            v-if="counterLabel"
            class="truncate"
          >
            {{ counterLabel }}
          </span>
        </div>
        <span
          v-if="completionLabel"
          class="font-medium flex items-center gap-1 shrink-0"
          :class="awarded ? 'text-yellow-500' : 'text-green-500'"
        >
          <CheckCircle class="h-3 w-3" aria-hidden="true" />
          {{ completionLabel }}
        </span>
      </div>
      <ProgressBar
        :progress="progress ?? 0"
        :target="target ?? 1"
        :state="progressState"
        :size="size === 'compact' ? 'xs' : 'sm'"
      />
    </div>

    <div
      v-if="(endsAt && !done) || hasAchievement"
      class="mt-3 flex items-center gap-3 text-xs flex-wrap"
    >
      <Deadline v-if="endsAt && !done" :ends-at="endsAt" />
      <span
        v-if="hasAchievement"
        class="inline-flex items-center gap-1 text-purple-500"
      >
        <Star class="h-3 w-3" aria-hidden="true" />
        + ачивка
      </span>
    </div>
  </component>
</template>
