<script setup lang="ts">
import type { TelegramUser } from '@/services/auth'
import { useWindowSize } from '@vueuse/core'
import { BurgerIcon, Button, CloseIcon } from 'itx-ui-kit'
import { computed, onUnmounted, ref, watch } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import Logo from '@/assets/itx-logo.svg'
import { useScrollHeader } from '@/composables/useScrollHeader.ts'
import { useToken } from '@/composables/useToken'
import { useUser } from '@/composables/useUser'
import Navigation from '../sections/Navigation.vue'
import TelegramAuth from './TelegramAuth.vue'

const tgUser = useUser()
const tgToken = useToken()
const yandexMetrika = useYandexMetrika()

function setUser(user: TelegramUser, token: string) {
  tgUser.value = user
  tgToken.value = token
}

function trackPlatformClick() {
  yandexMetrika.reachGoal('platform_redirect_click', {
    location: 'header',
    isAuthenticated: !!tgUser.value,
  } as any)
}

const { isScrolled } = useScrollHeader(10)
const { width } = useWindowSize()

const isMobile = computed(() => {
  return width.value < 1080
})

const isMenuOpen = ref(false)

function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
}

watch(isMenuOpen, (open) => {
  if (open) {
    document.body.style.overflow = 'hidden'
  }
  else {
    document.body.style.overflow = ''
  }
})

onUnmounted(() => {
  document.body.style.overflow = ''
})
</script>

<template>
  <header
    class="sticky top-0 z-50 w-full transition-all duration-500"
    :class="{ 'bg-background/0.9 backdrop-blur-[25px]': isScrolled, 'bg-transparent backdrop-blur-0': !isScrolled }"
  >
    <div class="container px-6 md:px-10 flex gap-5 h-[var(--header-height)] items-center justify-between md:justify-between">
      <div class="flex items-center gap-4 lg:gap-16">
        <a
          href="/"
          class="flex items-center gap-2 font-bold text-xl"
        >
          <div class=" w-16  sm:w-[88px] md:w-24">
            <Logo />
          </div>
        </a>
        <Navigation v-if="!isMobile" />
      </div>
      <div
        v-if="!isMobile"
        class="flex justify-end "
      >
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
      </div>
      <div v-else>
        <BurgerIcon
          v-if="!isMenuOpen"
          class="cursor-pointer hover:opacity-75 transition-opacity"
          @click="toggleMenu"
        />
        <CloseIcon
          v-else
          class="cursor-pointer size-12 hover:opacity-75 transition-opacity"
          @click="toggleMenu"
        />
      </div>
    </div>
  </header>
  <transition name="fade">
    <div
      v-if="isMenuOpen"
      class="px-6 md:px-10 inset-0 fixed bg-background z-40 flex flex-col justify-center items-center"
    >
      <Navigation
        class="flex flex-col px-6 md:px-10 justify-center items-center navigation"
        @click="toggleMenu"
      />
      <div class="absolute bottom-11 sm:bottom-16 sm:h-16 px-6 md:px-10">
        <TelegramAuth
          v-if="!tgUser"
          class="h-full"
          @auth="setUser"
        />
        <Button
          v-else
          class="h-full inline-block"
          as="a"
          href="/platform"
          @click="trackPlatformClick"
        >
          Перейти в платформу
        </Button>
      </div>
    </div>
  </transition>
</template>

<style scoped>
::v-deep(.navigation *) {
  font-size: 48px;
  @media (width <= 640px) {
    font-size: 36px;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
