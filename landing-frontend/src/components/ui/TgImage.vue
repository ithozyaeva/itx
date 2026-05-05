<script setup lang="ts">
import { computed, ref } from 'vue'
import AvatarPlaceholderIcon from '@/components/ui/AvatarPlaceholderIcon.vue'

const props = defineProps<{
  username: string
  avatarUrl?: string
}>()

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
    loading="lazy"
    style="aspect-ratio: 1 / 1;"
    @load="handleLoad"
    @error="isError = true"
  >
  <AvatarPlaceholderIcon
    v-else
    class="flex-shrink-0"
  />
</template>
