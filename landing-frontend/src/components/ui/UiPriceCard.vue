<script setup lang="ts">
import { Check } from 'lucide-vue-next'
import { computed } from 'vue'
import UiButton from './UiButton.vue'

const props = defineProps<{
  name: string
  price: number
  oldPrice?: number
  features: string[]
  link: string
  variant?: 'default' | 'highlighted'
  tierIndex?: string
  tierLabel?: string
}>()

const discount = computed(() => {
  if (!props.oldPrice)
    return null
  const old = +props.oldPrice
  const current = +props.price
  if (old <= 0 || current >= old)
    return null
  return Math.round((old - current) * 100 / old)
})
</script>

<template>
  <div
    class="price-card group relative flex flex-col border transition-all duration-300"
    :class="variant === 'highlighted'
      ? 'bg-accent text-[#0b0d0c] border-accent shadow-[0_30px_80px_-30px_hsl(var(--accent)/0.5)] md:-translate-y-3'
      : 'bg-background/80 text-foreground border-accent/20 hover:border-accent/50'"
  >
    <!-- corner brackets -->
    <span
      class="absolute top-0 left-0 w-3 h-3 border-t-2 border-l-2"
      :class="variant === 'highlighted' ? 'border-[#0b0d0c]' : 'border-accent'"
    />
    <span
      class="absolute top-0 right-0 w-3 h-3 border-t-2 border-r-2"
      :class="variant === 'highlighted' ? 'border-[#0b0d0c]' : 'border-accent'"
    />
    <span
      class="absolute bottom-0 left-0 w-3 h-3 border-b-2 border-l-2"
      :class="variant === 'highlighted' ? 'border-[#0b0d0c]' : 'border-accent'"
    />
    <span
      class="absolute bottom-0 right-0 w-3 h-3 border-b-2 border-r-2"
      :class="variant === 'highlighted' ? 'border-[#0b0d0c]' : 'border-accent'"
    />

    <!-- header -->
    <div
      class="px-6 md:px-7 pt-6 pb-5 border-b"
      :class="variant === 'highlighted' ? 'border-[#0b0d0c]/15' : 'border-accent/15'"
    >
      <div
        class="flex items-center justify-between font-mono text-[11px] tracking-[0.1em] uppercase mb-4"
        :class="variant === 'highlighted' ? 'text-[#0b0d0c]/70' : 'text-foreground/50'"
      >
        <span>[ уровень.{{ tierIndex }} ]</span>
        <span
          v-if="variant === 'highlighted'"
          class="px-2 py-0.5 bg-[#0b0d0c] text-accent"
        >выбор</span>
        <span v-else-if="tierLabel">{{ tierLabel }}</span>
      </div>
      <h3
        class="font-display uppercase text-2xl md:text-3xl leading-tight"
        :class="variant === 'highlighted' ? 'text-[#0b0d0c]' : 'text-accent'"
      >
        {{ name }}
      </h3>

      <div class="mt-5 flex items-baseline gap-2">
        <span
          v-if="oldPrice"
          class="font-mono text-sm line-through"
          :class="variant === 'highlighted' ? 'text-[#0b0d0c]/50' : 'text-foreground/30'"
        >
          {{ oldPrice }}₽
        </span>
        <span
          class="font-display text-4xl md:text-5xl leading-none"
          :class="variant === 'highlighted' ? 'text-[#0b0d0c]' : 'text-foreground'"
        >
          {{ price }}
        </span>
        <span
          class="font-display text-2xl md:text-3xl"
          :class="variant === 'highlighted' ? 'text-[#0b0d0c]/80' : 'text-foreground/70'"
        >
          ₽
        </span>
        <span
          class="font-mono text-xs self-end pb-1"
          :class="variant === 'highlighted' ? 'text-[#0b0d0c]/60' : 'text-foreground/40'"
        >
          /мес
        </span>
        <span
          v-if="discount"
          class="ml-auto font-mono text-[11px] px-2 py-1 border tracking-wider"
          :class="variant === 'highlighted'
            ? 'border-[#0b0d0c] text-[#0b0d0c]'
            : 'border-term-amber text-term-amber'"
        >
          −{{ discount }}%
        </span>
      </div>
    </div>

    <!-- features -->
    <ul class="flex-1 px-6 md:px-7 py-6 space-y-3">
      <li
        v-for="(feature, index) in features"
        :key="index"
        class="flex items-start gap-3 text-sm md:text-[15px] leading-snug"
        :class="variant === 'highlighted' ? 'text-[#0b0d0c]' : 'text-foreground/80'"
      >
        <Check
          class="w-4 h-4 mt-0.5 shrink-0"
          :class="variant === 'highlighted' ? 'text-[#0b0d0c]' : 'text-accent'"
          stroke-width="3"
        />
        <span>{{ feature }}</span>
      </li>
    </ul>

    <!-- footer -->
    <div class="px-6 md:px-7 pb-6">
      <UiButton
        as="a"
        :href="link"
        target="_blank"
        :variant="variant === 'highlighted' ? 'dark-filled' : 'filled'"
        class="w-full text-center"
      >
        Подписаться →
      </UiButton>
    </div>
  </div>
</template>

<style scoped>
.price-card {
  border-radius: 2px;
  min-width: 280px;
  min-height: 480px;
}
.price-card :deep(.ui-button) {
  width: 100%;
  border-radius: 2px !important;
}
</style>
