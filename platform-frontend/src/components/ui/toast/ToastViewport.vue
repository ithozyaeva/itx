<script setup lang="ts">
import type { ToastViewportProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import { ToastViewport } from 'reka-ui'
import { computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<ToastViewportProps & { class?: HTMLAttributes['class'] }>()

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
  <ToastViewport v-bind="delegatedProps" :class="cn('fixed top-[env(safe-area-inset-top)] z-[100] flex max-h-[100dvh] w-full flex-col-reverse p-4 sm:top-auto sm:bottom-[env(safe-area-inset-bottom)] sm:right-[env(safe-area-inset-right)] sm:flex-col md:max-w-[420px]', props.class)" />
</template>
