<script setup lang="ts">
import type { Component } from 'vue'
import { useWindowSize } from '@vueuse/core'
import { computed } from 'vue'
import {
  GradesIcon,
  GradesMobileIcon,
  MeetingsIcon,
  MeetingsMobileIcon,
  TopicsIcon,
  TopicsMobileIcon,
  VacanciesIcon,
  VacanciesMobileIcon,
} from '@/assets/icons'

import FeatureItem from '@/components/ui/FeatureItem.vue'
import Typography from '@/components/ui/UiTypography.vue'

interface Feature {
  title: string
  description: string
  icon: Component
}

const { width } = useWindowSize()

const isMobile = computed(() => width.value < 600)

const features: Feature[] = [
  {
    icon: isMobile.value ? TopicsMobileIcon : TopicsIcon,
    title: 'Обсуждаем технологии, карьеру и жизнь',
    description: 'Участники делятся опытом, разбирают кейсы — от вайбкодинга до переговоров по офферу.',
  },
  {
    icon: isMobile.value ? GradesMobileIcon : GradesIcon,
    title: 'Специалисты всех грейдов и направлений',
    description: 'Разработчики, тестировщики, аналитики, DevOps — от junior до lead.',
  },
  {
    icon: isMobile.value ? VacanciesMobileIcon : VacanciesIcon,
    title: 'Вакансии, рефералки и подготовка к собеседованиям',
    description: 'Вакансии от участников, реферальные ссылки, помощь с резюме и мок-интервью.',
  },
  {
    icon: isMobile.value ? MeetingsMobileIcon : MeetingsIcon,
    title: 'Еженедельные встречи и обучение',
    description: 'Английский клуб, книжный клуб, лекции по архитектуре, AI и soft skills.',
  },

]
</script>

<template>
  <section
    id="why"
    class="w-full pt-16 md:pt-24 lg:pt-32"
  >
    <div class="container px-6 md:px-10 flex flex-col gap-10">
      <div class="flex flex-col gap-5 text-center lg:flex-row lg:flex-auto lg:text-start justify-between">
        <Typography
          variant="h2"
          as="h2"
          class="text-accent basis-2/5"
        >
          Зачем вступать в IT-ХОЗЯЕВА?
        </Typography>
        <Typography
          variant="body-xl"
          as="p"
          class="basis-1/2"
        >
          IT-ХОЗЯЕВА — это IT-сообщество 250+ специалистов на стыке технологий и искусственного интеллекта: менторство, вайбкодинг, нетворкинг, подготовка к собеседованиям, обмен вакансиями и еженедельные встречи.
        </Typography>
      </div>
      <div class="mx-auto grid grid-cols-1  md:grid-cols-1 lg:grid-cols-2 gap-5 flex-auto">
        <FeatureItem
          v-for="feature in features"
          :key="feature.title"
          :title="feature.title"
          :description="feature.description"
          :icon="feature.icon"
        />
      </div>
    </div>
  </section>
</template>
