<script setup lang="ts">
import type { Component } from 'vue'

type TypographyVariant = 'h1' | 'h2' | 'h3' | 'h4' | 'body-xl' | 'body-l' | 'body-m' | 'body-s' | 'body-xs' | 'title' | 'label' | 'price' | 'button-text' | 'name-text' | 'date'

const props = withDefaults(defineProps<{
  as?: keyof HTMLElementTagNameMap | Component
  variant?: TypographyVariant
  class?: string
}>(), {
  as: 'span',
  variant: 'body-m',
})

const variantClasses: Record<TypographyVariant, string> = {
  'h1': 'text-3xl md:text-4xl font-bold tracking-tight uppercase',
  'h2': 'text-2xl font-bold tracking-tight uppercase',
  'h3': 'text-lg font-semibold',
  'h4': 'text-base font-semibold',
  'body-xl': 'text-xl',
  'body-l': 'text-lg',
  'body-m': 'text-base',
  'body-s': 'text-sm',
  'body-xs': 'text-xs',
  'title': 'text-xl font-semibold',
  'label': 'text-sm font-medium',
  'price': 'text-2xl font-bold',
  'button-text': 'text-sm font-medium',
  'name-text': 'text-base font-medium',
  'date': 'text-sm text-muted-foreground',
}
</script>

<template>
  <component :is="as" :class="[variantClasses[variant], props.class]">
    <slot />
  </component>
</template>
