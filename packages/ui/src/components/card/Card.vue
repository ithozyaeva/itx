<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '../../lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div
    :class="
      cn(
        'rounded-sm border bg-card text-card-foreground transition-all terminal-card',
        props.class,
      )
    "
  >
    <slot />
  </div>
</template>
