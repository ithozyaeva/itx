<script setup lang="ts">
import type { CommunityEvent } from '@/services/events'
import { Calendar as CalendarIcon } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import TgImage from '@/components/ui/TgImage.vue'
import Button from '@/components/ui/UiButton.vue'
import Popover from '@/components/ui/UiPopover.vue'
import Typography from '@/components/ui/UiTypography.vue'
import { useGoogleCalendar } from '@/composables/useGoogleCalendar.ts'
import { eventService } from '@/services/events'

const SLASH_RE = /\//g

const events = ref<Record<'new' | 'old', CommunityEvent[]>>({
  new: [],
  old: [],
})

function formatEventDate(dateString: string): string {
  const utcDate = new Date(dateString)
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

function formatStamp(dateString: string): string {
  const d = new Date(dateString)
  const fmt = new Intl.DateTimeFormat('en-GB', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
    timeZone: 'Europe/Moscow',
  }).format(d)
  return fmt.replace(',', '').replace(SLASH_RE, '-')
}

async function loadEvents() {
  events.value.old = await eventService.getOld()
  events.value.new = await eventService.getNext()
}

const visibleEvents = computed(() => {
  return events.value.new.slice(0, 2)
})

const hasFutureEvents = computed(() => events.value.new.length > 0)

const { openInGoogleCalendar } = useGoogleCalendar()
const yandexMetrika = useYandexMetrika()

function handleGoogleCalendarClick(eventTitle: string) {
  yandexMetrika.reachGoal('calendar_google_add', { event: eventTitle } as any)
}

function handleICSClick(eventTitle: string) {
  yandexMetrika.reachGoal('calendar_ics_add', { event: eventTitle } as any)
}

onMounted(loadEvents)
</script>

<template>
  <section
    v-if="events.new.length > 0 || events.old.length > 0"
    id="meets"
    class="w-full pt-20 md:pt-32 lg:pt-40"
  >
    <div class="container px-6 md:px-10">
      <div class="reveal">
        <SectionHeader
          index="03"
          path="~/community/events.log"
          title="Расписание"
          subtitle="Еженедельные онлайн-встречи. Скачивай в календарь — пропустить будет нельзя."
        />
      </div>

      <Typography
        v-if="!hasFutureEvents"
        variant="body-l"
        as="p"
        class="text-muted-foreground mt-12 font-mono text-sm"
      >
        &gt; no scheduled events — следите за обновлениями
      </Typography>

      <div class="mt-12 md:mt-16">
        <div class="border border-accent/20 bg-background/60 backdrop-blur-sm">
          <!-- log header -->
          <div class="flex items-center gap-2 px-5 py-3 border-b border-accent/20 bg-accent/5">
            <div class="flex gap-1.5">
              <span class="w-2.5 h-2.5 rounded-full bg-term-magenta/70" />
              <span class="w-2.5 h-2.5 rounded-full bg-term-amber/70" />
              <span class="w-2.5 h-2.5 rounded-full bg-accent/70" />
            </div>
            <span class="ml-3 font-mono text-[11px] text-foreground/50 tracking-wider">
              tail -f /var/log/community/events.log
            </span>
          </div>

          <!-- event entries -->
          <ul>
            <li
              v-for="event in visibleEvents"
              :key="event.title"
              class="group relative grid grid-cols-1 md:grid-cols-[auto_1fr_auto] gap-6 md:gap-8 px-5 md:px-8 py-6 md:py-8 border-b border-accent/10 hover:bg-accent/[0.03] transition-colors last:border-b-0"
            >
              <!-- timestamp col -->
              <div class="flex md:flex-col gap-3 md:gap-1 items-center md:items-start font-mono text-xs">
                <span class="text-accent">[{{ formatStamp(event.date) }}]</span>
                <Popover
                  :offset="12"
                  placement="top-start"
                >
                  <template #trigger>
                    <button
                      type="button"
                      class="flex items-center gap-1.5 text-foreground/60 hover:text-accent transition-colors"
                    >
                      <CalendarIcon class="w-3.5 h-3.5" />
                      <span>+add</span>
                    </button>
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
                <span
                  v-if="event.isRepeating && event.repeatPeriod"
                  class="text-term-amber/80 text-[10px] uppercase tracking-wider md:mt-1"
                >
                  recurring
                </span>
              </div>

              <!-- content col -->
              <div class="flex flex-col gap-3">
                <h3 class="font-display uppercase text-xl md:text-2xl leading-tight text-foreground group-hover:text-accent transition-colors">
                  {{ event.title }}
                </h3>
                <div class="flex flex-wrap items-center gap-x-3 gap-y-1 font-mono text-[11px] uppercase tracking-wider text-foreground/50">
                  <span class="text-accent/70">// tags:</span>
                  <span
                    v-for="tag in event.eventTags"
                    :key="tag.id"
                    class="text-foreground/70"
                  >
                    #{{ tag.name }}
                  </span>
                </div>
                <div class="text-xs font-mono text-foreground/50 md:hidden">
                  {{ formatEventDate(event.date) }} МСК
                </div>
              </div>

              <!-- hosts col -->
              <div class="flex flex-col gap-2 md:min-w-[180px]">
                <span class="font-mono text-[10px] uppercase tracking-wider text-foreground/40">
                  hosts
                </span>
                <div class="flex flex-col gap-2">
                  <div
                    v-for="host in event.hosts"
                    :key="host.id"
                    class="flex items-center gap-3"
                  >
                    <TgImage
                      :username="host.tg"
                      :avatar-url="host.avatarUrl"
                      class="w-9 h-9 shrink-0"
                      width="36"
                      height="36"
                    />
                    <div class="flex flex-col">
                      <span class="text-sm text-foreground">{{ host.firstName }} {{ host.lastName }}</span>
                      <span class="font-mono text-[10px] text-foreground/50">@{{ host.tg }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <Button
        v-if="events.new.length > 2"
        variant="filled"
        as="a"
        href="/platform/events"
        class="flex mx-auto mt-12"
      >
        Все события
      </Button>
    </div>
  </section>
</template>
