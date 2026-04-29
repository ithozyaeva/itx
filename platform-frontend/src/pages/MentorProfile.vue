<script setup lang="ts">
import type { MentorWithReviews } from '@/services/mentors'
import { ArrowLeft, Loader2, Lock } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import ErrorState from '@/components/common/ErrorState.vue'
import ReviewForm from '@/components/mentors/ReviewForm.vue'
import { Typography } from '@/components/ui/typography'
import { isUserSubscribed } from '@/composables/useUser'
import { formatShortDate } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { mentorsService } from '@/services/mentors'

const isSubscribed = isUserSubscribed()

const route = useRoute()
const mentor = ref<MentorWithReviews | null>(null)
const isLoading = ref(true)
const loadError = ref<string | null>(null)

async function loadMentor() {
  isLoading.value = true
  loadError.value = null
  try {
    const id = Number(route.params.id)
    if (Number.isNaN(id))
      throw new Error('Некорректный ID ментора')
    mentor.value = await mentorsService.getById(id)
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

onMounted(loadMentor)
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/mentor-profile
    </div>
    <RouterLink to="/mentors" class="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-6">
      <ArrowLeft class="h-4 w-4" />
      Назад к менторам
    </RouterLink>

    <div v-if="isLoading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <ErrorState v-else-if="loadError" :message="loadError" @retry="loadMentor" />

    <div v-else-if="mentor" class="space-y-6">
      <div class="bg-card rounded-sm terminal-card border p-6">
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
          <span
            v-for="tag in mentor.profTags"
            :key="tag.id"
            class="inline-flex items-center rounded-full border border-accent/30 px-2 py-0.5 text-xs text-accent"
          >
            {{ tag.title }}
          </span>
        </div>
        <a
          v-if="mentor.tg && isSubscribed"
          :href="`https://t.me/${mentor.tg}`"
          target="_blank"
          rel="noopener noreferrer"
          class="text-sm text-primary underline"
        >
          @{{ mentor.tg }}
        </a>
        <RouterLink
          v-else-if="mentor.tg"
          to="/tariffs"
          class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-accent transition-colors"
        >
          <Lock class="h-3.5 w-3.5" />
          Telegram ментора — по подписке
        </RouterLink>
      </div>

      <div v-if="mentor.contacts?.length && isSubscribed" class="bg-card rounded-sm terminal-card border p-6">
        <Typography variant="h3" as="h2" class="mb-4">
          Контакты
        </Typography>
        <div class="space-y-2">
          <div v-for="contact in mentor.contacts" :key="contact.id" class="flex items-center gap-2 text-sm min-w-0">
            <a :href="contact.link" target="_blank" rel="noopener noreferrer" class="text-accent underline break-all">
              {{ contact.link }}
            </a>
          </div>
        </div>
      </div>

      <div v-else-if="mentor.contacts?.length" class="bg-card/50 rounded-sm terminal-card border border-dashed p-6">
        <Typography variant="h3" as="h2" class="mb-2 flex items-center gap-2">
          <Lock class="h-5 w-5" />
          Контакты ментора
        </Typography>
        <p class="text-sm text-muted-foreground mb-3">
          Контакты для прямой связи доступны участникам с подпиской.
        </p>
        <RouterLink
          to="/tariffs"
          class="inline-block text-sm text-accent hover:underline"
        >
          Открыть тарифы →
        </RouterLink>
      </div>

      <div v-if="mentor.services?.length" class="bg-card rounded-sm terminal-card border p-6">
        <Typography variant="h3" as="h2" class="mb-4">
          Услуги
        </Typography>
        <div class="space-y-2">
          <div
            v-for="service in mentor.services"
            :key="service.id"
            class="flex justify-between items-start gap-2 border-b border-border pb-2 last:border-0"
          >
            <span class="break-words min-w-0">{{ service.name }}</span>
            <span v-if="service.price" class="text-sm text-muted-foreground shrink-0">{{ service.price }} ₽</span>
          </div>
        </div>
      </div>

      <div v-if="mentor.reviews?.length" class="bg-card rounded-sm terminal-card border p-6">
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
              <span>{{ formatShortDate(review.date) }}</span>
            </div>
            <p class="text-sm">
              {{ review.text }}
            </p>
          </div>
        </div>
      </div>

      <div v-if="mentor.services?.length && isSubscribed" class="bg-card rounded-sm terminal-card border p-6">
        <ReviewForm
          :mentor-id="mentor.id"
          :services="mentor.services"
          @submitted="loadMentor"
        />
      </div>
    </div>
  </div>
</template>
