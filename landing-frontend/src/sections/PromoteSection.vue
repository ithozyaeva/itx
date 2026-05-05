<script setup lang="ts">
import type { TelegramUser } from '@/services/auth.ts'
import { defineAsyncComponent, onMounted, onUnmounted, ref } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import TelegramAuth from '@/components/TelegramAuth.vue'
import Button from '@/components/ui/UiButton.vue'
import { useCountUp } from '@/composables/useCountUp'
import { useMagneticHover } from '@/composables/useMagneticHover'
import { useToken } from '@/composables/useToken.ts'
import { useUser } from '@/composables/useUser.ts'

// HeroConstellation — декоративный canvas, скрыт до lg-брейкпойнта.
// Грузим лениво, чтобы не утяжелять main chunk и не конкурировать с LCP.
const HeroConstellation = defineAsyncComponent(() =>
  import('@/components/HeroConstellation.vue'),
)

const tgUser = useUser()
const tgToken = useToken()
const yandexMetrika = useYandexMetrika()

function setUser(user: TelegramUser, token: string) {
  tgUser.value = user
  tgToken.value = token
}

function trackPlatformClick() {
  yandexMetrika.reachGoal('platform_redirect_click', {
    location: 'promote_section',
    isAuthenticated: !!tgUser.value,
  } as any)
}

function trackHeroPricingClick() {
  yandexMetrika.reachGoal('pricing_teaser_click', {
    location: 'hero',
  } as any)
}

const phrases = [
  'менторство_',
  'вайбкодинг_',
  'ИИ-практики_',
  'нетворкинг_',
  'собеседования_',
]
const typed = ref('')
const phraseIdx = ref(0)
let timer: number | undefined
let deleting = false

function tick() {
  const current = phrases[phraseIdx.value]
  if (!deleting) {
    typed.value = current.slice(0, typed.value.length + 1)
    if (typed.value === current) {
      deleting = true
      timer = window.setTimeout(tick, 1600)
      return
    }
    timer = window.setTimeout(tick, 65 + Math.random() * 40)
  }
  else {
    typed.value = current.slice(0, typed.value.length - 1)
    if (typed.value === '') {
      deleting = false
      phraseIdx.value = (phraseIdx.value + 1) % phrases.length
      timer = window.setTimeout(tick, 180)
      return
    }
    timer = window.setTimeout(tick, 35)
  }
}

const uptime = ref('00:00:00')
let uptimeTimer: number | undefined
const started = Date.now()
function updateUptime() {
  const sec = Math.floor((Date.now() - started) / 1000)
  const h = String(Math.floor(sec / 3600)).padStart(2, '0')
  const m = String(Math.floor((sec % 3600) / 60)).padStart(2, '0')
  const s = String(sec % 60).padStart(2, '0')
  uptime.value = `${h}:${m}:${s}`
}

// Count-up stats
const stat1 = useCountUp(250, 2000)
const stat2 = useCountUp(60, 1800)
const stat3 = useCountUp(7, 1200)

// Magnetic hover on CTA wrapper
const ctaRef = ref<HTMLElement | null>(null)
const { x: magX, y: magY } = useMagneticHover(ctaRef, 0.3)

// Hero entrance animation
const heroReady = ref(false)

onMounted(() => {
  tick()
  updateUptime()
  uptimeTimer = window.setInterval(updateUptime, 1000)

  // Staggered entrance
  requestAnimationFrame(() => {
    heroReady.value = true
  })

  // Start count-ups after a delay
  setTimeout(() => {
    stat1.start()
    stat2.start()
    stat3.start()
  }, 600)
})

onUnmounted(() => {
  if (timer)
    clearTimeout(timer)
  if (uptimeTimer)
    clearInterval(uptimeTimer)
})

const companies = ['Яндекс', 'Tinkoff', 'VK', 'Ozon', 'Wildberries', 'Авито', 'X5', 'Сбер', 'Kaspersky', 'JetBrains', 'Selectel', 'Тинькофф', 'Альфа', 'МТС']
</script>

