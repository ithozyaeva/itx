<script setup lang="ts">
import type { CasinoBetResult, CasinoFeedItem, CasinoStats } from '@/models/casino'
import { CircleDot, Dices, Loader2, RotateCw, TrendingDown, TrendingUp } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { casinoService } from '@/services/casino'
import { handleError } from '@/services/errorService'

const stats = ref<CasinoStats | null>(null)
const feed = ref<CasinoFeedItem[]>([])
const isLoading = ref(true)
const isPlaying = ref(false)
const lastResult = ref<CasinoBetResult | null>(null)
const showResult = ref(false)
const activeGame = ref<'coin' | 'dice' | 'wheel'>('coin')

const betAmount = ref(50)
const quickBets = [10, 25, 50, 100, 200]

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

const diceWinChance = computed(() => {
  return diceDirection.value === 'over'
    ? 100 - diceTarget.value
    : diceTarget.value
})

const wheelSegments = [
  { multiplier: 0, color: '#dc2626', label: 'x0' },
  { multiplier: 0.5, color: '#f97316', label: 'x0.5' },
  { multiplier: 1, color: '#eab308', label: 'x1' },
  { multiplier: 0, color: '#dc2626', label: 'x0' },
  { multiplier: 1.5, color: '#22c55e', label: 'x1.5' },
  { multiplier: 0.5, color: '#f97316', label: 'x0.5' },
  { multiplier: 3, color: '#a855f7', label: 'x3' },
  { multiplier: 0, color: '#dc2626', label: 'x0' },
  { multiplier: 1, color: '#eab308', label: 'x1' },
  { multiplier: 0.5, color: '#f97316', label: 'x0.5' },
  { multiplier: 2, color: '#3b82f6', label: 'x2' },
  { multiplier: 1.5, color: '#22c55e', label: 'x1.5' },
]

const coinFlipping = ref(false)
const coinResult = ref<'heads' | 'tails' | null>(null)
const coinRotation = ref(0)
const coinShowResult = ref(false)

const diceRolling = ref(false)
const diceResultValue = ref<number | null>(null)
const diceShowResult = ref(false)
const diceResultWin = ref(false)
const diceResultTarget = ref(0)
const diceResultDirection = ref<'over' | 'under'>('over')

// Two vertical reels (0-9 each): tens and units
// 10 panels in a ring, each 36deg apart
const diceTensAngle = ref(0)
const diceUnitsAngle = ref(0)

const wheelRotation = ref(0)
const isWheelSpinning = ref(false)

