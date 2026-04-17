import { readFileSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)
const distDir = join(__dirname, '..', 'dist')
const indexHtmlPath = join(distDir, 'index.html')
const baseHtml = readFileSync(indexHtmlPath, 'utf-8')

const SEO_SHELL_RE = /<div id="app">[\s\S]*?<\/div>\s*<script type="module"/

const BASE_URL = 'https://ithozyaeva.ru'

const SHELL_STYLE = `<style>
        .seo-fallback { max-width: 840px; margin: 0 auto; padding: 24px; font: 16px/1.5 system-ui, sans-serif; color: #d4d4d4; background: #0b0d0c; }
        .seo-fallback h1 { font-size: 32px; margin: 0 0 12px; color: #5eead4; }
        .seo-fallback h2 { font-size: 22px; margin: 28px 0 10px; color: #5eead4; }
        .seo-fallback h3 { font-size: 17px; margin: 16px 0 6px; }
        .seo-fallback a { color: #5eead4; }
        .seo-fallback nav { margin-bottom: 16px; font-size: 14px; color: #8a8a8a; }
      </style>`

/**
 * Каждый route — отдельный prerendered HTML.
 * title / description / canonical / og-теги подменяются в статическом HTML,
 * чтобы бот без JS видел правильную мету.
 * SEO-shell внутри #app — эквивалент реального контента страницы для JS-less ботов.
 * Vue 3 CSR mount замещает содержимое #app при бутe, поэтому пользователь видит настоящий SPA.
 */