<template>
  <section
    class="relative w-full md:min-h-[calc(100svh-var(--header-height))] flex flex-col md:justify-between overflow-hidden"
  >
    <!-- Decorative grid / crosshair -->
    <div
      aria-hidden="true"
      class="pointer-events-none absolute inset-0 opacity-[0.12]"
      style="background-image: radial-gradient(circle at 80% 30%, hsl(var(--accent)) 0, transparent 55%), radial-gradient(circle at 20% 80%, #ffb547 0, transparent 45%);"
    />
    <HeroConstellation />

    <!-- Terminal status bar -->
    <div
      class="container px-6 md:px-10 pt-8 md:pt-12 relative z-10 transition-all duration-700"
      :class="heroReady ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-4'"
    >
      <div class="flex items-center justify-between font-mono text-[11px] md:text-xs tracking-[0.08em] text-foreground/50">
        <div class="flex items-center gap-3 md:gap-6">
          <div class="flex items-center gap-1.5">
            <span class="w-2 h-2 rounded-full bg-accent shadow-[0_0_8px_hsl(var(--accent))] animate-pulse" />
            <span class="text-accent">ОНЛАЙН</span>
          </div>
          <span class="flex items-center gap-1">
            <span class="text-accent/40">↑</span>
            <span class="text-foreground/70 tabular-nums">{{ uptime }}</span>
          </span>
          <span class="hidden md:inline">v4.2.1</span>
        </div>
        <div class="flex items-center gap-3 md:gap-6">
          <span class="hidden sm:inline">250+ участников</span>
          <span class="flex items-center gap-1.5">
            <span
              aria-hidden="true"
              class="signal-bars"
            >
              <span /><span /><span /><span />
            </span>
            <span>msk/ru</span>
          </span>
        </div>
      </div>
    </div>

    <!-- Main hero content -->
    <div class="container px-6 md:px-10 md:flex-1 md:flex md:items-center pt-10 pb-10 md:py-16 relative z-10">
      <div class="w-full max-w-5xl">
        <!-- Prompt line -->
        <div
          class="font-mono text-xs md:text-sm text-foreground/50 mb-5 md:mb-6 transition-all duration-700 delay-100"
          :class="heroReady ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-6'"
        >
          <div>
            <span class="text-accent">community@ithozyaeva</span>:<span class="text-term-amber">~</span>$
            <span class="text-foreground/80">./добро_пожаловать<span class="hidden md:inline"> --новый пользователь</span><span class="md:hidden text-foreground/30"> \</span></span>
          </div>
          <div class="md:hidden mt-1 text-foreground/80">
            <span class="text-foreground/30">&gt;</span> --новый пользователь
          </div>
        </div>

        <!-- Main title with glitch -->
        <div
          class="flex items-baseline gap-2 md:gap-3 mb-1 transition-all duration-700 delay-200"
          :class="heroReady ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'"
        >
          <span class="font-mono text-accent/60 text-sm md:text-base">&gt;</span>
          <h1
            class="glitch-hover font-display uppercase text-[38px] xs:text-[46px] sm:text-[68px] md:text-[96px] lg:text-[130px] leading-[0.85] tracking-tight text-accent"
            data-text="IT-ХОЗЯЕВА"
          >
            IT-ХОЗЯЕВА
          </h1>
        </div>

        <div
          class="pl-4 md:pl-6 border-l border-accent/30 mt-6 md:mt-8 max-w-3xl transition-all duration-700 delay-300"
          :class="heroReady ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-8'"
        >
          <h2 class="font-display uppercase text-[22px] sm:text-[32px] md:text-[42px] leading-[1.05] text-foreground">
            Закрытое сообщество<br>
            <span class="text-foreground/40">IT-специалистов</span>
          </h2>

          <div class="mt-5 md:mt-7 font-mono text-sm md:text-base text-foreground/75 flex flex-wrap items-center gap-x-2">
            <span class="text-accent/60">&gt;&gt;</span>
            <span class="text-foreground/60">запуск:</span>
            <span class="text-term-amber min-w-[1ch]">{{ typed }}</span>
            <span class="inline-block w-[0.6ch] h-[1.1em] bg-term-amber align-middle animate-type-caret" />
          </div>

          <p class="mt-4 md:mt-5 text-base md:text-lg text-foreground/70 leading-relaxed max-w-2xl">
            250+ специалистов из Яндекса, Тинькофф, VK и стартапов.
            Менторство, ИИ, нетворкинг и подготовка к собеседованиям.
          </p>
        </div>

        <!-- CTA with magnetic hover -->
        <div
          ref="ctaRef"
          class="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:gap-6 mt-8 md:mt-12 transition-all duration-700 delay-[400ms]"
          :class="heroReady ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          :style="{ transform: heroReady ? `translate(${magX}px, ${magY}px)` : undefined }"
        >
          <div class="cta-glow rounded-[50px]">
            <TelegramAuth
              v-if="!tgUser || !tgToken"
              @auth="setUser"
            />
            <Button
              v-else
              variant="filled"
              as="a"
              href="/platform"
              rel="noopener noreferrer"
              @click="trackPlatformClick"
            >
              Перейти в платформу
            </Button>
          </div>
          <a
            href="#tariffs"
            class="font-mono text-xs md:text-sm text-foreground/60 hover:text-accent transition-colors flex items-center gap-2 group"
            @click="trackHeroPricingClick"
          >
            <span class="text-accent group-hover:animate-pulse">$</span>
            <span>./тарифы --от_520₽</span>
            <span class="inline-block transition-transform group-hover:translate-x-0.5">→</span>
          </a>
        </div>

        <!-- Stats strip with count-up -->
        <div
          class="mt-10 md:mt-16 grid grid-cols-3 max-w-2xl divide-x divide-accent/15 border-t border-b border-accent/15 transition-all duration-700 delay-500"
          :class="heroReady ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
        >
          <div class="py-4 pr-4">
            <div class="font-display text-2xl md:text-4xl text-accent tabular-nums">
              {{ stat1.value.value }}+
            </div>
            <div class="font-mono text-[10px] md:text-xs text-foreground/50 uppercase mt-1 tracking-widest">
              участников
            </div>
          </div>
          <div class="py-4 px-4">
            <div class="font-display text-2xl md:text-4xl text-accent tabular-nums">
              {{ stat2.value.value }}+
            </div>
            <div class="font-mono text-[10px] md:text-xs text-foreground/50 uppercase mt-1 tracking-widest">
              менторов
            </div>
          </div>
          <div class="py-4 pl-4">
            <div class="font-display text-2xl md:text-4xl text-accent tabular-nums">
              {{ stat3.value.value }}×
            </div>
            <div class="font-mono text-[10px] md:text-xs text-foreground/50 uppercase mt-1 tracking-widest">
              встреч в месяц.
            </div>
          </div>
        </div>

        <!-- Mobile scroll cue (terminal-flavored) -->
        <a
          href="#why"
          class="md:hidden mt-7 flex items-center gap-2 font-mono text-[10px] tracking-[0.2em] uppercase text-foreground/40 group"
        >
          <span class="text-accent/60">$</span>
          <span class="text-foreground/60">./след_секция</span>
          <span class="text-accent inline-block animate-bounce-soft">↓</span>
          <span class="flex-1 border-t border-dashed border-accent/15 ml-1" />
        </a>
      </div>
    </div>

    <!-- Marquee of companies -->
    <div class="relative z-10 border-t border-accent/15 bg-background/40 backdrop-blur-sm overflow-hidden py-4 group/marquee">
      <div class="flex items-center gap-3 font-mono text-[10px] md:text-xs tracking-[0.2em] uppercase">
        <div class="px-4 md:px-10 shrink-0 text-accent/70">
          // команды из:
        </div>
        <div class="flex-1 overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_8%,black_92%,transparent)]">
          <div class="flex gap-10 animate-marquee group-hover/marquee:[animation-play-state:paused] whitespace-nowrap will-change-transform">
            <span
              v-for="(c, i) in [...companies, ...companies]"
              :key="i"
              class="text-foreground/60 hover:text-accent transition-colors cursor-default"
            >
              {{ c }} <span class="text-accent/30 ml-10">◇</span>
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
@media (min-width: 420px) {
  .xs\:text-\[46px\] {
    font-size: 46px;
  }
}
.tabular-nums {
  font-variant-numeric: tabular-nums;
}

/* Cellular-style signal bars in status bar */
.signal-bars {
  display: inline-flex;
  align-items: flex-end;
  gap: 1.5px;
  height: 0.7em;
}
.signal-bars > span {
  display: block;
  width: 2px;
  background: hsl(var(--accent));
  opacity: 0.85;
  border-radius: 1px;
}
.signal-bars > span:nth-child(1) { height: 25%; }
.signal-bars > span:nth-child(2) { height: 50%; }
.signal-bars > span:nth-child(3) { height: 75%; }
.signal-bars > span:nth-child(4) { height: 100%; opacity: 0.35; }
.signal-bars > span:nth-child(4) {
  animation: signal-pulse 2.4s ease-in-out infinite;
}
@keyframes signal-pulse {
  0%, 100% { opacity: 0.35; }
  50% { opacity: 0.85; }
}

/* Soft bounce for mobile scroll cue */
@keyframes bounce-soft {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(3px); }
}
.animate-bounce-soft {
  animation: bounce-soft 1.8s ease-in-out infinite;
}
</style>
