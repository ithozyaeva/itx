<script setup lang="ts">
import Card from '@/components/ui/Card.vue'

interface Question {
  title: string
  answer: string
}

const questions: Question[] = [
  {
    title: 'Как присоединиться к сообществу?',
    answer: ' Для присоединения к сообществу IT-ХОЗЯЕВА необходимо выбрать подходящий тарифный план, заполнить форму регистрации и оплатить подписку. После этого вы получите доступ к закрытой платформе и всем материалам согласно выбранному плану.',
  },
  {
    title: 'Могу ли я сменить тарифный план?',
    answer: ' Да, вы можете изменить свой тарифный план в любое время. При переходе на более дорогой план доплата рассчитывается пропорционально оставшемуся времени подписки. При переходе на более дешевый план изменения вступят в силу со следующего периода оплаты.',
  },
  {
    title: 'Какие активности доступны в сообществе?',
    answer: ' В сообществе IT-ХОЗЯЕВА доступны различные активности: общение, регулярные оффлайн встречи, онлайн встречи технического книжного клуба, киберспортивные турниры, совместное обсуждение новых IT-трендов и обучение.',
  },
]
</script>

<template>
  <section id="FAQ" class="w-full py-12 md:py-24 lg:py-32">
    <div class="container px-4 md:px-6">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="space-y-2">
          <div class="inline-block rounded-lg bg-primary px-3 py-1 text-sm text-primary-foreground">
            FAQ
          </div>
          <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">
            Часто задаваемые вопросы
          </h2>
          <p class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
            Ответы на популярные вопросы о нашем сообществе
          </p>
        </div>
      </div>
      <div class="mx-auto grid max-w-5xl gap-6 py-12">
        <Card v-for="question in questions" :key="question.answer">
          <template #header>
            <div class="text-2xl font-semibold leading-none tracking-tight">
              {{ question.title }}
            </div>
          </template>
          <template #content>
            <p class="text-muted-foreground">
              {{ question.answer }}
            </p>
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>
