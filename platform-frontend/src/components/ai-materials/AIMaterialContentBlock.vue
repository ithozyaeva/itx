<script setup lang="ts">
import type { AIMaterial } from '@/models/aiMaterial'
import { Check, Copy, ExternalLink } from 'lucide-vue-next'
import { ref } from 'vue'

const props = defineProps<{ item: AIMaterial }>()

const copied = ref(false)
let copyTimeout: ReturnType<typeof setTimeout> | null = null

async function copyContent() {
  const text = props.item.contentType === 'agent'
    ? props.item.agentConfig
    : props.item.promptBody
  if (!text)
    return
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    if (copyTimeout)
      clearTimeout(copyTimeout)
    copyTimeout = setTimeout(() => (copied.value = false), 1800)
  }
  catch {
    // clipboard API недоступен (HTTP, старый браузер) — кнопка просто молчит,
    // пользователь увидит текст и сможет выделить вручную.
  }
}
</script>

<template>
  <div class="space-y-3">
    <!-- Prompt — большой <pre> с кнопкой Copy -->
    <div v-if="item.contentType === 'prompt'" class="relative">
      <button
        type="button"
        class="absolute right-2 top-2 inline-flex items-center gap-1 rounded-sm bg-background/80 backdrop-blur border border-border px-2 py-1 text-xs text-muted-foreground hover:text-foreground"
        :aria-label="copied ? 'Скопировано' : 'Скопировать промт'"
        @click="copyContent"
      >
        <component :is="copied ? Check : Copy" class="h-3.5 w-3.5" />
        {{ copied ? 'Скопировано' : 'Копировать' }}
      </button>
      <pre
        class="rounded-sm border border-border bg-muted/40 p-4 pr-24 text-xs font-mono whitespace-pre-wrap break-words max-h-[60vh] overflow-auto"
      >{{ item.promptBody }}</pre>
    </div>

    <!-- Link — большая кнопка-ссылка -->
    <a
      v-else-if="item.contentType === 'link' && item.externalUrl"
      :href="item.externalUrl"
      target="_blank"
      rel="noopener noreferrer"
      class="inline-flex items-center gap-2 rounded-sm border border-border bg-card px-4 py-3 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors break-all"
    >
      <ExternalLink class="h-4 w-4 shrink-0" />
      <span class="truncate">{{ item.externalUrl }}</span>
    </a>

    <!-- Agent — JSON/YAML конфиг с кнопкой Copy -->
    <div v-else-if="item.contentType === 'agent'" class="relative">
      <button
        type="button"
        class="absolute right-2 top-2 inline-flex items-center gap-1 rounded-sm bg-background/80 backdrop-blur border border-border px-2 py-1 text-xs text-muted-foreground hover:text-foreground"
        :aria-label="copied ? 'Скопировано' : 'Скопировать конфиг'"
        @click="copyContent"
      >
        <component :is="copied ? Check : Copy" class="h-3.5 w-3.5" />
        {{ copied ? 'Скопировано' : 'Копировать' }}
      </button>
      <pre
        class="rounded-sm border border-border bg-muted/40 p-4 pr-24 text-xs font-mono whitespace-pre-wrap break-words max-h-[60vh] overflow-auto"
      >{{ item.agentConfig }}</pre>
    </div>
  </div>
</template>
