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
  <footer class="w-full md:py-0 bg-accent rounded-t-[50px]">
    <div class="flex px-6 md:px-10 container  flex-col lg:flex-row gap-12 lg:gap-9 pt-8 pb-11">
      <div class="flex flex-col gap-5 basis-1/2">
        <Typography
          variant="h3"
          as="h3"
          class="text-2xl text-background"
        >
          Вступай в сообщество и становись IT-хозяином
        </Typography>
        <TelegramAuth
          v-if="!tgUser"
          variant="dark-filled"
          @auth="setUser"
        />
        <Button
          v-else
          variant="dark-filled"
          as="a"
          class="block w-fit "
          href="/platform"
          rel="noopener noreferrer"
          @click="trackPlatformClick"
        >
          Перейти в платформу
        </Button>
      </div>

      <div class="flex flex-col gap-10 lg:gap-14 basis-1/2 text-background">
        <Typography>
          По всем вопросам:<br>
          <Typography
            as="a"
            variant="body-l"
            href="https://t.me/jointimer"
            target="_blank"
            rel="noopener noreferrer"
            class="underline"
            @click="trackJointimerClick"
          >
            @jointimer
          </Typography>
        </Typography>
        <div class="flex flex-col gap-1 align-bottom">
          <Typography
            as="a"
            variant="body-s"
            href="/privacy"
            target="_blank"
            rel="noopener noreferrer"
            class="underline"
          >
            Политика обработки пользовательских данных
          </Typography>
          <Typography
            as="span"
            variant="body-s"
            class="flex items-center"
          >
            ©{{ new Date().getFullYear() }}, IT-ХОЗЯЕВА. Все права защищены
          </Typography>
        </div>
      </div>
    </div>
  </footer>
</template>
