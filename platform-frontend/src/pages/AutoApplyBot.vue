<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { Bot, ExternalLink } from 'lucide-vue-next'
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

    <div class="bg-card rounded-3xl border p-6 md:p-8 space-y-6">
      <div class="flex items-center gap-4">
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

      <div class="space-y-3 text-sm leading-relaxed">
        <p>
          Всем участникам сообщества IT-X доступен бот для автоматических откликов на вакансии на hh.ru.
          Бот поможет вам сэкономить время на рутинном поиске работы и автоматически откликаться на подходящие вакансии.
        </p>
        <p class="text-muted-foreground">
          Просто перейдите по ссылке ниже и начните пользоваться ботом прямо сейчас.
        </p>
      </div>

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
  </div>
</template>