async function fetchData() {
  isLoading.value = true
  try {
    const [s, f] = await Promise.all([
      casinoService.getStats(),
      casinoService.getFeed(),
    ])
    stats.value = s
    feed.value = f ?? []
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

function delay(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

const canPlay = computed(() => betAmount.value >= 10 && betAmount.value <= 200 && !isPlaying.value)

async function playGame(action: () => Promise<CasinoBetResult>, delayMs = 1500) {
  if (!canPlay.value)
    return
  isPlaying.value = true
  showResult.value = false
  lastResult.value = null
  try {
    const [result] = await Promise.all([action(), delay(delayMs)])
    lastResult.value = result
    showResult.value = true
    if (stats.value) {
      stats.value.balance = result.balance
    }
    casinoService.getFeed().then(f => feed.value = f ?? [])
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isPlaying.value = false
  }
}

function parseCoinOutcome(result: CasinoBetResult): 'heads' | 'tails' {
  const raw = result.result
  if (raw === 'heads' || raw === 'tails')
    return raw
  try {
    const parsed = JSON.parse(raw)
    return parsed.outcome as 'heads' | 'tails'
  }
  catch {
    return 'heads'
  }
}

function playCoinFlip(choice: 'heads' | 'tails') {
  coinFlipping.value = true
  coinResult.value = null
  coinShowResult.value = false
  // Random number of full spins (3-6)
  const fullSpins = (3 + Math.floor(Math.random() * 4)) * 360
  playGame(async () => {
    const result = await casinoService.coinFlip(betAmount.value, choice)
    const outcome = parseCoinOutcome(result)
    coinResult.value = outcome
    // Calculate target: heads = 0deg, tails = 180deg (mod 360)
    const targetFace = outcome === 'tails' ? 180 : 0
    // Normalize current position to find how far we need to go
    const currentMod = coinRotation.value % 360
    const correction = targetFace - currentMod
    coinRotation.value += fullSpins + correction
    // Wait for spin to finish
    await delay(1800)
    coinFlipping.value = false
    // Show big result face
    await delay(200)
    coinShowResult.value = true
    return result
  })
}

function rollReelTo(digit: number): number {
  // Each digit is 36deg apart. Target angle = digit * 36.
  // Add 3-5 full spins for visual effect.
  const fullSpins = (3 + Math.floor(Math.random() * 3)) * 360
  return fullSpins + digit * 36
}

function playDiceRoll() {
  diceRolling.value = true
  diceResultValue.value = null
  diceShowResult.value = false
  diceResultTarget.value = diceTarget.value
  diceResultDirection.value = diceDirection.value
  // Start reels spinning to random positions
  diceTensAngle.value = rollReelTo(Math.floor(Math.random() * 10))
  diceUnitsAngle.value = rollReelTo(Math.floor(Math.random() * 10))
  playGame(async () => {
    const result = await casinoService.diceRoll(betAmount.value, diceTarget.value, diceDirection.value)
    const raw = result.result
    if (raw.startsWith('{')) {
      try {
        const parsed = JSON.parse(raw)
        diceResultValue.value = parsed.roll
      }
      catch {
        diceResultValue.value = Math.floor(Math.random() * 100)
      }
    }
    else {
      diceResultValue.value = Number.parseInt(raw) || Math.floor(Math.random() * 100)
    }
    diceResultWin.value = result.profit > 0
    // Roll reels to correct digits
    const val = diceResultValue.value!
    diceTensAngle.value = rollReelTo(Math.floor(val / 10))
    diceUnitsAngle.value = rollReelTo(val % 10)
    await delay(2200)
    diceRolling.value = false
    await delay(300)
    diceShowResult.value = true
    return result
  }, 2800)
}

function findWheelSegmentAngle(multiplier: number): number {
  // Find all segments matching this multiplier, pick one randomly
  const matches = wheelSegments
    .map((seg, i) => ({ seg, i }))
    .filter(({ seg }) => seg.multiplier === multiplier)
  if (matches.length === 0)
    return 0
  const picked = matches[Math.floor(Math.random() * matches.length)]
  // Each segment is 30deg. Center of segment i is at (i * 30 + 15) degrees.
  // To land pointer (top = 0deg) on this segment, wheel must rotate so that
  // the segment center aligns with 0. Wheel rotation is clockwise, so:
  // targetAngle = 360 - (segmentIndex * 30 + 15)
  return 360 - (picked.i * 30 + 15)
}

function playWheel() {
  if (isWheelSpinning.value)
    return
  isWheelSpinning.value = true
  // Start with a fast initial spin (visual only, will be corrected)
  const tempRotation = wheelRotation.value + 720
  wheelRotation.value = tempRotation
  playGame(async () => {
    try {
      const result = await casinoService.wheelSpin(betAmount.value)
      // Calculate correct landing angle based on server result
      const targetAngle = findWheelSegmentAngle(result.multiplier)
      const fullSpins = (3 + Math.floor(Math.random() * 2)) * 360
      wheelRotation.value = tempRotation + fullSpins + targetAngle - (tempRotation % 360)
      return result
    }
    catch (error) {
      isWheelSpinning.value = false
      throw error
    }
    finally {
      setTimeout(() => {
        isWheelSpinning.value = false
      }, 4000)
    }
  }, 4000)
}

function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  return d.toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function gameLabel(game: string) {
  const labels: Record<string, string> = {
    coin_flip: 'Монетка',
    dice_roll: 'Кости',
    wheel: 'Колесо',
  }
  return labels[game] ?? game
}

function formatResult(result: string, game: string) {
  // Handle old JSON-format results from before migration
  if (result.startsWith('{')) {
    try {
      const parsed = JSON.parse(result)
      if (game === 'coin_flip') {
        const choiceLabel = parsed.choice === 'heads' ? 'Орёл' : 'Решка'
        const outcomeLabel = parsed.outcome === 'heads' ? 'Орёл' : 'Решка'
        return `${choiceLabel} → ${outcomeLabel}`
      }
      if (game === 'dice_roll') {
        return `${parsed.target} ${parsed.direction === 'over' ? '↑' : '↓'} → ${parsed.roll}`
      }
      return parsed.outcome ?? parsed.result ?? result
    }
    catch {
      return result
    }
  }
  // New string format results
  if (game === 'coin_flip') {
    return result === 'heads' ? 'Орёл' : result === 'tails' ? 'Решка' : result
  }
  return result
}

function gameIcon(game: string) {
  const icons: Record<string, string> = {
    coin_flip: 'coin',
    dice_roll: 'dice',
    wheel: 'wheel',
  }
  return icons[game] ?? 'coin'
}

onMounted(() => fetchData())
</script>

<template>
  <div class="casino-page">
    <!-- Atmospheric background -->
    <div class="casino-bg" />

    <div class="container mx-auto px-4 py-6 md:py-8 relative z-10">
      <!-- Header with balance -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="casino-title">
            Мини-игры
          </h1>
          <p class="text-sm text-muted-foreground mt-1">
            Испытай удачу
          </p>
        </div>
        <div
          v-if="stats"
          class="balance-chip"
        >
          <div class="balance-glow" />
          <span class="balance-label">Баланс</span>
          <span class="balance-value">{{ stats.balance }}</span>
          <span class="balance-currency">б.</span>
        </div>
      </div>

      <div
        v-if="isLoading"
        class="flex justify-center py-20"
      >
        <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
      </div>

      <template v-else>
        <!-- Result overlay -->
        <Transition
          enter-active-class="result-enter-active"
          enter-from-class="result-enter-from"
          enter-to-class="result-enter-to"
          leave-active-class="result-leave-active"
          leave-to-class="result-leave-to"
        >
          <div
            v-if="showResult && lastResult"
            class="result-banner"
            :class="{
              'result-win': lastResult.profit > 0,
              'result-lose': lastResult.profit < 0,
              'result-draw': lastResult.profit === 0,
            }"
            @click="showResult = false"
          >
            <div class="result-amount">
              <span v-if="lastResult.profit > 0">+{{ lastResult.payout }}</span>
              <span v-else-if="lastResult.profit === 0">0</span>
              <span v-else>{{ lastResult.profit }}</span>
              <span class="result-currency">б.</span>
            </div>
            <div class="result-details">
              {{ gameLabel(lastResult.game) }}: {{ formatResult(lastResult.result, lastResult.game) }}
              <span v-if="lastResult.multiplier > 0" class="result-multi">x{{ lastResult.multiplier }}</span>
            </div>
          </div>
        </Transition>

        <!-- Bet controls -->
        <div class="bet-strip">
          <span class="bet-label">Ставка</span>
          <div class="bet-chips">
            <button
              v-for="q in quickBets"
              :key="q"
              class="bet-chip"
              :class="{ 'bet-chip-active': betAmount === q }"
              @click="betAmount = q"
            >
              {{ q }}
            </button>
          </div>
          <div class="bet-input-wrap">
            <input
              v-model.number="betAmount"
              type="number"
              :min="10"
              :max="200"
              class="bet-input"
            >
            <span class="bet-input-suffix">б.</span>
          </div>
        </div>

        <!-- Game selector (mobile) -->
        <div class="game-nav">
          <button
            class="game-nav-btn"
            :class="{ active: activeGame === 'coin' }"
            @click="activeGame = 'coin'"
          >
            <CircleDot class="h-4 w-4" />
            Монетка
          </button>
          <button
            class="game-nav-btn"
            :class="{ active: activeGame === 'dice' }"
            @click="activeGame = 'dice'"
          >
            <Dices class="h-4 w-4" />
            Кости
          </button>
          <button
            class="game-nav-btn"
            :class="{ active: activeGame === 'wheel' }"
            @click="activeGame = 'wheel'"
          >
            <RotateCw class="h-4 w-4" />
            Колесо
          </button>
        </div>

        <!-- Games grid -->
        <div class="games-grid">
          <!-- Coin Flip -->
          <div
            class="game-card"
            :class="{ 'game-card-hidden': activeGame !== 'coin' }"
          >
            <div class="game-card-glow game-card-glow-amber" />
            <div class="game-header">
              <div class="game-icon game-icon-amber">
                <CircleDot class="h-5 w-5" />
              </div>
              <div>
                <h3 class="game-title">
                  Монетка
                </h3>
                <p class="game-subtitle">
                  Угадай сторону — x1.9
                </p>
              </div>
            </div>

            <div class="coin-visual">
              <div
                v-if="!coinShowResult"
                class="coin-scene"
              >
                <div
                  class="coin-3d"
                  :style="{ transform: `rotateY(${coinRotation}deg) rotateX(8deg)` }"
                >
                  <div class="coin-face-front">
                    <div class="coin-inner-ring" />
                    <span class="coin-symbol">O</span>
                    <span class="coin-label-text">ОРЁЛ</span>
                  </div>
                  <div class="coin-face-back">
                    <div class="coin-inner-ring" />
                    <span class="coin-symbol">P</span>
                    <span class="coin-label-text">РЕШКА</span>
                  </div>
                  <div class="coin-edge" />
                </div>
                <div class="coin-shadow" :class="{ 'coin-shadow-flip': coinFlipping }" />
              </div>

              <!-- Big result reveal -->
              <Transition
                enter-active-class="coin-reveal-enter-active"
                enter-from-class="coin-reveal-enter-from"
                leave-active-class="coin-reveal-leave-active"
                leave-to-class="coin-reveal-leave-to"
              >
                <div
                  v-if="coinShowResult && coinResult"
                  class="coin-result-reveal"
                  :class="coinResult === 'heads' ? 'coin-result-heads' : 'coin-result-tails'"
                  @click="coinShowResult = false"
                >
                  <div class="coin-result-face">
                    <div class="coin-result-ring" />
                    <span class="coin-result-symbol">{{ coinResult === 'heads' ? 'O' : 'P' }}</span>
                    <span class="coin-result-label">{{ coinResult === 'heads' ? 'ОРЁЛ' : 'РЕШКА' }}</span>
                  </div>
                </div>
              </Transition>
            </div>

            <div class="coin-actions">
              <button
                class="coin-btn coin-btn-heads"
                :disabled="!canPlay"
                @click="playCoinFlip('heads')"
              >
                <Loader2
                  v-if="isPlaying && activeGame === 'coin'"
                  class="h-4 w-4 animate-spin"
                />
                <span v-else class="coin-btn-icon">O</span>
                Орёл
              </button>
              <button
                class="coin-btn coin-btn-tails"
                :disabled="!canPlay"
                @click="playCoinFlip('tails')"
              >
                <Loader2
                  v-if="isPlaying && activeGame === 'coin'"
                  class="h-4 w-4 animate-spin"
                />
                <span v-else class="coin-btn-icon">P</span>
                Решка
              </button>
            </div>
          </div>

          <!-- Dice Roll -->
          <div
            class="game-card"
            :class="{ 'game-card-hidden': activeGame !== 'dice' }"
          >
            <div class="game-card-glow game-card-glow-blue" />
            <div class="game-header">
              <div class="game-icon game-icon-blue">
                <Dices class="h-5 w-5" />
              </div>
              <div>
                <h3 class="game-title">
                  Кости
                </h3>
                <p class="game-subtitle">
                  Выбери число и направление
                </p>
              </div>
            </div>

            <div class="dice-visual">
              <!-- Two slot reels (tens + units) -->
              <div
                v-if="!diceShowResult"
                class="dice-reels-row"
              >
                <div class="dice-reel-window">
                  <div
                    class="dice-reel"
                    :style="{ transform: `rotateX(-${diceTensAngle}deg)` }"
                  >
                    <div
                      v-for="d in 10"
                      :key="`t${d}`"
                      class="dice-reel-panel"
                      :style="{ transform: `rotateX(${(d - 1) * 36}deg) translateZ(50px)` }"
                    >
                      {{ d - 1 }}
                    </div>
                  </div>
                </div>
                <div class="dice-reel-window">
                  <div
                    class="dice-reel"
                    :style="{ transform: `rotateX(-${diceUnitsAngle}deg)` }"
                  >
                    <div
                      v-for="d in 10"
                      :key="`u${d}`"
                      class="dice-reel-panel"
                      :style="{ transform: `rotateX(${(d - 1) * 36}deg) translateZ(50px)` }"
                    >
                      {{ d - 1 }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Result reveal -->
              <Transition
                enter-active-class="dice-reveal-enter-active"
                enter-from-class="dice-reveal-enter-from"
                leave-active-class="dice-reveal-leave-active"
                leave-to-class="dice-reveal-leave-to"
              >
                <div
                  v-if="diceShowResult && diceResultValue !== null"
                  class="dice-result-reveal"
                  :class="diceResultWin ? 'dice-result-win' : 'dice-result-lose'"
                  @click="diceShowResult = false"
                >
                  <div class="dice-result-big-number">
                    {{ diceResultValue }}
                  </div>
                  <div class="dice-result-condition">
                    {{ diceResultDirection === 'over' ? '>' : '<' }} {{ diceResultTarget }}
                  </div>
                </div>
              </Transition>
            </div>

            <div class="dice-controls">
              <div class="dice-direction">
                <button
                  class="dir-btn"
                  :class="{ 'dir-btn-active': diceDirection === 'under' }"
                  @click="diceDirection = 'under'"
                >
                  <TrendingDown class="h-3.5 w-3.5" />
                  Меньше
                </button>
                <button
                  class="dir-btn"
                  :class="{ 'dir-btn-active': diceDirection === 'over' }"
                  @click="diceDirection = 'over'"
                >
                  <TrendingUp class="h-3.5 w-3.5" />
                  Больше
                </button>
              </div>

              <div class="dice-slider-wrap">
                <div class="dice-slider-track">
                  <div
                    class="dice-slider-fill"
                    :style="{
                      left: diceDirection === 'under' ? '0%' : `${diceTarget}%`,
                      width: diceDirection === 'under' ? `${diceTarget}%` : `${100 - diceTarget}%`,
                    }"
                  />
                  <div
                    class="dice-slider-marker"
                    :style="{ left: `${diceTarget}%` }"
                  >
                    <span class="dice-slider-value">{{ diceTarget }}</span>
                  </div>
                </div>
                <input
                  v-model.number="diceTarget"
                  type="range"
                  min="2"
                  max="98"
                  class="dice-range"
                >
              </div>

              <div class="dice-stats">
                <div class="dice-stat">
                  <span class="dice-stat-label">Шанс</span>
                  <span class="dice-stat-value">{{ diceWinChance }}%</span>
                </div>
                <div class="dice-stat">
                  <span class="dice-stat-label">Множитель</span>
                  <span class="dice-stat-value dice-stat-highlight">x{{ diceMultiplier }}</span>
                </div>
              </div>
            </div>

            <Button
              size="lg"
              class="w-full game-play-btn"
              :disabled="!canPlay"
              @click="playDiceRoll"
            >
              <Loader2
                v-if="isPlaying && activeGame === 'dice'"
                class="h-4 w-4 animate-spin"
              />
              Бросить кости
            </Button>
          </div>

          <!-- Wheel -->
          <div
            class="game-card"
            :class="{ 'game-card-hidden': activeGame !== 'wheel' }"
          >
            <div class="game-card-glow game-card-glow-purple" />
            <div class="game-header">
              <div class="game-icon game-icon-purple">
                <RotateCw class="h-5 w-5" />
              </div>
              <div>
                <h3 class="game-title">
                  Колесо фортуны
                </h3>
                <p class="game-subtitle">
                  Множитель от x0 до x3
                </p>
              </div>
            </div>

            <div class="wheel-visual">
              <div class="wheel-pointer" />
              <div
                class="wheel-body"
                :style="{ transform: `rotate(${wheelRotation}deg)` }"
                :class="{ 'wheel-spinning': isWheelSpinning }"
              >
                <svg
                  viewBox="0 0 200 200"
                  class="wheel-svg"
                >
                  <g
                    v-for="(seg, i) in wheelSegments"
                    :key="i"
                  >
                    <path
                      :d="`M100,100 L${100 + 95 * Math.cos((i * 30 - 90) * Math.PI / 180)},${100 + 95 * Math.sin((i * 30 - 90) * Math.PI / 180)} A95,95 0 0,1 ${100 + 95 * Math.cos(((i + 1) * 30 - 90) * Math.PI / 180)},${100 + 95 * Math.sin(((i + 1) * 30 - 90) * Math.PI / 180)} Z`"
                      :fill="seg.color"
                      :opacity="0.8"
                      stroke="rgba(0,0,0,0.3)"
                      stroke-width="0.5"
                    />
                    <text
                      :x="100 + 65 * Math.cos(((i + 0.5) * 30 - 90) * Math.PI / 180)"
                      :y="100 + 65 * Math.sin(((i + 0.5) * 30 - 90) * Math.PI / 180)"
                      text-anchor="middle"
                      dominant-baseline="central"
                      fill="white"
                      font-size="10"
                      font-weight="700"
                      :transform="`rotate(${(i + 0.5) * 30}, ${100 + 65 * Math.cos(((i + 0.5) * 30 - 90) * Math.PI / 180)}, ${100 + 65 * Math.sin(((i + 0.5) * 30 - 90) * Math.PI / 180)})`"
                    >
                      {{ seg.label }}
                    </text>
                  </g>
                  <circle
                    cx="100"
                    cy="100"
                    r="20"
                    fill="hsl(var(--card))"
                    stroke="rgba(255,255,255,0.1)"
                    stroke-width="1"
                  />
                </svg>
              </div>
            </div>

            <Button
              size="lg"
              class="w-full game-play-btn"
              :disabled="!canPlay || isWheelSpinning"
              @click="playWheel"
            >
              <RotateCw
                class="h-4 w-4"
                :class="{ 'animate-spin': isWheelSpinning }"
              />
              Крутить колесо
            </Button>
          </div>
        </div>

        <!-- History -->
        <div
          v-if="feed.length > 0"
          class="history-section"
        >
          <h3 class="history-title">
            Последние ставки
          </h3>
          <div class="history-list">
            <div
              v-for="bet in feed"
              :key="bet.id"
              class="history-item"
              :class="{
                'history-win': bet.profit > 0,
                'history-lose': bet.profit < 0,
              }"
            >
              <div class="history-icon-wrap">
                <CircleDot
                  v-if="gameIcon(bet.game) === 'coin'"
                  class="h-3.5 w-3.5"
                />
                <Dices
                  v-else-if="gameIcon(bet.game) === 'dice'"
                  class="h-3.5 w-3.5"
                />
                <RotateCw
                  v-else
                  class="h-3.5 w-3.5"
                />
              </div>
              <div class="history-info">
                <span class="history-game">
                  {{ bet.memberUsername ? `@${bet.memberUsername}` : bet.memberFirstName }}
                </span>
                <span class="history-date">{{ gameLabel(bet.game) }} · {{ formatDate(bet.createdAt) }}</span>
              </div>
              <div class="history-bet">
                {{ bet.betAmount }} б.
              </div>
              <div class="history-result-col">
                {{ formatResult(bet.result, bet.game) }}
              </div>
              <div class="history-profit">
                <Badge
                  :variant="bet.profit > 0 ? 'default' : bet.profit === 0 ? 'secondary' : 'destructive'"
                  class="history-badge"
                >
                  {{ bet.profit > 0 ? '+' : '' }}{{ bet.profit }}
                </Badge>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
/* ======= ATMOSPHERE ======= */
.casino-page {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
}

.casino-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
  background:
    radial-gradient(ellipse at 20% 0%, hsl(151 60% 54% / 0.06) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 100%, hsl(45 80% 60% / 0.04) 0%, transparent 50%),
    radial-gradient(ellipse at 50% 50%, hsl(var(--background)) 0%, hsl(var(--background)) 100%);
  pointer-events: none;
}

