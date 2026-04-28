<script setup lang="ts">
import { useYandexMetrika } from 'yandex-metrika-vue3'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import PriceCard from '@/components/ui/UiPriceCard.vue'
import { useScrollReveal } from '@/composables/useScrollReveal'

interface Tariff {
  name: string
  description: string
  price: number
  oldPrice?: number
  features: string[]
  link: string
  isPopular?: boolean
  tierIndex: string
  tierLabel: string
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
    description: 'Старт в ИТ-сообществе',
    price: 520,
    tierIndex: '01',
    tierLabel: 'базовый',
    features: [
      'Доступ ко всем обучающим материалам',
      'ИТ-чаты по направлениям и уровням',
      'Еженедельные вебинары и практикумы',
      'Вакансии и реферальные ссылки от участников',
    ],
    link: 'https://boosty.to/jointime/purchase/3150816',
  },
  {
    name: 'ХОЗЯИН',
    description: 'Максимум возможностей для роста в ИТ',
    price: 1000,
    oldPrice: 2000,
    isPopular: true,
    tierIndex: '02',
    tierLabel: 'про',
    features: [
      'Все возможности тарифа «Бригадир»',
      'Приоритетная поддержка и разбор резюме',
      'Влияние на темы встреч и контент',
      'Доступ к базе менторов и таблице экспертов',
      'Закрытые ИИ-беседы: AI-X и База стародубцева',
    ],
    link: 'https://boosty.to/jointime/purchase/3150814',
  },
  {
    name: 'МАСТЕР',
    description: 'Персональное менторство и продвижение',
    price: 5200,
    tierIndex: '03',
    tierLabel: 'макс',
    features: [
      'Все возможности тарифа «ХОЗЯИН»',
      'Размещение рекламы ваших ресурсов',
      'Верхняя позиция в таблице менторов с лычкой',
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
    class="w-full pt-20 md:pt-32 lg:pt-40"
  >
    <div class="container px-6 md:px-10">
      <div class="reveal">
        <SectionHeader
          index="04"
          path="~/community/access.sh"
          title="Уровни доступа"
          subtitle="Подписка через Boosty. План можно повысить или понизить в любой момент."
        />
      </div>

      <div class="grid pt-14 md:pt-16 md:grid-cols-2 lg:grid-cols-3 gap-5 md:gap-6 items-stretch">
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
            :tier-index="tariff.tierIndex"
            :tier-label="tariff.tierLabel"
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
