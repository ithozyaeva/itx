<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
</script>

<template>
  <section id="mentors" class="w-full py-12 md:py-24 lg:py-32 bg-muted">
    <div class="container px-4 md:px-6">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="space-y-2">
          <div class="inline-block rounded-lg bg-primary px-3 py-1 text-sm text-primary-foreground">
            Менторы
          </div>
          <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">
            Наши эксперты
          </h2>
          <p class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
            Выберите идеального ментора для себя: от бесплатных наставников до звёзд IT. Каждый из них — активный
            участник сообщества IT-хозяева и практикующий профессионал.
          </p>
        </div>
        <div class="m-auto max-w-6xl pt-12">
          <a href="https://airtable.com/appYqSv5VbbroB0mq/shr4GWqncLHvGqOCp" target="_blank">
            <Button>Выбрать ментора</Button>
          </a>
        </div>
      </div>
    </div>
  </section>
</template>
