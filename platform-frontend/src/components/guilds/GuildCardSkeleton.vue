<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
</script>

<template>
  <div class="rounded-2xl border bg-card border-border overflow-hidden">
    <Skeleton class="h-2 rounded-none" />
    <div class="p-4">
      <Skeleton class="h-5 w-1/2 mb-2" />
      <div class="space-y-1.5 mb-3">
        <Skeleton class="h-3.5 w-full" />
        <Skeleton class="h-3.5 w-3/4" />
      </div>
      <div class="flex items-center gap-4 mb-3">
        <Skeleton class="h-4 w-24" />
        <Skeleton class="h-4 w-20" />
      </div>
      <div class="flex items-center gap-1.5 mb-3">
        <Skeleton class="h-5 w-5 rounded-full" />
        <Skeleton class="h-3.5 w-24" />
      </div>
      <Skeleton class="h-9 w-24 rounded-xl" />
    </div>
  </div>
</template>
