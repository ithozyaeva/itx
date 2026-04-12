<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import { Calendar, ChevronDown, Crown, Loader2, MapPin } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Typography } from '@/components/ui/typography'
import { useDictionary } from '@/composables/useDictionary'
import { getNextOccurrenceDate } from '@/composables/useEventOccurrence'
import { useGoogleCalendar } from '@/composables/useGoogleCalendar'
import { useUser } from '@/composables/useUser'
import { dateFormatter, wrapLinks } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'
import Button from '../ui/button/Button.vue'

const props = defineProps<{
  event: CommunityEvent
}>()
const event = ref(props.event)

const user = useUser()
const isMembersExpanded = ref(false)

const nextOccurrenceDate = computed(() => getNextOccurrenceDate(event.value))
const formattedDate = computed(() => dateFormatter.format(nextOccurrenceDate.value))

const isExclusive = computed(() => !!event.value.exclusiveChatId)
const isHost = computed(() => user.value ? (event.value.hosts ?? []).map(item => item.id).includes(user.value.id) : false)
const isMember = computed(() => user.value ? (event.value.members ?? []).map(item => item.id).includes(user.value.id) : false)
const isPassedEvent = computed(() => {
  if (event.value.isRepeating && event.value.repeatPeriod) {
    return !!(event.value.repeatEndDate && new Date(event.value.repeatEndDate) < new Date())
  }
  return new Date(event.value.date) < new Date()
})
const isFull = computed(() => event.value.maxParticipants > 0 && (event.value.members?.length ?? 0) >= event.value.maxParticipants)

// Форматирование информации о повторениях
const repeatInfo = computed(() => {
  if (!event.value.isRepeating || !event.value.repeatPeriod) {
    return null
  }

  const periodLabels: Record<string, string> = {
    DAILY: 'день',
    WEEKLY: 'неделя',
    MONTHLY: 'месяц',
    YEARLY: 'год',
  }

  const periodLabel = periodLabels[event.value.repeatPeriod] || event.value.repeatPeriod.toLowerCase()
  const interval = event.value.repeatInterval || 1

  let info = `Повторяется каждые ${interval} ${interval === 1 ? periodLabel : getPluralForm(periodLabel, interval)}`

  if (event.value.repeatEndDate) {
    const endDate = new Date(event.value.repeatEndDate)
    info += ` до ${dateFormatter.format(endDate)}`
  }

  return info
})

function getPluralForm(word: string, count: number): string {
  const forms: Record<string, string[]> = {
    день: ['дня', 'дней'],
    неделя: ['недели', 'недель'],
    месяц: ['месяца', 'месяцев'],
    год: ['года', 'лет'],
  }

  if (!forms[word])
    return word

  if (count % 10 === 1 && count % 100 !== 11) {
    return word
  }
  else if (count % 10 >= 2 && count % 10 <= 4 && (count % 100 < 10 || count % 100 >= 20)) {
    return forms[word][0]
  }
  else {
    return forms[word][1]
  }
}

function toggleMembers() {
  isMembersExpanded.value = !isMembersExpanded.value
}

const isApplying = ref(false)
const isDeclining = ref(false)

async function applyEvent(eventId: number) {
  isApplying.value = true
  try {
    event.value = await eventsService.applyEvent(eventId)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isApplying.value = false
  }
}
async function declineEvent(eventId: number) {
  isDeclining.value = true
  try {
    event.value = await eventsService.declineEvent(eventId)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isDeclining.value = false
  }
}

function getICS() {
  const link = document.createElement('a')
  link.href = `/api/events/ics?eventId=${event.value.id}`
  link.download = `${event.value.title}.ics`
  link.click()
}

function getHostInitials(host: { firstName: string, lastName: string }) {
  return `${host.firstName.charAt(0)}${host.lastName.charAt(0)}`.toUpperCase()
}

const { placeTypesObject } = useDictionary(['placeTypes'])
const { openInGoogleCalendar } = useGoogleCalendar()
</script>

