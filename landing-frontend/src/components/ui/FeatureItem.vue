<script setup lang="ts">
import type { Component } from 'vue'
import { Typography } from 'itx-ui-kit'

defineProps<{
  icon: Component
  title: string
  description: string
}>()
</script>

<template>
  <div class="flex flex-col items-center align-middle justify-between py-10 gap-4 bg-primary rounded-3xl ">
    <div class="flex shrink-0 items-center justify-center w-full px-5 lg:px-10 min-h-32 aspect-[2/1] sm:aspect-auto ">
      <component
        :is="icon"
        class="w-full "
      />
    </div>
    <div class="flex flex-col gap-3 items-center text-center px-5 sm:px-28 lg:px-16 ">
      <Typography
        variant="h3"
        as="h3"
        class="uppercase text-accent"
      >
        {{ title }}
      </Typography>
      <Typography
        variant="body-l"
        as="p"
      >
        {{ description }}
      </Typography>
    </div>
  </div>
</template>