.casino-bg::after {
  content: '';
  position: absolute;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.03'/%3E%3C/svg%3E");
  background-size: 200px;
  opacity: 0.5;
}

/* ======= TITLE ======= */
.casino-title {
  font-size: 1.75rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  background: linear-gradient(135deg, hsl(var(--foreground)), hsl(var(--muted-foreground)));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* ======= BALANCE ======= */
.balance-chip {
  position: relative;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  border-radius: 999px;
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
  overflow: hidden;
}

.balance-glow {
  position: absolute;
  inset: -1px;
  border-radius: inherit;
  background: linear-gradient(135deg, hsl(45 80% 60% / 0.15), hsl(151 60% 54% / 0.1));
  z-index: 0;
  pointer-events: none;
}

.balance-label {
  position: relative;
  font-size: 0.75rem;
  color: hsl(var(--muted-foreground));
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 500;
}

.balance-value {
  position: relative;
  font-size: 1.25rem;
  font-weight: 800;
  font-variant-numeric: tabular-nums;
  color: hsl(45 80% 65%);
}

.balance-currency {
  position: relative;
  font-size: 0.8rem;
  color: hsl(var(--muted-foreground));
  font-weight: 500;
}

/* ======= RESULT BANNER ======= */
.result-banner {
  position: relative;
  margin-bottom: 1.5rem;
  padding: 1.25rem 1.5rem;
  border-radius: 1rem;
  text-align: center;
  cursor: pointer;
  border: 1px solid;
  overflow: hidden;
}

.result-banner::before {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0.1;
  background: radial-gradient(ellipse at center, currentColor, transparent 70%);
}

.result-win {
  border-color: hsl(142 70% 45% / 0.4);
  background: hsl(142 70% 45% / 0.08);
  color: hsl(142 70% 55%);
}

.result-lose {
  border-color: hsl(0 70% 50% / 0.4);
  background: hsl(0 70% 50% / 0.08);
  color: hsl(0 70% 60%);
}

.result-draw {
  border-color: hsl(var(--border));
  background: hsl(var(--card));
  color: hsl(var(--muted-foreground));
}

.result-amount {
  font-size: 2rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  line-height: 1;
  margin-bottom: 4px;
}

.result-currency {
  font-size: 1.25rem;
  opacity: 0.7;
  margin-left: 2px;
}

.result-details {
  font-size: 0.8rem;
  opacity: 0.7;
}

.result-multi {
  font-weight: 700;
  opacity: 1;
}

.result-enter-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.result-leave-active {
  transition: all 0.2s ease-in;
}
.result-enter-from {
  opacity: 0;
  transform: translateY(-12px) scale(0.96);
}
.result-enter-to {
  opacity: 1;
  transform: translateY(0) scale(1);
}
.result-leave-to {
  opacity: 0;
  transform: scale(0.98);
}

/* ======= BET STRIP ======= */
.bet-strip {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-bottom: 1.5rem;
  border-radius: 1rem;
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
  flex-wrap: wrap;
}

.bet-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: hsl(var(--muted-foreground));
  white-space: nowrap;
}

