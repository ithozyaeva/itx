<script setup lang="ts">
import type { PopoverContentProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import { PopoverContent, PopoverPortal } from 'reka-ui'
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<PopoverContentProps & { class?: HTMLAttributes['class'] }>(), {
  sideOffset: 4,
  align: 'start',
})
</script>

<template>
  <PopoverPortal>
    <PopoverContent
      v-bind="props"
      :class="cn('z-50 rounded-xl border bg-popover p-3 text-popover-foreground shadow-md outline-none animate-in fade-in-0 zoom-in-95', props.class)"
    >
      <slot />
    </PopoverContent>
  </PopoverPortal>
</template>
