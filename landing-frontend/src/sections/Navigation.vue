<script setup lang="ts">
const menuItems = [
  { id: '02', label: 'менторы', href: '#mentors' },
  { id: '03', label: 'встречи', href: '#meets' },
  { id: '04', label: 'тарифы', href: '#tariffs' },
  { id: '05', label: 'вопросы', href: '#FAQ' },
]
</script>

<template>
  <nav class="flex gap-1 md:gap-3 lg:gap-6">
    <a
      v-for="item in menuItems"
      :key="item.href"
      :href="item.href"
      class="group relative px-2 py-2 lg:px-3 lg:py-2 flex items-center gap-2 text-foreground/70 hover:text-accent transition-colors"
    >
      <span class="font-mono text-[10px] opacity-40 group-hover:opacity-100 transition-opacity">
        [{{ item.id }}]
      </span>
      <span class="font-display uppercase text-sm lg:text-base tracking-wide">
        {{ item.label }}
      </span>
      <span class="absolute left-2 right-2 bottom-0.5 h-px bg-accent scale-x-0 origin-left group-hover:scale-x-100 transition-transform duration-300" />
    </a>
  </nav>
</template>
