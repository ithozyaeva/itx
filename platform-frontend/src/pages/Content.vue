<script setup lang="ts">
import { BookOpen } from 'lucide-vue-next'
</script>

<template>
  <div class="flex flex-col items-center justify-center py-24 px-4 text-center">
    <div class="rounded-full bg-muted p-6 mb-6">
      <BookOpen class="h-12 w-12 text-muted-foreground" />
    </div>
    <h1 class="text-2xl font-bold mb-2">
      Контент
    </h1>
    <p class="text-muted-foreground max-w-md">
      Раздел с эксклюзивным контентом сообщества скоро будет доступен.
      Здесь появятся записи стримов, статьи и обучающие материалы.
    </p>
  </div>
</template>
