<script setup lang="ts">
import type { Component } from 'vue'
import { Button } from '@/components/ui/button'

withDefaults(defineProps<{
  icon?: Component
  title?: string
  description?: string
  actionLabel?: string
  variant?: 'plain' | 'dashed'
  size?: 'md' | 'sm'
}>(), {
  variant: 'plain',
  size: 'md',
})

const emit = defineEmits<{
  action: []
}>()
</script>

<template>
  <div
    class="flex flex-col items-center justify-center text-center"
    :class="[
      variant === 'dashed' ? 'rounded-sm border border-dashed bg-muted/20' : '',
      size === 'sm' ? 'py-8 px-4' : 'py-12 px-4',
    ]"
  >
    <div
      v-if="icon"
      class="flex items-center justify-center rounded-sm bg-muted mb-4"
      :class="size === 'sm' ? 'w-12 h-12' : 'w-16 h-16'"
    >
      <component
        :is="icon"
        class="text-muted-foreground"
        :class="size === 'sm' ? 'h-6 w-6' : 'h-8 w-8'"
        aria-hidden="true"
      />
    </div>
    <h3
      v-if="title"
      class="font-semibold mb-1"
      :class="size === 'sm' ? 'text-sm' : 'text-base'"
    >
      {{ title }}
    </h3>
    <p
      v-if="description"
      class="text-sm text-muted-foreground max-w-sm mb-4"
    >
      {{ description }}
    </p>
    <slot />
    <Button
      v-if="actionLabel"
      size="sm"
      @click="emit('action')"
    >
      {{ actionLabel }}
    </Button>
  </div>
</template>
