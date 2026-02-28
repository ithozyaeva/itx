<script setup lang="ts">
import type { Notification } from '@/services/notifications'
import { Badge } from '@/components/ui/badge'
import { handleError } from '@/services/errorService'
import { notificationService } from '@/services/notifications'
import { Bell, CheckCheck } from 'lucide-vue-next'
import { onMounted, onUnmounted, ref } from 'vue'

const notifications = ref<Notification[]>([])
const unreadCount = ref(0)
const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
let pollInterval: ReturnType<typeof setInterval> | null = null

async function fetchUnreadCount() {
  try {
    const result = await notificationService.getUnreadCount()
    unreadCount.value = result.count
  }
  catch {
    // Silently fail for polling
  }
}

async function fetchNotifications() {
  try {
    notifications.value = await notificationService.getAll()
  }
  catch (error) {
    handleError(error)
  }
}

async function toggleDropdown() {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    await fetchNotifications()
  }
}

async function markAsRead(notification: Notification) {
  if (notification.read)
    return
  try {
    await notificationService.markAsRead(notification.id)
    notification.read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }
  catch (error) {
    handleError(error)
  }
}

async function markAllAsRead() {
  try {
    await notificationService.markAllAsRead()
    notifications.value.forEach(n => n.read = true)
    unreadCount.value = 0
  }
  catch (error) {
    handleError(error)
  }
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1)
    return 'только что'
  if (minutes < 60)
    return `${minutes} мин. назад`
  if (hours < 24)
    return `${hours} ч. назад`
  if (days < 7)
    return `${days} дн. назад`
  return date.toLocaleDateString('ru-RU')
}

onMounted(() => {
  fetchUnreadCount()
  pollInterval = setInterval(fetchUnreadCount, 30000)
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  if (pollInterval)
    clearInterval(pollInterval)
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="dropdownRef" class="relative">
    <button class="relative p-1 rounded hover:bg-secondary cursor-pointer" @click="toggleDropdown">
      <Bell class="h-5 w-5" />
      <Badge
        v-if="unreadCount > 0"
        class="absolute -top-1 -right-1 h-4 min-w-4 px-1 text-[10px] flex items-center justify-center"
      >
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </Badge>
    </button>

    <div
      v-if="isOpen"
      class="absolute right-0 top-full mt-2 w-80 bg-popover border rounded-lg shadow-lg z-50 overflow-hidden"
    >
      <div class="flex items-center justify-between px-4 py-3 border-b">
        <span class="font-semibold text-sm">Уведомления</span>
        <button
          v-if="unreadCount > 0"
          class="text-xs text-muted-foreground hover:text-foreground cursor-pointer flex items-center gap-1"
          @click="markAllAsRead"
        >
          <CheckCheck class="h-3 w-3" />
          Прочитать все
        </button>
      </div>

      <div class="max-h-80 overflow-y-auto">
        <div v-if="notifications.length === 0" class="px-4 py-8 text-center text-sm text-muted-foreground">
          Нет уведомлений
        </div>

        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="px-4 py-3 border-b last:border-b-0 cursor-pointer hover:bg-muted/50 transition-colors"
          :class="{ 'bg-muted/30': !notification.read }"
          @click="markAsRead(notification)"
        >
          <div class="flex items-start gap-2">
            <div v-if="!notification.read" class="w-2 h-2 rounded-full bg-primary mt-1.5 shrink-0" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium truncate">
                {{ notification.title }}
              </p>
              <p v-if="notification.body" class="text-xs text-muted-foreground mt-0.5 line-clamp-2">
                {{ notification.body }}
              </p>
              <p class="text-xs text-muted-foreground mt-1">
                {{ formatDate(notification.createdAt) }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