.bet-chips {
  display: flex;
  gap: 6px;
}

.bet-chip {
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  background: hsl(var(--secondary));
  color: hsl(var(--secondary-foreground));
  border: 1px solid transparent;
  transition: all 0.15s ease;
  cursor: pointer;
}

.bet-chip:hover {
  background: hsl(var(--accent) / 0.15);
  border-color: hsl(var(--accent) / 0.3);
}

.bet-chip-active {
  background: hsl(var(--accent) / 0.2);
  border-color: hsl(var(--accent) / 0.5);
  color: hsl(var(--accent));
  box-shadow: 0 0 12px hsl(var(--accent) / 0.15);
}

.bet-input-wrap {
  position: relative;
  margin-left: auto;
}

.bet-input {
  width: 80px;
  padding: 6px 28px 6px 12px;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  background: hsl(var(--background));
  border: 1px solid hsl(var(--border));
  color: hsl(var(--foreground));
  outline: none;
  -moz-appearance: textfield;
}

.bet-input::-webkit-outer-spin-button,
.bet-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.bet-input:focus {
  border-color: hsl(var(--accent) / 0.5);
  box-shadow: 0 0 0 2px hsl(var(--accent) / 0.1);
}

.bet-input-suffix {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.75rem;
  color: hsl(var(--muted-foreground));
  pointer-events: none;
}

