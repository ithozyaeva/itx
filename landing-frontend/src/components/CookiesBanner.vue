<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import { Button, Typography } from 'itx-ui-kit'
import { computed } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'

const cookiesAccepted = useLocalStorage<boolean>('cookies_accepted', false)
const yandexMetrika = useYandexMetrika()

const isVisible = computed(() => !cookiesAccepted.value)

function acceptCookies() {
  yandexMetrika.reachGoal('cookies_banner_accept')
  cookiesAccepted.value = true
}
</script>

<template>
  <Transition name="banner-slide">
    <div
      v-if="isVisible"
      class="fixed bottom-0 left-0 right-0 z-40"
    >
      <div class="mx-5 my-3">
        <div
          class="max-w-[1400px] mx-auto bg-[#1D2723]  bg-opacity-95 border-2 rounded-3xl border-[#2B3D36]/20 shadow-lg"
        >
          <div class="container px-6 md:px-10 py-4 md:py-6">
            <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
              <div class="flex-1">
                <Typography
                  variant="body-m"
                  as="p"
                  class="text-foreground"
                >
                  Для повышения удобства работы с сайтом ithozyaeva.ru использует файлы cookie и рекомендательные технологии. В них содержатся данные о прошлых посещениях сайта. Если вы не хотите, чтобы эти данные обрабатывались, отключите cookie в настройках браузера.
                </Typography>
              </div>
              <Button
                as="button"
                variant="filled"
                class="shrink-0 w-full md:w-fit"
                @click="acceptCookies"
              >
                Понятно
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.banner-slide-enter-active {
  transition: transform 0.3s ease-out, opacity 0.3s ease-out;
}

.banner-slide-leave-active {
  transition: transform 0.3s ease-in, opacity 0.3s ease-in;
}

.banner-slide-enter-from {
  transform: translateY(100%);
  opacity: 0;
}

.banner-slide-leave-to {
  transform: translateY(100%);
  opacity: 0;
}

.banner-slide-enter-to,
.banner-slide-leave-from {
  transform: translateY(0);
  opacity: 1;
}
</style>
