<script setup lang="ts">
import { useAttrs } from 'vue'

defineOptions({
  inheritAttrs: false,
})

defineProps<{
  as?: string
  variant?: 'default' | 'active'
}>()

const attrs = useAttrs()
</script>

<template>
  <component
    :is="as || 'a'"
    class="ui-menu-item"
    :class="[variant || 'default']"
    v-bind="attrs"
  >
    <span class="ui-menu-item-text">
      <slot />
    </span>
  </component>
</template>

<style scoped>
.ui-menu-item-text {
  font-family: var(--button-font), sans-serif;
  font-weight: var(--button-weight);
  font-size: var(--button-size);
  line-height: var(--button-line-height);
  letter-spacing: var(--button-letter-spacing);
  text-transform: var(--button-transform);
}

.ui-menu-item {
  color: var(--color-light-grey);
  cursor: pointer;
  transition: var(--transition-default);
  text-decoration: none;
  display: inline-block;
}

.ui-menu-item.default:hover {
  color: var(--color-green-700);
}

.ui-menu-item.default:focus {
  color: var(--color-light-grey);
  border: 4px solid var(--color-white);
  border-radius: 8px;
  outline: none;
}

.ui-menu-item.default:active {
  color: var(--color-green-500);
}

.ui-menu-item.active {
  color: var(--color-green-700);
  text-decoration: underline;
}

.ui-menu-item.active:focus {
  border: 4px solid var(--color-white);
  border-radius: 8px;
  outline: none;
}
</style>
