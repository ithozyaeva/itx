<script setup lang="ts">
import logo from '@/assets/logo.png'
import Button from '@/components/ui/Button.vue'
import Laptop from '~icons/lucide/laptop'
import MessageSquare from '~icons/lucide/message-square'
import Users from '~icons/lucide/users'
</script>

<template>
  <section class="w-full py-12 md:py-24 lg:py-32 bg-gradient-to-b from-background to-muted">
    <div class="container px-4 md:px-6 space-y-10 xl:space-y-16">
      <div class="grid gap-4 px-4 sm:px-6 md:px-10 md:grid-cols-2 md:gap-16">
        <div class="flex flex-col justify-center space-y-4">
          <div class="space-y-2">
            <h1 class="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl lg:text-6xl">
              Станьте хозяином в мире IT
            </h1>
            <p class="max-w-[600px] text-muted-foreground md:text-xl">
              Присоединяйтесь к сообществу профессионалов, общайтесь, делитесь опытом и развивайтесь вместе
            </p>
          </div>
          <div class="flex flex-col gap-2 min-[400px]:flex-row">
            <a href="#invite">
              <Button target-id="pricing" class="px-8" to="invite">
                Присоединиться
              </Button>
            </a>
            <a href="#why">
              <Button target-id="features" to="more" type="outline">
                Узнать больше
              </Button>
            </a>
          </div>
        </div>
        <div class="flex items-center justify-center">
          <img
            :src="logo" alt="IT-ХОЗЯЕВА сообщество" class="rounded-lg object-contain"
          >
        </div>
      </div>
      <div class="mx-auto grid items-center gap-4 px-4 sm:px-6 md:px-10 md:grid-cols-3 md:gap-8">
        <div class="flex flex-col items-center space-y-2 border rounded-lg p-4">
          <Users class="h-8 w-8 text-primary" />
          <h3 class="text-xl font-bold">
            200+
          </h3>
          <p class="text-center text-muted-foreground">
            Активных участников
          </p>
        </div>
        <div class="flex flex-col items-center space-y-2 border rounded-lg p-4">
          <MessageSquare class="h-8 w-8 text-primary" />
          <h3 class="text-xl font-bold">
            1000+
          </h3>
          <p class="text-center text-muted-foreground">
            Сообщений ежедневно
          </p>
        </div>
        <div class="flex flex-col items-center space-y-2 border rounded-lg p-4">
          <Laptop class="h-8 w-8 text-primary" />
          <h3 class="text-xl font-bold">
            20+
          </h3>
          <p class="text-center text-muted-foreground">
            Эксклюзивных материалов
          </p>
        </div>
      </div>
    </div>
  </section>
</template>
