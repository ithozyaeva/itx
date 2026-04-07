<script setup lang="ts">
import { onBeforeMount, ref } from 'vue'
import OnboardingOverlay from '@/components/common/OnboardingOverlay.vue'
import { Toaster } from '@/components/ui/toast'
import { useOnboarding } from '@/composables/useOnboarding'
import { startSSE, stopSSE } from '@/composables/useSSE'
import { useToken } from '@/composables/useToken'
import { useUser } from '@/composables/useUser'
import { authService } from '@/services/auth'
import { handleError } from '@/services/errorService'
import Layout from './components/layout/Layout.vue'

const { start: startOnboarding } = useOnboarding()

const tg_user = useUser()
const tg_token = useToken()
const isLoading = ref(false)
onBeforeMount(() => {
  // Инициализация темы при запуске приложения
  const savedTheme = localStorage.getItem('theme')

  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark')
  }
  else {
    document.documentElement.classList.remove('dark')
  }

  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token') || tg_token.value
  if (token) {
    isLoading.value = true
    authService
      .authenticate(token)
      .then(({ user, token: authToken }) => {
        tg_user.value = { ...tg_user.value, ...user }
        tg_token.value = authToken
        window.history.replaceState({}, document.title, window.location.pathname)
        startSSE()
        // Запускаем онбординг с задержкой, чтобы дать DOM срендериться
        // (особенно важно для Safari/Telegram WebView)
        setTimeout(startOnboarding, 1500)
      })
      .catch((error) => {
        stopSSE()
        tg_user.value = null
        tg_token.value = null
        handleError(error)
      })
      .finally(() => {
        isLoading.value = false
      })
  }
  else if (tg_user.value) {
    startSSE()
  }
  else if (!tg_user.value && !import.meta.env.DEV) {
    window.location.pathname = '/'
  }
})
</script>

<template>
  <Toaster />
  <OnboardingOverlay />
  <div v-if="!isLoading" class="min-h-screen flex flex-col">
    <Layout>
      <router-view v-if="tg_user" v-slot="{ Component }">
        <Transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </Transition>
      </router-view>
    </Layout>
  </div>
</template>

<style scoped>
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.2s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
