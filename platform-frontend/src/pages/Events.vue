<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import type { EventSearchFilters } from '@/services/events'
import { Typography } from 'itx-ui-kit'
import { CalendarX, Loader2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EventCard from '@/components/events/EventCard.vue'
import EventFilters from '@/components/events/EventFilters.vue'
import { Button } from '@/components/ui/button'
import { useCardReveal } from '@/composables/useCardReveal'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'

const PAGE_SIZE = 10

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const pastEvents = ref<CommunityEvent[]>([])
const futureEvents = ref<CommunityEvent[]>([])
const pastTotal = ref(0)
const futureTotal = ref(0)
const isLoading = ref(false)
const isLoadingMorePast = ref(false)
const isLoadingMoreFuture = ref(false)
const currentFilters = ref<EventSearchFilters>({})

async function loadEvents(filters?: EventSearchFilters) {
  if (filters)
    currentFilters.value = filters
  isLoading.value = true
  try {
    const [pastResult, futureResult] = await Promise.all([
      eventsService.searchOld(PAGE_SIZE, 0, currentFilters.value),
      eventsService.searchNext(PAGE_SIZE, 0, currentFilters.value),
    ])
    pastEvents.value = pastResult.items
    pastTotal.value = pastResult.total
    futureEvents.value = futureResult.items
    futureTotal.value = futureResult.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

async function loadMorePast() {
  isLoadingMorePast.value = true
  try {
    const result = await eventsService.searchOld(PAGE_SIZE, pastEvents.value.length, currentFilters.value)
    pastEvents.value.push(...result.items)
    pastTotal.value = result.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoadingMorePast.value = false
  }
}

async function loadMoreFuture() {
  isLoadingMoreFuture.value = true
  try {
    const result = await eventsService.searchNext(PAGE_SIZE, futureEvents.value.length, currentFilters.value)
    futureEvents.value.push(...result.items)
    futureTotal.value = result.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoadingMoreFuture.value = false
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

    <div v-if="isLoading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <div>
        <Typography variant="h3" as="h2" class="mb-4">
          Архив событий
        </Typography>
        <div v-if="pastEvents.length === 0" class="flex flex-col items-center gap-2 py-8 text-muted-foreground">
          <CalendarX class="h-10 w-10" />
          <p>Нет архивных событий</p>
        </div>
        <template v-else>
          <div class="space-y-4">
            <EventCard
              v-for="event in pastEvents"
              :key="event.id"
              :event="event"
            />
          </div>
          <div v-if="pastEvents.length < pastTotal" class="mt-4 flex justify-center">
            <Button
              variant="outline"
              :disabled="isLoadingMorePast"
              @click="loadMorePast"
            >
              <Loader2 v-if="isLoadingMorePast" class="mr-2 h-4 w-4 animate-spin" />
              Показать ещё
            </Button>
          </div>
        </template>
      </div>
      <div>
        <Typography variant="h3" as="h2" class="mb-4">
          Предстоящие события
        </Typography>
        <div v-if="futureEvents.length === 0" class="flex flex-col items-center gap-2 py-8 text-muted-foreground">
          <CalendarX class="h-10 w-10" />
          <p>Нет предстоящих событий</p>
        </div>
        <template v-else>
          <div class="space-y-4">
            <EventCard
              v-for="event in futureEvents"
              :key="event.id"
              :event="event"
            />
          </div>
          <div v-if="futureEvents.length < futureTotal" class="mt-4 flex justify-center">
            <Button
              variant="outline"
              :disabled="isLoadingMoreFuture"
              @click="loadMoreFuture"
            >
              <Loader2 v-if="isLoadingMoreFuture" class="mr-2 h-4 w-4 animate-spin" />
              Показать ещё
            </Button>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