const routes = [
  {
    path: 'privacy.html',
    canonical: `${BASE_URL}/privacy`,
    title: 'Политика конфиденциальности | IT-ХОЗЯЕВА',
    description: 'Политика в отношении обработки персональных данных IT-ХОЗЯЕВА. Узнайте, как мы обрабатываем и защищаем ваши персональные данные.',
    ogType: 'article',
    jsonLd: {
      '@context': 'https://schema.org',
      '@type': 'BreadcrumbList',
      'itemListElement': [
        { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': `${BASE_URL}/` },
        { '@type': 'ListItem', 'position': 2, 'name': 'Политика конфиденциальности', 'item': `${BASE_URL}/privacy` },
      ],
    },
    shellHtml: `${SHELL_STYLE}
      <div class="seo-fallback">
        <nav aria-label="breadcrumb"><a href="/">Главная</a> › Политика конфиденциальности</nav>
        <h1>Политика в отношении обработки персональных данных</h1>
        <p>Настоящая политика обработки персональных данных составлена в соответствии с требованиями Федерального закона от 27.07.2006 № 152-ФЗ «О персональных данных» и определяет порядок обработки персональных данных и меры по обеспечению безопасности персональных данных, предпринимаемые Стародубцевым Александром Владимировичем (далее — Оператор).</p>
        <h2>Основные разделы</h2>
        <ul>
          <li>Общие положения и цели обработки персональных данных.</li>
          <li>Правовые основания и категории субъектов персональных данных.</li>
          <li>Порядок сбора, хранения, передачи и прекращения обработки.</li>
          <li>Права субъектов персональных данных.</li>
          <li>Контактные данные оператора для обращений.</li>
        </ul>
        <p>Полный текст доступен на странице после загрузки интерфейса. Для возврата на главную перейдите по ссылке <a href="/">ithozyaeva.ru</a>.</p>
      </div>`,
  },
  {
    path: 'mentors.html',
    canonical: `${BASE_URL}/mentors`,
    title: 'База менторов — senior-разработчики, AI-инженеры, тимлиды | IT-ХОЗЯЕВА',
    description: 'База из 60+ IT-менторов: senior-разработчики, тимлиды, архитекторы и AI-инженеры из Яндекса, Tinkoff, VK. Помогают с карьерой, сменой грейда, переходом в AI и подготовкой к собеседованиям.',
    ogType: 'website',
    jsonLd: {
      '@context': 'https://schema.org',
      '@type': 'BreadcrumbList',
      'itemListElement': [
        { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': `${BASE_URL}/` },
        { '@type': 'ListItem', 'position': 2, 'name': 'Менторы', 'item': `${BASE_URL}/mentors` },
      ],
    },
    shellHtml: `${SHELL_STYLE}
      <div class="seo-fallback">
        <nav aria-label="breadcrumb"><a href="/">Главная</a> › Менторы</nav>
        <h1>База менторов IT-ХОЗЯЕВА</h1>
        <p>60+ IT-менторов: senior-разработчики, тимлиды, архитекторы и AI-инженеры из топовых российских IT-компаний — Яндекс, Tinkoff, VK, Ozon, Сбер, Kaspersky, JetBrains. Помогаем расти по грейду, сменить направление, подготовиться к собеседованию и разобраться в AI и вайбкодинге.</p>
        <h2>Кто в базе</h2>
        <p>Senior и lead разработчики, архитекторы, AI-инженеры, DevOps, QA-лиды. Специалисты, которые реально собеседуют и принимают в топовых командах.</p>
        <h2>С чем помогают</h2>
        <ul>
          <li>Карьерное планирование и смена грейда.</li>
          <li>Ревью резюме и мок-интервью.</li>
          <li>Подготовка к системному дизайну и технической секции.</li>
          <li>Разбор PR, архитектурные консультации.</li>
          <li>Переход в AI-разработку, смена стека.</li>
        </ul>
        <h2>Как получить доступ</h2>
        <p>Доступ к базе менторов открывается на тарифе ХОЗЯИН — <a href="https://boosty.to/jointime/purchase/3150814" rel="noopener">от 1000 ₽/мес через Boosty</a>. Получаете полную таблицу контактов, расписание групповых сессий и возможность забронировать индивидуальную консультацию.</p>
        <p><a href="/">← Вернуться на главную</a></p>
      </div>`,
  },
  {
    path: 'vibe-coding.html',
    canonical: `${BASE_URL}/vibe-coding`,
    title: 'Вайбкодинг (vibe coding) — что это, как научиться и где практиковать | IT-ХОЗЯЕВА',
    description: 'Вайбкодинг — разработка в паре с LLM, где программист описывает результат, а код пишет и правит AI. Разбираем термин, инструменты (Cursor, Claude Code, Windsurf), практики и где освоить в IT-сообществе.',
    ogType: 'article',
    jsonLd: {
      '@context': 'https://schema.org',
      '@type': 'BreadcrumbList',
      'itemListElement': [
        { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': `${BASE_URL}/` },
        { '@type': 'ListItem', 'position': 2, 'name': 'Вайбкодинг', 'item': `${BASE_URL}/vibe-coding` },
      ],
    },
    shellHtml: `${SHELL_STYLE}
      <div class="seo-fallback">
        <nav aria-label="breadcrumb"><a href="/">Главная</a> › Вайбкодинг</nav>
        <h1>Вайбкодинг: что это и как научиться</h1>
        <p>Вайбкодинг (от англ. vibe coding) — стиль разработки, при котором программист описывает намерение и результат естественным языком, а большую часть кода пишет и правит LLM. Термин ввёл Андрей Карпаты в феврале 2025 — «ты не столько пишешь код, сколько задаёшь вайб».</p>
        <h2>Как это работает</h2>
        <p>Разработчик открывает IDE с агентом (Cursor, Claude Code, Windsurf) и даёт задачу в чате. Модель читает структуру репозитория, вносит изменения в несколько файлов, запускает линтер и тесты, показывает diff. Важное отличие от обычного автокомплита — агент действует автономно.</p>
        <h2>Инструменты</h2>
        <ul>
          <li><strong>Cursor</strong> — форк VS Code с Claude/GPT агентом в интерфейсе.</li>
          <li><strong>Claude Code</strong> — официальный CLI от Anthropic с долгими автономными сессиями.</li>
          <li><strong>Windsurf</strong> — агентный IDE от Codeium с многоступенчатым планированием.</li>
          <li><strong>Aider</strong> — open-source CLI, работает с любым API.</li>
          <li><strong>GitHub Copilot Workspace</strong> — агент внутри GitHub для работы по issue → PR.</li>
        </ul>
        <h2>Что меняется в работе</h2>
        <p>Рутина — CRUD, тесты, рефакторинг, перевод кода между языками — уходит к AI почти полностью. Критичным становится умение быстро проверить и отвергнуть неверное решение модели, держать архитектуру в голове и грамотно декомпозировать задачу.</p>
        <h2>Как начать</h2>
        <ol>
          <li>Установить один из агентных IDE (Cursor — самый быстрый старт).</li>
          <li>Выбрать свежую модель — Claude Opus 4.x или Sonnet 4.x.</li>
          <li>Взять реальную задачу, не учебную: багфикс, небольшая фича, миграция.</li>
          <li>Писать промпты с контекстом: куда положить файл, какие соглашения проекта, что не трогать.</li>
          <li>Каждый diff просматривать глазами.</li>
        </ol>
        <h2>Где практиковать в сообществе</h2>
        <p>В <a href="/">IT-ХОЗЯЕВА</a> еженедельно проходят онлайн-воркшопы по vibe coding: разбор реальных задач участников в Cursor и Claude Code, обсуждение промптов, подход к ревью AI-кода. Доступ к воркшопам — по подписке от 520 ₽/мес.</p>
        <p><a href="/">← Вернуться на главную</a></p>
      </div>`,
  },
]

const TITLE_RE = /<title>[^<]*<\/title>/
const DESCRIPTION_RE = /<meta name="description" content="[^"]*">/
const CANONICAL_RE = /<link rel="canonical" href="[^"]*">/
const OG_TYPE_RE = /<meta property="og:type" content="[^"]*">/
const OG_TITLE_RE = /<meta property="og:title" content="[^"]*">/
const OG_DESCRIPTION_RE = /<meta property="og:description" content="[^"]*">/
const OG_URL_RE = /<meta property="og:url" content="[^"]*">/
const TWITTER_TITLE_RE = /<meta name="twitter:title" content="[^"]*">/
const TWITTER_DESCRIPTION_RE = /<meta name="twitter:description" content="[^"]*">/
const HEAD_CLOSE_RE = /<\/head>/

