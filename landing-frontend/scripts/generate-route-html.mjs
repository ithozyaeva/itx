import { mkdirSync, readFileSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { createJiti } from 'jiti'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)
const distDir = join(__dirname, '..', 'dist')
const indexHtmlPath = join(distDir, 'index.html')
const baseHtml = readFileSync(indexHtmlPath, 'utf-8')

const jiti = createJiti(import.meta.url, { interopDefault: true })
const { articles, renderBlocksToHtml } = await jiti.import('../src/content/articles/index.ts')

// SPA-shell в собранном index.html выглядит так:
//   <div id="app">
//     <style>...</style>
//     <div class="seo-fallback">...</div>
//   </div>
// vite после билда переносит <script type="module"> в <head>, поэтому
// исторический паттерн "<div id=\"app\">...<script type=\"module\"" перестал матчиться,
// и на /mentors, /vibe-coding и т.д. оставался homepage-shell.
const SEO_SHELL_RE = /<div id="app">[\s\S]*?<\/div>\s*<\/div>/

const BASE_URL = 'https://ithozyaeva.ru'

const SHELL_STYLE = `<style>
        .seo-fallback { max-width: 840px; margin: 0 auto; padding: 24px; font: 16px/1.5 system-ui, sans-serif; color: #d4d4d4; background: #0b0d0c; }
        .seo-fallback h1 { font-size: 32px; margin: 0 0 12px; color: #5eead4; }
        .seo-fallback h2 { font-size: 22px; margin: 28px 0 10px; color: #5eead4; }
        .seo-fallback h3 { font-size: 17px; margin: 16px 0 6px; }
        .seo-fallback a { color: #5eead4; }
        .seo-fallback nav { margin-bottom: 16px; font-size: 14px; color: #8a8a8a; }
        .seo-fallback .article-meta { font-size: 12px; color: #8a8a8a; margin-bottom: 8px; text-transform: uppercase; letter-spacing: 0.08em; }
        .seo-fallback .article-note { border-left: 2px solid #5eead4; padding: 12px 16px; background: rgba(94,234,212,0.05); margin: 16px 0; }
        .seo-fallback .article-cta a { display: inline-block; padding: 10px 16px; background: #5eead4; color: #0b0d0c; text-decoration: none; font-weight: 700; text-transform: uppercase; }
        .seo-fallback ul, .seo-fallback ol { padding-left: 24px; }
        .seo-fallback li { margin: 6px 0; }
      </style>`

