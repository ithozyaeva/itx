<script setup lang="ts">
import StructuredData from '@/components/StructuredData.vue'
import Button from '@/components/ui/UiButton.vue'
import { usePageMeta } from '@/composables/useMeta'
import { useRevealObserver } from '@/composables/useRevealObserver'

useRevealObserver()

usePageMeta({
  title: 'Вайбкодинг (vibe coding) — что это, как научиться и где практиковать',
  description: 'Вайбкодинг — разработка в паре с языковой ИИ-моделью (LLM), где программист описывает результат, а код пишет и правит ИИ. Разбираем термин, инструменты (Cursor, Claude Code, Windsurf), практики и где освоить в ИТ-сообществе.',
  url: 'https://ithozyaeva.ru/vibe-coding',
  type: 'article',
})

const structuredData = [
  {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    'itemListElement': [
      { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': 'https://ithozyaeva.ru/' },
      { '@type': 'ListItem', 'position': 2, 'name': 'Вайбкодинг', 'item': 'https://ithozyaeva.ru/vibe-coding' },
    ],
  },
  {
    '@context': 'https://schema.org',
    '@type': 'Article',
    'headline': 'Вайбкодинг (vibe coding) — что это, как научиться и где практиковать',
    'description': 'Разбор термина «вайбкодинг» (vibe coding), инструменты и практики разработки в паре с ИИ.',
    'author': { '@type': 'Organization', 'name': 'IT-ХОЗЯЕВА', 'url': 'https://ithozyaeva.ru' },
    'publisher': {
      '@type': 'Organization',
      'name': 'IT-ХОЗЯЕВА',
      'url': 'https://ithozyaeva.ru',
      'logo': { '@type': 'ImageObject', 'url': 'https://ithozyaeva.ru/og-image.png' },
    },
    'mainEntityOfPage': 'https://ithozyaeva.ru/vibe-coding',
    'inLanguage': 'ru',
  },
  {
    '@context': 'https://schema.org',
    '@type': 'FAQPage',
    'mainEntity': [
      {
        '@type': 'Question',
        'name': 'Что такое вайбкодинг (vibe coding)?',
        'acceptedAnswer': {
          '@type': 'Answer',
          'text': 'Вайбкодинг — стиль разработки, при котором программист описывает намерение и результат естественным языком, а большую часть кода пишет и правит LLM (Claude, GPT). Термин ввёл Андрей Карпаты в 2025 — «ты не столько пишешь код, сколько задаёшь вайб».',
        },
      },
      {
        '@type': 'Question',
        'name': 'Какие инструменты используют для вайбкодинга?',
        'acceptedAnswer': {
          '@type': 'Answer',
          'text': 'Cursor, Claude Code, Windsurf, GitHub Copilot, Aider — IDE-агенты с автономной правкой файлов, запуском тестов и разрешением ошибок. Ключ — модель класса Claude 4 / GPT-5 и чёткие промпты с контекстом репозитория.',
        },
      },
      {
        '@type': 'Question',
        'name': 'Заменит ли вайбкодинг разработчиков?',
        'acceptedAnswer': {
          '@type': 'Answer',
          'text': 'Нет — он меняет распределение усилий. Рутина (CRUD, тесты, рефакторинг) уходит к AI, а критичными становятся архитектура, ревью кода, интеграции и способность быстро проверить и отвергнуть неверное решение модели.',
        },
      },
      {
        '@type': 'Question',
        'name': 'Где практиковать вайбкодинг в России?',
        'acceptedAnswer': {
          '@type': 'Answer',
          'text': 'В сообществе IT-ХОЗЯЕВА проходят еженедельные онлайн-воркшопы по vibe coding: разбор реальных задач в Cursor и Claude Code, обсуждение промптов, подход к ревью AI-кода. Доступ — по подписке от 520 ₽/мес.',
        },
      },
    ],
  },
]
</script>

<template>
  <StructuredData :data="structuredData" />
  <article class="w-full pt-10 md:pt-16 lg:pt-20 pb-20 md:pb-32">
    <div class="container px-6 md:px-10">
      <nav
        class="font-mono text-xs text-foreground/50 mb-6"
        aria-label="breadcrumb"
      >
        <a
          href="/"
          class="hover:text-accent transition-colors"
        >~/</a>
        <span class="mx-2 text-accent/40">/</span>
        <span class="text-foreground">vibe-coding</span>
      </nav>

      <header class="max-w-3xl">
        <h1 class="font-display uppercase text-[28px] sm:text-[40px] md:text-[56px] lg:text-[64px] leading-[0.95] tracking-tight text-accent">
          Вайбкодинг: что это и как научиться
        </h1>
        <p class="mt-6 md:mt-8 text-base md:text-lg text-foreground/75 leading-relaxed">
          Вайбкодинг (от англ. <em>vibe coding</em>) — стиль разработки, при котором программист описывает намерение
          и результат естественным языком, а большую часть кода пишет и правит LLM. Термин ввёл Андрей Карпаты
          в феврале 2025: «ты не столько пишешь код, сколько задаёшь вайб — говоришь, чего хочешь, и смотришь, что получилось».
        </p>
      </header>

      <div class="mt-14 md:mt-20 max-w-3xl space-y-14 md:space-y-20">
        <section>
          <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-4">
            Как это работает
          </h2>
          <p class="text-base md:text-lg text-foreground/75 leading-relaxed">
            Разработчик открывает IDE с агентом (Cursor, Claude Code, Windsurf) и даёт задачу в чате — «добавь фильтр
            по городу на странице списка вакансий, обнови тесты». Модель читает структуру репозитория, вносит изменения
            в несколько файлов, запускает линтер и тесты, показывает diff. Если что-то не так — программист пишет
            уточнение или правит сам.
          </p>
          <p class="mt-4 text-base md:text-lg text-foreground/75 leading-relaxed">
            Важное отличие от обычного автокомплита: агент действует <strong class="text-foreground">автономно</strong> —
            сам решает, какие файлы трогать, какие команды запустить, как проверить результат.
          </p>
        </section>

        <section>
          <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-4">
            Инструменты
          </h2>
          <ul class="space-y-3 text-base md:text-lg text-foreground/75 leading-relaxed">
            <li><strong class="text-foreground">Cursor</strong> — форк VS Code с Claude/GPT агентом в интерфейсе. Подходит для большинства задач, хорошо понимает структуру монорепо.</li>
            <li><strong class="text-foreground">Claude Code</strong> — официальный CLI от Anthropic. Умеет долгие автономные сессии, управление через файловые хуки и MCP-серверы.</li>
            <li><strong class="text-foreground">Windsurf</strong> — альтернативный агентный IDE от Codeium с акцентом на многоступенчатое планирование.</li>
            <li><strong class="text-foreground">Aider</strong> — open-source CLI-клиент, работает с локальными моделями или через API любого провайдера.</li>
            <li><strong class="text-foreground">GitHub Copilot Workspace</strong> — агент внутри GitHub, удобно для работы по issue → PR.</li>
          </ul>
        </section>

        <section>
          <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-4">
            Что меняется в работе
          </h2>
          <p class="text-base md:text-lg text-foreground/75 leading-relaxed">
            Рутина — CRUD, обвязка, тесты, рефакторинг, перевод кода между языками — уходит к AI почти полностью.
            Время освобождается, но критичным становится умение <strong class="text-foreground">быстро проверить и отвергнуть</strong>
            неверное решение модели, держать архитектуру в голове и грамотно декомпозировать задачу.
          </p>
          <p class="mt-4 text-base md:text-lg text-foreground/75 leading-relaxed">
            Разработчик, который раньше писал код 80% времени, теперь 80% времени <em>читает</em> код — свой,
            командный, AI-генерированный.
          </p>
        </section>

        <section>
          <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-4">
            Как начать
          </h2>
          <ol class="space-y-3 text-base md:text-lg text-foreground/75 leading-relaxed list-decimal pl-6">
            <li>Установить один из агентных IDE (Cursor — самый быстрый старт).</li>
            <li>Выбрать свежую модель — Claude Opus 4.x или Sonnet 4.x через API или подписку.</li>
            <li>Взять <em>реальную</em> задачу, не учебную: багфикс, небольшая фича, миграция.</li>
            <li>Писать промпты с контекстом: куда положить файл, какие соглашения проекта, что точно не трогать.</li>
            <li>Каждый diff просматривать глазами. На этом этапе все срезанные углы и вылезут.</li>
          </ol>
        </section>

        <section>
          <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-4">
            Где практиковать в сообществе
          </h2>
          <p class="text-base md:text-lg text-foreground/75 leading-relaxed">
            В <a
              href="/"
              class="text-accent hover:underline"
            >IT-ХОЗЯЕВА</a> еженедельно проходят онлайн-воркшопы по vibe coding: разбор реальных задач участников
            в Cursor и Claude Code, обсуждение промптов, подход к ревью AI-кода, закрытые AI-беседы
            и практические сессии. Доступ к воркшопам — по подписке от 520 ₽/мес.
          </p>
          <div class="mt-8">
            <Button
              as="a"
              href="/"
              variant="filled"
            >
              Вступить в сообщество
            </Button>
          </div>
        </section>

        <section>
          <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-4">
            Частые вопросы
          </h2>
          <div class="space-y-6">
            <div>
              <h3 class="font-display uppercase text-lg text-foreground mb-2">
                Заменит ли вайбкодинг разработчиков?
              </h3>
              <p class="text-base text-foreground/70 leading-relaxed">
                Нет — он меняет распределение усилий. Рутина уходит к AI, критичными становятся архитектура, ревью
                и способность быстро отвергнуть неверное решение модели.
              </p>
            </div>
            <div>
              <h3 class="font-display uppercase text-lg text-foreground mb-2">
                Нужно ли знать программирование, чтобы вайбкодить?
              </h3>
              <p class="text-base text-foreground/70 leading-relaxed">
                Для простых скриптов и прототипов — нет. Для продакшна нужна экспертиза: вы должны уметь прочитать
                diff и понять, что именно делает модель, иначе быстро накопите технический долг.
              </p>
            </div>
            <div>
              <h3 class="font-display uppercase text-lg text-foreground mb-2">
                Какая модель лучше для vibe coding?
              </h3>
              <p class="text-base text-foreground/70 leading-relaxed">
                На 2026 год для долгих агентных сессий — Claude Opus / Sonnet 4.x. Для быстрых правок и автокомплита
                — GPT или локальные open-source модели поменьше.
              </p>
            </div>
          </div>
        </section>
      </div>
    </div>
  </article>
</template>
