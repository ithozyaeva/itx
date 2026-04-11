<script setup lang="ts">
defineProps<{
  as?: string
  variant?: 'filled' | 'stroke' | 'dark-filled'
  disabled?: boolean
}>()
</script>

<template>
  <component
    :is="as || 'button'"
    :disabled="disabled"
    class="ui-button"
    :class="[variant || 'filled', { disabled }]"
  >
    <span class="ui-button-text">
      <slot />
    </span>
  </component>
</template>

<style scoped>
.ui-button-text {
  font-family: var(--button-font), sans-serif;
  font-weight: var(--button-weight);
  font-size: var(--button-size);
  line-height: var(--button-line-height);
  letter-spacing: var(--button-letter-spacing);
  text-transform: var(--button-transform);
}

.ui-button {
  border-radius: var(--radius-button);
  color: var(--color-green-black-700);
  transition: var(--transition-default);
  box-sizing: border-box;
  cursor: pointer;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  border: 4px solid transparent;
  align-content: center;
  min-height: 48px;
  padding: 11px 26px;
}

.ui-button.filled {
  background-color: var(--color-green-700);
}

.ui-button.filled:not(.disabled):hover {
  background-color: var(--color-white);
}

.ui-button.filled:not(.disabled):focus-visible {
  outline: 2px solid var(--color-green-700);
  outline-offset: 2px;
}

.ui-button.filled:not(.disabled):active {
  background-color: var(--color-grey);
  border-color: transparent;
}

.ui-button.filled.disabled {
  background-color: var(--color-green-black-400);
  color: var(--color-green-black-600);
  cursor: not-allowed;
}

.ui-button.stroke {
  border: 2px solid var(--color-green-700);
  color: var(--color-green-700);
  background-color: transparent;
}

.ui-button.stroke:not(.disabled):hover {
  border-color: var(--color-white);
  color: var(--color-white);
}

.ui-button.stroke:not(.disabled):focus-visible {
  outline: 2px solid var(--color-green-700);
  outline-offset: 2px;
}

.ui-button.stroke:not(.disabled):active {
  border-color: var(--color-grey);
  color: var(--color-grey);
}

.ui-button.stroke.disabled {
  border-color: var(--color-green-black-400);
  color: var(--color-green-black-400);
  cursor: not-allowed;
}

.ui-button.dark-filled {
  background-color: var(--color-green-black-700);
  color: var(--color-green-700);
}

.ui-button.dark-filled:not(.disabled):hover {
  background-color: var(--color-green-black-400);
}

.ui-button.dark-filled:not(.disabled):focus-visible {
  outline: 2px solid var(--color-green-700);
  outline-offset: 2px;
}

.ui-button.dark-filled:not(.disabled):active {
  background-color: var(--color-green-black-500);
  border-color: transparent;
}

.ui-button.dark-filled.disabled {
  background-color: var(--color-green-black-600);
  color: var(--color-green-black-400);
  cursor: not-allowed;
}
</style>
