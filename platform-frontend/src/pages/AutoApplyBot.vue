<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import {
  Bot,
  CheckCircle,
  Clock,
  ExternalLink,
  Filter,
  MessageSquare,
  Rocket,
  Settings,
  Zap,
} from 'lucide-vue-next'
import { ref } from 'vue'

const faqItems = ref([
  {
    question: 'Бот бесплатный?',
    answer: 'Да, бот полностью бесплатен для всех участников сообщества IT-X.',
    open: false,
  },
  {
    question: 'Нужно ли предоставлять логин и пароль от hh.ru?',
    answer: 'Нет, бот работает через API hh.ru и авторизуется через токен. Ваши учётные данные в безопасности.',
    open: false,
  },
  {
    question: 'Можно ли настроить фильтры вакансий?',
    answer: 'Да, вы можете указать желаемую должность, зарплату, город, формат работы и другие параметры.',
    open: false,
  },
  {
    question: 'Как часто бот проверяет новые вакансии?',
    answer: 'Бот проверяет новые вакансии каждые 30 минут и автоматически откликается на подходящие.',
    open: false,
  },
])

const steps = [
  { icon: MessageSquare, title: 'Откройте бота', description: 'Перейдите в Telegram и запустите бота' },
  { icon: Settings, title: 'Настройте профиль', description: 'Укажите желаемую должность, зарплату и город' },
  { icon: Filter, title: 'Задайте фильтры', description: 'Выберите формат работы, опыт и другие параметры' },
  { icon: Rocket, title: 'Запустите автоотклик', description: 'Бот будет откликаться на подходящие вакансии за вас' },
]

const features = [
  { icon: Zap, title: 'Автоматические отклики', description: 'Бот откликается на новые вакансии без вашего участия' },
  { icon: Filter, title: 'Умные фильтры', description: 'Точная настройка параметров поиска вакансий' },
  { icon: Clock, title: 'Работает 24/7', description: 'Не пропустите ни одной подходящей вакансии' },
  { icon: CheckCircle, title: 'Сопроводительные письма', description: 'Генерация персонализированных откликов' },
]

function toggleFaq(index: number) {
  faqItems.value[index].open = !faqItems.value[index].open
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-3xl">
    <Typography
      variant="h2"
      as="h1"
      class="mb-6"
    >
      Бот для автооткликов
    </Typography>

    <!-- Hero card -->
    <div class="bg-card rounded-3xl border p-6 md:p-8 mb-6">
      <div class="flex items-center gap-4 mb-4">
        <div class="flex items-center justify-center w-14 h-14 rounded-2xl bg-primary/10 text-primary shrink-0">
          <Bot :size="28" />
        </div>
        <div>
          <h3 class="text-lg font-semibold">
            Roaster Resume Bot
          </h3>
          <p class="text-sm text-muted-foreground">
            Автоматические отклики на вакансии hh.ru
          </p>
        </div>
      </div>

      <p class="text-sm leading-relaxed mb-6">
        Бот помогает участникам IT-X экономить время на рутинном поиске работы.
        Настройте параметры один раз — и бот будет автоматически откликаться на подходящие вакансии.
      </p>

      <a
        href="https://t.me/roaster_resume_bot"
        target="_blank"
        rel="noopener noreferrer"
        class="inline-flex items-center gap-2 rounded-xl bg-primary text-primary-foreground px-5 py-2.5 text-sm font-medium hover:bg-primary/90 transition-colors"
      >
        Открыть бота в Telegram
        <ExternalLink :size="14" />
      </a>
    </div>

    <!-- How it works -->
    <div class="mb-6">
      <h2 class="text-base font-semibold mb-4">
        Как это работает
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div
          v-for="(step, index) in steps"
          :key="index"
          class="flex gap-3 rounded-2xl border bg-card border-border p-4"
        >
          <div class="flex items-center justify-center w-10 h-10 rounded-xl bg-primary/10 text-primary shrink-0">
            <span class="text-sm font-bold">{{ index + 1 }}</span>
          </div>
          <div>
            <h4 class="text-sm font-medium mb-0.5">
              {{ step.title }}
            </h4>
            <p class="text-xs text-muted-foreground">
              {{ step.description }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Features -->
    <div class="mb-6">
      <h2 class="text-base font-semibold mb-4">
        Возможности
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div
          v-for="feature in features"
          :key="feature.title"
          class="flex gap-3 rounded-2xl border bg-card border-border p-4"
        >
          <div class="flex items-center justify-center w-10 h-10 rounded-xl bg-muted shrink-0">
            <component
              :is="feature.icon"
              class="h-5 w-5 text-muted-foreground"
            />
          </div>
          <div>
            <h4 class="text-sm font-medium mb-0.5">
              {{ feature.title }}
            </h4>
            <p class="text-xs text-muted-foreground">
              {{ feature.description }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- FAQ -->
    <div class="mb-6">
      <h2 class="text-base font-semibold mb-4">
        Частые вопросы
      </h2>
      <div class="space-y-2">
        <div
          v-for="(item, index) in faqItems"
          :key="index"
          class="rounded-2xl border bg-card border-border overflow-hidden"
        >
          <button
            class="w-full flex items-center justify-between p-4 text-left text-sm font-medium hover:bg-muted/50 transition-colors"
            :aria-expanded="item.open"
            @click="toggleFaq(index)"
          >
            {{ item.question }}
            <span
              class="text-muted-foreground transition-transform duration-200"
              :class="{ 'rotate-180': item.open }"
            >
              &#9662;
            </span>
          </button>
          <div
            v-if="item.open"
            class="px-4 pb-4 text-sm text-muted-foreground"
          >
            {{ item.answer }}
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom CTA -->
    <div class="text-center py-6">
      <a
        href="https://t.me/roaster_resume_bot"
        target="_blank"
        rel="noopener noreferrer"
        class="inline-flex items-center gap-2 rounded-xl bg-primary text-primary-foreground px-6 py-3 text-sm font-medium hover:bg-primary/90 transition-colors"
      >
        <Bot :size="18" />
        Начать пользоваться ботом
        <ExternalLink :size="14" />
      </a>
    </div>
  </div>
</template>
