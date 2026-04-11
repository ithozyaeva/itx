<script setup lang="ts">
defineProps<{
  as?: string
  variant?: 'default' | 'active'
  disabled?: boolean
}>()
</script>

<template>
  <component
    :is="as || 'button'"
    class="ui-tag"
    :class="[variant || 'default', { disabled }]"
    :disabled="disabled"
    :aria-disabled="disabled"
  >
    <span class="ui-tag-text">
      <slot />
    </span>
  </component>
</template>

<style scoped>
.ui-tag-text {
  font-family: var(--body-l-font), sans-serif;
  font-weight: var(--body-l-weight);
  font-size: var(--body-l-size);
  line-height: var(--body-l-line-height);
  letter-spacing: var(--body-l-letter-spacing);
}

.ui-tag {
  border-radius: var(--radius-default);
  color: var(--color-green-700);
  transition: var(--transition-default);
  box-sizing: border-box;
  cursor: pointer;
  border: none;
  outline: 4px solid transparent;
  padding: 8px 20px;
}

.ui-tag.default {
  outline: 1px solid var(--color-green-700);
  background-color: transparent;
  pointer-events: none;
}

.ui-tag.default:hover:not(.disabled) {
  outline-color: var(--color-white);
  color: var(--color-white);
}

.ui-tag.default:focus:not(.disabled) {
  outline: 4px solid var(--color-white);
}

.ui-tag.default:active:not(.disabled) {
  outline-color: var(--color-grey);
  color: var(--color-grey);
}

.ui-tag.active {
  background-color: var(--color-green-700);
  color: var(--color-green-black-700);
  outline: none;
}

.ui-tag.active:hover:not(.disabled) {
  background-color: var(--color-white);
}

.ui-tag.active:focus:not(.disabled) {
  outline: 4px solid var(--color-white);
}

.ui-tag.active:active:not(.disabled) {
  background-color: var(--color-grey);
  border-color: transparent;
}

.ui-tag.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.ui-tag.disabled.default {
  outline-color: var(--color-green-700);
  color: var(--color-green-700);
}

.ui-tag.disabled.active {
  background-color: var(--color-green-700);
  color: var(--color-green-black-700);
}
</style>
