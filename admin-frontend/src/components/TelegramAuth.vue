<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { loginWithTelegram } from '@/services/authService'
import { handleError } from '@/services/errorService'

const emit = defineEmits<{
  (e: 'authSuccess'): void
}>()

const isLoading = ref(false)

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token')

  if (token) {
    handleTelegramAuth(token)
    window.history.replaceState({}, document.title, window.location.pathname)
  }
})

async function handleTelegramAuth(token: string) {
  isLoading.value = true
  try {
    const response = await loginWithTelegram(token)
    if (response) {
      emit('authSuccess')
    }
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

function openTelegramBot() {
  window.open(`https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_NAME}?start=from_admin`, '_blank')
}
</script>

<template>
  <div class="telegram-auth">
    <Button :disabled="isLoading" @click="openTelegramBot">
      {{ isLoading ? 'Authenticating...' : 'Зайти через ТГ' }}
    </Button>
    <div v-if="isLoading" class="loading mt-2 text-sm text-muted-foreground">
      Authenticating...
    </div>
  </div>
</template>
