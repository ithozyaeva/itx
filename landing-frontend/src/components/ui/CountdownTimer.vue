<script setup lang="ts">
import { useNow } from '@vueuse/core'
import { computed, ref } from 'vue'

const props = defineProps<{
  deadline: string
  label?: string
}>()

const now = useNow({ interval: 1000 })
const deadlineMs = new Date(props.deadline).getTime()

const remaining = computed(() => Math.max(0, deadlineMs - now.value.getTime()))
const isExpired = computed(() => remaining.value <= 0)
const isDismissed = ref(false)

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
    v-if="!isExpired && !isDismissed"
    class="fixed bottom-4 right-4 z-50 flex flex-col items-center gap-1.5 rounded-xl border border-accent/20 bg-background/95 px-4 py-3 shadow-lg backdrop-blur-sm"
  >
    <button
      class="absolute -top-2 -right-2 flex h-5 w-5 items-center justify-center rounded-full bg-muted text-muted-foreground text-xs hover:bg-accent hover:text-accent-foreground transition-colors"
      @click="isDismissed = true"
    >
      &times;
    </button>
    <a
      href="#tariffs"
      class="flex flex-col items-center gap-1.5 no-underline"
    >
      <span class="text-xs font-medium text-muted-foreground">
        {{ label || 'ХОЗЯИН -50% — осталось:' }}
      </span>
      <div class="flex items-center gap-1">
        <div class="flex flex-col items-center rounded-md bg-accent/10 px-2 py-1 min-w-[40px]">
          <span class="text-base font-bold tabular-nums text-accent leading-none">{{ days }}</span>
          <span class="text-[9px] uppercase text-muted-foreground">дн</span>
        </div>
        <span class="text-sm font-bold text-accent">:</span>
        <div class="flex flex-col items-center rounded-md bg-accent/10 px-2 py-1 min-w-[40px]">
          <span class="text-base font-bold tabular-nums text-accent leading-none">{{ pad(hours) }}</span>
          <span class="text-[9px] uppercase text-muted-foreground">час</span>
        </div>
        <span class="text-sm font-bold text-accent">:</span>
        <div class="flex flex-col items-center rounded-md bg-accent/10 px-2 py-1 min-w-[40px]">
          <span class="text-base font-bold tabular-nums text-accent leading-none">{{ pad(minutes) }}</span>
          <span class="text-[9px] uppercase text-muted-foreground">мин</span>
        </div>
        <span class="text-sm font-bold text-accent">:</span>
        <div class="flex flex-col items-center rounded-md bg-accent/10 px-2 py-1 min-w-[40px]">
          <span class="text-base font-bold tabular-nums text-accent leading-none">{{ pad(seconds) }}</span>
          <span class="text-[9px] uppercase text-muted-foreground">сек</span>
        </div>
      </div>
    </a>
  </div>
</template>
