import mdx from '@astrojs/mdx'
import sitemap from '@astrojs/sitemap'
import tailwind from '@astrojs/tailwind'
import { defineConfig } from 'astro/config'

export default defineConfig({
  site: 'https://ithozyaeva.ru',
  output: 'static',
  build: {
    format: 'file',
  },
  integrations: [
    tailwind({ applyBaseStyles: false }),
    mdx(),
    sitemap({
      changefreq: 'weekly',
      lastmod: new Date(),
      serialize(item) {
        // Tune per-route priority/changefreq based on URL pattern.
        if (item.url === 'https://ithozyaeva.ru/') {
          item.priority = 1.0
          item.changefreq = 'weekly'
        }
        else if (item.url === 'https://ithozyaeva.ru/mentors' || item.url === 'https://ithozyaeva.ru/articles') {
          item.priority = 0.9
          item.changefreq = 'weekly'
        }
        else if (item.url.includes('/articles/')) {
          item.priority = 0.7
          item.changefreq = 'monthly'
        }
        else if (item.url === 'https://ithozyaeva.ru/privacy') {
          item.priority = 0.3
          item.changefreq = 'monthly'
        }
        else {
          item.priority = 0.8
          item.changefreq = 'monthly'
        }
        return item
      },
    }),
  ],
})
