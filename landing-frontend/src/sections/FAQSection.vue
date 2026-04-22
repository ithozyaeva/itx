<script setup lang="ts">
import { ref } from 'vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'

interface Question {
  title: string
  answer: string
}

const questions: Question[] = [
  {
    title: 'Как вступить в сообщество IT-ХОЗЯЕВА?',
    answer: 'Выберите тариф (Бригадир, ХОЗЯИН или KING), оформите подписку на Boosty и получите доступ к закрытой платформе, чатам, базе менторов и всем материалам вашего тарифа. Подписка от 520 ₽/мес.',
  },
  {
    title: 'Можно ли сменить тариф или отменить подписку?',
    answer: 'Да. Тарифный план можно повысить или понизить в любой момент через Boosty. При повышении доплата рассчитывается пропорционально, при понижении — изменения вступают со следующего периода оплаты.',
  },
  {
    title: 'Какие активности проходят в сообществе?',
    answer: 'Еженедельные онлайн-встречи: английский клуб, технический книжный клуб, лекции по архитектуре, AI и вайбкодингу (vibe coding). Воркшопы по искусственному интеллекту, нетворкинг-сессии, мок-интервью, разбор резюме, обмен вакансиями и реферальными ссылками.',
  },
  {
    title: 'Для кого подходит IT-ХОЗЯЕВА?',
    answer: 'Для IT-специалистов любого грейда и направления: разработчики, AI-инженеры, тестировщики, аналитики, DevOps, менеджеры. Для тех, кто интересуется искусственным интеллектом, вайбкодингом и хочет расти через менторство и нетворкинг в IT-комьюнити. Участники из Яндекса, Тинькофф, VK — от junior до lead.',
  },
]

const openIndex = ref<number | null>(0)
function toggle(i: number) {
  openIndex.value = openIndex.value === i ? null : i
}
</script>

<template>
  <section
    id="FAQ"
    class="w-full py-20 md:py-32 lg:py-40"
  >
    <div class="container px-6 md:px-10">
      <div class="reveal">
        <SectionHeader
          index="05"
          path="~/community/faq.txt"
          title="Частые вопросы"
          subtitle="Ответы на популярные вопросы о сообществе IT-ХОЗЯЕВА."
        />
      </div>

      <div class="mt-12 md:mt-16 border border-accent/20 bg-background/60 backdrop-blur-sm reveal">
        <div
          v-for="(q, i) in questions"
          :key="q.title"
          class="border-b border-accent/15 last:border-b-0"
        >
          <button
            type="button"
            class="w-full flex items-start gap-4 md:gap-6 text-left px-5 md:px-8 py-5 md:py-7 group transition-colors hover:bg-accent/[0.04]"
            :aria-expanded="openIndex === i"
            @click="toggle(i)"
          >
            <span
              class="font-mono text-xs md:text-sm mt-1 shrink-0 transition-transform duration-300"
              :class="openIndex === i ? 'text-accent rotate-90' : 'text-accent/60'"
            >
              &gt;
            </span>
            <div class="flex-1 flex items-start justify-between gap-4">
              <div class="flex flex-col gap-1.5">
                <span class="font-mono text-[10px] md:text-xs text-foreground/40 tracking-[0.12em]">
                  Q.{{ String(i + 1).padStart(2, '0') }}
                </span>
                <h3
                  class="font-display uppercase text-base md:text-lg lg:text-xl leading-tight transition-colors"
                  :class="openIndex === i ? 'text-accent' : 'text-foreground group-hover:text-accent'"
                >
                  {{ q.title }}
                </h3>
              </div>
              <span
                class="font-mono text-xs shrink-0 mt-1"
                :class="openIndex === i ? 'text-accent' : 'text-foreground/40'"
              >
                {{ openIndex === i ? '[−]' : '[+]' }}
              </span>
            </div>
          </button>
          <div
            class="grid transition-[grid-template-rows] duration-300"
            :style="{ gridTemplateRows: openIndex === i ? '1fr' : '0fr' }"
          >
            <div class="overflow-hidden">
              <div class="px-5 md:px-8 pb-6 md:pb-8 pl-14 md:pl-[66px]">
                <div class="font-mono text-[11px] text-accent/60 mb-2">
                  &gt; answer:
                </div>
                <p class="text-sm md:text-base text-foreground/75 leading-relaxed max-w-3xl">
                  {{ q.answer }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