<template>
  <div
    data-reveal
    class="rounded-sm border p-4 transition-all duration-200 flex flex-col gap-2 terminal-card"
    :class="isExclusive
      ? 'bg-gradient-to-br from-amber-50/80 to-yellow-50/50 dark:from-amber-950/30 dark:to-yellow-950/20 border-amber-300/60 dark:border-amber-600/40 hover:shadow-lg hover:shadow-amber-200/30 dark:hover:shadow-amber-900/20'
      : 'bg-card hover:shadow-md'"
  >
    <!-- Exclusive badge -->
    <div v-if="isExclusive" class="flex items-center gap-1.5 text-amber-600 dark:text-amber-400">
      <Crown class="h-4 w-4" />
      <span class="text-xs font-semibold uppercase tracking-wide">{{ event.exclusiveChatTitle || 'Эксклюзив' }}</span>
    </div>

    <!-- Date with calendar popover -->
    <Popover v-if="!isPassedEvent">
      <PopoverTrigger>
        <div class="flex items-center gap-2 text-accent text-sm cursor-pointer">
          <Calendar class="shrink-0" />
          <span>{{ formattedDate }} ({{ event.timezone || 'UTC' }})</span>
        </div>
      </PopoverTrigger>
      <PopoverContent>
        <div class="flex flex-col gap-1 text-sm">
          <button
            class="text-left hover:text-accent transition-colors py-1"
            @click="openInGoogleCalendar(event)"
          >
            + Google Calendar
          </button>
          <button
            class="text-left hover:text-accent transition-colors py-1"
            @click="getICS"
          >
            + iCalendar
          </button>
        </div>
      </PopoverContent>
    </Popover>
    <div v-else class="flex items-center gap-2 text-muted-foreground text-sm">
      <Calendar class="shrink-0" />
      <span>{{ formattedDate }} ({{ event.timezone || 'UTC' }})</span>
    </div>

    <!-- Repeat info -->
    <span v-if="repeatInfo" class="text-xs text-muted-foreground italic">
      {{ repeatInfo }}
    </span>

    <!-- Title -->
    <Typography
      variant="h4"
      as="h3"
    >
      {{ event.title }}
    </Typography>

    <!-- Tags -->
    <div class="flex flex-wrap items-center gap-1.5">
      <span class="inline-flex items-center rounded-full border border-accent/30 px-2 py-0.5 text-xs text-accent">
        {{ placeTypesObject[event.placeType] }}
      </span>
      <span
        v-if="event.eventType !== 'ONLINE' && !!event.customPlaceType"
        class="inline-flex items-center rounded-full border border-accent/30 px-2 py-0.5 text-xs text-accent"
      >
        {{ event.customPlaceType }}
      </span>
      <span
        v-for="tag in event.eventTags"
        :key="tag.id"
        class="inline-flex items-center rounded-full border border-accent/30 px-2 py-0.5 text-xs text-accent"
      >
        {{ tag.name }}
      </span>
    </div>

    <!-- Description -->
    <p class="text-muted-foreground text-sm break-words line-clamp-2">
      {{ event.description }}
    </p>

    <!-- Place/location -->
    <div class="flex items-center gap-2 text-sm">
      <MapPin class="shrink-0 h-4 w-4" />
      <span v-if="event.placeType === 'OFFLINE'" class="break-all line-clamp-1">{{ event.place }}</span>
      <p v-if="event.placeType === 'ONLINE'" class="break-all line-clamp-1" v-html="wrapLinks(event.place)" />
      <p v-if="event.placeType === 'HYBRID'" class="break-all line-clamp-1" v-html="wrapLinks(event.place)" />
    </div>

    <!-- Hosts with avatars -->
    <div class="flex items-center gap-2 text-sm">
      <span class="font-medium shrink-0">Спикеры:</span>
      <div class="flex items-center gap-2 flex-wrap">
        <div
          v-for="host in event.hosts"
          :key="host.id"
          class="flex items-center gap-1.5"
        >
          <img
            v-if="host.avatarUrl"
            :src="host.avatarUrl"
            :alt="`${host.firstName} ${host.lastName}`"
            class="w-8 h-8 rounded-full object-cover"
          >
          <div
            v-else
            class="w-8 h-8 rounded-full bg-accent text-white flex items-center justify-center text-xs font-medium"
          >
            {{ getHostInitials(host) }}
          </div>
          <span>{{ host.firstName }} {{ host.lastName }}</span>
        </div>
      </div>
    </div>

    <!-- Video/recording links -->
    <div v-if="event.videoLink" class="flex items-center gap-2 text-sm">
      <span class="shrink-0">Трансляция:</span>
      <a :href="event.videoLink" target="_blank" class="underline break-all line-clamp-1">
        {{ event.videoLink }}
      </a>
    </div>
    <div v-if="event.recordingUrl" class="flex items-center gap-2 text-sm">
      <span class="shrink-0">Запись:</span>
      <a :href="event.recordingUrl" target="_blank" class="underline break-all text-accent line-clamp-1">
        {{ event.recordingUrl }}
      </a>
    </div>

    <!-- Members expandable -->
    <div class="flex flex-col">
      <button
        class="flex items-center gap-2 text-left hover:text-accent transition-colors"
        @click="toggleMembers"
      >
        <span class="text-sm font-medium">Участники ({{ event.members.length }}{{ event.maxParticipants > 0 ? `/${event.maxParticipants}` : '' }})</span>
        <ChevronDown
          class="w-4 h-4 transition-transform duration-200"
          :class="{ 'rotate-180': isMembersExpanded }"
        />
      </button>
      <div
        class="grid transition-all duration-200 ease-in-out"
        :class="[
          isMembersExpanded
            ? 'grid-rows-[1fr] opacity-100'
            : 'grid-rows-[0fr] opacity-0',
        ]"
      >
        <div class="overflow-hidden">
          <div class="mt-1 space-y-1 text-muted-foreground text-sm break-words">
            <span>{{ event.members.map(member => `${member.firstName} ${member.lastName}`).join(', ') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Action button -->
    <div v-if="!isPassedEvent" class="self-end">
      <ConfirmDialog
        v-if="isMember"
        title="Отменить участие?"
        description="Вы будете исключены из списка участников события."
        confirm-label="Отменить участие"
        @confirm="declineEvent(event.id)"
      >
        <template #trigger>
          <Button :disabled="isDeclining">
            <Loader2 v-if="isDeclining" class="h-4 w-4 animate-spin mr-1" />
            Отменить участие
          </Button>
        </template>
      </ConfirmDialog>
      <Button v-if="!isMember && !isHost" :disabled="isApplying || isFull" @click="applyEvent(event.id)">
        <Loader2 v-if="isApplying" class="h-4 w-4 animate-spin mr-1" />
        {{ isFull ? 'Мест нет' : 'Участвую!' }}
      </Button>
    </div>
  </div>
</template>
