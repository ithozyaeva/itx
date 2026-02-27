<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import type { EventSearchFilters } from '@/services/events'
import EventCard from '@/components/events/EventCard.vue'
import EventFilters from '@/components/events/EventFilters.vue'
import { useCardReveal } from '@/composables/useCardReveal'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'
import { Typography } from 'itx-ui-kit'
import { onMounted, ref } from 'vue'

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const pastEvents = ref<CommunityEvent[]>([])
const futureEvents = ref<CommunityEvent[]>([])
const currentFilters = ref<EventSearchFilters>({})

async function loadEvents(filters?: EventSearchFilters) {
  if (filters)
    currentFilters.value = filters
  try {
    pastEvents.value = (await eventsService.searchOld(30, 0, currentFilters.value)).items
    futureEvents.value = (await eventsService.searchNext(30, 0, currentFilters.value)).items
  }
  catch (error) {
    handleError(error)
  }
}

onMounted(() => loadEvents())
</script>

<template>
  <div ref="containerRef" class="container mx-auto px-4 py-8">
    <Typography variant="h2" as="h1" class="mb-6">
      События сообщества
    </Typography>

    <EventFilters class="mb-6" @change="loadEvents" />

    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <div>
        <Typography variant="h3" as="h2" class="mb-4">
          Архив событий
        </Typography>
        <div v-if="pastEvents.length === 0" class="text-muted-foreground">
          Нет архивных событий
        </div>
        <div v-else class="space-y-4">
          <EventCard
            v-for="event in pastEvents"
            :key="event.id"
            :event="event"
          />
        </div>
      </div>
      <div>
        <Typography variant="h3" as="h2" class="mb-4">
          Предстоящие события
        </Typography>
        <div v-if="futureEvents.length === 0" class="text-muted-foreground">
          Нет предстоящих событий
        </div>
        <div v-else class="space-y-4">
          <EventCard
            v-for="event in futureEvents"
            :key="event.id"
            :event="event"
          />
        </div>
      </div>
    </div>
  </div>
</template>
