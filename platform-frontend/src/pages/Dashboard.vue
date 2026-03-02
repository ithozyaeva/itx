<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import { Bot, ExternalLink, Loader2, Lock, Play } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EventsSection from '@/components/dashboard/EventsSection.vue'
import GreetingCard from '@/components/dashboard/GreetingCard.vue'
import NearestEvent from '@/components/dashboard/NearestEvent.vue'
import ReferralsSidebar from '@/components/dashboard/ReferralsSidebar.vue'
import { useUserLevel } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'

const nearestEvent = ref<CommunityEvent | null>(null)
const isLoading = ref(false)
const { levelIndex } = useUserLevel()
const showVideoEmbed = ref(false)

async function loadNearestEvent() {
  isLoading.value = true
  try {
    const result = await eventsService.searchNext(1, 0)
    nearestEvent.value = result.items[0] ?? null
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

onMounted(() => loadNearestEvent())
</script>

<template>
  <div class="flex gap-6 px-4 py-6 md:py-8">
    <!-- Main content -->
    <div class="flex-1 min-w-0">
      <GreetingCard />

      <div
        v-if="isLoading"
        class="flex justify-center py-12"
      >
        <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
      </div>

      <NearestEvent
        v-else-if="nearestEvent"
        :event="nearestEvent"
        class="mt-6"
      />

      <EventsSection :skip-first="!!nearestEvent" />

      <!-- Latest YouTube video -->
      <div class="mt-6 rounded-3xl border bg-card overflow-hidden">
        <div
          class="relative w-full overflow-hidden"
          style="padding-bottom: 56.25%;"
        >
          <iframe
            v-if="showVideoEmbed"
            class="absolute inset-0 w-full h-full"
            src="https://www.youtube.com/embed/H6cWhHG_KBQ?autoplay=1"
            title="Деплой для глупеньких [ IT-X: Mornings ]"
            frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen
          />
          <button
            v-else
            class="absolute inset-0 w-full h-full cursor-pointer"
            @click="showVideoEmbed = true"
          >
            <img
              src="https://img.youtube.com/vi/H6cWhHG_KBQ/mqdefault.jpg"
              alt="Деплой для глупеньких [ IT-X: Mornings ]"
              class="w-full h-full object-cover"
            >
            <div class="absolute inset-0 bg-black/20 hover:bg-black/30 transition-colors flex items-center justify-center">
              <div class="w-14 h-14 rounded-full bg-white/90 flex items-center justify-center">
                <Play class="h-6 w-6 text-foreground ml-0.5" fill="currentColor" />
              </div>
            </div>
          </button>
        </div>
        <div class="p-5 flex items-center justify-between">
          <div>
            <p class="font-medium text-sm">
              Деплой для глупеньких [ IT-X: Mornings ]
            </p>
            <p class="text-xs text-muted-foreground mt-0.5">
              Последнее видео сообщества
            </p>
          </div>
          <RouterLink
            to="/content"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors whitespace-nowrap"
          >
            Все видео →
          </RouterLink>
        </div>
      </div>

      <!-- Auto-apply bot -->
      <div class="mt-6 rounded-3xl border bg-card p-5 flex items-center gap-4">
        <div class="flex items-center justify-center w-11 h-11 rounded-xl bg-primary/10 text-primary shrink-0">
          <Bot :size="22" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="font-medium text-sm">
            Бот для автооткликов на hh.ru
          </p>
          <p class="text-xs text-muted-foreground mt-0.5">
            Доступен всем участникам сообщества
          </p>
        </div>
        <a
          href="https://t.me/roaster_resume_bot"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-1.5 rounded-xl bg-primary text-primary-foreground px-4 py-2 text-xs font-medium hover:bg-primary/90 transition-colors shrink-0"
        >
          Открыть
          <ExternalLink :size="12" />
        </a>
      </div>

      <!-- Locked content teaser -->
      <div
        v-if="levelIndex < 2"
        class="mt-6 rounded-3xl border border-dashed border-muted-foreground/30 bg-muted/30 p-6 text-center"
      >
        <Lock class="h-8 w-8 mx-auto text-muted-foreground/50 mb-2" />
        <p class="text-muted-foreground font-medium">
          Больше контента доступно на уровне «Хозяин»
        </p>
        <p class="text-sm text-muted-foreground/60 mt-1">
          Повышайте уровень, чтобы разблокировать эксклюзивные материалы
        </p>
      </div>
    </div>

    <!-- Right sidebar (desktop only) -->
    <div class="hidden lg:block w-80 shrink-0">
      <ReferralsSidebar />
    </div>
  </div>
</template>
