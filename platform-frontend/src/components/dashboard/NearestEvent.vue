<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import { Loader2, Radio } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { useUser } from '@/composables/useUser'
import { dateFormatter } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'

const props = defineProps<{
  event: CommunityEvent
}>()

const event = ref(props.event)
const user = useUser()
const isApplying = ref(false)

const formattedDate = computed(() => dateFormatter.format(new Date(event.value.date)))
const isMember = computed(() => user.value ? event.value.members.map(m => m.id).includes(user.value.id) : false)
const isHost = computed(() => user.value ? event.value.hosts.map(h => h.id).includes(user.value.id) : false)
const isLive = computed(() => {
  const now = new Date()
  const eventDate = new Date(event.value.date)
  const diffMs = now.getTime() - eventDate.getTime()
  return diffMs >= 0 && diffMs < 2 * 60 * 60 * 1000
})

async function applyEvent() {
  isApplying.value = true
  try {
    event.value = await eventsService.applyEvent(event.value.id)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isApplying.value = false
  }
}
</script>

<template>
  <div class="rounded-3xl border bg-card p-6">
    <div class="flex items-center gap-2 mb-4">
      <h2 class="text-lg font-semibold">
        Ближайшее событие
      </h2>
      <span
        v-if="isLive"
        class="inline-flex items-center gap-1 rounded-full bg-red-500 px-2.5 py-0.5 text-xs font-medium text-white"
      >
        <Radio class="h-3 w-3" />
        LIVE
      </span>
    </div>
    <div class="flex flex-col gap-3">
      <h3 class="text-xl font-bold">
        {{ event.title }}
      </h3>
      <p class="text-muted-foreground text-sm line-clamp-2">
        {{ event.description }}
      </p>
      <div class="flex flex-wrap gap-4 text-sm text-muted-foreground">
        <span>{{ formattedDate }}</span>
        <span v-if="event.hosts.length">
          Ведущий: {{ event.hosts.map(h => `${h.firstName} ${h.lastName}`).join(', ') }}
        </span>
      </div>
      <div class="flex items-center gap-3 mt-2">
        <Button
          v-if="!isMember && !isHost"
          :disabled="isApplying"
          @click="applyEvent"
        >
          <Loader2
            v-if="isApplying"
            class="h-4 w-4 animate-spin mr-1"
          />
          Записаться
        </Button>
        <span
          v-if="isMember"
          class="text-sm text-green-600 font-medium"
        >
          Вы записаны
        </span>
        <span class="text-sm text-muted-foreground">
          {{ event.members.length }}{{ event.maxParticipants > 0 ? `/${event.maxParticipants}` : '' }} участников
        </span>
      </div>
    </div>
  </div>
</template>
