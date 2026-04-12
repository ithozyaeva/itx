<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
</script>

<template>
  <div class="bg-card rounded-sm border p-4 flex flex-col gap-2">
    <Skeleton class="h-6 w-3/4 rounded-lg" />
    <div class="flex gap-1">
      <Skeleton class="h-5 w-16 rounded-full" />
      <Skeleton class="h-5 w-20 rounded-full" />
    </div>
    <div class="space-y-1.5">
      <Skeleton class="h-3.5 w-full" />
      <Skeleton class="h-3.5 w-5/6" />
    </div>
    <div class="flex items-center gap-2 mt-1">
      <Skeleton class="h-4 w-4" />
      <Skeleton class="h-3.5 w-32" />
    </div>
    <div class="flex items-center gap-2">
      <Skeleton class="h-4 w-4" />
      <Skeleton class="h-3.5 w-24" />
    </div>
    <div class="flex items-center gap-2">
      <Skeleton class="h-4 w-4" />
      <Skeleton class="h-3.5 w-28" />
    </div>
    <Skeleton class="h-9 w-32 rounded-xl mt-2" />
  </div>
</template>
