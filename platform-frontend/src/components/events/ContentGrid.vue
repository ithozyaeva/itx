<script setup lang="ts">
import { ExternalLink, Play } from 'lucide-vue-next'
import { ref } from 'vue'
import { openLink } from '@/composables/useTelegramWebApp'

interface Video {
  id: string
  title: string
  date: string
}

const videos: Video[] = [
  { id: 'H6cWhHG_KBQ', title: 'Деплой для глупеньких [ IT-X: Mornings ]', date: '2026-02-10' },
  { id: '5x4BKfrQhrY', title: 'Что с работой в 2026: QA с HR сообщества', date: '2026-02-08' },
  { id: '05EXlY1q-Kc', title: 'Postgres для настоящих слонов [ IT-X: Mornings ]', date: '2026-01-27' },
  { id: 'Aiy6rwQNrds', title: 'Делаем SSG-утилиту на Rust [ IT-X: Mornings ]', date: '2025-12-18' },
  { id: 'NZhTyuJWVJE', title: 'Фриланс: опыт выживания [ IT-X: Mornings ]', date: '2025-12-18' },
  { id: '-UA56Roynpg', title: 'Docker: База. Часть 1 [ IT-X: Mornings ]', date: '2025-12-16' },
  { id: '_J090s_jeOk', title: 'Валентин Ким. Как найти работу фронтом в 2025', date: '2025-07-24' },
  { id: 'SO9Xn_bF1zU', title: 'Василий Кузенков. База по ИИ', date: '2025-07-23' },
  { id: '7tYCbNyIun4', title: 'Владимир Балун. Системный дизайн', date: '2025-07-10' },
  { id: 'w-TETYEhzxs', title: 'Альтернативный способ заработка в IT. Как стать ментором?', date: '2025-02-22' },
]

const activeVideo = ref<string | null>(null)

const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})

function formatDate(dateStr: string) {
  return dateFormatter.format(new Date(dateStr))
}
</script>

<template>
  <div>
    <div class="flex items-center justify-end mb-4">
      <a
        href="https://www.youtube.com/@joindev"
        target="_blank"
        rel="noopener noreferrer"
        class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
        @click.prevent="openLink('https://www.youtube.com/@joindev')"
      >
        YouTube канал
        <ExternalLink :size="14" />
      </a>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="video in videos"
        :key="video.id"
        class="rounded-sm border bg-card overflow-hidden group"
      >
        <div
          class="relative w-full overflow-hidden"
          style="padding-bottom: 56.25%;"
        >
          <iframe
            v-if="activeVideo === video.id"
            class="absolute inset-0 w-full h-full"
            :src="`https://www.youtube.com/embed/${video.id}?autoplay=1`"
            :title="video.title"
            frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen
          />
          <button
            v-else
            class="absolute inset-0 w-full h-full cursor-pointer"
            @click="activeVideo = video.id"
          >
            <img
              :src="`https://img.youtube.com/vi/${video.id}/mqdefault.jpg`"
              :alt="video.title"
              class="w-full h-full object-cover"
            >
            <div class="absolute inset-0 bg-black/20 group-hover:bg-black/30 transition-colors flex items-center justify-center">
              <div class="w-12 h-12 rounded-full bg-background/90 flex items-center justify-center">
                <Play class="h-5 w-5 text-foreground ml-0.5" fill="currentColor" />
              </div>
            </div>
          </button>
        </div>
        <div class="p-4">
          <p class="font-medium text-sm line-clamp-2">
            {{ video.title }}
          </p>
          <p class="text-xs text-muted-foreground mt-1">
            {{ formatDate(video.date) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
