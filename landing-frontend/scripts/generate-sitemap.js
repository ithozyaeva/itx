import { writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const BASE_URL = 'https://ithozyaeva.ru'
const routes = [
  { path: '/', priority: '1.0', changefreq: 'weekly' },
  { path: '/privacy', priority: '0.5', changefreq: 'monthly' },
]

const today = new Date().toISOString().split('T')[0]

const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${routes.map(route => `  <url>
    <loc>${BASE_URL}${route.path}</loc>
    <lastmod>${today}</lastmod>
    <changefreq>${route.changefreq}</changefreq>
    <priority>${route.priority}</priority>
  </url>`).join('\n')}
</urlset>
`

const outputPath = join(__dirname, '../public/sitemap.xml')
writeFileSync(outputPath, sitemap, 'utf-8')
console.log('✅ Sitemap generated:', outputPath)
