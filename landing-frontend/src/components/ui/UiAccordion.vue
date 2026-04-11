<script setup lang="ts">
import { ref } from 'vue'
import UiTypography from './UiTypography.vue'

defineProps<{
  title: string
  content?: string
  defaultOpen?: boolean
}>()

const isOpen = ref(false)

function toggle() {
  isOpen.value = !isOpen.value
}
</script>

<template>
  <div
    class="accordion"
    :class="{ open: isOpen }"
  >
    <button
      type="button"
      class="accordion-header"
      :aria-expanded="isOpen"
      @click="toggle"
    >
      <UiTypography
        variant="h4"
        as="h4"
        class="accordion-title"
      >
        {{ title }}
      </UiTypography>
      <div class="accordion-icon">
        <svg
          v-if="!isOpen"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 40 40"
        >
          <path
            stroke="currentColor"
            stroke-linecap="round"
            stroke-width="2"
            d="M20 5v30M35 20H5"
          />
        </svg>
        <svg
          v-else
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 40 40"
        >
          <path
            stroke="currentColor"
            stroke-linecap="round"
            stroke-width="2"
            d="M31 9 9 31M31 31 9 9"
          />
        </svg>
      </div>
    </button>
    <div
      class="accordion-wrapper"
      :class="{ 'accordion-wrapper--open': isOpen }"
    >
      <div class="accordion-content">
        <UiTypography variant="body-l">
          <slot>{{ content }}</slot>
        </UiTypography>
      </div>
    </div>
  </div>
</template>

<style scoped>
.accordion {
  box-sizing: border-box;
  background-color: var(--color-green-black-600);
  border-radius: var(--radius-accordion);
  cursor: pointer;
  transition: var(--transition-default);
  padding: 16px 24px;
}

.accordion:hover {
  background-color: var(--color-green-black-500);
}

.accordion:hover .accordion-icon {
  background-color: var(--color-green-black-400);
}

.accordion-header {
  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;
  justify-content: space-between;
  align-items: center;
  gap: 14px;
  display: flex;
  width: 100%;
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  font: inherit;
  color: inherit;
  text-align: left;
}

.accordion-title {
  color: var(--color-green-700);
  transition: color 0.2s;
}

.accordion-icon {
  border-radius: var(--radius-circle);
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  transition: color 0.2s;
}

.accordion-icon svg {
  width: 100%;
  height: 100%;
}

.accordion-wrapper {
  display: grid;
  grid-template-rows: 0fr;
  transition: grid-template-rows 0.25s;
}

.accordion-wrapper--open {
  grid-template-rows: 1fr;
}

.accordion-content {
  color: var(--color-light-grey);
  overflow: hidden;
}
</style>
