<script setup lang="ts">
import { PriceCard, Typography } from 'itx-ui-kit'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import { useScrollReveal } from '@/composables/useScrollReveal'

interface Tariff {
  name: string
  description: string
  price: number
  oldPrice?: number
  features: string[]
  link: string
  isPopular?: boolean
}

const yandexMetrika = useYandexMetrika()

function handleSubscriptionClick(tariffName: string, link: string) {
  yandexMetrika.reachGoal('subscription_click', {
    tariff: tariffName,
  } as any)
  yandexMetrika.extLink(link, { title: `Подписка ${tariffName}` })
}

const tariffs: Tariff[] = [
  {
    name: 'Бригадир',
    description: 'Старт в IT-сообществе',
    price: 520,
    features: [
      'Доступ ко всем обучающим материалам',
      'IT-чаты по направлениям и грейдам',
      'Еженедельные вебинары и воркшопы',
      'Вакансии и рефералки от участников',
    ],
    link: 'https://boosty.to/jointime/purchase/3150816',
  },
  {
    name: 'ХОЗЯИН',
    description: 'Максимум возможностей для роста в IT',
    price: 2000,
    isPopular: true,
    features: [
      'Все возможности тарифа «Бригадир»',
      'Приоритетная поддержка и ревью резюме',
      'Влияние на темы встреч и контент',
      'Доступ к базе менторов и таблице экспертов',
    ],
    link: 'https://boosty.to/jointime/purchase/3150814',
  },
  {
    name: 'KING',
    description: 'Персональное менторство и продвижение',
    price: 5200,
    features: [
      'Все возможности тарифа «ХОЗЯИН»',
      'Размещение рекламы ваших ресурсов',
      'Топ-позиция в таблице менторов с лычкой',
      'Личная консультация по карьере или технологиям',
    ],
    link: 'https://boosty.to/jointime/purchase/967625',
  },
]

const { containerRef, isVisible } = useScrollReveal({ threshold: 0.05 })
</script>

<template>
  <section
    id="tariffs"
    ref="containerRef"
    class="w-full pt-12 md:pt-24 lg:pt-32"
  >
    <div class="container px-6 md:px-10">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="space-y-2">
          <Typography
            variant="h2"
            as="h2"
            class="text-accent"
          >
            Тарифы и подписка
          </Typography>
          <Typography
            variant="body-xl"
            as="p"
            class="max-w-[540px]"
          >
            Выберите подходящий тариф — подписка оформляется через Boosty, можно сменить план в любой момент
          </Typography>
        </div>
      </div>
      <div class="grid pt-12 lg:grid-cols-3 gap-5">
        <div
          v-for="(tariff, index) in tariffs"
          :key="tariff.name"
          class="contents"
          @click="handleSubscriptionClick(tariff.name, tariff.link)"
        >
          <PriceCard
            :name="tariff.name"
            :price="tariff.price"
            :old-price="tariff.oldPrice"
            :features="tariff.features"
            :link="tariff.link"
            :variant="tariff.isPopular ? 'highlighted' : 'default'"
            :class="[
              isVisible
                ? tariff.isPopular ? 'animate-card-reveal-highlight' : 'animate-card-reveal'
                : 'opacity-0',
            ]"
            :style="isVisible ? { animationDelay: `${index * 150}ms` } : undefined"
          />
        </div>
      </div>
    </div>
  </section>
</template>
