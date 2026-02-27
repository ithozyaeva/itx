<script setup lang="ts">
import { useToken } from '@/composables/useToken'
import { useUser } from '@/composables/useUser'
import { authService } from '@/services/auth'
import { onBeforeMount, ref } from 'vue'
import Layout from './components/layout/Layout.vue'

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
      })
      .catch((error) => {
        console.error('Authentication failed:', error)
      })
      .finally(() => {
        isLoading.value = false
      })
  }
  else if (!tg_user.value && !import.meta.env.DEV) {
    window.location.pathname = '/'
  }
})
</script>

<template>
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
