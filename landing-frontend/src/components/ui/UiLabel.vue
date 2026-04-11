<script setup lang="ts">
defineProps<{
  as?: string
}>()
</script>

<template>
  <component
    :is="as || 'span'"
    class="ui-label"
  >
    <label class="ui-label-text">
      <slot />
    </label>
  </component>
</template>

<style scoped>
.ui-label-text {
  font-family: var(--label-font), sans-serif;
  font-weight: var(--label-weight);
  font-size: var(--label-size);
  line-height: var(--label-line-height);
  letter-spacing: var(--label-letter-spacing);
}

.ui-label {
  border-radius: var(--radius-button);
  color: var(--color-green-500);
  box-sizing: border-box;
  background-color: rgba(var(--color-green-500-rgb), 0.1);
  padding: 4px 12px;
}
</style>
