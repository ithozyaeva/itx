<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { X } from 'lucide-vue-next'
import { ref } from 'vue'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
  save: [text: string]
}>()

const reviewText = ref('')

function handleSave() {
  emit('save', reviewText.value)
  reviewText.value = ''
}

function handleClose() {
  emit('close')
  reviewText.value = ''
}

function handleBackdropClick(event: MouseEvent) {
  if (event.target === event.currentTarget) {
    handleClose()
  }
}
</script>

<template>
  <Transition name="modal-fade">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center backdrop-blur-sm"
      @click="handleBackdropClick"
    >
      <Transition name="modal-scale">
        <div v-if="isOpen" class="bg-card text-card-foreground rounded-3xl p-6 w-full max-w-md relative shadow-xl">
          <button
            class="absolute right-4 top-4 text-muted-foreground hover:text-foreground cursor-pointer"
            @click="handleClose"
          >
            <X class="h-6 w-6" />
          </button>

          <Typography variant="h3" as="h2" class="mb-4">
            Добавить отзыв
          </Typography>

          <div class="mb-4">
            <label for="review" class="block text-sm font-medium text-muted-foreground mb-2">
              Текст отзыва
            </label>
            <textarea
              id="review"
              v-model="reviewText"
              rows="4"
              class="w-full px-3 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="Введите ваш отзыв..."
            />
          </div>

          <div class="flex justify-end gap-3">
            <button
              class="px-4 py-2 border border-input rounded-full hover:bg-secondary transition duration-300 cursor-pointer"
              @click="handleClose"
            >
              Отменить
            </button>
            <button
              class="px-4 py-2 bg-accent text-accent-foreground rounded-full hover:bg-accent/90 transition duration-300 cursor-pointer"
              :disabled="!reviewText.trim()"
              @click="handleSave"
            >
              Сохранить
            </button>
          </div>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<style scoped>
/* Анимация фона */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

/* Анимация самого окна */
.modal-scale-enter-active {
  transition: all 0.3s ease-out;
}

.modal-scale-leave-active {
  transition: all 0.2s ease-in;
}

.modal-scale-enter-from,
.modal-scale-leave-to {
  transform: scale(0.95);
  opacity: 0;
}

.modal-scale-enter-to,
.modal-scale-leave-from {
  transform: scale(1);
  opacity: 1;
}
</style>
