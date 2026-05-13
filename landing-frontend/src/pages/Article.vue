<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import StructuredData from '@/components/StructuredData.vue'
import { usePageMeta } from '@/composables/useMeta'
import { useRevealObserver } from '@/composables/useRevealObserver'
import { articles, getArticleBySlug, renderBlocksToHtml } from '@/content/articles'

useRevealObserver()

const route = useRoute()

const article = computed(() => {
  const slug = String(route.params.slug || '')
  return getArticleBySlug(slug)
})

const bodyHtml = computed(() => article.value ? renderBlocksToHtml(article.value.body) : '')

const related = computed(() => {
  if (!article.value)
    return []
  return articles.filter(a => a.slug !== article.value!.slug).slice(0, 3)
})

const canonicalUrl = computed(() => article.value ? `https://ithozyaeva.ru/articles/${article.value.slug}` : 'https://ithozyaeva.ru/articles')

usePageMeta({
  title: article.value?.title ?? 'Статья не найдена',
  description: article.value?.description ?? 'Запрашиваемая статья не найдена.',
  url: canonicalUrl.value,
  type: 'article',
  noindex: !article.value,
})

const structuredData = computed(() => {
  if (!article.value)
    return []
  const a = article.value
  const data: object[] = [
    {
      '@context': 'https://schema.org',
      '@type': 'BreadcrumbList',
      'itemListElement': [
        { '@type': 'ListItem', 'position': 1, 'name': 'Главная', 'item': 'https://ithozyaeva.ru/' },
        { '@type': 'ListItem', 'position': 2, 'name': 'Статьи', 'item': 'https://ithozyaeva.ru/articles' },
        { '@type': 'ListItem', 'position': 3, 'name': a.breadcrumb, 'item': canonicalUrl.value },
      ],
    },
    {
      '@context': 'https://schema.org',
      '@type': 'Article',
      'headline': a.h1,
      'description': a.description,
      'datePublished': a.publishedAt,
      'dateModified': a.updatedAt ?? a.publishedAt,
      'author': { '@type': 'Organization', 'name': 'IT-ХОЗЯЕВА', 'url': 'https://ithozyaeva.ru' },
      'publisher': {
        '@type': 'Organization',
        'name': 'IT-ХОЗЯЕВА',
        'url': 'https://ithozyaeva.ru',
        'logo': { '@type': 'ImageObject', 'url': 'https://ithozyaeva.ru/og-image.png' },
      },
      'mainEntityOfPage': canonicalUrl.value,
      'inLanguage': 'ru',
    },
  ]
  if (a.faq?.length) {
    data.push({
      '@context': 'https://schema.org',
      '@type': 'FAQPage',
      'mainEntity': a.faq.map(f => ({
        '@type': 'Question',
        'name': f.q,
        'acceptedAnswer': { '@type': 'Answer', 'text': f.a },
      })),
    })
  }
  return data
})
</script>

