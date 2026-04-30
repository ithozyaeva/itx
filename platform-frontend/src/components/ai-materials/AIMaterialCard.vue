<script setup lang="ts">
import type { AIMaterial } from '@/models/aiMaterial'
import { Bookmark, ExternalLink, FileCode2, Heart, MessageCircle, Sparkles } from 'lucide-vue-next'
import { computed } from 'vue'
import { AI_MATERIAL_KIND_LABELS } from '@/models/aiMaterial'

const props = defineProps<{ item: AIMaterial }>()

const contentIcon = computed(() => {
  switch (props.item.contentType) {
    case 'link':
      return ExternalLink
    case 'agent':
      return FileCode2
    default:
      return Sparkles
  }
})

function authorName(): string {
  const a = props.item.author
  if (!a)
    return 'Аноним'
  const name = [a.firstName, a.lastName].filter(Boolean).join(' ')
  return name || (a.tg ? `@${a.tg}` : 'Аноним')
}
</script>

<template>
  <article
    class="group flex flex-col gap-3 rounded-sm border border-border bg-card p-4 transition-colors hover:border-accent terminal-card"
  >
    <header class="flex items-start gap-2">
      <span
        class="shrink-0 inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium bg-accent/15 text-accent"
      >
        <component :is="contentIcon" class="h-3 w-3" />
        {{ AI_MATERIAL_KIND_LABELS[item.materialKind] }}
      </span>
      <h3 class="font-medium text-sm leading-tight line-clamp-2 flex-1">
        {{ item.title }}
      </h3>
    </header>

    <p class="text-xs text-muted-foreground line-clamp-3 min-h-[3rem]">
      {{ item.summary }}
    </p>

    <div v-if="item.tags.length" class="flex flex-wrap gap-1">
      <span
        v-for="t in item.tags.slice(0, 4)"
        :key="t"
        class="px-1.5 py-0.5 rounded-full text-[10px] font-mono bg-muted text-muted-foreground"
      >#{{ t }}</span>
      <span
        v-if="item.tags.length > 4"
        class="px-1.5 py-0.5 rounded-full text-[10px] font-mono bg-muted text-muted-foreground"
      >+{{ item.tags.length - 4 }}</span>
    </div>

    <footer class="mt-auto flex items-center justify-between text-xs text-muted-foreground">
      <span class="truncate">{{ authorName() }}</span>
      <div class="flex items-center gap-2">
        <span class="inline-flex items-center gap-0.5">
          <Heart class="h-3 w-3" :class="item.liked ? 'fill-red-500 text-red-500' : ''" />
          {{ item.likesCount }}
        </span>
        <span class="inline-flex items-center gap-0.5">
          <Bookmark class="h-3 w-3" :class="item.bookmarked ? 'fill-current' : ''" />
          {{ item.bookmarksCount }}
        </span>
        <span class="inline-flex items-center gap-0.5">
          <MessageCircle class="h-3 w-3" />
          {{ item.commentsCount }}
        </span>
      </div>
    </footer>
  </article>
</template>
