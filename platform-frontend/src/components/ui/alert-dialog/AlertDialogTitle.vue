<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { AlertDialogTitle } from 'reka-ui'
import { cn } from '@/lib/utils'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <AlertDialogTitle :class="cn('text-lg font-semibold', props.class)">
    <slot />
  </AlertDialogTitle>
</template>
