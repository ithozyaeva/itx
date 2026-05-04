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
  progressLabel?: string // кастомный текст слева от прогресса (напр. «3 / 7 дней подряд»)
  // Состояния
  done?: boolean // выполнено (по факту)
  awarded?: boolean // награда зачислена
  // Иконка
  icon?: Component
  iconTone?: IconTone
  // Дедлайн
  endsAt?: string
  // Доп. бейдж сверху над прогрессом (для дейликов: tier)
  pillLabel?: string
  pillClasses?: string
  // Маркер «+ ачивка» (для челленджей)
  hasAchievement?: boolean
  // Навигация
  to?: RouteLocationRaw
  // Стиль
  variant?: 'card' | 'row'
  size?: 'compact' | 'normal'
  showDoneStrike?: boolean // зачёркивать заголовок при done (для action-карточек)
}>(), {
  done: false,
  awarded: false,
  iconTone: 'orange',
  variant: 'card',
  size: 'normal',
  hasAchievement: false,
  showDoneStrike: false,
})

const showProgressBlock = computed(() => typeof props.progress === 'number' && typeof props.target === 'number' && props.target > 0)

const completionLabel = computed(() => {
  if (!props.done && !props.awarded)
    return null
  return props.awarded ? 'Награда зачислена' : 'Выполнено'
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
      <div class="flex items-center justify-between text-xs text-muted-foreground mb-1.5">
        <span
          v-if="pillLabel"
          class="px-1.5 py-0.5 rounded text-[10px] font-medium"
          :class="pillClasses"
        >
          {{ pillLabel }}
        </span>
        <span v-else>
          {{ progressLabel ?? `${progress} / ${target}` }}
        </span>
        <span
          v-if="completionLabel"
          class="font-medium flex items-center gap-1"
          :class="awarded ? 'text-yellow-500' : 'text-green-500'"
        >
          <CheckCircle class="h-3 w-3" aria-hidden="true" />
          {{ completionLabel }}
        </span>
        <span v-else-if="!pillLabel">
          {{ progressLabel ? `${progress} / ${target}` : `${Math.min(100, Math.round(((progress ?? 0) / Math.max(target ?? 1, 1)) * 100))}%` }}
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

    <slot name="footer" />
  </component>
</template>
