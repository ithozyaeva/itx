<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
</script>

<template>
  <div class="rounded-2xl border bg-card border-border overflow-hidden">
    <Skeleton class="w-full h-48 rounded-none" />
    <div class="p-4">
      <div class="flex items-start justify-between gap-2 mb-2">
        <Skeleton class="h-4 w-2/3" />
        <Skeleton class="h-5 w-16 rounded-full shrink-0" />
      </div>
      <Skeleton class="h-4 w-20 mb-2" />
      <Skeleton class="h-5 w-16 rounded-full mb-2" />
      <Skeleton class="h-3.5 w-32 mb-1" />
      <Skeleton class="h-3.5 w-28 mb-3" />
      <Skeleton class="h-8 w-24 rounded-lg mt-3" />
    </div>
  </div>
</template>
