import { useHead, useSeoMeta } from '@unhead/vue'

const SITE_BASE = {
  name: 'IT-ХОЗЯЕВА',
  description: 'Закрытое сообщество IT-специалистов, в котором ценят опыт, юмор и профессионализм',
  url: 'https://ithozyaeva.ru',
  image: 'https://ithozyaeva.ru/og-image.png',
  locale: 'ru_RU',
} as const

export interface PageMeta {
  title?: string
  description?: string
  keywords?: string
  image?: string
  url?: string
  type?: 'website' | 'article'
  noindex?: boolean
  nofollow?: boolean
}

export function usePageMeta(meta: PageMeta) {
  const title = meta.title
    ? `${meta.title} | ${SITE_BASE.name}`
    : SITE_BASE.name

  const description = meta.description || SITE_BASE.description
  const keywords = meta.keywords || 'IT сообщество, программисты, разработчики, IT специалисты, закрытое сообщество, менторство, карьера в IT'
  const image = meta.image || SITE_BASE.image
  const url = meta.url || SITE_BASE.url
  const type = meta.type || 'website'

  const robots = meta.noindex || meta.nofollow
    ? `${meta.noindex ? 'noindex' : 'index'}, ${meta.nofollow ? 'nofollow' : 'follow'}`
    : 'index, follow'

  useSeoMeta({
    title,
    description,
    ogTitle: title,
    ogDescription: description,
    ogImage: image,
    ogUrl: url,
    ogType: type,
    ogLocale: SITE_BASE.locale,
    ogSiteName: SITE_BASE.name,
    twitterCard: 'summary_large_image',
    twitterTitle: title,
    twitterDescription: description,
    twitterImage: image,
    robots,
  })

  useHead({
    meta: [
      { name: 'keywords', content: keywords },
      { name: 'author', content: SITE_BASE.name },
      { name: 'ai:description', content: description },
      { name: 'ai:keywords', content: keywords },
      { name: 'theme-color', content: '#000000' },
      { name: 'apple-mobile-web-app-title', content: SITE_BASE.name },
    ],
    link: [
      { rel: 'canonical', href: url },
    ],
  })
}
