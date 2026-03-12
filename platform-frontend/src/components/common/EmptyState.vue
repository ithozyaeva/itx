<script setup lang="ts">
import type { Component } from 'vue'
import { Button } from '@/components/ui/button'

defineProps<{
  icon: Component
  title: string
  description?: string
  actionLabel?: string
}>()

const emit = defineEmits<{
  action: []
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-12 px-4 text-center">
    <div class="flex items-center justify-center w-16 h-16 rounded-2xl bg-muted mb-4">
      <component
        :is="icon"
        class="h-8 w-8 text-muted-foreground"
      />
    </div>
    <h3 class="text-base font-semibold mb-1">
      {{ title }}
    </h3>
    <p
      v-if="description"
      class="text-sm text-muted-foreground max-w-sm mb-4"
    >
      {{ description }}
    </p>
    <Button
      v-if="actionLabel"
      size="sm"
      @click="emit('action')"
    >
      {{ actionLabel }}
    </Button>
  </div>
</template>
