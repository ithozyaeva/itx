<script setup lang="ts">
import type { Component } from 'vue'

withDefaults(defineProps<{
  as?: keyof HTMLElementTagNameMap | Component
  variant?: 'active' | 'default'
  disabled?: boolean
}>(), {
  as: 'span',
  variant: 'default',
  disabled: false,
})
</script>

<template>
  <component
    :is="as"
    class="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors"
    :class="{
      'border-accent/30 text-accent': variant === 'active',
      'border-border text-foreground': variant === 'default',
      'opacity-50 pointer-events-none': disabled,
    }"
  >
    <slot />
  </component>
</template>
