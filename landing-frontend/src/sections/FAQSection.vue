<script setup lang="ts">
import { Accordion, Typography } from 'itx-ui-kit'

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
    answer: 'Еженедельные онлайн-встречи: английский клуб, технический книжный клуб, лекции по архитектуре и карьерному росту. Также — оффлайн-митапы, мок-интервью, разбор резюме, обмен вакансиями и реферальными ссылками.',
  },
  {
    title: 'Для кого подходит IT-ХОЗЯЕВА?',
    answer: 'Для IT-специалистов любого грейда и направления: разработчики, тестировщики, аналитики, DevOps, менеджеры. В сообществе участники из Яндекса, Тинькофф, VK и других компаний — от junior до lead.',
  },
]
</script>

<template>
  <section
    id="FAQ"
    class="w-full py-12 md:py-24 lg:py-32"
  >
    <div class="container px-6 md:px-10 space-y-9">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="flex flex-col gap-5">
          <Typography
            variant="h2"
            as="h2"
            class="text-accent"
          >
            Частые вопросы
          </Typography>
          <Typography
            variant="body-xl"
            as="p"
            class="max-w-[540px]"
          >
            Ответы на популярные вопросы о сообществе IT-ХОЗЯЕВА
          </Typography>
        </div>
      </div>
      <div class="grid grid-cols-1 lg:grid-cols-2 items-start gap-5">
        <Accordion
          v-for="question in questions"
          :key="question.title"
          :title="question.title"
          :content="question.answer"
        />
      </div>
    </div>
  </section>
</template>
