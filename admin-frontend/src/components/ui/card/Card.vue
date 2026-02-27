<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div
    :class="
      cn(
        'rounded-3xl border bg-card text-card-foreground shadow hover:shadow-md transition-shadow',
        props.class,
      )
    "
  >
    <slot />
  </div>
</template>