/* ======= GAME NAV (mobile) ======= */
.game-nav {
  display: flex;
  gap: 4px;
  margin-bottom: 1rem;
  padding: 4px;
  border-radius: 0.75rem;
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
}

@media (min-width: 768px) {
  .game-nav {
    display: none;
  }
}

.game-nav-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 0.5rem;
  font-size: 0.8rem;
  font-weight: 500;
  color: hsl(var(--muted-foreground));
  transition: all 0.15s ease;
}

.game-nav-btn.active {
  background: hsl(var(--secondary));
  color: hsl(var(--foreground));
  font-weight: 600;
}

/* ======= GAMES GRID ======= */
.games-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 2rem;
}

@media (max-width: 768px) {
  .games-grid {
    grid-template-columns: 1fr;
  }
  .game-card-hidden {
    display: none;
  }
}

/* ======= GAME CARD ======= */
.game-card {
  position: relative;
  padding: 1.25rem;
  border-radius: 1.25rem;
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.game-card-glow {
  position: absolute;
  top: -40px;
  right: -40px;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  opacity: 0.07;
  pointer-events: none;
  filter: blur(40px);
}

.game-card-glow-amber { background: hsl(45 90% 55%); }
.game-card-glow-blue { background: hsl(217 90% 60%); }
.game-card-glow-purple { background: hsl(270 80% 60%); }

.game-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 1.25rem;
}

