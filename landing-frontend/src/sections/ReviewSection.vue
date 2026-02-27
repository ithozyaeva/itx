<script setup lang="ts">
import type { TelegramUser } from '@/services/auth'
import { AlertCircle } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import Card from '@/components/ui/Card.vue'
import TgImage from '@/components/ui/TgImage.vue'

interface Review {
  avatarPath: string
  author: TelegramUser
  occupation: string
  text: string
}

function reloadPage() {
  window.location.reload()
}

const reviews = ref<Review[]>([])
const loading = ref<boolean>(false)
const error = ref()

async function search() {
  try {
    error.value = false
    loading.value = true
    const response = await fetch('/api/review-on-community')

    if (!response.ok) {
      throw new Error('Не удалось загрузить данные о отзывах')
    }

    reviews.value = await response.json()
  }
  catch (err) {
    console.error(`Произошла ошибка при запросе отзывов - ${err}`)
    error.value = err
  }
  finally {
    loading.value = false
  }
}

onMounted(search)
</script>

<template>
  <section
    id="reviews"
    class="w-full py-12 md:py-24 lg:py-32 bg-muted"
  >
    <div
      v-if="loading"
      class="text-center py-16"
    >
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto" />
      <p class="mt-4 text-gray-600">
        Загрузка отзывов...
      </p>
    </div>
    <div class="container px-4 md:px-6">
      <div class="flex flex-col items-center justify-center space-y-4 text-center">
        <div class="space-y-2">
          <div class="inline-block rounded-lg bg-primary px-3 py-1 text-sm text-primary-foreground">
            Отзывы
          </div>
          <h2 class="text-3xl font-bold tracking-tighter sm:text-5xl">
            Что говорят участники
          </h2>
          <p class="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
            Узнайте, как наше сообщество помогает IT-специалистам достигать новых высот
          </p>
        </div>
      </div>
      <div class="mx-auto grid max-w-5xl gap-6 py-12 md:grid-cols-2 md:gap-8">
        <div
          v-if="error && !reviews.length"
          class="text-center py-10"
        >
          <div class="flex justify-center mb-4">
            <AlertCircle class="h-12 w-12 text-red-500" />
          </div>
          <p class="text-red-500 mb-2">
            {{ error }}
          </p>
          <button
            class="mt-4 px-4 py-2 bg-blue-500 text-white rounded"
            @click="reloadPage"
          >
            Попробовать снова
          </button>
        </div>

        <Card
          v-for="review in reviews"
          :key="review.author.id"
        >
          <template #content>
            <div class="pt-6">
              <div class="flex items-start gap-4">
                <TgImage
                  :username="review.author.tg"
                  class="rounded-full"
                  width="60"
                  height="60"
                />
                <div class="grid gap-1">
                  <h3 class="font-bold">
                    {{ review.author.firstName ?? "" }} {{ review.author.lastName ?? "" }}
                  </h3>
                  <p class="text-sm text-muted-foreground">
                    {{ review.author.tg }}
                  </p>
                </div>
              </div>
              <div class="mt-4">
                <p class="text-muted-foreground">
                  {{ review.text }}
                </p>
              </div>
            </div>
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>
