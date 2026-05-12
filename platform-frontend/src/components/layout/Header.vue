<script setup lang="ts">
import { LogOut, Menu } from 'lucide-vue-next'
import NotificationDropdown from '@/components/NotificationDropdown.vue'
import ThemeToggle from '@/components/ui/theme-toggle.vue'
import { useSidebar } from '@/composables/useSidebar'
import { stopSSE } from '@/composables/useSSE'
import { getTelegramWebApp, isMiniApp } from '@/composables/useTelegramWebApp'
import { useUser } from '@/composables/useUser'
import { stopProactiveRefresh } from '@/services/api'
import { authService } from '@/services/auth'

const user = useUser()
const { toggleSidebar } = useSidebar()

// authService.logout инвалидирует токен серверно; stopSSE/stopProactiveRefresh
// гасят фоновые соединения, чтобы не дёргали мёртвый токен в окно между
// инвалидацией и навигацией на лендинг.
// Внутри Mini App навигация на '/' триггерит auto-relogin: App.vue видит
// insideMiniApp + !tg_token и молча обменивает (всё ещё валидный) initData
// на новый токен — кнопка «Выйти» становится no-op. В Mini App закрываем
// окно через tg.close(); из браузера — редирект на лендинг как раньше.
async function logout() {
  stopProactiveRefresh()
  stopSSE()
  await authService.logout()
  user.value = null
  if (isMiniApp()) {
    getTelegramWebApp()?.close?.()
    return
  }
  window.location.pathname = '/'
}
</script>

<template>
  <header
    class="sticky top-0 z-40 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
  >
    <div class="container mx-auto px-3 sm:px-4 flex h-14 sm:h-16 items-center justify-between gap-2">
      <button
        class="p-2 rounded-sm bg-background border border-border md:hidden"
        aria-label="Открыть меню"
        @click="toggleSidebar"
      >
        <Menu class="w-6 h-6" />
      </button>
      <a
        href="/"
        class="flex items-center gap-2 font-bold text-xl"
      >
        <svg
          class="w-20"
          viewBox="0 0 238 104"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M39.7609 14.8145V29.6038H54.6269H69.4929V65.8376V102.072H83.6156H97.7383V65.8376V29.6038H112.604H127.47V14.8145V0.0251393H83.6156H39.7609V14.8145ZM0 52.5273V102.072H14.866H29.732V52.5273V0H14.866H0V52.5273Z"
            fill="currentColor"
          />
          <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M143.432 9.72996L133.671 19.4585L149.948 35.73L166.223 52L149.948 68.27L133.671 84.5415L143.432 94.2701L153.191 104L169.514 87.7747L185.836 71.5509L202.157 87.7747L218.48 104L228.239 94.2701L238 84.5415L221.723 68.27L205.448 52L221.723 35.73L238 19.4585L228.239 9.72996L218.48 0L202.157 16.2253L185.836 32.4491L169.514 16.2253L153.191 0L143.432 9.72996Z"
            fill="currentColor"
          />
        </svg>
      </a>
      <span class="hidden md:flex items-center gap-1.5 font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
        <span class="w-1.5 h-1.5 rounded-full bg-accent shadow-[0_0_6px_hsl(var(--accent))]" />
        online
      </span>
      <div class="flex-1" />
      <div
        v-if="user"
        class="flex items-center gap-3"
      >
        <ThemeToggle />
        <NotificationDropdown />
        <span class="mr-1 text-sm hidden sm:inline">
          {{ [user.firstName, user.lastName?.[0]].filter(Boolean).join(' ') }}
        </span>
        <button
          aria-label="Выйти"
          @click="logout"
        >
          <LogOut
            class="h-5 w-5 cursor-pointer text-muted-foreground hover:text-foreground transition-colors"
          />
        </button>
      </div>
    </div>
  </header>
</template>
