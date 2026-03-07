<script setup lang="ts">
import type { CasinoBetResult, CasinoStats } from '@/models/casino'
import { Typography } from 'itx-ui-kit'
import { CircleDot, Dices, Loader2, RotateCw } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { casinoService } from '@/services/casino'
import { handleError } from '@/services/errorService'

const stats = ref<CasinoStats | null>(null)
const history = ref<CasinoBetResult[]>([])
const isLoading = ref(true)
const isPlaying = ref(false)
const lastResult = ref<CasinoBetResult | null>(null)
const showResult = ref(false)

const betAmount = ref(50)
const quickBets = [10, 25, 50, 100, 200]

// Dice game state
const diceTarget = ref(50)
const diceDirection = ref<'over' | 'under'>('over')

const diceMultiplier = computed(() => {
  const chance = diceDirection.value === 'over'
    ? (100 - diceTarget.value) / 100
    : diceTarget.value / 100
  if (chance <= 0)
    return 0
  return Math.round((0.97 / chance) * 100) / 100
})

async function fetchData() {
  isLoading.value = true
  try {
    const [s, h] = await Promise.all([
      casinoService.getStats(),
      casinoService.getHistory(),
    ])
    stats.value = s
    history.value = h ?? []
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

async function playGame(action: () => Promise<CasinoBetResult>) {
  if (isPlaying.value)
    return
  isPlaying.value = true
  showResult.value = false
  lastResult.value = null
  try {
    const result = await action()
    lastResult.value = result
    showResult.value = true
    if (stats.value) {
      stats.value.balance = result.balance
    }
    history.value = [result, ...history.value.slice(0, 19)]
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isPlaying.value = false
  }
}

function playCoinFlip(choice: 'heads' | 'tails') {
  playGame(() => casinoService.coinFlip(betAmount.value, choice))
}

function playDiceRoll() {
  playGame(() => casinoService.diceRoll(betAmount.value, diceTarget.value, diceDirection.value))
}

function playWheel() {
  playGame(() => casinoService.wheelSpin(betAmount.value))
}

function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  return d.toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function gameLabel(game: string) {
  const labels: Record<string, string> = {
    'coin-flip': 'Монетка',
    'dice-roll': 'Кости',
    'wheel': 'Колесо',
  }
  return labels[game] ?? game
}

onMounted(() => fetchData())
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Казино
      </Typography>
      <div
        v-if="stats"
        class="flex items-center gap-2 rounded-xl border border-border bg-card px-4 py-2"
      >
        <span class="text-sm text-muted-foreground">Баланс:</span>
        <span class="font-bold text-lg">{{ stats.balance }} б.</span>
      </div>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <!-- Bet amount controls -->
      <div class="rounded-2xl border border-border bg-card p-4 mb-6">
        <div class="flex flex-wrap items-center gap-3">
          <span class="text-sm font-medium text-muted-foreground">Ставка:</span>
          <Input
            v-model="betAmount"
            type="number"
            :min="1"
            class="w-24"
          />
          <div class="flex gap-1.5">
            <Button
              v-for="q in quickBets"
              :key="q"
              :variant="betAmount === q ? 'default' : 'outline'"
              size="sm"
              @click="betAmount = q"
            >
              {{ q }}
            </Button>
          </div>
        </div>
      </div>

      <!-- Result banner -->
      <Transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 -translate-y-2 scale-95"
        enter-to-class="opacity-100 translate-y-0 scale-100"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="showResult && lastResult"
          class="mb-6 rounded-2xl border p-5 text-center"
          :class="lastResult.profit > 0
            ? 'border-green-500/30 bg-green-500/10'
            : lastResult.profit === 0
              ? 'border-border bg-card'
              : 'border-red-500/30 bg-red-500/10'"
        >
          <div class="text-2xl font-bold mb-1">
            <span v-if="lastResult.profit > 0" class="text-green-500">
              +{{ lastResult.payout }} б.
            </span>
            <span v-else-if="lastResult.profit === 0" class="text-muted-foreground">
              0 б.
            </span>
            <span v-else class="text-red-500">
              {{ lastResult.profit }} б.
            </span>
          </div>
          <div class="text-sm text-muted-foreground">
            {{ gameLabel(lastResult.game) }}: {{ lastResult.result }}
            <span v-if="lastResult.multiplier > 0"> (x{{ lastResult.multiplier }})</span>
          </div>
        </div>
      </Transition>

      <!-- Game tabs -->
      <Tabs default-value="coin">
        <TabsList class="mb-6 w-full justify-start">
          <TabsTrigger value="coin">
            <CircleDot class="h-4 w-4 mr-1.5" />
            Монетка
          </TabsTrigger>
          <TabsTrigger value="dice">
            <Dices class="h-4 w-4 mr-1.5" />
            Кости
          </TabsTrigger>
          <TabsTrigger value="wheel">
            <RotateCw class="h-4 w-4 mr-1.5" />
            Колесо
          </TabsTrigger>
          <TabsTrigger value="history">
            История
          </TabsTrigger>
        </TabsList>

        <!-- Coin Flip -->
        <TabsContent value="coin">
          <div class="rounded-2xl border border-border bg-card p-6">
            <h3 class="text-lg font-semibold mb-2">
              Подбрось монетку
            </h3>
            <p class="text-sm text-muted-foreground mb-6">
              Угадай сторону. Выигрыш x2.
            </p>
            <div class="flex gap-3 justify-center">
              <Button
                size="lg"
                :disabled="isPlaying"
                @click="playCoinFlip('heads')"
              >
                <Loader2
                  v-if="isPlaying"
                  class="h-4 w-4 animate-spin"
                />
                Орёл
              </Button>
              <Button
                size="lg"
                variant="outline"
                :disabled="isPlaying"
                @click="playCoinFlip('tails')"
              >
                <Loader2
                  v-if="isPlaying"
                  class="h-4 w-4 animate-spin"
                />
                Решка
              </Button>
            </div>
          </div>
        </TabsContent>

        <!-- Dice Roll -->
        <TabsContent value="dice">
          <div class="rounded-2xl border border-border bg-card p-6">
            <h3 class="text-lg font-semibold mb-2">
              Кости
            </h3>
            <p class="text-sm text-muted-foreground mb-6">
              Выбери число и направление. Чем меньше шанс, тем больше множитель.
            </p>
            <div class="space-y-4">
              <div class="flex items-center gap-4">
                <div class="flex gap-2">
                  <Button
                    :variant="diceDirection === 'under' ? 'default' : 'outline'"
                    size="sm"
                    @click="diceDirection = 'under'"
                  >
                    Меньше
                  </Button>
                  <Button
                    :variant="diceDirection === 'over' ? 'default' : 'outline'"
                    size="sm"
                    @click="diceDirection = 'over'"
                  >
                    Больше
                  </Button>
                </div>
                <Input
                  v-model="diceTarget"
                  type="number"
                  :min="1"
                  :max="99"
                  class="w-20"
                />
                <span class="text-sm text-muted-foreground">
                  Множитель: <span class="font-bold text-foreground">x{{ diceMultiplier }}</span>
                </span>
              </div>
              <input
                v-model.number="diceTarget"
                type="range"
                min="1"
                max="99"
                class="w-full accent-primary"
              >
              <div class="flex justify-between text-xs text-muted-foreground">
                <span>1</span>
                <span>{{ diceTarget }}</span>
                <span>99</span>
              </div>
              <Button
                size="lg"
                class="w-full"
                :disabled="isPlaying"
                @click="playDiceRoll"
              >
                <Loader2
                  v-if="isPlaying"
                  class="h-4 w-4 animate-spin"
                />
                Бросить кости
              </Button>
            </div>
          </div>
        </TabsContent>

        <!-- Wheel -->
        <TabsContent value="wheel">
          <div class="rounded-2xl border border-border bg-card p-6">
            <h3 class="text-lg font-semibold mb-2">
              Колесо фортуны
            </h3>
            <p class="text-sm text-muted-foreground mb-6">
              Крути колесо и получи случайный множитель от x0 до x5.
            </p>
            <Button
              size="lg"
              class="w-full"
              :disabled="isPlaying"
              @click="playWheel"
            >
              <Loader2
                v-if="isPlaying"
                class="h-4 w-4 animate-spin mr-1"
              />
              <RotateCw
                v-else
                class="h-4 w-4 mr-1"
              />
              Крутить колесо
            </Button>
          </div>
        </TabsContent>

        <!-- History -->
        <TabsContent value="history">
          <div class="rounded-2xl border border-border bg-card">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Дата</TableHead>
                  <TableHead>Игра</TableHead>
                  <TableHead>Ставка</TableHead>
                  <TableHead>Выбор</TableHead>
                  <TableHead>Результат</TableHead>
                  <TableHead class="text-right">
                    Выигрыш
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow
                  v-if="history.length === 0"
                >
                  <TableCell
                    colspan="6"
                    class="text-center text-muted-foreground py-8"
                  >
                    Пока нет ставок
                  </TableCell>
                </TableRow>
                <TableRow
                  v-for="bet in history"
                  :key="bet.id"
                >
                  <TableCell class="text-muted-foreground text-xs">
                    {{ formatDate(bet.createdAt) }}
                  </TableCell>
                  <TableCell>{{ gameLabel(bet.game) }}</TableCell>
                  <TableCell>{{ bet.betAmount }} б.</TableCell>
                  <TableCell>{{ bet.betChoice }}</TableCell>
                  <TableCell>{{ bet.result }}</TableCell>
                  <TableCell class="text-right">
                    <Badge
                      :variant="bet.profit > 0 ? 'default' : bet.profit === 0 ? 'secondary' : 'destructive'"
                    >
                      {{ bet.profit > 0 ? '+' : '' }}{{ bet.profit }} б.
                    </Badge>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </TabsContent>
      </Tabs>
    </template>
  </div>
</template>
