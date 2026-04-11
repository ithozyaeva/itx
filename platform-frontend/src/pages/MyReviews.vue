<script setup lang="ts">
import type { ReviewOnCommunity } from '@/services/reviews'
import { Loader2, MessageSquare, Pencil, Plus, Trash2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import ReviewModal from '@/components/ReviewModal.vue'
import { Button } from '@/components/ui/button'
import { Typography } from '@/components/ui/typography'
import { handleError } from '@/services/errorService'
import { reviewService } from '@/services/reviews'

const reviews = ref<ReviewOnCommunity[]>([])
const isLoading = ref(false)
const loadError = ref<string | null>(null)
const editingId = ref<number | null>(null)
const editText = ref('')
const isModalOpen = ref(false)
const isSaving = ref(false)

async function loadReviews() {
  isLoading.value = true
  loadError.value = null
  try {
    reviews.value = await reviewService.getMyReviews()
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

async function handleCreateReview(text: string) {
  if (isSaving.value)
    return
  isSaving.value = true
  try {
    await reviewService.createReview(text)
    isModalOpen.value = false
    await loadReviews()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSaving.value = false
  }
}

function startEdit(review: ReviewOnCommunity) {
  editingId.value = review.id
  editText.value = review.text
}

function cancelEdit() {
  editingId.value = null
  editText.value = ''
}

async function saveEdit(id: number) {
  if (isSaving.value)
    return
  isSaving.value = true
  try {
    await reviewService.updateReview(id, editText.value)
    cancelEdit()
    await loadReviews()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSaving.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await reviewService.deleteReview(id)
    reviews.value = reviews.value.filter(r => r.id !== id)
  }
  catch (error) {
    handleError(error)
  }
}

const statusLabels: Record<string, string> = {
  DRAFT: 'На модерации',
  APPROVED: 'Опубликован',
}

onMounted(loadReviews)
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Мои отзывы
      </Typography>
      <Button
        @click="isModalOpen = true"
      >
        <Plus class="h-4 w-4 mr-1" />
        Добавить отзыв
      </Button>
    </div>

    <div v-if="isLoading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="loadReviews"
    />

    <EmptyState
      v-else-if="reviews.length === 0"
      :icon="MessageSquare"
      title="Отзывов пока нет"
      description="Ваши отзывы о сообществе появятся здесь"
      action-label="Добавить отзыв"
      @action="isModalOpen = true"
    />

    <div v-else class="space-y-4">
      <div
        v-for="review in reviews"
        :key="review.id"
        class="bg-card border border-border rounded-2xl p-4"
      >
        <div class="flex items-center justify-between mb-2">
          <span
            class="text-xs px-2 py-1 rounded-full"
            :class="review.status === 'APPROVED' ? 'bg-green-500/10 text-green-600' : 'bg-yellow-500/10 text-yellow-600'"
          >
            {{ statusLabels[review.status] }}
          </span>
          <span class="text-xs text-muted-foreground">
            {{ new Date(review.date).toLocaleDateString() }}
          </span>
        </div>

        <template v-if="editingId === review.id">
          <textarea
            v-model="editText"
            rows="3"
            class="w-full px-3 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring resize-none text-sm"
          />
          <div class="flex gap-2 mt-2">
            <Button
              size="sm"
              :disabled="isSaving"
              @click="saveEdit(review.id)"
            >
              Сохранить
            </Button>
            <Button
              size="sm"
              variant="outline"
              @click="cancelEdit"
            >
              Отмена
            </Button>
          </div>
        </template>

        <template v-else>
          <p class="text-sm mb-3">
            {{ review.text }}
          </p>
          <div class="flex gap-2">
            <Button
              size="sm"
              variant="ghost"
              @click="startEdit(review)"
            >
              <Pencil class="h-4 w-4 mr-1" />
              Редактировать
            </Button>
            <Button
              size="sm"
              variant="ghost"
              class="text-destructive hover:text-destructive"
              @click="handleDelete(review.id)"
            >
              <Trash2 class="h-4 w-4 mr-1" />
              Удалить
            </Button>
          </div>
        </template>
      </div>
    </div>
    <ReviewModal
      :is-open="isModalOpen"
      @close="isModalOpen = false"
      @save="handleCreateReview"
    />
  </div>
</template>
