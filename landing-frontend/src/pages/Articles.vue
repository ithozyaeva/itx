<script setup lang="ts">
import StructuredData from '@/components/StructuredData.vue'
import { usePageMeta } from '@/composables/useMeta'
import { useRevealObserver } from '@/composables/useRevealObserver'
import { articles } from '@/content/articles'

useRevealObserver()

usePageMeta({
  title: 'Статьи о вайбкодинге, менторстве и карьере в IT',
  description: 'Практические статьи о vibe coding, AI-инжиниринге, менторстве и подготовке к IT-собеседованиям. Опыт участников сообщества IT-ХОЗЯЕВА.',
  url: 'https://ithozyaeva.ru/articles',
  type: 'website',
})

const structuredData = [
  {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    'itemListElement': [
      { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': 'https://ithozyaeva.ru/' },
      { '@type': 'ListItem', 'position': 2, 'name': 'Статьи', 'item': 'https://ithozyaeva.ru/articles' },
    ],
  },
  {
    '@context': 'https://schema.org',
    '@type': 'CollectionPage',
    'name': 'Статьи IT-ХОЗЯЕВА',
    'url': 'https://ithozyaeva.ru/articles',
    'description': 'Статьи о vibe coding, AI-инжиниринге, менторстве и подготовке к IT-собеседованиям.',
    'isPartOf': { '@type': 'WebSite', 'name': 'IT-ХОЗЯЕВА', 'url': 'https://ithozyaeva.ru' },
    'mainEntity': {
      '@type': 'ItemList',
      'itemListElement': articles.map((a, i) => ({
        '@type': 'ListItem',
        'position': i + 1,
        'url': `https://ithozyaeva.ru/articles/${a.slug}`,
        'name': a.h1,
      })),
    },
  },
]
</script>

<template>
  <StructuredData :data="structuredData" />
  <section class="w-full pt-10 md:pt-16 lg:pt-20 pb-20 md:pb-32">
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
        <span class="text-foreground">articles</span>
      </nav>

      <header class="max-w-4xl">
        <h1 class="font-display uppercase text-[32px] sm:text-[48px] md:text-[64px] lg:text-[80px] leading-[0.9] tracking-tight text-accent">
          Статьи
        </h1>
        <p class="mt-6 md:mt-8 text-base md:text-lg text-foreground/75 leading-relaxed max-w-3xl">
          Практика vibe coding, путь в AI-инжиниринг, выбор ментора, подготовка к собеседованиям и обзор IT-сообществ России — статьи на стыке инструментов и карьеры. Опыт участников IT-ХОЗЯЕВА и общие наблюдения по индустрии.
        </p>
      </header>

      <div class="mt-12 md:mt-16 grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-8 max-w-6xl">
        <router-link
          v-for="article in articles"
          :key="article.slug"
          :to="`/articles/${article.slug}`"
          class="group block border border-accent/15 bg-background/95 p-6 md:p-8 hover:border-accent/50 transition-colors"
        >
          <div class="font-mono text-[11px] text-foreground/40 uppercase tracking-widest mb-4 flex items-center gap-2">
            <span class="w-1.5 h-1.5 rounded-full bg-accent/60" />
            <time>{{ article.publishedAt }}</time>
            <span
              v-if="article.tags?.length"
              class="text-foreground/30"
            >|</span>
            <span
              v-if="article.tags?.length"
              class="text-foreground/50"
            >{{ article.tags.slice(0, 2).join(' · ') }}</span>
          </div>
          <h2 class="font-display uppercase text-xl md:text-2xl text-accent leading-tight group-hover:underline underline-offset-4">
            {{ article.h1 }}
          </h2>
          <p class="mt-4 text-sm md:text-base text-foreground/70 leading-relaxed">
            {{ article.excerpt }}
          </p>
          <div class="mt-6 font-mono text-xs text-accent/70 group-hover:text-accent transition-colors">
            читать →
          </div>
        </router-link>
      </div>
    </div>
  </section>
</template>