/**
 * Базовые route — отдельные prerendered HTML.
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
        <p>В <a href="/">IT-ХОЗЯЕВА</a> еженедельно проходят онлайн-воркшопы по vibe coding: разбор реальных задач участников в Cursor и Claude Code, обсуждение промптов и проверка кода, написанного ИИ. Доступ к воркшопам — по подписке от 520 ₽/мес. Подробнее — в <a href="/articles">статьях</a>.</p>
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

const AMP_RE = /&/g
const QUOT_RE = /"/g

function escapeAttr(value) {
  return value.replace(AMP_RE, '&amp;').replace(QUOT_RE, '&quot;')
}

function buildHtml({ canonical, title, description, ogType, jsonLd, shellHtml }) {
  const jsonLdScripts = Array.isArray(jsonLd) ? jsonLd : [jsonLd]
  const jsonLdHtml = jsonLdScripts
    .map(o => `<script type="application/ld+json">\n${JSON.stringify(o, null, 2)}\n    </script>`)
    .join('\n    ')

  return baseHtml
    .replace(TITLE_RE, `<title>${title}</title>`)
    .replace(DESCRIPTION_RE, `<meta name="description" content="${escapeAttr(description)}">`)
    .replace(CANONICAL_RE, `<link rel="canonical" href="${canonical}">`)
    .replace(OG_TYPE_RE, `<meta property="og:type" content="${ogType}">`)
    .replace(OG_TITLE_RE, `<meta property="og:title" content="${escapeAttr(title)}">`)
    .replace(OG_DESCRIPTION_RE, `<meta property="og:description" content="${escapeAttr(description)}">`)
    .replace(OG_URL_RE, `<meta property="og:url" content="${canonical}">`)
    .replace(TWITTER_TITLE_RE, `<meta name="twitter:title" content="${escapeAttr(title)}">`)
    .replace(TWITTER_DESCRIPTION_RE, `<meta name="twitter:description" content="${escapeAttr(description)}">`)
    .replace(HEAD_CLOSE_RE, `${jsonLdHtml}\n  </head>`)
    .replace(SEO_SHELL_RE, `<div id="app">\n      ${shellHtml}\n    </div>`)
}

function buildArticlesIndexShell() {
  const items = articles
    .map(a => `<li><a href="/articles/${a.slug}"><strong>${a.h1}</strong></a> — ${a.excerpt}</li>`)
    .join('\n          ')
  return `${SHELL_STYLE}
      <div class="seo-fallback">
        <nav aria-label="breadcrumb"><a href="/">Главная</a> › Статьи</nav>
        <h1>Статьи IT-ХОЗЯЕВА</h1>
        <p>Практика vibe coding, путь в AI-инжиниринг, выбор ментора, подготовка к собеседованиям и обзор IT-сообществ России — статьи на стыке инструментов и карьеры.</p>
        <ul>
          ${items}
        </ul>
        <p><a href="/">← Вернуться на главную</a></p>
      </div>`
}

function buildArticleShell(article) {
  const body = renderBlocksToHtml(article.body)
  const tagsLine = article.tags?.length ? ` · ${article.tags.join(' · ')}` : ''
  const faqHtml = article.faq?.length
    ? `<h2>Частые вопросы</h2>${article.faq.map(f => `<h3>${f.q}</h3><p>${f.a}</p>`).join('')}`
    : ''
  return `${SHELL_STYLE}
      <div class="seo-fallback">
        <nav aria-label="breadcrumb"><a href="/">Главная</a> › <a href="/articles">Статьи</a> › ${article.breadcrumb}</nav>
        <div class="article-meta">${article.publishedAt}${tagsLine}</div>
        <h1>${article.h1}</h1>
        <p>${article.lead}</p>
        ${body}
        ${faqHtml}
        <p><a href="/articles">← Все статьи</a></p>
      </div>`
}

const articlesIndexRoute = {
  path: 'articles.html',
  canonical: `${BASE_URL}/articles`,
  title: 'Статьи о вайбкодинге, менторстве и карьере в IT | IT-ХОЗЯЕВА',
  description: 'Практические статьи о vibe coding, AI-инжиниринге, менторстве и подготовке к IT-собеседованиям. Опыт участников сообщества IT-ХОЗЯЕВА.',
  ogType: 'website',
  jsonLd: [
    {
      '@context': 'https://schema.org',
      '@type': 'BreadcrumbList',
      'itemListElement': [
        { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': `${BASE_URL}/` },
        { '@type': 'ListItem', 'position': 2, 'name': 'Статьи', 'item': `${BASE_URL}/articles` },
      ],
    },
    {
      '@context': 'https://schema.org',
      '@type': 'CollectionPage',
      'name': 'Статьи IT-ХОЗЯЕВА',
      'url': `${BASE_URL}/articles`,
      'description': 'Статьи о vibe coding, AI-инжиниринге, менторстве и подготовке к IT-собеседованиям.',
      'isPartOf': { '@type': 'WebSite', 'name': 'IT-ХОЗЯЕВА', 'url': BASE_URL },
      'mainEntity': {
        '@type': 'ItemList',
        'itemListElement': articles.map((a, i) => ({
          '@type': 'ListItem',
          'position': i + 1,
          'url': `${BASE_URL}/articles/${a.slug}`,
          'name': a.h1,
        })),
      },
    },
  ],
  shellHtml: buildArticlesIndexShell(),
}

const articleRoutes = articles.map(article => ({
  path: `articles/${article.slug}.html`,
  canonical: `${BASE_URL}/articles/${article.slug}`,
  title: `${article.title}`.includes('IT-ХОЗЯЕВА') ? article.title : `${article.title} | IT-ХОЗЯЕВА`,
  description: article.description,
  ogType: 'article',
  jsonLd: [
    {
      '@context': 'https://schema.org',
      '@type': 'BreadcrumbList',
      'itemListElement': [
        { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': `${BASE_URL}/` },
        { '@type': 'ListItem', 'position': 2, 'name': 'Статьи', 'item': `${BASE_URL}/articles` },
        { '@type': 'ListItem', 'position': 3, 'name': article.breadcrumb, 'item': `${BASE_URL}/articles/${article.slug}` },
      ],
    },
    {
      '@context': 'https://schema.org',
      '@type': 'Article',
      'headline': article.h1,
      'description': article.description,
      'datePublished': article.publishedAt,
      'dateModified': article.updatedAt ?? article.publishedAt,
      'author': { '@type': 'Organization', 'name': 'IT-ХОЗЯЕВА', 'url': BASE_URL },
      'publisher': {
        '@type': 'Organization',
        'name': 'IT-ХОЗЯЕВА',
        'url': BASE_URL,
        'logo': { '@type': 'ImageObject', 'url': `${BASE_URL}/og-image.png` },
      },
      'mainEntityOfPage': `${BASE_URL}/articles/${article.slug}`,
      'inLanguage': 'ru',
    },
    ...(article.faq?.length
      ? [{
          '@context': 'https://schema.org',
          '@type': 'FAQPage',
          'mainEntity': article.faq.map(f => ({
            '@type': 'Question',
            'name': f.q,
            'acceptedAnswer': { '@type': 'Answer', 'text': f.a },
          })),
        }]
      : []),
  ],
  shellHtml: buildArticleShell(article),
}))

const allRoutes = [...routes, articlesIndexRoute, ...articleRoutes]

for (const route of allRoutes) {
  const outPath = join(distDir, route.path)
  mkdirSync(dirname(outPath), { recursive: true })
  writeFileSync(outPath, buildHtml(route), 'utf-8')
  console.log(`✅ Generated: ${outPath}`)
}
