<script setup lang="ts">
import { onBeforeMount, ref } from 'vue'
import OnboardingOverlay from '@/components/common/OnboardingOverlay.vue'
import NpsWidget from '@/components/NpsWidget.vue'
import { Toaster } from '@/components/ui/toast'
import { useOnboarding } from '@/composables/useOnboarding'
import { startSSE, stopSSE } from '@/composables/useSSE'
import { useToken } from '@/composables/useToken'
import { useUser } from '@/composables/useUser'
import { startProactiveRefresh, stopProactiveRefresh } from '@/services/api'
import { authService } from '@/services/auth'
import { handleError } from '@/services/errorService'
import { profileService } from '@/services/profile'
import Layout from './components/layout/Layout.vue'

const { start: startOnboarding } = useOnboarding()

const tg_user = useUser()
const tg_token = useToken()
const isLoading = ref(false)
// sessionExpired — взводим, когда у юзера в localStorage был tg_token, но
// /me + refresh вернули 401 (токен инвалидирован — например, массовой
// миграцией в #325, либо протух). Раньше в этом случае молча редиректили на
// лендинг, и из-за того что tg_token не чистился, юзер попадал в петлю
// «платформа на полсекунды → лендинг → клик → платформа на полсекунды».
// Теперь apiClient чистит токены сам, а мы показываем человеческий экран
// с кнопкой в бот.
const sessionExpired = ref(false)
const botUrl = `https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_NAME}?start=from_site`
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
  const urlToken = urlParams.get('token')
  // Если токен и user уже в localStorage (юзер пришёл с лендинга), повторно
  // /api/auth/telegram не дёргаем — иначе при флакающем бэке мы рискуем
  // обнулить только что положенный auth-state из-за одной 5xx-ошибки.
  // /me ниже всё равно подтянет свежий профиль.
  if (urlToken) {
    isLoading.value = true
    authService
      .authenticate(urlToken)
      .then(({ user, token: authToken }) => {
        tg_user.value = { ...tg_user.value, ...user }
        tg_token.value = authToken
        window.history.replaceState({}, document.title, window.location.pathname)
        startSSE()
        startProactiveRefresh()
        // Авторизация возвращает не все поля профиля (нет subscriptionTier и т.п.)
        // — дёргаем /me чтобы перетереть tg_user актуальной версией.
        profileService.getMe().catch(() => {})
        // Запускаем онбординг с задержкой, чтобы дать DOM срендериться
        // (особенно важно для Safari/Telegram WebView)
        setTimeout(startOnboarding, 1500)
      })
      .catch((error) => {
        stopSSE()
        stopProactiveRefresh()
        tg_user.value = null
        tg_token.value = null
        handleError(error)
        // Битый одноразовый токен в URL — показываем экран «сессия истекла»
        // вместо редиректа на лендинг (откуда юзер тут же снова приедет).
        sessionExpired.value = true
      })
      .finally(() => {
        isLoading.value = false
      })
  }
  else if (tg_token.value) {
    startSSE()
    startProactiveRefresh()
    // tg_user мог протухнуть (TTL 1 час) — getMe сам поднимет его в localStorage
    // через apiClient (подставит tg_token в header). Пока ждём — isLoading,
    // чтобы не мигать голым Layout без router-view.
    if (!tg_user.value)
      isLoading.value = true
    profileService.getMe().finally(() => {
      isLoading.value = false
      // Если /me + refresh не помогли — apiClient уже почистил tg_token.
      // Показываем экран «сессия истекла» с кнопкой в бот, а не редирект:
      // юзер видит причину и может перелогиниться одним кликом.
      if (!tg_user.value)
        sessionExpired.value = true
    })
  }
  else if (!import.meta.env.DEV) {
    window.location.pathname = '/'
  }
})
</script>

<template>
  <Toaster />
  <OnboardingOverlay />
  <div v-if="sessionExpired" class="min-h-screen flex items-center justify-center px-6">
    <div class="max-w-md w-full text-center space-y-6">
      <h1 class="text-2xl font-semibold">
        Сессия истекла
      </h1>
      <p class="text-foreground/70">
        Войди снова через Telegram-бота — он откроет платформу одним кликом.
      </p>
      <a
        :href="botUrl"
        class="inline-block px-6 py-3 rounded-md bg-accent text-accent-foreground font-medium hover:opacity-90 transition-opacity"
      >
        Войти через Telegram
      </a>
      <div>
        <a href="/" class="text-sm text-foreground/50 hover:text-foreground/70 underline">
          Вернуться на главную
        </a>
      </div>
    </div>
  </div>
  <div v-else-if="!isLoading" class="min-h-screen flex flex-col">
    <Layout>
      <router-view v-if="tg_user" v-slot="{ Component }">
        <Transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </Transition>
      </router-view>
    </Layout>
    <NpsWidget v-if="tg_user" />
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
