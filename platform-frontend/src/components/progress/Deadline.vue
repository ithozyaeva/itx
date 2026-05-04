<script setup lang="ts">
import { Calendar } from 'lucide-vue-next'
import { computed } from 'vue'
import { deadlineUrgency, formatDeadline } from '@/lib/progressFormat'

const props = defineProps<{
  endsAt: string
}>()

const label = computed(() => formatDeadline(props.endsAt))
const urgency = computed(() => deadlineUrgency(props.endsAt))

const colorClass = computed(() => {
  switch (urgency.value) {
    case 'expired':
    case 'critical':
      return 'text-red-500'
    case 'warning':
      return 'text-orange-500'
    default:
      return 'text-muted-foreground'
  }
})
</script>

<template>
  <span
    class="inline-flex items-center gap-1 text-xs"
    :class="colorClass"
  >
    <Calendar class="h-3 w-3" aria-hidden="true" />
    {{ label }}
  </span>
</template>
