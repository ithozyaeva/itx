<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import type { EventSearchFilters } from '@/services/events'
import { Typography } from 'itx-ui-kit'
import { Calendar, CalendarX, List, Loader2 } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import CalendarView from '@/components/events/CalendarView.vue'
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
const viewMode = ref<'list' | 'calendar'>('list')
const loadError = ref<string | null>(null)
const calendarEvents = ref<CommunityEvent[]>([])
const allEvents = computed(() => viewMode.value === 'calendar' ? calendarEvents.value : [...futureEvents.value, ...pastEvents.value])

async function loadEvents(filters?: EventSearchFilters) {
  if (filters)
    currentFilters.value = filters
  isLoading.value = true
  loadError.value = null
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
    loadError.value = (await handleError(error)).message
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

const isLoadingCalendar = ref(false)

async function loadCalendarEvents() {
  isLoadingCalendar.value = true
  try {
    const [past, future] = await Promise.all([
      eventsService.searchOld(200, 0, currentFilters.value),
      eventsService.searchNext(200, 0, currentFilters.value),
    ])
    calendarEvents.value = [...future.items, ...past.items]
  }
  catch (error) {
    calendarEvents.value = [...futureEvents.value, ...pastEvents.value]
    handleError(error)
  }
  finally {
    isLoadingCalendar.value = false
  }
}

watch(viewMode, (mode) => {
  if (mode === 'calendar')
    loadCalendarEvents()
})

onMounted(() => loadEvents())
</script>

<template>
  <div ref="containerRef" class="container mx-auto px-4 py-6 md:py-8">
    <div class="flex items-center justify-between mb-4 md:mb-6">
      <Typography variant="h2" as="h1">
        События сообщества
      </Typography>
      <div class="flex gap-1 bg-muted rounded-lg p-0.5">
        <button
          class="p-2 rounded-md transition-colors"
          :class="viewMode === 'list' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'"
          aria-label="Список"
          @click="viewMode = 'list'"
        >
          <List class="h-4 w-4" />
        </button>
        <button
          class="p-2 rounded-md transition-colors"
          :class="viewMode === 'calendar' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'"
          aria-label="Календарь"
          @click="viewMode = 'calendar'"
        >
          <Calendar class="h-4 w-4" />
        </button>
      </div>
    </div>

    <EventFilters class="mb-4 md:mb-6" @change="loadEvents" />

    <div v-if="isLoading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <ErrorState v-else-if="loadError" :message="loadError" @retry="loadEvents()" />

    <template v-else-if="viewMode === 'list'">
      <div class="flex flex-col md:grid md:grid-cols-2 gap-6 md:gap-8">
        <div>
          <Typography variant="h3" as="h2" class="mb-4">
            Предстоящие события
          </Typography>
          <EmptyState
            v-if="futureEvents.length === 0"
            :icon="CalendarX"
            title="Нет предстоящих событий"
            description="Следите за обновлениями — новые события появляются регулярно"
          />
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
        <div>
          <Typography variant="h3" as="h2" class="mb-4">
            Архив событий
          </Typography>
          <EmptyState
            v-if="pastEvents.length === 0"
            :icon="CalendarX"
            title="Нет архивных событий"
          />
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
      </div>
    </template>

    <div v-else-if="viewMode === 'calendar'" class="rounded-2xl border bg-card border-border p-4">
      <div v-if="isLoadingCalendar" class="flex justify-center py-12">
        <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
      <CalendarView v-else :events="allEvents" />
    </div>
  </div>
</template>
