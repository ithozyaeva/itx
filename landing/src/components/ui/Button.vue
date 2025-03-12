<script setup lang="ts">
import type { VariantProps } from 'cva'
import { cva } from 'cva'

withDefaults(defineProps<{
  type?: ButtonProps['type']
}>(), {
  type: 'primary',
})

const button = cva('inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 h-11 rounded-md px-8', {
  variants: {
    type: {
      primary: 'bg-primary text-primary-foreground hover:bg-primary/90',
      secondary:
    'bg-secondary text-secondary-foreground hover:bg-secondary/80',
      outline:
    'border border-input bg-background hover:bg-accent hover:text-accent-foreground',
    },
  },
})

type ButtonProps = VariantProps<typeof button>
</script>

<template>
  <button
    :class="button({ type })"
  >
    <slot />
  </button>
</template>
