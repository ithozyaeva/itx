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
  <footer class="w-full border-t py-6">
    <div class="container mx-auto px-4 flex flex-col gap-6">
      <div class="grid grid-cols-2 sm:grid-cols-3 gap-x-8 gap-y-2">
        <RouterLink
          v-for="link in links"
          :key="link.path"
          :to="link.path"
          class="text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          {{ link.label }}
        </RouterLink>
      </div>
      <p class="text-center text-sm text-muted-foreground">
        © {{ new Date().getFullYear() }} IT-ХОЗЯЕВА. Все права защищены.
      </p>
    </div>
  </footer>
</template>
