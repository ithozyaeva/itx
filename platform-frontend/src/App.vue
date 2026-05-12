<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { onBeforeMount, ref } from 'vue'
import OnboardingOverlay from '@/components/common/OnboardingOverlay.vue'
import NpsWidget from '@/components/NpsWidget.vue'
import ReferralWelcome from '@/components/ReferralWelcome.vue'
import { Toaster } from '@/components/ui/toast'
import { useOnboarding } from '@/composables/useOnboarding'
import { startSSE, stopSSE } from '@/composables/useSSE'
import { getTelegramWebApp, initTelegramWebApp, isMiniApp } from '@/composables/useTelegramWebApp'
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
const insideMiniApp = isMiniApp()
// sessionExpired — взводим, когда у юзера в localStorage был tg_token, но
// /me + refresh вернули 401 (токен инвалидирован — например, массовой
// миграцией в #325, либо протух). Раньше в этом случае молча редиректили на
// лендинг, и из-за того что tg_token не чистился, юзер попадал в петлю
// «платформа на полсекунды → лендинг → клик → платформа на полсекунды».
// Теперь apiClient чистит токены сам, а мы показываем человеческий экран
// с кнопкой в бот.
const sessionExpired = ref(false)
const botUrl = `https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_NAME}?start=from_site`
// loginViaMiniApp — забираем initData у Telegram-клиента и обмениваем его
// на tg_token. Бэк валидирует HMAC, поэтому подделать ничего нельзя.
// Используется и при первом заходе из чата с ботом, и как fallback при
// протухшем токене (вместо экрана «Сессия истекла»).
async function loginViaMiniApp(): Promise<boolean> {
  const tg = getTelegramWebApp()
  if (!tg || !tg.initData)
    return false
  try {
    const { user, token: authToken } = await authService.authenticateWebApp(tg.initData)
    tg_user.value = { ...tg_user.value, ...user }
    tg_token.value = authToken
    startSSE()
    startProactiveRefresh()
    profileService.getMe().catch(() => {})
    setTimeout(startOnboarding, 1500)
    return true
  }
  catch (error) {
    stopSSE()
    stopProactiveRefresh()
    tg_user.value = null
    tg_token.value = null
    handleError(error)
    return false
  }
}

onBeforeMount(async () => {
  // Инициализация темы при запуске приложения
  const savedTheme = localStorage.getItem('theme')

  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark')
  }
  else {
    document.documentElement.classList.remove('dark')
  }

  // Внутри Telegram сразу даём знать клиенту, что мы готовы (иначе чёрный
  // экран до первой отрисовки) и разворачиваемся на полный viewport.
  if (insideMiniApp)
    initTelegramWebApp()

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
    profileService.getMe().finally(async () => {
      // Если /me + refresh не помогли — apiClient уже почистил tg_token.
      // Внутри miniapp молча перелогиниваемся через initData; снаружи
      // (десктоп-браузер) — показываем экран «сессия истекла» с кнопкой
      // в бот, чтобы юзер видел причину и мог перелогиниться одним кликом.
      if (!tg_user.value) {
        if (insideMiniApp && await loginViaMiniApp()) {
          isLoading.value = false
          return
        }
        sessionExpired.value = true
      }
      isLoading.value = false
    })
  }
  else if (insideMiniApp) {
    // Свежий заход через menu button бота: токена ещё нет, но Telegram уже
    // подсунул нам подписанный initData — обмениваем его на сессию.
    isLoading.value = true
    const ok = await loginViaMiniApp()
    isLoading.value = false
    if (!ok)
      sessionExpired.value = true
  }
  else if (!import.meta.env.DEV) {
    window.location.pathname = '/'
  }
})
</script>

<template>
  <Toaster />
  <OnboardingOverlay />
  <div v-if="sessionExpired" class="min-h-[calc(100dvh-var(--safe-y))] flex items-center justify-center px-6">
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
  <div v-else-if="isLoading" class="min-h-[calc(100dvh-var(--safe-y))] flex items-center justify-center" aria-live="polite" aria-busy="true">
    <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    <span class="sr-only">Загрузка…</span>
  </div>
  <div v-else class="min-h-[calc(100dvh-var(--safe-y))] flex flex-col">
    <Layout>
      <router-view v-if="tg_user" v-slot="{ Component }">
        <Transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </Transition>
      </router-view>
    </Layout>
    <NpsWidget v-if="tg_user" />
    <ReferralWelcome v-if="tg_user" />
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
