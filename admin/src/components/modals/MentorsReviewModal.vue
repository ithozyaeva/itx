<script setup lang="ts">
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import MentorsReviewForm from '../forms/MentorsReviewForm.vue'

defineProps<{
  isOpen: boolean
  reviewId?: number
}>()

const emit = defineEmits(['update:isOpen', 'saved'])

function handleSaved() {
  emit('saved')
  emit('update:isOpen', false)
}

function handleCancel() {
  emit('update:isOpen', false)
}
</script>

<template>
  <Dialog :open="isOpen" @update:open="$emit('update:isOpen', $event)">
    <DialogContent class="sm:max-w-[700px]">
      <DialogHeader>
        <DialogTitle>{{ reviewId ? 'Редактирование отзыва' : 'Добавление отзыва' }}</DialogTitle>
      </DialogHeader>
      <MentorsReviewForm :review-id="reviewId" @saved="handleSaved" @cancel="handleCancel" />
    </DialogContent>
  </Dialog>
</template>