.game-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  flex-shrink: 0;
}

.game-icon-amber {
  background: hsl(45 80% 55% / 0.15);
  color: hsl(45 80% 60%);
}
.game-icon-blue {
  background: hsl(217 80% 55% / 0.15);
  color: hsl(217 80% 65%);
}
.game-icon-purple {
  background: hsl(270 70% 55% / 0.15);
  color: hsl(270 70% 65%);
}

.game-title {
  font-size: 1rem;
  font-weight: 700;
  line-height: 1.2;
}

.game-subtitle {
  font-size: 0.75rem;
  color: hsl(var(--muted-foreground));
  margin-top: 1px;
}

.game-play-btn {
  margin-top: auto;
}

/* ======= COIN FLIP 3D ======= */
.coin-visual {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 1rem 0;
  position: relative;
  min-height: 130px;
}

.coin-scene {
  position: relative;
  width: 96px;
  height: 96px;
  perspective: 600px;
  transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
}

.coin-3d {
  width: 100%;
  height: 100%;
  position: relative;
  transform-style: preserve-3d;
  transition: transform 1.8s cubic-bezier(0.12, 0.8, 0.2, 1);
}

.coin-face-front,
.coin-face-back {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  backface-visibility: hidden;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
}

.coin-face-front {
  background: linear-gradient(145deg, hsl(45 75% 58%), hsl(38 85% 42%));
  box-shadow:
    0 4px 20px hsl(45 80% 50% / 0.3),
    inset 0 2px 6px hsl(45 90% 80% / 0.4),
    inset 0 -3px 6px hsl(35 80% 30% / 0.3);
}

.coin-face-front::after {
  content: '';
  position: absolute;
  top: 8%;
  left: 15%;
  width: 35%;
  height: 20%;
  border-radius: 50%;
  background: linear-gradient(180deg, hsl(45 80% 80% / 0.4), transparent);
  pointer-events: none;
}

.coin-face-back {
  background: linear-gradient(145deg, hsl(220 15% 50%), hsl(220 12% 38%));
  box-shadow:
    0 4px 20px hsl(220 20% 40% / 0.3),
    inset 0 2px 6px hsl(220 15% 70% / 0.3),
    inset 0 -3px 6px hsl(220 20% 20% / 0.4);
  transform: rotateY(180deg);
}

.coin-face-back::after {
  content: '';
  position: absolute;
  top: 8%;
  left: 15%;
  width: 35%;
  height: 20%;
  border-radius: 50%;
  background: linear-gradient(180deg, hsl(220 20% 70% / 0.3), transparent);
  pointer-events: none;
}

.coin-inner-ring {
  position: absolute;
  inset: 5px;
  border-radius: 50%;
  border: 2px solid currentColor;
  opacity: 0.2;
  pointer-events: none;
}

.coin-symbol {
  font-size: 1.75rem;
  font-weight: 900;
  line-height: 1;
}

.coin-face-front .coin-symbol {
  color: hsl(35 50% 18%);
  text-shadow: 0 1px 0 hsl(45 60% 70% / 0.5);
}

.coin-face-back .coin-symbol {
  color: hsl(220 10% 80%);
  text-shadow: 0 1px 0 hsl(220 20% 30% / 0.5);
}

.coin-label-text {
  font-size: 0.5rem;
  font-weight: 700;
  letter-spacing: 0.15em;
  text-transform: uppercase;
}

.coin-face-front .coin-label-text {
  color: hsl(35 40% 28%);
}

.coin-face-back .coin-label-text {
  color: hsl(220 10% 70%);
}

.coin-edge {
  position: absolute;
  top: 50%;
  left: 0;
  width: 100%;
  height: 6px;
  transform: translateY(-50%) translateZ(-3px);
  background: linear-gradient(90deg, hsl(38 60% 35%), hsl(45 70% 50%), hsl(38 60% 35%));
  border-radius: 50%;
}

.coin-shadow {
  position: absolute;
  bottom: -12px;
  left: 50%;
  transform: translateX(-50%);
  width: 70px;
  height: 14px;
  border-radius: 50%;
  background: radial-gradient(ellipse, hsl(0 0% 0% / 0.15), transparent 70%);
  transition: all 0.3s ease;
}

