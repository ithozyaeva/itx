<script setup lang="ts">
import Code from '~icons/ic/twotone-code'
import Navigation from '../sections/Navigation.vue'
</script>

<template>
  <header
    class="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
  >
    <div class="container flex h-16 items-center justify-between">
      <div class="flex items-center gap-2 font-bold text-xl">
        <Code class="h-6 w-6" />
        <span>IT-ХОЗЯЕВА</span>
      </div>
      <Navigation />
      <div class="w-20" />
    </div>
  </header>
</template>
