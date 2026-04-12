<script setup lang="ts">
const props = defineProps<{
  class?: string
}>()
</script>

<template>
  <div
    class="animate-pulse rounded-md bg-muted"
    :class="props.class"
  />
</template>