.coin-shadow-flip {
  animation: coinShadowBounce 1.8s cubic-bezier(0.12, 0.8, 0.2, 1);
}

@keyframes coinShadowBounce {
  0% { transform: translateX(-50%) scale(1); opacity: 1; }
  25% { transform: translateX(-50%) scale(0.4); opacity: 0.2; }
  50% { transform: translateX(-50%) scale(0.35); opacity: 0.15; }
  80% { transform: translateX(-50%) scale(0.85); opacity: 0.9; }
  100% { transform: translateX(-50%) scale(1); opacity: 1; }
}

/* ======= COIN RESULT REVEAL ======= */
.coin-result-reveal {
  display: flex;
  justify-content: center;
  cursor: pointer;
}

.coin-result-face {
  position: relative;
  width: 88px;
  height: 88px;
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
}

.coin-result-heads .coin-result-face {
  background: linear-gradient(145deg, hsl(45 75% 58%), hsl(38 85% 42%));
  box-shadow:
    0 6px 30px hsl(45 80% 50% / 0.4),
    0 0 40px hsl(45 80% 50% / 0.15),
    inset 0 2px 6px hsl(45 90% 80% / 0.4),
    inset 0 -3px 6px hsl(35 80% 30% / 0.3);
}

.coin-result-tails .coin-result-face {
  background: linear-gradient(145deg, hsl(220 15% 50%), hsl(220 12% 38%));
  box-shadow:
    0 6px 30px hsl(220 20% 40% / 0.4),
    0 0 40px hsl(220 30% 50% / 0.15),
    inset 0 2px 6px hsl(220 15% 70% / 0.3),
    inset 0 -3px 6px hsl(220 20% 20% / 0.4);
}

.coin-result-ring {
  position: absolute;
  inset: 5px;
  border-radius: 50%;
  border: 2px solid currentColor;
  opacity: 0.2;
}

.coin-result-symbol {
  font-size: 2rem;
  font-weight: 900;
  line-height: 1;
}

.coin-result-heads .coin-result-symbol {
  color: hsl(35 50% 18%);
  text-shadow: 0 1px 0 hsl(45 60% 70% / 0.5);
}

.coin-result-tails .coin-result-symbol {
  color: hsl(220 10% 80%);
  text-shadow: 0 1px 0 hsl(220 20% 30% / 0.5);
}

.coin-result-label {
  font-size: 0.55rem;
  font-weight: 700;
  letter-spacing: 0.15em;
}

.coin-result-heads .coin-result-label {
  color: hsl(35 40% 28%);
}

.coin-result-tails .coin-result-label {
  color: hsl(220 10% 70%);
}

.coin-reveal-enter-active {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.coin-reveal-enter-from {
  opacity: 0;
  transform: scale(0.3) rotateY(90deg);
}
.coin-reveal-leave-active {
  transition: all 0.25s ease-in;
}
.coin-reveal-leave-to {
  opacity: 0;
  transform: scale(0.8);
}

.coin-actions {
  display: flex;
  gap: 8px;
  margin-top: auto;
}

.coin-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 0.75rem;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.coin-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.coin-btn-icon {
  font-weight: 800;
  font-size: 1rem;
}

.coin-btn-heads {
  background: hsl(45 70% 50% / 0.15);
  color: hsl(45 80% 60%);
  border: 1px solid hsl(45 70% 50% / 0.25);
}

.coin-btn-heads:hover:not(:disabled) {
  background: hsl(45 70% 50% / 0.25);
  border-color: hsl(45 70% 50% / 0.4);
  box-shadow: 0 0 20px hsl(45 80% 50% / 0.1);
}

.coin-btn-tails {
  background: hsl(var(--secondary));
  color: hsl(var(--secondary-foreground));
  border: 1px solid hsl(var(--border));
}

.coin-btn-tails:hover:not(:disabled) {
  background: hsl(var(--secondary));
  border-color: hsl(var(--muted-foreground) / 0.3);
}

/* ======= DICE SLOT REELS ======= */
.dice-visual {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0.75rem 0 0.25rem;
  position: relative;
  min-height: 110px;
}

.dice-reels-row {
  display: flex;
  gap: 6px;
  align-items: center;
  justify-content: center;
}

.dice-reel-window {
  width: 48px;
  height: 56px;
  perspective: 300px;
  overflow: hidden;
  border-radius: 10px;
  background: linear-gradient(145deg, hsl(217 55% 22%), hsl(220 50% 16%));
  border: 1px solid hsl(217 40% 32% / 0.6);
  box-shadow:
    0 4px 12px hsl(217 60% 15% / 0.4),
    inset 0 1px 3px hsl(217 60% 40% / 0.1),
    inset 0 -1px 3px hsl(220 50% 8% / 0.4);
  position: relative;
}

.dice-reel-window::before,
.dice-reel-window::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  height: 14px;
  z-index: 2;
  pointer-events: none;
}

.dice-reel-window::before {
  top: 0;
  background: linear-gradient(to bottom, hsl(220 50% 16%), transparent);
}

.dice-reel-window::after {
  bottom: 0;
  background: linear-gradient(to top, hsl(220 50% 16%), transparent);
}

.dice-reel {
  width: 100%;
  height: 100%;
  position: relative;
  transform-style: preserve-3d;
  transition: transform 2s cubic-bezier(0.12, 0.8, 0.2, 1);
}

.dice-reel-panel {
  position: absolute;
  width: 48px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  font-weight: 900;
  font-variant-numeric: tabular-nums;
  color: hsl(217 80% 75%);
  text-shadow: 0 0 12px hsl(217 80% 60% / 0.3);
  backface-visibility: hidden;
}

/* ======= DICE RESULT REVEAL ======= */
.dice-result-reveal {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
}

