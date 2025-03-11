<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  icon: Component
  title: string
  description: string
}>()
</script>

<template>
  <div class="flex items-start gap-4">
    <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary text-primary-foreground">
      <component :is="icon" class="h-5 w-5" />
    </div>
    <div class="space-y-2">
      <h3 class="text-xl font-bold">
        {{ title }}
      </h3>
      <p class="text-muted-foreground">
        {{ description }}
      </p>
    </div>
  </div>
</template>
