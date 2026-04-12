<script setup lang="ts">
import type { NotificationSettings } from '@/models/profile'
import { Bell, Loader2 } from 'lucide-vue-next'
import { onMounted, reactive, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Typography } from '@/components/ui/typography'
import { handleError } from '@/services/errorService'
import { notificationSettingsService } from '@/services/notificationSettings'

const isLoading = ref(true)
const isSaving = ref(false)
const hasChanges = ref(false)

const settings = reactive({
  muteAll: false,
  newEvents: true,
  remindWeek: true,
  remindDay: true,
  remindHour: true,
  eventStart: true,
  eventUpdates: true,
  eventCancelled: true,
})

let originalSettings: Partial<NotificationSettings> = {}

const toggleItems = [
  { key: 'newEvents' as const, label: 'Новые события', description: 'Уведомления о новых событиях' },
  { key: 'remindWeek' as const, label: 'Напоминание за неделю', description: 'За ~7 дней до события' },
  { key: 'remindDay' as const, label: 'Напоминание за день', description: 'За ~1 день до события' },
  { key: 'remindHour' as const, label: 'Напоминание за час', description: 'За ~1 час до события' },
  { key: 'eventStart' as const, label: 'Начало события', description: 'Когда событие начинается' },
  { key: 'eventUpdates' as const, label: 'Изменения событий', description: 'Когда событие обновляется' },
  { key: 'eventCancelled' as const, label: 'Отмена событий', description: 'Когда событие отменяется' },
]

function toggleMuteAll() {
  settings.muteAll = !settings.muteAll
  checkChanges()
}

onMounted(async () => {
  try {
    const data = await notificationSettingsService.get()
    if (data) {
      Object.assign(settings, {
        muteAll: data.muteAll,
        newEvents: data.newEvents,
        remindWeek: data.remindWeek,
        remindDay: data.remindDay,
        remindHour: data.remindHour,
        eventStart: data.eventStart,
        eventUpdates: data.eventUpdates,
        eventCancelled: data.eventCancelled,
      })
      originalSettings = { ...settings }
    }
  }
  catch (err) {
    handleError(err)
  }
  finally {
    isLoading.value = false
  }
})

function toggleSetting(key: keyof typeof settings) {
  settings[key] = !settings[key]
  checkChanges()
}

function checkChanges() {
  hasChanges.value = settings.muteAll !== originalSettings.muteAll
    || toggleItems.some(item => settings[item.key] !== originalSettings[item.key])
}

async function handleSubmit() {
  isSaving.value = true
  try {
    const result = await notificationSettingsService.update(settings)
    if (result) {
      originalSettings = { ...settings }
      hasChanges.value = false
    }
  }
  catch (err) {
    handleError(err)
  }
  finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="p-6 md:p-8 bg-card backdrop-blur-lg border border-border shadow-lg rounded-sm terminal-card">
    <div class="flex relative flex-col items-center space-y-4">
      <div class="flex items-center gap-2">
        <Bell class="h-5 w-5 text-muted-foreground" />
        <Typography
          variant="h3"
          as="h5"
        >
          Уведомления в Telegram
        </Typography>
      </div>

      <div
        v-if="isLoading"
        class="flex justify-center py-4"
      >
        <Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
      </div>

      <div
        v-else
        class="w-full space-y-3"
      >
        <div
          class="flex items-center justify-between py-2 px-3 rounded-xl bg-destructive/5 border border-destructive/20 cursor-pointer"
          @click="toggleMuteAll"
        >
          <div class="flex flex-col">
            <span class="text-sm font-medium">Отключить все уведомления</span>
            <span class="text-xs text-muted-foreground">Включая уведомления от бота в Telegram</span>
          </div>
          <button
            type="button"
            role="switch"
            :aria-checked="settings.muteAll"
            class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            :class="settings.muteAll ? 'bg-destructive' : 'bg-input'"
            @click.stop="toggleMuteAll"
          >
            <span
              class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-background shadow-lg ring-0 transition duration-200 ease-in-out"
              :class="settings.muteAll ? 'translate-x-5' : 'translate-x-0'"
            />
          </button>
        </div>

        <div
          v-for="item in toggleItems"
          :key="item.key"
          class="flex items-center justify-between py-2 px-3 rounded-xl transition-colors"
          :class="settings.muteAll ? 'opacity-40 pointer-events-none' : 'hover:bg-accent/50 cursor-pointer'"
          @click="toggleSetting(item.key)"
        >
          <div class="flex flex-col">
            <span class="text-sm font-medium">{{ item.label }}</span>
            <span class="text-xs text-muted-foreground">{{ item.description }}</span>
          </div>
          <button
            type="button"
            role="switch"
            :aria-checked="settings[item.key]"
            class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            :class="settings[item.key] ? 'bg-accent' : 'bg-input'"
            @click.stop="toggleSetting(item.key)"
          >
            <span
              class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-background shadow-lg ring-0 transition duration-200 ease-in-out"
              :class="settings[item.key] ? 'translate-x-5' : 'translate-x-0'"
            />
          </button>
        </div>

        <Button
          v-if="hasChanges"
          class="w-full mt-2 cursor-pointer transition duration-300 gap-2"
          :disabled="isSaving"
          @click="handleSubmit"
        >
          <Loader2
            v-if="isSaving"
            class="h-4 w-4 animate-spin"
          />
          Сохранить
        </Button>
      </div>
    </div>
  </div>
</template>
