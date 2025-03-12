<template>
  <div class="rounded-lg border bg-card text-card-foreground shadow-sm">
    <div v-if="$slots.header" class="flex flex-col space-y-1.5 p-6">
      <slot name="header" />
    </div>
    <div v-if="$slots.content" class="p-6 pt-0">
      <slot name="content" />
    </div>
    <div v-if="$slots.footer" class="flex items-center p-6 pt-0">
      <slot name="footer" />
    </div>
  </div>
</template>
