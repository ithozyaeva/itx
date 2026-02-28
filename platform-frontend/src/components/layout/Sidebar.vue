<script setup lang="ts">
import { CloseIcon, Typography } from 'itx-ui-kit'
import { Shield } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { useSidebar } from '@/composables/useSidebar'
import { isUserAdmin, isUserSubscribed, useUser, useUserLevel } from '@/composables/useUser'
import { reviewService } from '@/services/reviews'
import ReviewModal from '../ReviewModal.vue'

const { sidebarItems, isOpen, toggleSidebar } = useSidebar()
const route = useRoute()
const router = useRouter()
const user = useUser()
const { level, levelIndex, maxLevel } = useUserLevel()

function isActive(path: string) {
  if (path === '/')
    return route.path === '/'
  return route.path.startsWith(path)
}

function navigateTo(path: string) {
  router.push(path)
  isOpen.value = false
}

const isModalOpen = ref(false)
const isSubscribed = isUserSubscribed()
const isAdmin = isUserAdmin()

async function handleSaveReview(text: string) {
  await reviewService
    .createReview(text)
    .finally(() => {
      isModalOpen.value = false
    })
}
</script>

<template>
  <div>
    <div
      class="fixed md:static h-screen border-r border-border bg-primary text-primary-foreground transition-all duration-300 w-full md:w-56 z-40"
      :class="[
        isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0',
      ]"
    >
      <div class="flex flex-col justify-between items-center pb-2 h-full">
        <div class="w-full">
          <div class="flex items-center justify-between p-4 md:hidden">
            <Typography
              variant="h4"
              as="h2"
            >
              Меню
            </Typography>
            <button
              class="p-2"
              @click="toggleSidebar"
            >
              <CloseIcon class="w-6 h-6" />
            </button>
          </div>
          <div class="flex-1 overflow-y-auto py-4">
            <ul class="space-y-1">
              <li
                v-for="item in sidebarItems"
                :key="item.path"
              >
                <Button
                  variant="ghost"
                  class="w-full justify-start py-2 cursor-pointer text-primary-foreground hover:bg-accent hover:text-accent-foreground"
                  :class="[
                    isActive(item.path) ? 'bg-accent text-accent-foreground' : '',
                  ]"
                  @click="navigateTo(item.path)"
                >
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
        <div class="flex flex-col items-center gap-2 w-full px-3">
          <a
            v-if="isAdmin"
            href="/admin"
            class="flex items-center gap-2 cursor-pointer border border-border/20 rounded-full px-4 py-1 text-primary-foreground hover:bg-accent hover:text-accent-foreground transition duration-300 active:scale-95"
          >
            <Shield class="h-4 w-4" />
            <span>Админ-панель</span>
          </a>
          <button
            v-if="isSubscribed"
            class="cursor-pointer border border-border/20 rounded-full px-4 py-1 text-primary-foreground hover:bg-accent hover:text-accent-foreground transition duration-300 active:scale-95"
            @click="isModalOpen = true"
          >
            Добавить отзыв
          </button>

          <!-- User profile card -->
          <div
            v-if="user"
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
                <span class="text-sm font-medium truncate">{{ user.firstName }} {{ user.lastName?.[0] ?? '' }}</span>
                <span class="text-xs text-primary-foreground/60">{{ level }}</span>
              </div>
            </div>
            <div class="mt-2">
              <div class="flex justify-between text-xs text-primary-foreground/60 mb-1">
                <span>Уровень {{ levelIndex + 1 }}/{{ maxLevel + 1 }}</span>
              </div>
              <div class="h-1.5 w-full rounded-full bg-primary-foreground/20">
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
