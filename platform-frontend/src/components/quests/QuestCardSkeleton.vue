<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
</script>

<template>
  <div class="rounded-sm border bg-card p-5">
    <div class="flex items-start gap-3">
      <Skeleton class="w-10 h-10 rounded-xl shrink-0" />
      <div class="flex-1 min-w-0">
        <div class="flex items-start justify-between gap-2">
          <Skeleton class="h-4 w-2/3" />
          <Skeleton class="h-5 w-12 rounded-full shrink-0" />
        </div>
        <Skeleton class="h-3 w-full mt-1.5" />
      </div>
    </div>
    <div class="mt-4">
      <div class="flex items-center justify-between mb-1.5">
        <Skeleton class="h-3 w-28" />
        <Skeleton class="h-3 w-8" />
      </div>
      <Skeleton class="w-full h-1.5 rounded-full" />
    </div>
    <div class="flex items-center gap-3 mt-3">
      <Skeleton class="h-3 w-16" />
      <Skeleton class="h-3 w-20" />
    </div>
  </div>
</template>
