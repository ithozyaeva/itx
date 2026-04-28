<script setup lang="ts">
import StructuredData from '@/components/StructuredData.vue'
import { usePageMeta } from '@/composables/useMeta'
import { useRevealObserver } from '@/composables/useRevealObserver'
import MentorsSection from '@/sections/MentorsSection.vue'

useRevealObserver()

usePageMeta({
  title: 'База менторов — старшие разработчики, ИИ-инженеры, тимлиды',
  description: 'База из 60+ ИТ-менторов: старшие (senior) разработчики, тимлиды, архитекторы и ИИ-инженеры из Яндекса, Тинькофф, VK. Помогают с карьерой, сменой уровня, переходом в ИИ и подготовкой к собеседованиям.',
  url: 'https://ithozyaeva.ru/mentors',
  type: 'website',
})

const structuredData = [
  {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    'itemListElement': [
      { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': 'https://ithozyaeva.ru/' },
      { '@type': 'ListItem', 'position': 2, 'name': 'Наставники', 'item': 'https://ithozyaeva.ru/mentors' },
    ],
  },
  {
    '@context': 'https://schema.org',
    '@type': 'Service',
    'serviceType': 'Наставничество в ИТ',
    'provider': { '@type': 'Organization', 'name': 'IT-ХОЗЯЕВА', 'url': 'https://ithozyaeva.ru' },
    'areaServed': { '@type': 'Country', 'name': 'Россия' },
    'description': 'Индивидуальные консультации со старшими разработчиками, архитекторами и ИИ-инженерами. Карьерное менторство, разбор резюме, подготовка к собеседованиям, смена уровня, переход в ИИ.',
    'offers': {
      '@type': 'Offer',
      'priceCurrency': 'RUB',
      'price': '1000',
      'description': 'Доступ к базе менторов на тарифе ХОЗЯИН от 1000 ₽/мес',
    },
  },
]
</script>

<template>
  <StructuredData :data="structuredData" />
  <section class="w-full pt-10 md:pt-16 lg:pt-20">
    <div class="container px-6 md:px-10">
      <nav
        class="font-mono text-xs text-foreground/50 mb-6"
        aria-label="breadcrumb"
      >
        <a
          href="/"
          class="hover:text-accent transition-colors"
        >~/</a>
        <span class="mx-2 text-accent/40">/</span>
        <span class="text-foreground">mentors</span>
      </nav>

      <div class="max-w-4xl">
        <h1 class="font-display uppercase text-[32px] sm:text-[48px] md:text-[64px] lg:text-[80px] leading-[0.9] tracking-tight text-accent">
          База менторов
        </h1>
        <p class="mt-6 md:mt-8 text-base md:text-lg text-foreground/75 leading-relaxed max-w-3xl">
          60+ ИТ-менторов: старшие разработчики, тимлиды, архитекторы и ИИ-инженеры из топовых российских ИТ-компаний —
          Яндекс, Тинькофф, VK, Ozon, Сбер, Kaspersky, JetBrains. Помогаем расти по уровням, сменить направление, подготовиться
          к собеседованию и разобраться в ИИ и вайбкодинге.
        </p>

        <div class="mt-10 md:mt-14 grid grid-cols-1 md:grid-cols-3 gap-px bg-accent/15 border border-accent/15">
          <div class="bg-background/95 p-6 md:p-7">
            <h2 class="font-display uppercase text-lg md:text-xl text-accent mb-2">
              Кто в базе
            </h2>
            <p class="text-sm md:text-base text-foreground/70 leading-relaxed">
              Старшие и ведущие разработчики, архитекторы, ИИ-инженеры, DevOps-инженеры, руководители тестирования. Ребята, которые реально собеседуют
              и принимают в топовых командах.
            </p>
          </div>
          <div class="bg-background/95 p-6 md:p-7">
            <h2 class="font-display uppercase text-lg md:text-xl text-accent mb-2">
              С чем помогают
            </h2>
            <p class="text-sm md:text-base text-foreground/70 leading-relaxed">
              Карьерное планирование, разбор резюме, тренировочные собеседования, подготовка к системному дизайну, разбор пул-реквестов, переход
              в ИИ-разработку, смена стека.
            </p>
          </div>
          <div class="bg-background/95 p-6 md:p-7">
            <h2 class="font-display uppercase text-lg md:text-xl text-accent mb-2">
              Как получить доступ
            </h2>
            <p class="text-sm md:text-base text-foreground/70 leading-relaxed">
              Оформить тариф
              <a
                href="https://boosty.to/jointime/purchase/3150814"
                class="text-accent hover:underline"
                rel="noopener"
              >ХОЗЯИН на Boosty</a>
              — от 1000 ₽/мес. Получаете полную таблицу контактов и расписание сессий.
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <MentorsSection />
</template>