.dice-result-win .dice-result-big-number {
  color: hsl(142 70% 55%);
  text-shadow: 0 0 30px hsl(142 70% 50% / 0.5), 0 0 60px hsl(142 70% 50% / 0.15);
}

.dice-result-lose .dice-result-big-number {
  color: hsl(0 70% 60%);
  text-shadow: 0 0 30px hsl(0 70% 55% / 0.5), 0 0 60px hsl(0 70% 55% / 0.15);
}

.dice-result-big-number {
  font-size: 2.75rem;
  font-weight: 900;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.03em;
  line-height: 1;
}

.dice-result-condition {
  font-size: 0.75rem;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: hsl(var(--muted-foreground));
  letter-spacing: 0.02em;
}

.dice-result-win .dice-result-condition {
  color: hsl(142 50% 50% / 0.7);
}

.dice-result-lose .dice-result-condition {
  color: hsl(0 50% 55% / 0.7);
}

.dice-reveal-enter-active {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.dice-reveal-enter-from {
  opacity: 0;
  transform: scale(0.3) translateY(15px);
}
.dice-reveal-leave-active {
  transition: all 0.25s ease-in;
}
.dice-reveal-leave-to {
  opacity: 0;
  transform: scale(0.8);
}

/* ======= DICE ======= */
.dice-controls {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 16px;
}

.dice-direction {
  display: flex;
  gap: 6px;
}

.dir-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  padding: 7px 12px;
  border-radius: 0.5rem;
  font-size: 0.8rem;
  font-weight: 500;
  background: hsl(var(--secondary));
  color: hsl(var(--muted-foreground));
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.15s ease;
}

.dir-btn-active {
  background: hsl(217 80% 55% / 0.15);
  color: hsl(217 80% 65%);
  border-color: hsl(217 80% 55% / 0.3);
}

/* Dice slider */
.dice-slider-wrap {
  position: relative;
  padding: 16px 0 8px;
}

.dice-slider-track {
  position: relative;
  height: 6px;
  border-radius: 3px;
  background: hsl(var(--secondary));
  overflow: visible;
}

.dice-slider-fill {
  position: absolute;
  top: 0;
  height: 100%;
  border-radius: 3px;
  background: linear-gradient(90deg, hsl(142 70% 45% / 0.6), hsl(151 60% 54% / 0.8));
  transition: left 0.1s ease, width 0.1s ease;
}

.dice-slider-marker {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  z-index: 2;
  transition: left 0.1s ease;
}

.dice-slider-value {
  display: block;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  background: hsl(217 80% 55%);
  color: white;
  white-space: nowrap;
  box-shadow: 0 2px 8px hsl(217 80% 55% / 0.3);
  transform: translateY(-14px);
}

.dice-range {
  position: absolute;
  top: 12px;
  left: 0;
  width: 100%;
  height: 20px;
  opacity: 0;
  cursor: pointer;
  z-index: 3;
  margin: 0;
}

.dice-stats {
  display: flex;
  gap: 12px;
}

.dice-stat {
  flex: 1;
  padding: 8px 12px;
  border-radius: 0.5rem;
  background: hsl(var(--background));
  text-align: center;
}

.dice-stat-label {
  display: block;
  font-size: 0.65rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: hsl(var(--muted-foreground));
  margin-bottom: 2px;
}

.dice-stat-value {
  font-size: 1rem;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
}

.dice-stat-highlight {
  color: hsl(217 80% 65%);
}

/* ======= WHEEL ======= */
.wheel-visual {
  display: flex;
  justify-content: center;
  padding: 1rem 0;
  position: relative;
}

.wheel-pointer {
  position: absolute;
  top: 0.75rem;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-top: 14px solid hsl(var(--foreground));
  z-index: 2;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3));
}

.wheel-body {
  width: 180px;
  height: 180px;
  transition: transform 4s cubic-bezier(0.17, 0.67, 0.12, 0.99);
}

.wheel-spinning {
  transition: transform 4s cubic-bezier(0.17, 0.67, 0.12, 0.99);
}

.wheel-svg {
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 4px 12px rgba(0,0,0,0.2));
}

/* ======= HISTORY ======= */
.history-section {
  margin-top: 0.5rem;
}

.history-title {
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: hsl(var(--muted-foreground));
  margin-bottom: 0.75rem;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 0.75rem;
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
  font-size: 0.8rem;
  transition: background 0.15s ease;
}

.history-item:hover {
  background: hsl(var(--secondary));
}

.history-win {
  border-left: 3px solid hsl(142 70% 45% / 0.5);
}

.history-lose {
  border-left: 3px solid hsl(0 70% 50% / 0.3);
}

.history-icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: hsl(var(--secondary));
  color: hsl(var(--muted-foreground));
  flex-shrink: 0;
}

.history-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
}

.history-game {
  font-weight: 600;
  color: hsl(var(--foreground));
}

.history-date {
  font-size: 0.7rem;
  color: hsl(var(--muted-foreground));
}

.history-bet {
  font-variant-numeric: tabular-nums;
  color: hsl(var(--muted-foreground));
  white-space: nowrap;
}

.history-result-col {
  font-variant-numeric: tabular-nums;
  color: hsl(var(--foreground));
  white-space: nowrap;
  min-width: 40px;
  text-align: center;
}

.history-profit {
  min-width: 60px;
  text-align: right;
}

.history-badge {
  font-variant-numeric: tabular-nums;
  font-weight: 700;
  font-size: 0.75rem;
}

@media (max-width: 640px) {
  .history-bet,
  .history-result-col {
    display: none;
  }
  .bet-chips {
    order: 3;
    flex-basis: 100%;
  }
  .bet-input-wrap {
    margin-left: 0;
  }
}
</style>
