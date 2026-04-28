<script setup lang="ts">
import { ChevronDown } from 'lucide-vue-next'
import { AccordionContent, AccordionHeader, AccordionItem, AccordionRoot, AccordionTrigger } from 'reka-ui'
import { Typography } from '@/components/ui/typography'

const sections = [
  {
    title: 'Что такое IT-X платформа?',
    content: 'IT-X — это закрытое сообщество ИТ-специалистов. Здесь можно участвовать в событиях, находить менторов, обмениваться опытом, зарабатывать баллы за активность и многое другое.',
  },
  {
    title: 'Как работает реферальная система?',
    content: 'Перейдите в раздел «Рефералки» через ваш профиль. Скопируйте вашу реферальную ссылку и поделитесь ей. Когда человек присоединится по вашей ссылке, вы получите баллы за создание реферала, а при его активности — дополнительные баллы за конверсию.',
  },
  {
    title: 'Как разместить резюме?',
    content: 'Перейдите в раздел «Резюме» через ваш профиль и загрузите PDF-файл с вашим резюме. Другие участники смогут его увидеть. Менторы и рекрутеры сообщества могут порекомендовать вас — это называется «зарефать».',
  },
  {
    title: 'Как зарабатывать баллы?',
    content: 'Баллы начисляются за различную активность: посещение событий, создание рефералов, заполнение профиля, еженедельную активность в чатах, выполнение заданий на бирже. Полный список — в разделе «Мои баллы».',
  },
  {
    title: 'Что такое события и как участвовать?',
    content: 'В разделе «События» публикуются мероприятия сообщества: созвоны, воркшопы, лекции. Нажмите «Участвую!» чтобы записаться. Можно добавить событие в Google Calendar или скачать ICS-файл.',
  },
  {
    title: 'Как найти ментора?',
    content: 'Перейдите в раздел «Менторы» — там список экспертов с описанием их опыта, навыков и услуг. Вы можете связаться с ментором через Telegram для консультации.',
  },
  {
    title: 'Что такое Барахолка?',
    content: 'Барахолка — это маркетплейс внутри сообщества. Можно размещать объявления о продаже техники, книг, курсов и других товаров, актуальных для ИТ-специалистов.',
  },
  {
    title: 'Что такое задания и квесты?',
    content: 'На «Бирже заданий» публикуются задачи от участников сообщества, за выполнение которых начисляются баллы. «Задания в чатах» — это автоматические квесты за активность в Telegram-чатах.',
  },
]
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-3xl">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/faq
    </div>
    <Typography
      variant="h2"
      as="h1"
      class="mb-6"
    >
      Часто задаваемые вопросы
    </Typography>

    <AccordionRoot
      type="multiple"
      class="space-y-3"
    >
      <AccordionItem
        v-for="(item, index) in sections"
        :key="index"
        :value="`item-${index}`"
        class="rounded-sm border bg-card overflow-hidden"
      >
        <AccordionHeader>
          <AccordionTrigger class="flex w-full items-center justify-between px-5 py-4 text-left text-sm font-medium hover:bg-muted/50 transition-colors group">
            {{ item.title }}
            <ChevronDown class="h-4 w-4 shrink-0 text-muted-foreground transition-transform duration-200 group-data-[state=open]:rotate-180" />
          </AccordionTrigger>
        </AccordionHeader>
        <AccordionContent class="overflow-hidden text-sm data-[state=closed]:animate-accordion-up data-[state=open]:animate-accordion-down">
          <div class="px-5 pb-4 text-muted-foreground leading-relaxed break-words">
            {{ item.content }}
          </div>
        </AccordionContent>
      </AccordionItem>
    </AccordionRoot>
  </div>
</template>
