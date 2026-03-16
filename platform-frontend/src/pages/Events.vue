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
import EventCardSkeleton from '@/components/events/EventCardSkeleton.vue'
import EventFilters from '@/components/events/EventFilters.vue'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { useCardReveal } from '@/composables/useCardReveal'
import { useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'

const user = useUser()
const filterMode = ref<'all' | 'my'>('all')

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
function isMyEvent(event: CommunityEvent) {
  if (!user.value)
    return false
  const userId = user.value.id
  return event.members?.some(m => m.id === userId) || event.hosts?.some(h => h.id === userId)
}

const filteredFutureEvents = computed(() =>
  filterMode.value === 'my' ? futureEvents.value.filter(isMyEvent) : futureEvents.value,
)
const filteredPastEvents = computed(() =>
  filterMode.value === 'my' ? pastEvents.value.filter(isMyEvent) : pastEvents.value,
)
const allEvents = computed(() => {
  const events = viewMode.value === 'calendar' ? calendarEvents.value : [...futureEvents.value, ...pastEvents.value]
  return filterMode.value === 'my' ? events.filter(isMyEvent) : events
})

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

    <div class="flex flex-wrap items-center gap-4 mb-4 md:mb-6">
      <EventFilters class="flex-1 min-w-0" @change="loadEvents" />
      <div
        v-if="user"
        class="flex gap-1 bg-muted rounded-lg p-0.5"
      >
        <button
          class="px-3 py-1.5 text-sm rounded-md transition-colors"
          :class="filterMode === 'all' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'"
          @click="filterMode = 'all'"
        >
          Все события
        </button>
        <button
          class="px-3 py-1.5 text-sm rounded-md transition-colors"
          :class="filterMode === 'my' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'"
          @click="filterMode = 'my'"
        >
          Мои события
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col md:grid md:grid-cols-2 gap-6 md:gap-8">
      <div class="space-y-4">
        <Skeleton class="h-7 w-48 rounded-lg mb-4" />
        <EventCardSkeleton v-for="i in 2" :key="`f-${i}`" />
      </div>
      <div class="space-y-4">
        <Skeleton class="h-7 w-36 rounded-lg mb-4" />
        <EventCardSkeleton v-for="i in 2" :key="`p-${i}`" />
      </div>
    </div>

    <ErrorState v-else-if="loadError" :message="loadError" @retry="loadEvents()" />

    <template v-else-if="viewMode === 'list'">
      <div class="flex flex-col md:grid md:grid-cols-2 gap-6 md:gap-8">
        <div>
          <Typography variant="h3" as="h2" class="mb-4">
            Предстоящие события
          </Typography>
          <EmptyState
            v-if="filteredFutureEvents.length === 0"
            :icon="CalendarX"
            :title="filterMode === 'my' ? 'Нет ваших предстоящих событий' : 'Нет предстоящих событий'"
            :description="filterMode === 'my' ? undefined : 'Следите за обновлениями — новые события появляются регулярно'"
          />
          <template v-else>
            <div class="space-y-4">
              <EventCard
                v-for="event in filteredFutureEvents"
                :key="event.id"
                :event="event"
              />
            </div>
            <div v-if="filterMode === 'all' && futureEvents.length < futureTotal" class="mt-4 flex justify-center">
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
            v-if="filteredPastEvents.length === 0"
            :icon="CalendarX"
            :title="filterMode === 'my' ? 'Нет ваших архивных событий' : 'Нет архивных событий'"
          />
          <template v-else>
            <div class="space-y-4">
              <EventCard
                v-for="event in filteredPastEvents"
                :key="event.id"
                :event="event"
              />
            </div>
            <div v-if="filterMode === 'all' && pastEvents.length < pastTotal" class="mt-4 flex justify-center">
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