function buildHtml({ canonical, title, description, ogType, jsonLd, shellHtml }) {
  const jsonLdScript = `<script type="application/ld+json">\n${JSON.stringify(jsonLd, null, 2)}\n    </script>\n  </head>`

  return baseHtml
    .replace(TITLE_RE, `<title>${title}</title>`)
    .replace(DESCRIPTION_RE, `<meta name="description" content="${description}">`)
    .replace(CANONICAL_RE, `<link rel="canonical" href="${canonical}">`)
    .replace(OG_TYPE_RE, `<meta property="og:type" content="${ogType}">`)
    .replace(OG_TITLE_RE, `<meta property="og:title" content="${title}">`)
    .replace(OG_DESCRIPTION_RE, `<meta property="og:description" content="${description}">`)
    .replace(OG_URL_RE, `<meta property="og:url" content="${canonical}">`)
    .replace(TWITTER_TITLE_RE, `<meta name="twitter:title" content="${title}">`)
    .replace(TWITTER_DESCRIPTION_RE, `<meta name="twitter:description" content="${description}">`)
    .replace(HEAD_CLOSE_RE, jsonLdScript)
    .replace(SEO_SHELL_RE, `<div id="app">\n      ${shellHtml}\n    </div>\n    <script type="module"`)
}

for (const route of routes) {
  const outPath = join(distDir, route.path)
  writeFileSync(outPath, buildHtml(route), 'utf-8')
  console.log(`✅ Generated: ${outPath}`)
}
