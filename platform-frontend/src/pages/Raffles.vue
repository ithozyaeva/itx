<script setup lang="ts">
import type { RaffleItem } from '@/models/raffle'
import { Typography } from 'itx-ui-kit'
import { Gift, Loader2, Ticket, Trophy } from 'lucide-vue-next'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useSSE } from '@/composables/useSSE'
import { useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { raffleService } from '@/services/raffles'

const user = useUser()
const items = ref<RaffleItem[]>([])
const isLoading = ref(true)
const buyingId = ref<number | null>(null)
const ticketCounts = ref<Record<number, number>>({})
const filterMode = ref<'all' | 'my'>('all')

const filteredItems = computed(() =>
  filterMode.value === 'my' ? items.value.filter(r => r.myTickets > 0) : items.value,
)
const activeRaffles = computed(() => filteredItems.value.filter(r => r.status === 'ACTIVE'))
const finishedRaffles = computed(() => filteredItems.value.filter(r => r.status === 'FINISHED'))

async function fetchRaffles() {
  isLoading.value = true
  try {
    items.value = await raffleService.getAll() ?? []
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

function getTicketCount(id: number) {
  return ticketCounts.value[id] ?? 1
}

async function buyTickets(id: number) {
  if (buyingId.value)
    return
  buyingId.value = id
  try {
    await raffleService.buyTickets(id, getTicketCount(id))
    ticketCounts.value[id] = 1
    await fetchRaffles()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    buyingId.value = null
  }
}

const now = ref(Date.now())
const tickTimer = setInterval(() => {
  now.value = Date.now()
}, 60000)
onUnmounted(() => clearInterval(tickTimer))

function timeLeft(endsAt: string) {
  const diff = new Date(endsAt).getTime() - now.value
  if (diff <= 0)
    return 'Завершён'
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(hours / 24)
  if (days > 0)
    return `${days} дн. ${hours % 24} ч.`
  const mins = Math.floor((diff % 3600000) / 60000)
  if (hours > 0)
    return `${hours} ч. ${mins} мин.`
  return `${mins} мин.`
}

function winnerName(r: RaffleItem) {
  return [r.winnerFirstName, r.winnerLastName].filter(Boolean).join(' ') || '—'
}

useSSE('raffles', () => fetchRaffles())

onMounted(() => {
  fetchRaffles()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Розыгрыши
      </Typography>
      <div
        v-if="user"
        class="flex gap-1 bg-muted rounded-lg p-0.5"
      >
        <button
          type="button"
          class="px-3 py-1.5 text-sm rounded-md transition-colors"
          :class="filterMode === 'all' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'"
          @click="filterMode = 'all'"
        >
          Все
        </button>
        <button
          type="button"
          class="px-3 py-1.5 text-sm rounded-md transition-colors"
          :class="filterMode === 'my' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'"
          @click="filterMode = 'my'"
        >
          Мои розыгрыши
        </button>
      </div>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <div
        v-if="filteredItems.length === 0"
        class="text-center py-12 text-muted-foreground"
      >
        <Gift class="h-12 w-12 mx-auto mb-3 opacity-50" />
        <p>{{ filterMode === 'my' ? 'Вы ещё не участвовали в розыгрышах' : 'Розыгрышей пока нет' }}</p>
      </div>

      <!-- Active -->
      <div
        v-if="activeRaffles.length"
        class="mb-8"
      >
        <h2 class="text-lg font-semibold mb-4">
          Активные
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="raffle in activeRaffles"
            :key="raffle.id"
            class="rounded-2xl border bg-card border-border p-5"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="min-w-0">
                <h3 class="font-semibold break-words">
                  {{ raffle.title }}
                </h3>
                <p
                  v-if="raffle.description"
                  class="text-sm text-muted-foreground mt-1"
                >
                  {{ raffle.description }}
                </p>
              </div>
              <Gift class="h-5 w-5 text-primary shrink-0" />
            </div>

            <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm mb-3">
              <div>
                <span class="text-muted-foreground">Приз:</span>
                <span class="font-medium ml-1">{{ raffle.prize }}</span>
              </div>
              <div>
                <span class="text-muted-foreground">Билет:</span>
                <span class="font-medium ml-1">{{ raffle.ticketCost }} б.</span>
              </div>
            </div>

            <div class="flex items-center gap-4 text-sm text-muted-foreground mb-4">
              <div class="flex items-center gap-1">
                <Ticket class="h-3.5 w-3.5" />
                <span>Продано: {{ raffle.totalTickets }}{{ raffle.maxTickets ? `/${raffle.maxTickets}` : '' }}</span>
              </div>
              <div>
                Ваших: {{ raffle.myTickets }}
              </div>
              <div class="ml-auto">
                {{ timeLeft(raffle.endsAt) }}
              </div>
            </div>

            <div class="flex items-center gap-2">
              <Select
                :model-value="String(getTicketCount(raffle.id))"
                @update:model-value="ticketCounts[raffle.id] = Number($event)"
              >
                <SelectTrigger class="w-24">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="n in 10"
                    :key="n"
                    :value="String(n)"
                  >
                    {{ n }} шт.
                  </SelectItem>
                </SelectContent>
              </Select>
              <button
                class="flex-1 px-4 py-1.5 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
                :disabled="buyingId === raffle.id"
                @click="buyTickets(raffle.id)"
              >
                <Loader2
                  v-if="buyingId === raffle.id"
                  class="h-4 w-4 animate-spin inline mr-1"
                />
                Купить ({{ raffle.ticketCost * getTicketCount(raffle.id) }} б.)
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Finished -->
      <div v-if="finishedRaffles.length">
        <h2 class="text-lg font-semibold mb-4">
          Завершённые
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="raffle in finishedRaffles"
            :key="raffle.id"
            class="rounded-2xl border bg-card border-border p-5 opacity-75"
          >
            <h3 class="font-semibold mb-1">
              {{ raffle.title }}
            </h3>
            <p class="text-sm text-muted-foreground mb-3">
              Приз: {{ raffle.prize }}
            </p>
            <div
              v-if="raffle.winnerId"
              class="flex items-center gap-2 text-sm"
            >
              <Trophy class="h-4 w-4 text-yellow-500" />
              <span class="text-muted-foreground">Победитель:</span>
              <router-link
                :to="`/members/${raffle.winnerId}`"
                class="font-medium hover:underline"
              >
                {{ winnerName(raffle) }}
              </router-link>
            </div>
            <p
              v-else
              class="text-sm text-muted-foreground"
            >
              Без участников
            </p>
            <div
              v-if="raffle.myTickets > 0"
              class="text-sm text-muted-foreground mt-1"
            >
              Ваших билетов: {{ raffle.myTickets }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
