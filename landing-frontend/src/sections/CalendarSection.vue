<script setup lang="ts">
import type { CommunityEvent } from '@/services/events'
import { Button, CalendarIcon, Label, Popover, Tag, Typography } from 'itx-ui-kit'
import { computed, onMounted, ref } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import Card from '@/components/ui/Card.vue'
import TgImage from '@/components/ui/TgImage.vue'
import { useGoogleCalendar } from '@/composables/useGoogleCalendar.ts'
import { eventService } from '@/services/events'

const neededEvents = ref<'new' | 'old'>('new')
const events = ref<Record<'new' | 'old', CommunityEvent[]>>({
  new: [],
  old: [],
})

function formatEventDate(dateString: string): string {
  const utcDate = new Date(dateString)

  // Конвертируем в МСК (UTC+3)
  const formatter = new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
    timeZone: 'Europe/Moscow',
  })

  return formatter.format(utcDate)
}

async function loadEvents() {
  events.value.old = await eventService.getOld()
  events.value.new = await eventService.getNext()

  if (events.value.new.length === 0) {
    neededEvents.value = 'old'
  }
}

const visibleCount = ref(4)

function showMore() {
  visibleCount.value += 4
}

const visibleEvents = computed(() => {
  return events.value[neededEvents.value].slice(0, visibleCount.value)
})

const hasFutureEvents = computed(() => events.value.new.length > 0)
const isFuture = computed(() => neededEvents.value === 'new')

const { openInGoogleCalendar } = useGoogleCalendar()
const yandexMetrika = useYandexMetrika()

function handleGoogleCalendarClick(eventTitle: string) {
  yandexMetrika.reachGoal('calendar_google_add', {
    event: eventTitle,
  } as any)
}

function handleICSClick(eventTitle: string) {
  yandexMetrika.reachGoal('calendar_ics_add', {
    event: eventTitle,
  } as any)
}

onMounted(loadEvents)
</script>

<template>
  <section
    v-if="events.new.length > 0 || events.old.length > 0"
    id="meets"
    class="w-full py-12 md:py-18 lg:pt-20 lg:pb-14 rounded-[50px] bg-primary mt-20 lg:mt-28"
  >
    <div class="container px-6 md:px-10">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="flex flex-col gap-5 items-center justify-center">
          <Typography
            variant="h2"
            as="h2"
            class="text-accent"
          >
            Мероприятия для участников
          </Typography>
          <Typography
            variant="body-xl"
            as="p"
            class="max-w-[540px]"
          >
            Эксперты нашего сообщества, которые готовы поделиться экспертизой
          </Typography>
        </div>

        <div class="flex space-x-4 rounded-full select-none w-fit">
          <Tag
            :disabled="!hasFutureEvents"
            :variant="isFuture ? 'active' : 'default'"
            @click="neededEvents = 'new'"
          >
            Будущие
          </Tag>
          <Tag
            :variant="!isFuture ? 'active' : 'default'"
            @click="neededEvents = 'old'"
          >
            Прошедшие
          </Tag>
        </div>
      </div>

      <div class="grid gap-3 pt-12 lg:grid-cols-2 lg:gap-5">
        <Card
          v-for="event in visibleEvents"
          :key="event.title"
          class="min-h-[253px]"
        >
          <template #header>
            <div class="flex flex-col gap-[14px]">
              <div class="flex items-center justify-between">
                <Popover
                  :offset="12"
                  placement="top-end"
                >
                  <template #trigger>
                    <div
                      class="flex space-x-2 items-center text-accent hover:opacity-75 transition-opacity"
                    >
                      <div class="flex flex-col">
                        <Typography
                          as="span"
                          variant="date"
                        >
                          {{ formatEventDate(event.date) }} (МСК)
                        </Typography>
                        <Typography
                          v-if="event.isRepeating && event.repeatPeriod"
                          as="span"
                          variant="label"
                          class="text-xs text-muted-foreground italic"
                        >
                          Повторяется
                        </Typography>
                      </div>
                      <CalendarIcon />
                    </div>
                  </template>
                  <template #content>
                    <ul class="flex flex-col gap-2">
                      <li
                        class="cursor-pointer hover:text-accent transition-colors"
                        @click="() => { handleGoogleCalendarClick(event.title); openInGoogleCalendar(event) }"
                      >
                        <Typography variant="body-m">
                          + Google Calendar
                        </Typography>
                      </li>
                      <li
                        class="cursor-pointer hover:text-accent transition-colors"
                        @click="() => { handleICSClick(event.title); eventService.getICS(event.id) }"
                      >
                        <Typography variant="body-m">
                          + iCalendar
                        </Typography>
                      </li>
                    </ul>
                  </template>
                </Popover>
              </div>
              <Typography
                variant="h4"
                as="h4"
              >
                {{ event.title }}
              </Typography>
              <div class="flex flex-wrap items-center gap-1">
                <Label
                  v-for="tag in event.eventTags"
                  :key="tag.id"
                >{{ tag.name }}</Label>
              </div>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-6 pt-4">
              <div
                v-for="host in event.hosts"
                :key="host.id"
                class="flex items-center space-x-4 "
              >
                <TgImage
                  :username="host.tg"
                  class="rounded-full"
                  width="48"
                  height="48"
                />
                <div>
                  <Typography
                    variant="name-text"
                    as="p"
                    class="font-semibold"
                  >
                    {{ host.firstName }} {{ host.lastName }}
                  </Typography>
                  <Typography
                    variant="label"
                    class="text-muted-foreground"
                  >
                    {{ host.tg }}
                  </Typography>
                </div>
              </div>
            </div>
          </template>
        </Card>
      </div>
    </div>
    <Button
      v-if="visibleCount < events[neededEvents].length"
      variant="filled"
      class="flex mx-auto mt-12"
      @click="showMore"
    >
      Показать больше
    </Button>
  </section>
</template>
