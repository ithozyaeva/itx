<script setup lang="ts">
import { computed, ref } from 'vue'
import AvatarPlaceholderIcon from '@/components/ui/AvatarPlaceholderIcon.vue'

const props = withDefaults(defineProps<{
  username: string
  avatarUrl?: string
  // eager=true для above-the-fold использования (LCP-критичные секции).
  // По умолчанию lazy, потому что реальные usages — в below-the-fold.
  eager?: boolean
}>(), {
  eager: false,
})

const isError = ref(false)

const imgSrc = computed(() => {
  if (props.avatarUrl)
    return props.avatarUrl
  return `https://t.me/i/userpic/160/${props.username}.jpg`
})

function handleLoad(event: Event) {
  if ((event.target as HTMLImageElement).naturalWidth <= 1) {
    isError.value = true
  }
}
</script>

<template>
  <img
    v-if="!isError"
    :src="imgSrc"
    :alt="`Аватар ${username}`"
    width="160"
    height="160"
    decoding="async"
    :loading="eager ? 'eager' : 'lazy'"
    :fetchpriority="eager ? 'high' : 'auto'"
    @load="handleLoad"
    @error="isError = true"
  >
  <AvatarPlaceholderIcon
    v-else
    class="flex-shrink-0"
  />
</template>
