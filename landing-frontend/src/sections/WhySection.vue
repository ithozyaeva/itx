<script setup lang="ts">
import { Braces, Network, Presentation, Target } from 'lucide-vue-next'
import SectionHeader from '@/components/ui/SectionHeader.vue'

interface Feature {
  id: string
  tag: string
  title: string
  description: string
  icon: any
}

const features: Feature[] = [
  {
    id: '01',
    tag: 'topics',
    title: 'Обсуждаем технологии, карьеру и жизнь',
    description: 'Участники делятся опытом, разбирают рабочие случаи — от вайбкодинга до переговоров о предложении работы.',
    icon: Braces,
  },
  {
    id: '02',
    tag: 'grades',
    title: 'Специалисты всех уровней и направлений',
    description: 'Разработчики, тестировщики, аналитики, DevOps-инженеры — от начинающих до руководителей.',
    icon: Target,
  },
  {
    id: '03',
    tag: 'jobs',
    title: 'Вакансии, реферальные ссылки и подготовка к собеседованиям',
    description: 'Вакансии от участников, реферальные приглашения, помощь с резюме и тренировочные собеседования.',
    icon: Network,
  },
  {
    id: '04',
    tag: 'meetings',
    title: 'Еженедельные встречи и обучение',
    description: 'Английский клуб, книжный клуб, лекции по архитектуре, искусственному интеллекту и гибким навыкам.',
    icon: Presentation,
  },
]
</script>

<template>
  <section
    id="why"
    class="w-full pt-20 md:pt-32 lg:pt-40"
  >
    <div class="container px-6 md:px-10">
      <div class="reveal">
        <SectionHeader
          index="01"
          path="~/community/why.md"
          title="Зачем вступать?"
          subtitle="Профессиональное сообщество из 250+ специалистов из самый разных технических сфер. Менторская поддержка, рефералки, вайбкодинг-практики, живой нетворкинг, стресс-полигон для собеседований, закрытые вакансии и регулярные офлайн/онлайн-встречи."
        />
      </div>

      <div class="mt-12 md:mt-16 grid grid-cols-1 md:grid-cols-2 gap-px bg-accent/15 reveal-stagger">
        <article
          v-for="(feature, i) in features"
          :key="feature.id"
          class="reveal-child trace-border relative group p-7 md:p-10 bg-background/95 hover:bg-accent/5 transition-colors duration-300 min-h-[280px] flex flex-col justify-between cursor-default"
        >
          <!-- Tag row -->
          <div class="flex items-start justify-between gap-4">
            <div class="flex items-center gap-3">
              <span class="font-mono text-xs text-accent/60 tracking-[0.12em]">
                [{{ feature.id }}]
              </span>
              <span class="font-mono text-xs text-foreground/40 tracking-[0.08em]">
                &lt;{{ feature.tag }}/&gt;
              </span>
            </div>
            <div class="shrink-0 w-12 h-12 md:w-14 md:h-14 flex items-center justify-center border border-accent/30 text-accent group-hover:border-accent group-hover:bg-accent/10 transition-all duration-300">
              <component
                :is="feature.icon"
                class="w-5 h-5 md:w-6 md:h-6"
                stroke-width="1.5"
              />
            </div>
          </div>

          <div class="mt-8 md:mt-12">
            <h3 class="font-display uppercase text-xl md:text-2xl lg:text-[28px] leading-[1.15] text-foreground group-hover:text-accent transition-colors duration-300">
              {{ feature.title }}
            </h3>
            <p class="mt-3 text-sm md:text-base text-foreground/60 leading-relaxed max-w-lg">
              {{ feature.description }}
            </p>
          </div>

          <!-- Corner accents -->
          <span class="absolute top-0 left-0 w-3 h-3 border-t-2 border-l-2 border-accent opacity-0 group-hover:opacity-100 transition-opacity" />
          <span class="absolute bottom-0 right-0 w-3 h-3 border-b-2 border-r-2 border-accent opacity-0 group-hover:opacity-100 transition-opacity" />

          <!-- Index watermark -->
          <span
            aria-hidden="true"
            class="pointer-events-none absolute -bottom-4 -right-2 font-display text-[120px] md:text-[180px] leading-none text-accent/[0.03] group-hover:text-accent/[0.06] transition-colors select-none"
          >
            {{ i + 1 }}
          </span>
        </article>
      </div>
    </div>
  </section>
</template>
