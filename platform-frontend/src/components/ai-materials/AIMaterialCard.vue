<script setup lang="ts">
import type { AIMaterial } from '@/models/aiMaterial'
import { ExternalLink, FileCode2, Sparkles } from 'lucide-vue-next'
import { computed } from 'vue'
import { AI_MATERIAL_KIND_LABELS } from '@/models/aiMaterial'
import AIMaterialReactions from './AIMaterialReactions.vue'

const props = defineProps<{ item: AIMaterial }>()

const emit = defineEmits<{
  'update:item': [v: AIMaterial]
}>()

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

function patch(field: keyof AIMaterial, value: AIMaterial[keyof AIMaterial]) {
  emit('update:item', { ...props.item, [field]: value })
}
</script>

<template>
  <article
    class="group flex flex-col gap-3 rounded-sm border border-border bg-card p-4 transition-colors hover:border-accent terminal-card h-full"
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
      <AIMaterialReactions
        :material-id="item.id"
        :liked="item.liked"
        :bookmarked="item.bookmarked"
        :likes-count="item.likesCount"
        :bookmarks-count="item.bookmarksCount"
        :comments-count="item.commentsCount"
        :stop-propagation="true"
        size="sm"
        @update:liked="(v: boolean) => patch('liked', v)"
        @update:bookmarked="(v: boolean) => patch('bookmarked', v)"
        @update:likes-count="(v: number) => patch('likesCount', v)"
        @update:bookmarks-count="(v: number) => patch('bookmarksCount', v)"
      />
    </footer>
  </article>
</template>
