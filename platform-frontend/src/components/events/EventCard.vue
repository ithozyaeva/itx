<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import { useDictionary } from '@/composables/useDictionary'
import { useUser } from '@/composables/useUser'
import { dateFormatter, wrapLinks } from '@/lib/utils'
import { eventsService } from '@/services/events'
import { CalendarIcon, Tag, Typography } from 'itx-ui-kit'
import { ChevronDown, MapPin } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import Button from '../ui/button/Button.vue'

const props = defineProps<{
  event: CommunityEvent
}>()
const event = ref(props.event)

const user = useUser()
const isMembersExpanded = ref(false)

const formattedDate = computed(() => dateFormatter.format(new Date(props.event.date)))

const isHost = computed(() => user.value ? event.value.hosts.map(item => item.id).includes(user.value.id) : false)
const isMember = computed(() => user.value ? event.value.members.map(item => item.id).includes(user.value.id) : false)
const isPassedEvent = computed(() => new Date(props.event.date) < new Date())

// Форматирование информации о повторениях
const repeatInfo = computed(() => {
  if (!props.event.isRepeating || !props.event.repeatPeriod) {
    return null
  }

  const periodLabels: Record<string, string> = {
    DAILY: 'день',
    WEEKLY: 'неделя',
    MONTHLY: 'месяц',
    YEARLY: 'год',
  }

  const periodLabel = periodLabels[props.event.repeatPeriod] || props.event.repeatPeriod.toLowerCase()
  const interval = props.event.repeatInterval || 1

  let info = `Повторяется каждые ${interval} ${interval === 1 ? periodLabel : getPluralForm(periodLabel, interval)}`

  if (props.event.repeatEndDate) {
    const endDate = new Date(props.event.repeatEndDate)
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

async function applyEvent(eventId: number) {
  event.value = await eventsService.applyEvent(eventId)
}
async function declineEvent(eventId: number) {
  event.value = await eventsService.declineEvent(eventId)
}

function getICS() {
  window.open(`${window.location.origin}/api/events/ics?eventId=${event.value.id}`, '_blank')
}

const { placeTypesObject } = useDictionary(['placeTypes'])
</script>

<template>
  <div data-reveal class="bg-card rounded-3xl border p-4 hover:shadow-md transition-shadow flex flex-col gap-1">
    <div class="flex justify-between items-start mb-2">
      <Typography variant="h4" as="h3">
        {{ event.title }}
      </Typography>
      <div class="space-x-2">
        <Tag>
          {{ placeTypesObject[event.placeType] }}
        </Tag>
        <Tag
          v-if="event.eventType !== 'ONLINE' && !!event.customPlaceType"
        >
          {{ event.customPlaceType }}
        </Tag>
        <Button v-if="!isPassedEvent" class="cursor-pointer" size="sm" variant="outline" @click="getICS">
          + Добавить в календарь
        </Button>
      </div>
    </div>
    <p class="text-muted-foreground mb-4">
      {{ event.description }}
    </p>
    <div class="space-y-2 text-sm">
      <div class="flex items-center gap-2">
        <CalendarIcon />
        <div class="flex flex-col">
          <span>{{ formattedDate }} ({{ event.timezone || 'UTC' }})</span>
          <span v-if="repeatInfo" class="text-xs text-muted-foreground italic">
            {{ repeatInfo }}
          </span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <MapPin />
        <span v-if="event.placeType === 'OFFLINE'">{{ event.place }}</span>
        <a v-if="event.placeType === 'ONLINE'" :href="event.place" target="_blank" class="underline">{{ event.place }}</a>
        <p v-if="event.placeType === 'HYBRID'" v-html="wrapLinks(event.place)" />
      </div>
      <div class="flex items-center gap-2">
        <span class="text-bold">Cспикеры: </span>
        <span>{{ event.hosts.map(host => `${host.firstName} ${host.lastName}`).join(', ') }}</span>
      </div>
      <div v-if="event.videoLink">
        <span>Запись встречи:</span>
        <a :href="event.videoLink" target="_blank" class="underline">
          {{ event.videoLink }}
        </a>
      </div>
    </div>
    <div class="flex flex-col">
      <button
        class="flex items-center gap-2 text-left hover:text-primary transition-colors"
        @click="toggleMembers"
      >
        <span class="text-sm text-bold">Участники ({{ event.members.length }})</span>
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
          <div class="mt-1 space-y-1 text-muted-foreground">
            <span>{{ event.members.map(member => `${member.firstName} ${member.lastName}`).join(', ') }}</span>
          </div>
        </div>
      </div>
    </div>
    <div v-if="!isPassedEvent" class="self-end">
      <Button v-if="isMember" @click="declineEvent(event.id)">
        Отменить участие
      </Button>
      <Button v-if="!isMember && !isHost" @click="applyEvent(event.id)">
        Участвую!
      </Button>
    </div>
  </div>
</template>
