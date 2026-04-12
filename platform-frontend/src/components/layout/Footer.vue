<script setup lang="ts">
import { RouterLink } from 'vue-router'

const links = [
  { label: 'Барахолка', path: '/marketplace' },
  { label: 'Менторы', path: '/mentors' },
  { label: 'Рефералки', path: '/referals' },
  { label: 'Рейтинг', path: '/leaderboard' },
  { label: 'Достижения', path: '/achievements' },
  { label: 'Резюме', path: '/resumes' },
]
</script>

<template>
  <footer class="w-full border-t py-4">
    <div class="container mx-auto px-4 flex flex-col gap-4">
      <div class="grid grid-cols-2 sm:grid-cols-3 gap-x-4 sm:gap-x-8 gap-y-2">
        <RouterLink
          v-for="link in links"
          :key="link.path"
          :to="link.path"
          class="font-mono text-xs text-muted-foreground hover:text-accent transition-colors tracking-wide"
        >
          ./{{ link.label.toLowerCase() }}
        </RouterLink>
      </div>
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 font-mono text-[11px] text-muted-foreground/70 tracking-wide border-t border-border pt-3">
        <span>© {{ new Date().getFullYear() }} IT-ХОЗЯЕВА // all rights reserved</span>
        <span class="flex items-center gap-2">
          <span class="w-1.5 h-1.5 rounded-full bg-accent" />
          status: operational
        </span>
      </div>
    </div>
  </footer>
</template>
