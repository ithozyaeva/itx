<script setup lang="ts">
import type { TelegramUser } from '@/services/auth.ts'
import { onMounted, onUnmounted, ref } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import TelegramAuth from '@/components/TelegramAuth.vue'
import Button from '@/components/ui/UiButton.vue'
import { useToken } from '@/composables/useToken.ts'
import { useUser } from '@/composables/useUser.ts'

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

// Typewriter effect for the tagline
const phrases = [
  'менторство_',
  'vibe_coding_',
  'AI-практики_',
  'нетворкинг_',
  'собесы_',
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

// Live "uptime" counter
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

onMounted(() => {
  tick()
  updateUptime()
  uptimeTimer = window.setInterval(updateUptime, 1000)
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
    class="relative w-full min-h-[100svh] md:min-h-[calc(100svh-var(--header-height))] flex flex-col justify-between overflow-hidden"
  >
    <!-- Decorative grid / crosshair -->
    <div
      aria-hidden="true"
      class="pointer-events-none absolute inset-0 opacity-[0.12]"
      style="background-image: radial-gradient(circle at 80% 30%, hsl(var(--accent)) 0, transparent 55%), radial-gradient(circle at 20% 80%, #ffb547 0, transparent 45%);"
    />
    <div
      aria-hidden="true"
      class="pointer-events-none absolute top-[18%] right-[6%] hidden lg:block"
    >
      <div class="w-[520px] h-[520px] rounded-full border border-accent/10 animate-[glitch-x_6s_steps(1)_infinite]">
        <div class="absolute inset-8 rounded-full border border-accent/15" />
        <div class="absolute inset-20 rounded-full border border-accent/20" />
        <div class="absolute inset-32 rounded-full border border-accent/25" />
        <div class="absolute top-1/2 left-0 right-0 h-px bg-accent/20" />
        <div class="absolute left-1/2 top-0 bottom-0 w-px bg-accent/20" />
      </div>
    </div>

    <!-- Terminal status bar -->
    <div class="container px-6 md:px-10 pt-8 md:pt-12 relative z-10">
      <div class="flex items-center justify-between font-mono text-[11px] md:text-xs tracking-[0.08em] text-foreground/50">
        <div class="flex items-center gap-4 md:gap-6">
          <div class="flex items-center gap-1.5">
            <span class="w-2 h-2 rounded-full bg-accent shadow-[0_0_8px_hsl(var(--accent))]" />
            <span class="text-accent">ONLINE</span>
          </div>
          <span class="hidden sm:inline">uptime: <span class="text-foreground/80">{{ uptime }}</span></span>
          <span class="hidden md:inline">v4.2.1</span>
        </div>
        <div class="flex items-center gap-4 md:gap-6">
          <span class="hidden sm:inline">250+ users</span>
          <span>msk/ru</span>
        </div>
      </div>
    </div>

    <!-- Main hero content -->
    <div class="container px-6 md:px-10 flex-1 flex items-center py-12 md:py-16 relative z-10">
      <div class="w-full max-w-5xl">
        <div class="font-mono text-xs md:text-sm text-foreground/50 mb-5 md:mb-6">
          <span class="text-accent">community@ithozyaeva</span>:<span class="text-term-amber">~</span>$
          <span class="text-foreground/80">./welcome --new-user</span>
        </div>

        <div class="flex items-baseline gap-2 md:gap-3 mb-1">
          <span class="font-mono text-accent/60 text-sm md:text-base">&gt;</span>
          <h1 class="font-display uppercase text-[38px] xs:text-[46px] sm:text-[68px] md:text-[96px] lg:text-[130px] leading-[0.85] tracking-tight text-accent">
            IT-ХОЗЯЕВА
          </h1>
        </div>

        <div class="pl-4 md:pl-6 border-l border-accent/30 mt-6 md:mt-8 max-w-3xl">
          <h2 class="font-display uppercase text-[22px] sm:text-[32px] md:text-[42px] leading-[1.05] text-foreground">
            Закрытое сообщество<br>
            <span class="text-foreground/40">IT-специалистов</span>
          </h2>

          <div class="mt-5 md:mt-7 font-mono text-sm md:text-base text-foreground/75 flex flex-wrap items-center gap-x-2">
            <span class="text-accent/60">&gt;&gt;</span>
            <span class="text-foreground/60">exec:</span>
            <span class="text-term-amber min-w-[1ch]">{{ typed }}</span>
            <span class="inline-block w-[0.6ch] h-[1.1em] bg-term-amber align-middle animate-type-caret" />
          </div>

          <p class="mt-4 md:mt-5 text-base md:text-lg text-foreground/70 leading-relaxed max-w-2xl">
            250+ специалистов из Яндекса, Tinkoff, VK и стартапов.
            Менторство, AI, нетворкинг и подготовка к собеседованиям.
          </p>
        </div>

        <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:gap-6 mt-8 md:mt-12">
          <TelegramAuth
            v-if="!tgUser"
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
          <a
            href="#why"
            class="font-mono text-xs md:text-sm text-foreground/50 hover:text-accent transition-colors flex items-center gap-2"
          >
            <span class="text-accent">$</span> ./scroll-down --see-more
          </a>
        </div>

        <!-- Stats strip -->
        <div class="mt-10 md:mt-16 grid grid-cols-3 max-w-2xl divide-x divide-accent/15 border-t border-b border-accent/15">
          <div class="py-4 pr-4">
            <div class="font-display text-2xl md:text-4xl text-accent">
              250+
            </div>
            <div class="font-mono text-[10px] md:text-xs text-foreground/50 uppercase mt-1 tracking-widest">
              участников
            </div>
          </div>
          <div class="py-4 px-4">
            <div class="font-display text-2xl md:text-4xl text-accent">
              60+
            </div>
            <div class="font-mono text-[10px] md:text-xs text-foreground/50 uppercase mt-1 tracking-widest">
              менторов
            </div>
          </div>
          <div class="py-4 pl-4">
            <div class="font-display text-2xl md:text-4xl text-accent">
              7×
            </div>
            <div class="font-mono text-[10px] md:text-xs text-foreground/50 uppercase mt-1 tracking-widest">
              встреч в мес.
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Marquee of companies -->
    <div class="relative z-10 border-t border-accent/15 bg-background/40 backdrop-blur-sm overflow-hidden py-4">
      <div class="flex items-center gap-3 font-mono text-[10px] md:text-xs tracking-[0.2em] uppercase">
        <div class="px-4 md:px-10 shrink-0 text-accent/70">
          // команды из:
        </div>
        <div class="flex-1 overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_8%,black_92%,transparent)]">
          <div class="flex gap-10 animate-marquee whitespace-nowrap will-change-transform">
            <span
              v-for="(c, i) in [...companies, ...companies]"
              :key="i"
              class="text-foreground/60 hover:text-accent transition-colors"
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
</style>
