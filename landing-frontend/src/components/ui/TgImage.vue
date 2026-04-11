<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import AvatarPlaceholderIcon from '@/components/ui/AvatarPlaceholderIcon.vue'

defineProps<{ username: string }>()

const imgRef = ref()
const isError = ref(false)
function handleLoad(event: Event) {
  if ((event.target as HTMLImageElement).naturalWidth <= 1) {
    isError.value = true
  }
}

function handleError() {
  isError.value = true
}

onMounted(() => {
  if (imgRef.value) {
    imgRef.value.addEventListener('load', handleLoad)
    imgRef.value.addEventListener('error', handleError)
  }
})

onUnmounted(() => {
  if (imgRef.value) {
    imgRef.value.removeEventListener('load', handleLoad)
    imgRef.value.removeEventListener('error', handleError)
  }
})
</script>

<template>
  <img
    v-if="!isError"
    ref="imgRef"
    :src="`https://t.me/i/userpic/160/${username}.jpg`"
    :alt="`Аватар ${username}`"
    @error="isError = true"
  >
  <AvatarPlaceholderIcon
    v-else
    class="flex-shrink-0"
  />
</template>
