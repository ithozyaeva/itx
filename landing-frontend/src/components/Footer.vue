<script setup lang="ts">
import type { TelegramUser } from '@/services/auth.ts'
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
    location: 'footer',
    isAuthenticated: !!tgUser.value,
  } as any)
}

function trackJointimerClick() {
  yandexMetrika.reachGoal('footer_jointimer_click', {
    location: 'footer',
  } as any)
  yandexMetrika.extLink('https://t.me/jointimer', { title: '@jointimer' })
}
</script>

<template>
  <footer class="w-full mt-24 md:mt-32 relative overflow-hidden">
    <!-- big CTA band -->
    <div class="relative bg-accent text-[#0b0d0c]">
      <!-- scanline overlay -->
      <div
        aria-hidden="true"
        class="pointer-events-none absolute inset-0 opacity-[0.06]"
        style="background-image: repeating-linear-gradient(to bottom, transparent 0, transparent 3px, #000 3px, #000 4px);"
      />

      <div class="container px-6 md:px-10 py-16 md:py-24 relative">
        <div class="flex items-center gap-3 font-mono text-[11px] md:text-xs tracking-[0.12em] uppercase mb-6">
          <span class="w-2 h-2 rounded-full bg-[#0b0d0c] animate-pulse" />
          <span class="text-[#0b0d0c]">$ ./join_now --final-call</span>
          <span class="flex-1 h-px bg-[#0b0d0c]/30" />
          <span class="hidden sm:inline text-[#0b0d0c]/70">[06] ~/join</span>
        </div>

        <h2 class="reveal font-display uppercase text-[38px] sm:text-[56px] md:text-[80px] lg:text-[104px] leading-[0.9] tracking-tight text-[#0b0d0c]">
          Стань<br>
          <span class="text-[#0b0d0c]/50">хозяином</span>
        </h2>

        <div class="mt-8 md:mt-12 flex flex-col md:flex-row md:items-center gap-5 md:gap-8">
          <TelegramAuth
            v-if="!tgUser"
            variant="dark-filled"
            @auth="setUser"
          />
          <Button
            v-else
            variant="dark-filled"
            as="a"
            class="block w-fit"
            href="/platform"
            rel="noopener noreferrer"
            @click="trackPlatformClick"
          >
            Перейти в платформу
          </Button>

          <div class="font-mono text-xs md:text-sm text-[#0b0d0c]/70">
            <div>
              &gt; support:
              <a
                href="https://t.me/jointimer"
                target="_blank"
                rel="noopener noreferrer"
                class="text-[#0b0d0c] underline underline-offset-4 hover:opacity-70 transition-opacity"
                @click="trackJointimerClick"
              >@jointimer</a>
            </div>
            <div class="mt-1">
              &gt; subscription: ~520₽/mo
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Terminal footer -->
    <div class="bg-background border-t border-accent/15">
      <div class="container px-6 md:px-10 py-8">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 font-mono text-[11px] md:text-xs text-foreground/50 tracking-wide">
          <div class="flex items-center gap-4 flex-wrap">
            <span class="text-accent">[●]</span>
            <span>©{{ new Date().getFullYear() }} IT-ХОЗЯЕВА</span>
            <span class="text-foreground/25">|</span>
            <a
              href="/mentors"
              class="hover:text-accent transition-colors underline underline-offset-4"
            >
              mentors
            </a>
            <span class="text-foreground/25">|</span>
            <a
              href="/vibe-coding"
              class="hover:text-accent transition-colors underline underline-offset-4"
            >
              vibe-coding
            </a>
            <span class="text-foreground/25">|</span>
            <a
              href="/privacy"
              class="hover:text-accent transition-colors underline underline-offset-4"
            >
              privacy.policy
            </a>
          </div>
          <div class="flex items-center gap-4 text-foreground/35">
            <span>build: 4.2.1</span>
            <span class="text-foreground/20">|</span>
            <span>region: ru-msk</span>
            <span class="text-foreground/20">|</span>
            <span>status: <span class="text-accent">operational</span></span>
          </div>
        </div>

        <!-- ASCII banner -->
        <div class="mt-6 font-mono text-[10px] leading-[1.2] text-accent/25 select-none overflow-hidden whitespace-pre hidden md:block">
          ────────────────────────────────────────────────────────────────────────────────────────────────
          &gt; end of transmission. stay curious. ship code. own it.
          ────────────────────────────────────────────────────────────────────────────────────────────────
        </div>
      </div>
    </div>
  </footer>
</template>
