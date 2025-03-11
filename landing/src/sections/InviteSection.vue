<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
</script>

<template>
  <section id="invite" class="w-full py-12 md:py-24 lg:py-32 bg-primary text-primary-foreground">
    <div class="container px-4 md:px-6">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div>
          <h2 class="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl mb-2">
            Готовы стать IT-ХОЗЯИНОМ?
          </h2>
          <p class="mx-auto max-w-[700px] text-primary-foreground/80 md:text-xl mb-4">
            Присоединяйтесь к сообществу профессионалов и начните свой путь к успеху в IT уже сегодня
          </p>
        </div>
        <div class="w-full max-w-sm text-black">
          <a href="https://boosty.to/jointime?isFromFeed=true" target="_blank">
            <Button type="outline">Присоединиться сейчас</Button>
          </a>
        </div>
      </div>
    </div>
  </section>
</template>
