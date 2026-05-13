import { writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { createJiti } from 'jiti'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const jiti = createJiti(import.meta.url, { interopDefault: true })
const { articles } = await jiti.import('../src/content/articles/index.ts')

const BASE_URL = 'https://ithozyaeva.ru'
const today = new Date().toISOString().split('T')[0]

const staticRoutes = [
  { path: '/', priority: '1.0', changefreq: 'weekly', lastmod: today },
  { path: '/mentors', priority: '0.9', changefreq: 'weekly', lastmod: today },
  { path: '/vibe-coding', priority: '0.8', changefreq: 'monthly', lastmod: today },
  { path: '/articles', priority: '0.8', changefreq: 'weekly', lastmod: today },
  { path: '/privacy', priority: '0.3', changefreq: 'monthly', lastmod: today },
]

const articleRoutes = articles.map(a => ({
  path: `/articles/${a.slug}`,
  priority: '0.7',
  changefreq: 'monthly',
  lastmod: a.updatedAt ?? a.publishedAt,
}))

const routes = [...staticRoutes, ...articleRoutes]

const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${routes.map(route => `  <url>
    <loc>${BASE_URL}${route.path}</loc>
    <lastmod>${route.lastmod}</lastmod>
    <changefreq>${route.changefreq}</changefreq>
    <priority>${route.priority}</priority>
  </url>`).join('\n')}
</urlset>
`

const outputPath = join(__dirname, '../public/sitemap.xml')
writeFileSync(outputPath, sitemap, 'utf-8')
console.log('✅ Sitemap generated:', outputPath)
