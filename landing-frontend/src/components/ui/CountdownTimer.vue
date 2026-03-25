<script setup lang="ts">
import { useNow } from '@vueuse/core'
import { computed } from 'vue'

const props = defineProps<{
  deadline: string
  label?: string
}>()

const now = useNow({ interval: 1000 })
const deadlineMs = new Date(props.deadline).getTime()

const remaining = computed(() => Math.max(0, deadlineMs - now.value.getTime()))
const isExpired = computed(() => remaining.value <= 0)

const days = computed(() => Math.floor(remaining.value / (1000 * 60 * 60 * 24)))
const hours = computed(() => Math.floor((remaining.value / (1000 * 60 * 60)) % 24))
const minutes = computed(() => Math.floor((remaining.value / (1000 * 60)) % 60))
const seconds = computed(() => Math.floor((remaining.value / 1000) % 60))

function pad(n: number): string {
  return String(n).padStart(2, '0')
}
</script>

<template>
  <div
    v-if="!isExpired"
    class="flex flex-col items-center gap-2"
  >
    <span class="text-sm font-medium text-muted-foreground">
      {{ label || 'Акция заканчивается через:' }}
    </span>
    <div class="flex items-center gap-1.5">
      <div class="flex flex-col items-center rounded-lg bg-accent/10 px-3 py-1.5 min-w-[52px]">
        <span class="text-xl font-bold tabular-nums text-accent">{{ days }}</span>
        <span class="text-[10px] uppercase text-muted-foreground">дн</span>
      </div>
      <span class="text-lg font-bold text-accent">:</span>
      <div class="flex flex-col items-center rounded-lg bg-accent/10 px-3 py-1.5 min-w-[52px]">
        <span class="text-xl font-bold tabular-nums text-accent">{{ pad(hours) }}</span>
        <span class="text-[10px] uppercase text-muted-foreground">час</span>
      </div>
      <span class="text-lg font-bold text-accent">:</span>
      <div class="flex flex-col items-center rounded-lg bg-accent/10 px-3 py-1.5 min-w-[52px]">
        <span class="text-xl font-bold tabular-nums text-accent">{{ pad(minutes) }}</span>
        <span class="text-[10px] uppercase text-muted-foreground">мин</span>
      </div>
      <span class="text-lg font-bold text-accent">:</span>
      <div class="flex flex-col items-center rounded-lg bg-accent/10 px-3 py-1.5 min-w-[52px]">
        <span class="text-xl font-bold tabular-nums text-accent">{{ pad(seconds) }}</span>
        <span class="text-[10px] uppercase text-muted-foreground">сек</span>
      </div>
    </div>
  </div>
</template>
