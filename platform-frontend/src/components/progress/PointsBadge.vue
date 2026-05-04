<script setup lang="ts">
withDefaults(defineProps<{
  amount: number
  earned?: boolean
  size?: 'sm' | 'md'
}>(), {
  earned: false,
  size: 'sm',
})
</script>

<template>
  <span
    class="font-bold rounded-full shrink-0"
    :class="[
      earned
        ? 'text-green-500 bg-green-500/10'
        : 'text-yellow-500 bg-yellow-500/10',
      size === 'sm' ? 'text-xs px-2 py-0.5' : 'text-sm px-2.5 py-1',
    ]"
  >
    <template v-if="!earned">+</template>{{ amount }}
  </span>
</template>
