<script setup lang="ts">
import type { TelegramUser } from '@/services/auth.ts'
import { Button, Typography } from 'itx-ui-kit'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import TelegramAuth from '@/components/TelegramAuth.vue'
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
</script>

<template>
  <section
    class="relative w-full py-16 md:py-10 lg:py-12 flex flex-col  min-h-[calc(100svh-72px)] md:min-h-fit md:h-auto justify-between items-center
    before:absolute before:content-['']  before:-z-10 before:w-[1230px] before:h-[836px]
    before:bg-[url('../assets/icons/crosses-bg.svg')] before:bg-no-repeat before:bg-contain
    lg:before:-top-2/3 lg:before:-right-[20%] lg:before:bg-right-top
    sm:before:-top-[280px] sm:before:-right-[480px] sm:before:scale-100 sm:before:opacity-100
    before:scale-[0.65] before:-top-[200px] before:opacity-50
    overflow-x-clip
"
  >
    <div class="container px-6 md:px-10 flex flex-col gap-8 items-center md:items-start justify-between h-full flex-1 ">
      <div class="flex flex-col gap-6 sm:gap-1 text-center md:text-left  sm:items-center md:items-start w-full  md:max-w-6xl">
        <Typography
          variant="h1"
          as="h1"
          class="uppercase text-accent"
        >
          IT-хозяева
        </Typography>
        <Typography
          variant="h1"
          as="h1"
          class="uppercase max-sm:text-2xl break-words hyphens-auto w-full max-w-full"
        >
          закрытое сообщество <span class="whitespace-nowrap">IT-специалистов</span>, в котором ценят опыт, юмор и профессио&shy;нализм
        </Typography>
      </div>
      <div class="flex justify-start ">
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
    </div>
  </section>
</template>
