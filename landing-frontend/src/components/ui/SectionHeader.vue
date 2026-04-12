<script setup lang="ts">
defineProps<{
  index: string
  path: string
  title: string
  subtitle?: string
  accentSubtitle?: boolean
}>()
</script>

<template>
  <div class="flex flex-col gap-5 md:gap-6">
    <div class="flex items-center gap-3 md:gap-4 overflow-hidden">
      <span class="font-mono text-[11px] md:text-xs text-accent/70 tracking-[0.15em] shrink-0">
        [{{ index }}]
      </span>
      <span class="font-mono text-[11px] md:text-xs text-foreground/50 tracking-[0.08em] shrink-0">
        {{ path }}
      </span>
      <span class="flex-1 h-px bg-gradient-to-r from-accent/40 via-accent/10 to-transparent" />
      <span class="font-mono text-[11px] md:text-xs text-term-amber/80 shrink-0 hidden sm:inline">
        ● READY
      </span>
    </div>
    <h2 class="font-display text-accent text-[32px] sm:text-[44px] md:text-[56px] lg:text-[72px] leading-[0.95] tracking-tight uppercase">
      {{ title }}
    </h2>
    <p
      v-if="subtitle"
      class="max-w-2xl text-base md:text-lg text-foreground/70 leading-relaxed"
      :class="{ 'text-accent': accentSubtitle }"
    >
      {{ subtitle }}
    </p>
  </div>
</template>
