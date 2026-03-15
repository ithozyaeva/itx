<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
</script>

<template>
  <div class="bg-card rounded-3xl border p-4 flex flex-col gap-2">
    <Skeleton class="h-6 w-2/3 rounded-lg" />
    <Skeleton class="h-4 w-1/2" />
    <Skeleton class="h-4 w-3/4" />
    <div class="flex flex-wrap gap-1 mt-1">
      <Skeleton class="h-5 w-14 rounded-full" />
      <Skeleton class="h-5 w-18 rounded-full" />
      <Skeleton class="h-5 w-12 rounded-full" />
    </div>
    <Skeleton class="h-4 w-24 mt-auto" />
  </div>
</template>
