<script setup lang="ts">
import type { Component } from 'vue'
import { CheckCircle } from 'lucide-vue-next'
import { computed } from 'vue'

export type IconTone = 'orange' | 'green' | 'blue' | 'yellow' | 'purple' | 'accent'

const props = withDefaults(defineProps<{
  icon: Component
  tone?: IconTone
  done?: boolean
  size?: 'sm' | 'md' | 'lg'
}>(), {
  tone: 'accent',
  done: false,
  size: 'md',
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm': return { box: 'w-8 h-8', icon: 'h-4 w-4' }
    case 'lg': return { box: 'w-14 h-14', icon: 'h-7 w-7' }
    default: return { box: 'w-10 h-10', icon: 'h-5 w-5' }
  }
})

const toneClasses = computed(() => {
  if (props.done)
    return { bg: 'bg-green-500/20', text: 'text-green-500' }
  switch (props.tone) {
    case 'orange': return { bg: 'bg-orange-500/10', text: 'text-orange-500' }
    case 'green': return { bg: 'bg-green-500/10', text: 'text-green-500' }
    case 'blue': return { bg: 'bg-blue-500/10', text: 'text-blue-500' }
    case 'yellow': return { bg: 'bg-yellow-500/10', text: 'text-yellow-500' }
    case 'purple': return { bg: 'bg-purple-500/10', text: 'text-purple-500' }
    default: return { bg: 'bg-accent/10', text: 'text-accent' }
  }
})

const iconComponent = computed(() => props.done ? CheckCircle : props.icon)
</script>

<template>
  <div
    class="flex items-center justify-center rounded-sm shrink-0"
    :class="[sizeClasses.box, toneClasses.bg]"
    aria-hidden="true"
  >
    <component
      :is="iconComponent"
      :class="[sizeClasses.icon, toneClasses.text]"
    />
  </div>
</template>
