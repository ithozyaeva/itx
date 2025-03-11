<script setup lang="ts">
import izede from '@/assets/izede.jpg'
import placeholderSvg from '@/assets/placeholder.svg'
import Card from '@/components/ui/Card.vue'

interface Review {
  avatarPath: string
  author: string
  occupation: string
  text: string
}

const reviews: Review[] = [
  {
    avatarPath: izede,
    author: 'Денис Чернов',
    occupation: 'Vue-евангелист',
    text: '"В сообществе царит атмосфера полной поддержки друг дуга, есть чаты на любой вкус. Разработчики самых разных уровней, от начинающих до матерых разработчиков. Клубы и активности идут бесконечным потоком, сообщество само генерирует контент."',
  },
  {
    avatarPath: placeholderSvg,
    author: 'Мария Иванова',
    occupation: 'CTO, Startup Hub',
    text: '"Сообщество IT-ХОЗЯЕВА стало для меня настоящим прорывом. Здесь я нашла единомышленников, менторов и партнеров для своих проектов. Рекомендую всем, кто хочет расти в IT-сфере."',
  },
  {
    avatarPath: placeholderSvg,
    author: 'Дмитрий Сидоров',
    occupation: 'Product Manager, Sber',
    text: '"Благодаря экспертам из IT-ХОЗЯЕВА я смог быстро перепрофилироваться из разработчика в продуктового менеджера. Материалы и поддержка сообщества бесценны для тех, кто хочет развиваться."',
  },
  {
    avatarPath: placeholderSvg,
    author: 'Елена Смирнова',
    occupation: 'DevOps Engineer, VK',
    text: '"Вступила в сообщество год назад, когда только начинала карьеру в DevOps. Сейчас я работаю в крупной компании и продолжаю учиться у лучших специалистов отрасли в IT-ХОЗЯЕВА."',
  },
]
</script>

<template>
  <section id="reviews" class="w-full py-12 md:py-24 lg:py-32 bg-muted">
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
        <Card v-for="review in reviews" :key="review.author">
          <template #content>
            <div class="pt-6">
              <div class="flex items-start gap-4">
                <img
                  :src="review.avatarPath" :alt="`Аватар ${review.author}`" class="rounded-full" width="60"
                  height="60"
                >
                <div class="grid gap-1">
                  <h3 class="font-bold">
                    {{ review.author }}
                  </h3>
                  <p class="text-sm text-muted-foreground">
                    {{ review.occupation }}
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
