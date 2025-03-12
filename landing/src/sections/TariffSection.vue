<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import Check from '~icons/lucide/check'

interface Tariff {
  name: string
  description: string
  price: number
  features: string[]
  link: string
  isPopular?: boolean
}

const tariffs: Tariff[] = [
  {
    name: 'Работяга',
    description: 'Для начинающих специалистов',
    price: 200,
    features: [
      'Доступ к основным материалам',
      'Участие в общих обсуждениях',
      'Еженедельные вебинары',
      'Базовая поддержка сообщества',
    ],
    link: 'https://boosty.to/jointime/purchase/3150818',
  },
  {
    name: 'Начальник отдела',
    description: 'Для активных IT-специалистов',
    price: 400,
    isPopular: true,
    features: [
      'Все преимущества предыдущих тарифов',
      'VPN',
      'Приоритетная техническая поддержка',
      'Возможность влиять на запись контента',
      'Возможность участия в таблице менторов',
    ],
    link: 'https://boosty.to/jointime/purchase/3150815',
  },
  {
    name: 'Проверенный ментор',
    description: 'Для лидеров и руководителей',
    price: 4000,
    features: [
      'Все преимущества предыдущих тарифов',
      'Реклама ресурсов',
      'Повышенная позиция в таблице менторов',
      'Лычка в таблице менторов',
      'Личная консультация по любой теме',
    ],
    link: 'https://boosty.to/jointime/purchase/967625',
  },
]
</script>

<template>
  <section id="tariffs" class="w-full py-12 md:py-24 lg:py-32">
    <div class="container px-4 md:px-6">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="space-y-2">
          <div class="inline-block rounded-lg bg-primary px-3 py-1 text-sm text-primary-foreground">
            Тарифы
          </div>
          <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">
            Выберите свой план
          </h2>
          <p class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
            Гибкие тарифы для разных потребностей и уровней профессионального развития
          </p>
        </div>
      </div>
      <div class="grid gap-6 pt-12 lg:grid-cols-3 lg:gap-8">
        <Card v-for="tariff in tariffs" :key="tariff.name" :class="{ 'border-primary': tariff.isPopular }">
          <template #header>
            <div v-if="tariff.isPopular" class="flex items-center justify-between">
              <div class="text-2xl font-semibold leading-none tracking-tight">
                {{ tariff.name }}
              </div>
              <div class="inline-block rounded-lg bg-primary px-3 py-1 text-xs text-primary-foreground">
                Популярный
              </div>
            </div>
            <div v-else class="text-2xl font-semibold leading-none tracking-tight">
              {{ tariff.name }}
            </div>
            <div class="text-sm text-[hsl(var(--muted-foreground))]">
              {{ tariff.description }}
            </div>
            <div class="pt-4">
              <span class="text-4xl font-bold">{{ tariff.price }}₽</span>
              <span class="text-muted-foreground"> / месяц</span>
            </div>
          </template>
          <template #content>
            <ul class="grid gap-2">
              <li v-for="(feature, index) in tariff.features" :key="index" class="flex items-center gap-2">
                <Check class="h-4 w-4 text-primary" />
                <span>{{ feature }}</span>
              </li>
            </ul>
          </template>
          <template #footer>
            <a :href="tariff.link" target="_blank">
              <Button class="w-full">Выбрать план</Button>
            </a>
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>
