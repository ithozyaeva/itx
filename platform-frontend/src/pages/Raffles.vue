<script setup lang="ts">
import type { RaffleItem, RaffleTicketSource } from '@/models/raffle'
import {
  Award,
  CalendarCheck,
  CalendarDays,
  CheckCircle2,
  Circle,
  Flame,
  Gift,
  ListChecks,
  Loader2,
  Star,
  Ticket,
  Trophy,
  Users,
} from 'lucide-vue-next'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import EmptyState from '@/components/common/EmptyState.vue'
import { TintedIcon } from '@/components/progress'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { useSSE } from '@/composables/useSSE'
import { useUser } from '@/composables/useUser'
import { dailiesService } from '@/services/dailies'
import { handleError } from '@/services/errorService'
import { raffleService } from '@/services/raffles'

interface RaffleSourceRow {
  key: RaffleTicketSource
  label: string
  hint: string
  icon: typeof CalendarCheck
}

// Источники, которые пользователь видит как «способы получить билет».
// purchase/legacy скрыты — они либо для manual-раффлов, либо историч.
const TICKET_SOURCES: RaffleSourceRow[] = [
  { key: 'check_in', label: 'Check-in', hint: 'Зайди и нажми «check-in»', icon: CalendarCheck },
  { key: 'daily_task', label: 'Дейлик', hint: 'Закрой любую из 5 дневных задач', icon: ListChecks },
  { key: 'all_dailies_bonus', label: 'Все 5 дейликов', hint: 'Закрой все задачи дня', icon: Star },
  { key: 'challenge', label: 'Челлендж', hint: 'Заверши недельный или месячный челлендж', icon: Award },
  { key: 'attend_event', label: 'Событие', hint: 'Посети любое событие', icon: CalendarDays },
]

const { toast } = useToast()

const user = useUser()
const items = ref<RaffleItem[]>([])
const isLoading = ref(true)
const buyingId = ref<number | null>(null)
const ticketCounts = ref<Record<number, number>>({})
const filterMode = ref<'all' | 'my'>('all')

const dailyRaffle = ref<RaffleItem | null>(null)

function hasSource(source: RaffleTicketSource) {
  return dailyRaffle.value?.mySources?.includes(source) ?? false
}

// Daily-раффл показывается отдельной секцией сверху, поэтому исключаем
// его из общего списка.
const manualItems = computed(() => items.value.filter(r => r.kind !== 'daily'))
const filteredItems = computed(() =>
  filterMode.value === 'my' ? manualItems.value.filter(r => r.myTickets > 0) : manualItems.value,
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

async function fetchDailyRaffle() {
  try {
    const resp = await dailiesService.getDailyRaffle()
    if (resp && 'id' in resp)
      dailyRaffle.value = resp as RaffleItem
    else
      dailyRaffle.value = null
  }
  catch {
    dailyRaffle.value = null
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
    toast({ title: 'Билеты куплены' })
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

useSSE('raffles', () => {
  fetchRaffles()
  fetchDailyRaffle()
})

onMounted(() => {
  fetchRaffles()
  fetchDailyRaffle()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/raffles
    </div>
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

    <!-- Ежедневный авто-розыгрыш — отдельная секция сверху -->
    <div
      v-if="dailyRaffle"
      class="rounded-sm border-2 border-orange-500/30 bg-gradient-to-br from-orange-500/5 via-card to-yellow-500/5 p-5 md:p-6 mb-6 terminal-card"
    >
      <div class="flex items-start justify-between gap-3 mb-3 flex-wrap">
        <div class="flex items-center gap-3 min-w-0">
          <TintedIcon :icon="Flame" tone="orange" size="lg" />
          <div class="min-w-0">
            <p class="font-mono text-[10px] uppercase tracking-widest text-orange-500/80 mb-0.5">
              // ежедневный авто-розыгрыш
            </p>
            <h3 class="font-bold text-lg leading-tight">
              {{ dailyRaffle.title }}
            </h3>
          </div>
        </div>
        <div class="text-right shrink-0">
          <div class="text-xs text-muted-foreground">
            До розыгрыша
          </div>
          <div class="font-mono font-bold text-base">
            {{ timeLeft(dailyRaffle.endsAt) }}
          </div>
        </div>
      </div>
      <p class="text-sm text-muted-foreground mb-4">
        {{ dailyRaffle.description || 'Случайный участник получает 100 баллов в 23:59 МСК.' }}
      </p>
      <div class="flex items-center justify-between gap-3 flex-wrap mb-4">
        <div class="flex items-center gap-4 text-sm">
          <div class="flex items-center gap-1.5">
            <Users class="h-4 w-4 text-muted-foreground" />
            <span class="font-medium">{{ dailyRaffle.totalTickets }}</span>
            <span class="text-muted-foreground">участн.</span>
          </div>
          <div class="flex items-center gap-1.5">
            <Ticket class="h-4 w-4 text-muted-foreground" />
            <span class="text-muted-foreground">Ваших:</span>
            <span class="font-medium">{{ dailyRaffle.myTickets }}</span>
          </div>
          <div class="flex items-center gap-1.5 text-yellow-500 font-medium">
            <Trophy class="h-4 w-4" />
            {{ dailyRaffle.prize }}
          </div>
        </div>
      </div>

      <div class="border-t border-orange-500/15 pt-4">
        <p class="font-mono text-[10px] uppercase tracking-widest text-orange-500/80 mb-3">
          // как получить билет
        </p>
        <ul class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-2">
          <li
            v-for="src in TICKET_SOURCES"
            :key="src.key"
            class="flex items-start gap-2.5 text-sm"
            :class="hasSource(src.key) ? 'text-foreground' : 'text-muted-foreground'"
          >
            <component
              :is="hasSource(src.key) ? CheckCircle2 : Circle"
              class="h-4 w-4 mt-0.5 shrink-0"
              :class="hasSource(src.key) ? 'text-green-500' : 'text-muted-foreground/40'"
            />
            <div class="min-w-0">
              <div class="font-medium leading-tight">
                {{ src.label }}
              </div>
              <div class="text-xs text-muted-foreground/80">
                {{ src.hint }}
              </div>
            </div>
          </li>
        </ul>
        <RouterLink
          v-if="!hasSource('check_in')"
          to="/progress?tab=today"
          class="mt-4 inline-flex items-center justify-center px-3 py-2 rounded-sm bg-primary text-primary-foreground text-xs font-medium hover:bg-primary/90 transition-colors min-h-[36px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Сделай check-in →
        </RouterLink>
      </div>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <EmptyState
        v-if="filteredItems.length === 0"
        :icon="Gift"
        variant="dashed"
        :title="filterMode === 'my' ? 'Вы ещё не участвовали в розыгрышах' : 'Розыгрышей пока нет'"
        :description="filterMode === 'my' ? 'Купите билет в активный розыгрыш или дождитесь нового.' : undefined"
      />

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
            class="rounded-sm border bg-card border-border terminal-card p-5"
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
                class="flex-1 px-4 py-1.5 rounded-sm bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
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
            class="rounded-sm border bg-card border-border terminal-card p-5 opacity-75"
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
              <RouterLink
                :to="`/members/${raffle.winnerId}`"
                class="font-medium hover:underline"
              >
                {{ winnerName(raffle) }}
              </RouterLink>
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
