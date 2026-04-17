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
    shellHtml: `<style>
        .seo-fallback { max-width: 840px; margin: 0 auto; padding: 24px; font: 16px/1.5 system-ui, sans-serif; color: #d4d4d4; background: #0b0d0c; }
        .seo-fallback h1 { font-size: 32px; margin: 0 0 12px; color: #5eead4; }
        .seo-fallback h2 { font-size: 22px; margin: 28px 0 10px; color: #5eead4; }
        .seo-fallback a { color: #5eead4; }
        .seo-fallback nav { margin-bottom: 16px; font-size: 14px; color: #8a8a8a; }
      </style>
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
