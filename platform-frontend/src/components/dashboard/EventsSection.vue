<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import { CalendarX, Loader2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EventCard from '@/components/events/EventCard.vue'
import { Button } from '@/components/ui/button'
import { useCardReveal } from '@/composables/useCardReveal'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'

const PAGE_SIZE = 5

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const events = ref<CommunityEvent[]>([])
const total = ref(0)
const isLoading = ref(false)
const isLoadingMore = ref(false)
const activeTab = ref('all')

const tabs = [
  { key: 'all', label: 'Все' },
  { key: 'review', label: 'Разборы' },
  { key: 'stream', label: 'Стримы' },
  { key: 'networking', label: 'Нетворкинг' },
]

async function loadEvents() {
  isLoading.value = true
  try {
    const result = await eventsService.searchNext(PAGE_SIZE, 0)
    events.value = result.items
    total.value = result.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

async function loadMore() {
  isLoadingMore.value = true
  try {
    const result = await eventsService.searchNext(PAGE_SIZE, events.value.length)
    events.value.push(...result.items)
    total.value = result.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoadingMore.value = false
  }
}

onMounted(() => loadEvents())
</script>

<template>
  <div
    ref="containerRef"
    class="mt-6"
  >
    <h2 class="text-lg font-semibold mb-4">
      События сообщества
    </h2>
    <div class="flex gap-2 mb-4 overflow-x-auto">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="rounded-full px-4 py-1.5 text-sm font-medium transition-colors whitespace-nowrap"
        :class="activeTab === tab.key
          ? 'bg-primary text-primary-foreground'
          : 'bg-muted text-muted-foreground hover:bg-muted/80'"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div
      v-else-if="events.length === 0"
      class="flex flex-col items-center gap-2 py-8 text-muted-foreground"
    >
      <CalendarX class="h-10 w-10" />
      <p>Нет предстоящих событий</p>
    </div>

    <div
      v-else
      class="space-y-4"
    >
      <EventCard
        v-for="event in events"
        :key="event.id"
        :event="event"
      />
      <div
        v-if="events.length < total"
        class="flex justify-center"
      >
        <Button
          variant="outline"
          :disabled="isLoadingMore"
          @click="loadMore"
        >
          <Loader2
            v-if="isLoadingMore"
            class="mr-2 h-4 w-4 animate-spin"
          />
          Показать ещё
        </Button>
      </div>
    </div>
  </div>
</template>
