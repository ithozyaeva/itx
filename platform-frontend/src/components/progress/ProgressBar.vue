<script setup lang="ts">
import { computed } from 'vue'
import { progressPct } from '@/lib/progressFormat'

const props = withDefaults(defineProps<{
  progress: number
  target: number
  state?: 'active' | 'done' | 'awarded'
  size?: 'xs' | 'sm'
  label?: string
}>(), {
  state: 'active',
  size: 'sm',
})

const pct = computed(() => progressPct(props.progress, props.target))

const trackHeight = computed(() => props.size === 'xs' ? 'h-1' : 'h-1.5')
const fillColor = computed(() => {
  if (props.state === 'done')
    return 'bg-green-500'
  if (props.state === 'awarded')
    return 'bg-yellow-500'
  return 'bg-accent'
})
</script>

<template>
  <div
    class="w-full rounded-full bg-muted overflow-hidden"
    :class="trackHeight"
    role="progressbar"
    :aria-valuenow="progress"
    :aria-valuemin="0"
    :aria-valuemax="target"
    :aria-label="label ?? `Прогресс ${progress} из ${target}`"
  >
    <div
      class="h-full rounded-full transition-all"
      :class="fillColor"
      :style="{ width: `${pct}%` }"
    />
  </div>
</template>
