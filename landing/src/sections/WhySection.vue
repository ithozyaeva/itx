<script setup lang="ts">
import type { Component } from 'vue'
import Button from '@/components/ui/Button.vue'
import FeatureItem from '@/components/ui/FeatureItem.vue'
import Database from '~icons/lucide/database'
import Globe from '~icons/lucide/globe'
import MessageSquare from '~icons/lucide/message-square'
import Shield from '~icons/lucide/shield'
import Star from '~icons/lucide/star'
import Users from '~icons/lucide/users'

interface Feature {
  title: string
  description: string
  icon: Component
}

const features: Feature[] = [
  {
    icon: Shield,
    title: 'Экспертная поддержка',
    description: 'Получите доступ к консультациям от ведущих специалистов отрасли',
  },
  {
    icon: Database,
    title: 'Эксклюзивные материалы',
    description: 'Библиотека курсов, вебинаров и документации, недоступная в открытом доступе',
  },
  {
    icon: Users,
    title: 'Сильное комьюнити',
    description: 'Нетворкинг с профессионалами и единомышленниками из IT-индустрии',
  },
  {
    icon: Globe,
    title: 'Карьерные возможности',
    description: 'Закрытая база вакансий и рекомендации для трудоустройства',
  },
  {
    icon: MessageSquare,
    title: 'Приватные обсуждения',
    description: 'Закрытые чаты для обсуждения актуальных тем и решения сложных задач',
  },
  {
    icon: Star,
    title: 'Персональный рост',
    description: 'Индивидуальные планы развития и менторство от лидеров индустрии',
  },
]
</script>

<template>
  <section id="why" class="w-full py-12 md:py-24 lg:py-32">
    <div class="container px-4 md:px-6">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="space-y-2">
          <div class="inline-block rounded-lg bg-primary px-3 py-1 text-sm text-primary-foreground">
            Преимущества
          </div>
          <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">
            Почему выбирают нас
          </h2>
          <p class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
            Наше сообщество предлагает уникальные возможности для профессионального роста и развития в IT-сфере
          </p>
        </div>
      </div>
      <div class="mx-auto grid max-w-5xl items-center gap-6 py-12 md:grid-cols-2 md:gap-10 gap-y-4 md:gap-y-4">
        <FeatureItem v-for="feature in features" :key="feature.title" :title="feature.title" :description="feature.description" :icon="feature.icon" />
      </div>
      <div class="flex justify-center flex-col sm:flex-row items-start gap-4">
        <a href="#tariffs">
          <Button class="px-8">
            Выбрать тариф
          </Button>
        </a>
        <a href="#reviews">
          <Button type="outline" class="px-8">
            Отзывы участников
          </Button>
        </a>
      </div>
    </div>
  </section>
</template>
