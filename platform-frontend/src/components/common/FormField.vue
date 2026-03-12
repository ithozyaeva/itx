<script setup lang="ts">
const props = defineProps<{
  label: string
  error?: string
  htmlFor?: string
  required?: boolean
}>()
</script>

<template>
  <div>
    <label
      :for="props.htmlFor"
      class="block text-sm font-medium mb-1"
    >
      {{ label }}
      <span
        v-if="required"
        class="text-destructive"
      >*</span>
    </label>
    <slot />
    <p
      v-if="error"
      class="mt-1 text-xs text-destructive"
      role="alert"
    >
      {{ error }}
    </p>
  </div>
</template>
