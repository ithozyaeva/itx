<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
</script>

<template>
  <div class="rounded-2xl p-4 border bg-card border-border">
    <div class="flex items-start justify-between gap-2 mb-2">
      <Skeleton class="h-4 w-2/3" />
      <Skeleton class="h-5 w-16 rounded-full shrink-0" />
    </div>
    <div class="space-y-1.5 mb-3">
      <Skeleton class="h-3 w-full" />
      <Skeleton class="h-3 w-4/5" />
    </div>
    <div class="flex items-center gap-1.5 mb-1">
      <Skeleton class="h-3.5 w-3.5" />
      <Skeleton class="h-3.5 w-28" />
    </div>
    <div class="flex items-center gap-1.5 mb-1">
      <Skeleton class="h-3.5 w-3.5" />
      <Skeleton class="h-3.5 w-24" />
    </div>
    <Skeleton class="h-8 w-28 rounded-lg mt-3" />
  </div>
</template>
