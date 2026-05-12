<script setup lang="ts">
import { Shield, X } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Typography } from '@/components/ui/typography'
import { useSidebar } from '@/composables/useSidebar'
import { canViewAdminPanel, hasMinTier, isUserSubscribed, useUser, useUserLevel } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { reviewService } from '@/services/reviews'
import ReviewModal from '../ReviewModal.vue'

const { sidebarGroups, isOpen, toggleSidebar } = useSidebar()
const isSubscribedRef = isUserSubscribed()

// Скрываем платные пункты меню для UNSUBSCRIBER и наоборот пункт «Тарифы»
// прячем у подписчиков. Группы без видимых пунктов целиком убираем,
// чтобы не оставались пустые «// Сообщество».
const visibleGroups = computed(() => {
  return sidebarGroups.value
    .map(group => ({
      ...group,
      items: group.items.filter((item) => {
        if (item.requiresSubscription && !isSubscribedRef.value)
          return false
        if (item.visibleFor === 'unsubscribed' && isSubscribedRef.value)
          return false
        if (item.requiresMinTier && !hasMinTier(item.requiresMinTier).value)
          return false
        return true
      }),
    }))
    .filter(group => group.items.length > 0)
})
const route = useRoute()
const router = useRouter()
const user = useUser()
const { level, levelIndex, maxLevel } = useUserLevel()

function isActive(path: string) {
  if (path === '/')
    return route.path === '/'
  return route.path === path || route.path.startsWith(`${path}/`)
}

function navigateTo(path: string) {
  router.push(path)
  isOpen.value = false
}

const isModalOpen = ref(false)
const isSubscribed = isSubscribedRef
const isAdmin = canViewAdminPanel()

async function handleSaveReview(text: string) {
  try {
    await reviewService.createReview(text)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isModalOpen.value = false
  }
}
</script>

<template>
  <div>
    <div
      data-onboarding="sidebar"
      class="fixed inset-y-0 left-0 md:static md:h-[calc(100dvh-var(--safe-y))] pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)] pl-[env(safe-area-inset-left)] pr-[env(safe-area-inset-right)] md:p-0 border-r border-sidebar-border bg-sidebar text-sidebar-foreground transition-all duration-300 w-full md:w-56 z-40"
      :class="[
        isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0',
      ]"
    >
      <div class="flex flex-col justify-between items-center pb-2 h-full overflow-y-auto">
        <div class="w-full">
          <div class="flex items-center justify-between px-4 py-3 md:hidden">
            <Typography
              variant="h4"
              as="h2"
            >
              Меню
            </Typography>
            <button
              class="p-2"
              aria-label="Закрыть меню"
              @click="toggleSidebar"
            >
              <X class="w-6 h-6" />
            </button>
          </div>
          <!-- Terminal system header -->
          <div class="px-4 py-3 border-b border-sidebar-border hidden md:flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-accent shadow-[0_0_6px_hsl(var(--accent))]" />
            <span class="font-mono text-xs font-medium tracking-wider text-sidebar-foreground/80">[SYS] IT-X</span>
          </div>
          <div class="flex-1 py-4">
            <div
              v-for="(group, groupIndex) in visibleGroups"
              :key="groupIndex"
              :class="groupIndex > 0 ? 'mt-4' : ''"
            >
              <p
                v-if="group.label"
                class="px-4 mb-1.5 text-[11px] font-mono font-medium uppercase tracking-wider text-sidebar-foreground/40"
              >
                // {{ group.label }}
              </p>
              <ul class="space-y-0.5">
                <li
                  v-for="(item, itemIndex) in group.items"
                  :key="item.path"
                  :data-onboarding="item.dataOnboarding"
                >
                  <Button
                    variant="ghost"
                    class="w-full justify-start py-2 cursor-pointer rounded-sm text-sidebar-foreground hover:bg-accent hover:text-accent-foreground"
                    :class="[
                      isActive(item.path) ? 'bg-accent text-accent-foreground' : '',
                    ]"
                    @click="navigateTo(item.path)"
                  >
                    <span class="font-mono text-[10px] text-sidebar-foreground/40 mr-1.5 w-4 shrink-0">{{ String(itemIndex + 1).padStart(2, '0') }}</span>
                    <component
                      :is="item.icon"
                      class="h-5 w-5 mr-2"
                    />
                    <span>{{ item.title }}</span>
                    <span
                      v-if="item.indicator"
                      class="ml-auto h-2 w-2 rounded-full bg-green-500"
                    />
                  </Button>
                </li>
              </ul>
            </div>
          </div>
        </div>
        <div class="flex flex-col items-center gap-2 w-full px-3 shrink-0">
          <a
            v-if="isAdmin"
            href="/admin"
            class="flex items-center gap-2 cursor-pointer border border-border/20 rounded-full px-4 py-1 text-sidebar-foreground hover:bg-accent hover:text-accent-foreground transition duration-300 active:scale-95"
          >
            <Shield class="h-4 w-4" />
            <span>Админ-панель</span>
          </a>
          <button
            v-if="isSubscribed"
            class="cursor-pointer border border-border/20 rounded-full px-4 py-1 text-sidebar-foreground hover:bg-accent hover:text-accent-foreground transition duration-300 active:scale-95"
            @click="isModalOpen = true"
          >
            Добавить отзыв
          </button>

          <!-- User profile card -->
          <div
            v-if="user"
            data-onboarding="profile"
            class="w-full mt-2 border-t border-border/20 pt-3 pb-1 cursor-pointer hover:opacity-80 transition-opacity"
            @click="navigateTo('/me')"
          >
            <div class="flex items-center gap-3">
              <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-full bg-accent">
                <img
                  v-if="user.avatarUrl"
                  :src="user.avatarUrl"
                  :alt="user.firstName"
                  class="h-full w-full object-cover"
                >
                <span
                  v-else
                  class="flex h-full w-full items-center justify-center text-sm font-medium text-accent-foreground"
                >
                  {{ user.firstName?.[0] }}{{ user.lastName?.[0] }}
                </span>
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-sm font-medium truncate">{{ [user.firstName, user.lastName?.[0]].filter(Boolean).join(' ') }}</span>
                <span class="text-xs text-sidebar-foreground/60">{{ level }}</span>
              </div>
            </div>
            <div class="mt-2">
              <div class="flex justify-between text-xs text-sidebar-foreground/60 mb-1">
                <span>Уровень {{ levelIndex + 1 }}/{{ maxLevel + 1 }}</span>
              </div>
              <div class="h-1.5 w-full rounded-full bg-sidebar-foreground/20">
                <div
                  class="h-full rounded-full bg-green-500 transition-all duration-500"
                  :style="{ width: `${((levelIndex) / maxLevel) * 100}%` }"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <ReviewModal
      :is-open="isModalOpen"
      @close="isModalOpen = false"
      @save="handleSaveReview"
    />
  </div>
</template>
