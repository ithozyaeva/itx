<script setup lang="ts">
import { isUserSubscribed, useUser } from '@/composables/useUser'
import { authService } from '@/services/auth'
import { reviewService } from '@/services/reviews'
import { LogOut } from 'lucide-vue-next'
import { ref } from 'vue'
import ReviewModal from './ReviewModal.vue'

const user = useUser()
const isModalOpen = ref(false)
const isSubscribed = isUserSubscribed()

function logout() {
  authService.clearAuthHeader()
  user.value = null
  window.location.pathname = '/'
}

async function handleSaveReview(text: string) {
  await reviewService
    .createReview(text)
    .finally(() => {
      isModalOpen.value = false
    })
}
</script>

<template>
  <header
    class="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
  >
    <div class="container mr-auto ml-auto px-2 flex h-16 items-center justify-between">
      <a href="/" class="flex items-center gap-2 font-bold text-xl">
        <svg class="w-20" viewBox="0 0 238 104" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path fill-rule="evenodd" clip-rule="evenodd" d="M39.7609 14.8145V29.6038H54.6269H69.4929V65.8376V102.072H83.6156H97.7383V65.8376V29.6038H112.604H127.47V14.8145V0.0251393H83.6156H39.7609V14.8145ZM0 52.5273V102.072H14.866H29.732V52.5273V0H14.866H0V52.5273Z" fill="currentColor" />
          <path fill-rule="evenodd" clip-rule="evenodd" d="M143.432 9.72996L133.671 19.4585L149.948 35.73L166.223 52L149.948 68.27L133.671 84.5415L143.432 94.2701L153.191 104L169.514 87.7747L185.836 71.5509L202.157 87.7747L218.48 104L228.239 94.2701L238 84.5415L221.723 68.27L205.448 52L221.723 35.73L238 19.4585L228.239 9.72996L218.48 0L202.157 16.2253L185.836 32.4491L169.514 16.2253L153.191 0L143.432 9.72996Z" fill="currentColor" />
        </svg>
      </a>
      <div class="flex-1" />
      <div v-if="user" class="flex items-center gap-3">
        <button
          v-if="isSubscribed"
          class="cursor-pointer border border-border rounded-full px-4 py-1 hover:bg-secondary transition duration-300 active:scale-95"
          @click="isModalOpen = true"
        >
          Добавить отзыв
        </button>
        <span class="mr-3">
          {{ `${user.firstName} ${user.lastName[0] ?? ""}` }}
        </span>
        <LogOut class="h-6 w-6 cursor-pointer" @click="logout" />
      </div>
    </div>
  </header>

  <ReviewModal :is-open="isModalOpen" @close="isModalOpen = false" @save="handleSaveReview" />
</template>
