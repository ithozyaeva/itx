<script setup lang="ts">
import type { Placement } from '@floating-ui/vue'
import { autoUpdate, flip, offset, shift, useFloating } from '@floating-ui/vue'
import { computed, onMounted, onUnmounted, ref, toRef } from 'vue'

const props = withDefaults(defineProps<{
  placement?: Placement
  offset?: number
}>(), {
  placement: 'bottom',
  offset: 8,
})

const isOpen = ref(false)
const triggerRef = ref<HTMLElement | null>(null)
const popoverRef = ref<HTMLElement | null>(null)

const placementRef = toRef(props, 'placement')
const offsetRef = toRef(props, 'offset')

const middleware = computed(() => [
  offset(offsetRef.value),
  flip(),
  shift({ padding: 5 }),
])

const { floatingStyles } = useFloating(triggerRef, popoverRef, {
  placement: placementRef,
  strategy: 'fixed',
  middleware,
  whileElementsMounted: autoUpdate,
})

function toggle() {
  isOpen.value = !isOpen.value
}

function handleClickOutside(e: MouseEvent) {
  if (
    popoverRef.value
    && triggerRef.value
    && !popoverRef.value.contains(e.target as Node)
    && !triggerRef.value.contains(e.target as Node)
  ) {
    isOpen.value = false
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <div class="popover-container">
    <div
      ref="triggerRef"
      class="trigger"
      @click="toggle"
    >
      <slot name="trigger" />
    </div>
    <Transition name="popover">
      <div
        v-if="isOpen"
        ref="popoverRef"
        class="popover"
        :style="floatingStyles"
      >
        <slot name="content" />
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.popover-container {
  display: inline-block;
}

.popover {
  background: var(--color-green-black-500);
  border: 1px solid var(--color-green-black-400);
  border-radius: var(--radius-card);
  z-index: 50;
  padding: 16px;
}

.popover-enter-active,
.popover-leave-active {
  transition: opacity var(--transition-popover);
}

.popover-enter-from,
.popover-leave-to {
  opacity: 0;
}
</style>
