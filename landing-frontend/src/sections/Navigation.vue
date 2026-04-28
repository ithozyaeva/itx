<script setup lang="ts">
interface MenuItem {
  id: string
  label: string
  href: string
  accent?: boolean
  badge?: string
}

const menuItems: MenuItem[] = [
  { id: '02', label: 'менторы', href: '#mentors' },
  { id: '03', label: 'встречи', href: '#meets' },
  { id: '04', label: 'тарифы', href: '#tariffs', accent: true, badge: 'от 520 ₽' },
  { id: '05', label: 'вопросы', href: '#FAQ' },
]
</script>

<template>
  <nav class="flex gap-0.5 md:gap-1 lg:gap-4">
    <a
      v-for="item in menuItems"
      :key="item.href"
      :href="item.href"
      class="group relative px-1.5 py-2 lg:px-3 lg:py-2 flex items-center gap-1.5 transition-colors"
      :class="item.accent ? 'text-accent hover:text-accent/80' : 'text-foreground/70 hover:text-accent'"
    >
      <span class="font-mono text-[10px] opacity-40 group-hover:opacity-100 transition-opacity hidden lg:inline">
        [{{ item.id }}]
      </span>
      <span class="font-display uppercase text-xs lg:text-base tracking-wide">
        {{ item.label }}
      </span>
      <span
        v-if="item.badge"
        class="hidden lg:inline font-mono text-[10px] tracking-wider px-1.5 py-0.5 border border-accent/40 bg-accent/10 text-accent"
      >
        {{ item.badge }}
      </span>
      <span class="absolute left-2 right-2 bottom-0.5 h-px bg-accent scale-x-0 origin-left group-hover:scale-x-100 transition-transform duration-300" />
    </a>
  </nav>
</template>
