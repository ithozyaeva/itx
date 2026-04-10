<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { getOccurrencesInMonth } from '@/composables/useEventOccurrence'
import EventCard from './EventCard.vue'

const props = defineProps<{
  events: CommunityEvent[]
}>()

const today = new Date()
const currentMonth = ref(today.getMonth())
const currentYear = ref(today.getFullYear())

const monthNames = [
  'Январь',
  'Февраль',
  'Март',
  'Апрель',
  'Май',
  'Июнь',
  'Июль',
  'Август',
  'Сентябрь',
  'Октябрь',
  'Ноябрь',
  'Декабрь',
]

const weekDays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const selectedDate = ref<string | null>(null)

interface CalendarDay {
  date: number
  month: number
  year: number
  isCurrentMonth: boolean
  isToday: boolean
  dateKey: string
}

const calendarDays = computed<CalendarDay[]>(() => {
  const year = currentYear.value
  const month = currentMonth.value

  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)

  // Monday = 0, Sunday = 6
  let startDow = firstDay.getDay() - 1
  if (startDow < 0)
    startDow = 6

  const days: CalendarDay[] = []

  // Previous month days
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = startDow - 1; i >= 0; i--) {
    const d = prevMonthLastDay - i
    const m = month === 0 ? 11 : month - 1
    const y = month === 0 ? year - 1 : year
    days.push({
      date: d,
      month: m,
      year: y,
      isCurrentMonth: false,
      isToday: false,
      dateKey: formatDateKey(y, m, d),
    })
  }

  // Current month days
  for (let d = 1; d <= lastDay.getDate(); d++) {
    const isToday = d === today.getDate() && month === today.getMonth() && year === today.getFullYear()
    days.push({
      date: d,
      month,
      year,
      isCurrentMonth: true,
      isToday,
      dateKey: formatDateKey(year, month, d),
    })
  }

  // Next month days to fill 6 rows
  const remaining = 42 - days.length
  for (let d = 1; d <= remaining; d++) {
    const m = month === 11 ? 0 : month + 1
    const y = month === 11 ? year + 1 : year
    days.push({
      date: d,
      month: m,
      year: y,
      isCurrentMonth: false,
      isToday: false,
      dateKey: formatDateKey(y, m, d),
    })
  }

  return days
})

const eventsByDate = computed(() => {
  const map = new Map<string, CommunityEvent[]>()
  for (const event of props.events) {
    if (event.isRepeating && event.repeatPeriod) {
      const occurrences = getOccurrencesInMonth(event, currentYear.value, currentMonth.value)
      for (const date of occurrences) {
        const key = formatDateKey(date.getFullYear(), date.getMonth(), date.getDate())
        if (!map.has(key))
          map.set(key, [])
        map.get(key)!.push({ ...event, date: date.toISOString() })
      }
    }
    else {
      const date = new Date(event.date)
      const key = formatDateKey(date.getFullYear(), date.getMonth(), date.getDate())
      if (!map.has(key))
        map.set(key, [])
      map.get(key)!.push(event)
    }
  }
  return map
})

const selectedDayEvents = computed(() => {
  if (!selectedDate.value)
    return []
  return eventsByDate.value.get(selectedDate.value) ?? []
})

function formatDateKey(year: number, month: number, date: number): string {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(date).padStart(2, '0')}`
}

function prevMonth() {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  }
  else {
    currentMonth.value--
  }
}

function nextMonth() {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  }
  else {
    currentMonth.value++
  }
}

function selectDay(day: CalendarDay) {
  const key = day.dateKey
  selectedDate.value = selectedDate.value === key ? null : key
}

function getEventCount(dateKey: string): number {
  return eventsByDate.value.get(dateKey)?.length ?? 0
}

function goToToday() {
  currentMonth.value = today.getMonth()
  currentYear.value = today.getFullYear()
}
</script>

<template>
  <div>
    <!-- Calendar header -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <h3 class="text-lg font-semibold">
          {{ monthNames[currentMonth] }} {{ currentYear }}
        </h3>
        <button
          class="text-xs text-muted-foreground hover:text-foreground transition-colors px-2 py-0.5 rounded-md hover:bg-muted"
          @click="goToToday"
        >
          Сегодня
        </button>
      </div>
      <div class="flex gap-1">
        <button
          class="p-1.5 rounded-lg hover:bg-muted transition-colors"
          aria-label="Предыдущий месяц"
          @click="prevMonth"
        >
          <ChevronLeft class="h-4 w-4" />
        </button>
        <button
          class="p-1.5 rounded-lg hover:bg-muted transition-colors"
          aria-label="Следующий месяц"
          @click="nextMonth"
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Weekday headers -->
    <div class="grid grid-cols-7 gap-px mb-1">
      <div
        v-for="day in weekDays"
        :key="day"
        class="text-center text-xs font-medium text-muted-foreground py-2"
      >
        {{ day }}
      </div>
    </div>

    <!-- Calendar grid -->
    <div class="grid grid-cols-7 gap-px bg-border rounded-xl overflow-hidden">
      <button
        v-for="(day, index) in calendarDays"
        :key="index"
        class="relative bg-card p-2 min-h-[3.5rem] text-center transition-colors hover:bg-muted/50"
        :class="{
          'opacity-40': !day.isCurrentMonth,
          'ring-2 ring-primary ring-inset': selectedDate === day.dateKey,
        }"
        @click="selectDay(day)"
      >
        <span
          class="text-sm inline-flex items-center justify-center w-6 h-6"
          :class="{
            'rounded-full bg-primary text-primary-foreground text-xs font-bold': day.isToday,
          }"
        >
          {{ day.date }}
        </span>
        <!-- Event dots -->
        <div
          v-if="getEventCount(day.dateKey) > 0"
          class="flex justify-center gap-0.5 mt-1"
        >
          <div
            v-for="n in Math.min(getEventCount(day.dateKey), 3)"
            :key="n"
            class="h-1.5 w-1.5 rounded-full bg-accent"
          />
          <span
            v-if="getEventCount(day.dateKey) > 3"
            class="text-[10px] text-muted-foreground leading-none"
          >
            +{{ getEventCount(day.dateKey) - 3 }}
          </span>
        </div>
      </button>
    </div>

    <!-- Selected day events -->
    <div
      v-if="selectedDate && selectedDayEvents.length > 0"
      class="mt-4 space-y-3"
    >
      <h4 class="text-sm font-medium text-muted-foreground">
        События за {{ selectedDate }}
      </h4>
      <EventCard
        v-for="event in selectedDayEvents"
        :key="event.id"
        :event="event"
      />
    </div>
    <div
      v-else-if="selectedDate && selectedDayEvents.length === 0"
      class="mt-4 text-center py-6 text-sm text-muted-foreground"
    >
      Нет событий в этот день
    </div>
  </div>
</template>
