<script setup lang="ts">
import type { MentorWithReviews } from '@/services/mentors'
import ReviewForm from '@/components/mentors/ReviewForm.vue'
import { mentorsService } from '@/services/mentors'
import { Tag, Typography } from 'itx-ui-kit'
import { ArrowLeft, Loader2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const mentor = ref<MentorWithReviews | null>(null)

async function loadMentor() {
  const id = Number(route.params.id)
  mentor.value = await mentorsService.getById(id)
}

onMounted(loadMentor)
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <RouterLink to="/mentors" class="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-6">
      <ArrowLeft class="h-4 w-4" />
      Назад к менторам
    </RouterLink>

    <div v-if="!mentor" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div v-else class="space-y-6">
      <div class="bg-card rounded-3xl border p-6">
        <Typography variant="h2" as="h1" class="mb-2">
          {{ mentor.firstName }} {{ mentor.lastName }}
        </Typography>
        <p v-if="mentor.occupation" class="text-muted-foreground mb-1">
          {{ mentor.occupation }}
        </p>
        <p v-if="mentor.experience" class="text-sm mb-3">
          {{ mentor.experience }}
        </p>
        <div v-if="mentor.profTags?.length" class="flex flex-wrap gap-1 mb-3">
          <Tag v-for="tag in mentor.profTags" :key="tag.id">
            {{ tag.title }}
          </Tag>
        </div>
        <a
          v-if="mentor.tg"
          :href="`https://t.me/${mentor.tg}`"
          target="_blank"
          class="text-sm text-primary underline"
        >
          @{{ mentor.tg }}
        </a>
      </div>

      <div v-if="mentor.contacts?.length" class="bg-card rounded-3xl border p-6">
        <Typography variant="h3" as="h2" class="mb-4">
          Контакты
        </Typography>
        <div class="space-y-2">
          <div v-for="contact in mentor.contacts" :key="contact.id" class="flex items-center gap-2 text-sm">
            <a :href="contact.link" target="_blank" class="text-primary underline">
              {{ contact.link }}
            </a>
          </div>
        </div>
      </div>

      <div v-if="mentor.services?.length" class="bg-card rounded-3xl border p-6">
        <Typography variant="h3" as="h2" class="mb-4">
          Услуги
        </Typography>
        <div class="space-y-2">
          <div
            v-for="service in mentor.services"
            :key="service.id"
            class="flex justify-between items-center border-b border-border pb-2 last:border-0"
          >
            <span>{{ service.name }}</span>
            <span v-if="service.price" class="text-sm text-muted-foreground">{{ service.price }} ₽</span>
          </div>
        </div>
      </div>

      <div v-if="mentor.reviews?.length" class="bg-card rounded-3xl border p-6">
        <Typography variant="h3" as="h2" class="mb-4">
          Отзывы
        </Typography>
        <div class="space-y-4">
          <div
            v-for="review in mentor.reviews"
            :key="review.id"
            class="border-b border-border pb-3 last:border-0"
          >
            <div class="flex items-center gap-2 text-sm text-muted-foreground mb-1">
              <span>{{ review.author }}</span>
              <span v-if="review.service">— {{ review.service.name }}</span>
              <span>{{ new Date(review.date).toLocaleDateString() }}</span>
            </div>
            <p class="text-sm">
              {{ review.text }}
            </p>
          </div>
        </div>
      </div>

      <div v-if="mentor.services?.length" class="bg-card rounded-3xl border p-6">
        <ReviewForm
          :mentor-id="mentor.id"
          :services="mentor.services"
          @submitted="loadMentor"
        />
      </div>
    </div>
  </div>
</template>