<template>
  <StructuredData
    v-if="article"
    :data="structuredData"
  />

  <article
    v-if="article"
    class="w-full pt-10 md:pt-16 lg:pt-20 pb-20 md:pb-32"
  >
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
        <a
          href="/articles"
          class="hover:text-accent transition-colors"
        >articles</a>
        <span class="mx-2 text-accent/40">/</span>
        <span class="text-foreground">{{ article.slug }}</span>
      </nav>

      <header class="max-w-3xl">
        <div class="font-mono text-[11px] text-foreground/40 uppercase tracking-widest mb-4 flex items-center gap-2">
          <span class="w-1.5 h-1.5 rounded-full bg-accent/60" />
          <time :datetime="article.publishedAt">{{ article.publishedAt }}</time>
          <span
            v-if="article.tags?.length"
            class="text-foreground/30"
          >|</span>
          <span
            v-if="article.tags?.length"
            class="text-foreground/50"
          >{{ article.tags.join(' · ') }}</span>
        </div>
        <h1 class="font-display uppercase text-[28px] sm:text-[40px] md:text-[56px] lg:text-[64px] leading-[0.95] tracking-tight text-accent">
          {{ article.h1 }}
        </h1>
        <p class="mt-6 md:mt-8 text-base md:text-lg text-foreground/75 leading-relaxed">
          {{ article.lead }}
        </p>
      </header>

      <!-- eslint-disable-next-line vue/no-v-html -->
      <div
        class="article-body mt-14 md:mt-20 max-w-3xl text-base md:text-lg text-foreground/75 leading-relaxed"
        v-html="bodyHtml"
      />

      <section
        v-if="article.faq?.length"
        class="mt-16 md:mt-24 max-w-3xl"
      >
        <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-6">
          Частые вопросы
        </h2>
        <div class="space-y-6">
          <div
            v-for="(item, i) in article.faq"
            :key="i"
          >
            <h3 class="font-display uppercase text-base md:text-lg text-foreground mb-2">
              {{ item.q }}
            </h3>
            <p class="text-sm md:text-base text-foreground/70 leading-relaxed">
              {{ item.a }}
            </p>
          </div>
        </div>
      </section>

      <section
        v-if="related.length"
        class="mt-20 md:mt-28 max-w-6xl"
      >
        <h2 class="font-display uppercase text-2xl md:text-3xl text-accent mb-8">
          Ещё статьи
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <router-link
            v-for="r in related"
            :key="r.slug"
            :to="`/articles/${r.slug}`"
            class="group block border border-accent/15 bg-background/95 p-5 md:p-6 hover:border-accent/50 transition-colors"
          >
            <h3 class="font-display uppercase text-base md:text-lg text-accent leading-tight group-hover:underline underline-offset-4">
              {{ r.h1 }}
            </h3>
            <p class="mt-3 text-sm text-foreground/65 leading-relaxed">
              {{ r.excerpt }}
            </p>
          </router-link>
        </div>
      </section>
    </div>
  </article>

  <section
    v-else
    class="w-full pt-20 pb-32"
  >
    <div class="container px-6 md:px-10">
      <h1 class="font-display uppercase text-4xl text-accent mb-4">
        Статья не найдена
      </h1>
      <p class="text-foreground/70">
        Возможно, ссылка устарела. <router-link
          to="/articles"
          class="text-accent hover:underline"
        >
          Посмотреть все статьи
        </router-link>.
      </p>
    </div>
  </section>
</template>

<style scoped>
.article-body :deep(h2) {
  margin-top: 3rem;
  margin-bottom: 1rem;
  color: hsl(var(--accent));
  font-family: var(--font-display, inherit);
  font-weight: 700;
  font-size: 1.5rem;
  line-height: 1.2;
  text-transform: uppercase;
  letter-spacing: -0.01em;
}

@media (min-width: 768px) {
  .article-body :deep(h2) {
    font-size: 1.875rem;
  }
}

.article-body :deep(h3) {
  margin-top: 2rem;
  margin-bottom: 0.75rem;
  color: hsl(var(--foreground));
  font-family: var(--font-display, inherit);
  font-weight: 700;
  font-size: 1.125rem;
  line-height: 1.3;
  text-transform: uppercase;
  letter-spacing: 0;
}

@media (min-width: 768px) {
  .article-body :deep(h3) {
    font-size: 1.25rem;
  }
}

.article-body :deep(p) {
  margin-top: 1rem;
  margin-bottom: 1rem;
}

.article-body :deep(p):first-child {
  margin-top: 0;
}

.article-body :deep(ul),
.article-body :deep(ol) {
  margin-top: 1rem;
  margin-bottom: 1rem;
  padding-left: 1.5rem;
}

.article-body :deep(ul) {
  list-style-type: disc;
}

.article-body :deep(ol) {
  list-style-type: decimal;
}

.article-body :deep(li) {
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}

.article-body :deep(strong) {
  color: hsl(var(--foreground));
  font-weight: 600;
}

.article-body :deep(em) {
  color: hsl(var(--foreground) / 0.85);
  font-style: italic;
}

.article-body :deep(a) {
  color: hsl(var(--accent));
  text-decoration: underline;
  text-underline-offset: 4px;
  text-decoration-thickness: 1px;
  transition: opacity 0.15s ease;
}

.article-body :deep(a):hover {
  opacity: 0.75;
}

.article-body :deep(.article-note) {
  margin-top: 1.5rem;
  margin-bottom: 1.5rem;
  padding: 1rem 1.25rem;
  border-left: 2px solid hsl(var(--accent) / 0.5);
  background: hsl(var(--accent) / 0.05);
  font-size: 0.95em;
}

.article-body :deep(.article-cta) {
  margin-top: 2rem;
  margin-bottom: 2rem;
}

.article-body :deep(.article-cta a) {
  display: inline-block;
  padding: 0.75rem 1.25rem;
  border: 1px solid hsl(var(--accent));
  background: hsl(var(--accent));
  color: #0b0d0c;
  text-decoration: none;
  font-family: var(--font-display, inherit);
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 0.04em;
  transition: opacity 0.15s ease;
}

.article-body :deep(.article-cta a):hover {
  opacity: 0.85;
}
</style>
