<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import { Smile, X } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { feedbackService } from '@/services/feedbackService'

const COOLDOWN_MS = 30 * 24 * 60 * 60 * 1000 // 30 дней

const dismissedAt = useLocalStorage<number | null>('nps_dismissed_at', null)
const isOpen = ref(false)
const score = ref<number | null>(null)
const comment = ref('')
const isSubmitting = ref(false)

const { toast } = useToast()

const isHidden = computed(() => {
  if (!dismissedAt.value)
    return false
  return Date.now() - dismissedAt.value < COOLDOWN_MS
})

function open() {
  if (isHidden.value)
    return
  isOpen.value = true
  score.value = null
  comment.value = ''
}

// dismiss — пользователь явно отказался (кнопка «Отменить») → 30 дней не дёргаем.
function dismiss() {
  isOpen.value = false
  dismissedAt.value = Date.now()
}

// closeSoft — закрытие по крестику или фону, без long-cooldown'а;
// чтобы случайный клик не выключал виджет на месяц.
function closeSoft() {
  isOpen.value = false
}

function handleBackdropClick(event: MouseEvent) {
  if (event.target === event.currentTarget) {
    closeSoft()
  }
}

async function submit() {
  if (score.value === null || isSubmitting.value)
    return
  isSubmitting.value = true
  try {
    await feedbackService.submit(score.value, comment.value)
    toast({ title: 'Спасибо за отзыв' })
    dismissedAt.value = Date.now()
    isOpen.value = false
  }
  catch {
    // handleError в сервисе уже показал ошибку
  }
  finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <button
    v-if="!isHidden"
    type="button"
    aria-label="Оставить отзыв"
    class="fixed bottom-[calc(1.5rem+env(safe-area-inset-bottom))] right-[calc(1.5rem+env(safe-area-inset-right))] z-40 size-12 rounded-full bg-accent text-accent-foreground shadow-lg flex items-center justify-center hover:scale-105 transition-transform cursor-pointer"
    @click="open"
  >
    <Smile class="size-6" />
  </button>

  <Transition name="nps-fade">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center backdrop-blur-sm p-2 sm:p-4"
      @click="handleBackdropClick"
    >
      <Transition name="nps-scale">
        <div v-if="isOpen" class="bg-card text-card-foreground rounded-sm p-4 sm:p-6 w-full max-w-md relative shadow-xl">
          <button
            type="button"
            class="absolute right-4 top-4 text-muted-foreground hover:text-foreground cursor-pointer"
            aria-label="Закрыть"
            @click="closeSoft"
          >
            <X class="size-6" />
          </button>

          <Typography variant="h3" as="h2" class="mb-1">
            Оцените платформу
          </Typography>
          <p class="text-sm text-muted-foreground mb-4">
            Насколько вероятно, что вы порекомендуете нас другу?
          </p>

          <!-- 11 кнопок в один ряд при любой ширине: на узких экранах
               получаются маленькие, но не переносятся (NPS-шкала визуально
               важнее размера кнопок). min-h обеспечивает достаточный
               touch-target. -->
          <div class="grid grid-cols-11 gap-0.5 sm:gap-1 mb-2">
            <button
              v-for="n in 11"
              :key="n - 1"
              type="button"
              :aria-label="`Оценка ${n - 1}`"
              class="min-h-10 px-0 text-xs sm:text-sm rounded-sm border border-input transition-colors cursor-pointer"
              :class="score === n - 1
                ? 'bg-accent text-accent-foreground border-accent'
                : 'hover:bg-secondary'"
              @click="score = n - 1"
            >
              {{ n - 1 }}
            </button>
          </div>
          <div class="flex justify-between text-[10px] sm:text-xs text-muted-foreground mb-4">
            <span>Точно нет</span>
            <span>Точно да</span>
          </div>

          <div class="mb-4">
            <label for="nps-comment" class="block text-sm font-medium text-muted-foreground mb-2">
              Комментарий (необязательно)
            </label>
            <textarea
              id="nps-comment"
              v-model="comment"
              rows="3"
              class="w-full px-3 py-2 border border-input rounded-sm bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="Что можно улучшить?"
            />
          </div>

          <div class="flex justify-end gap-3">
            <button
              type="button"
              class="px-4 py-2 border border-input rounded-full hover:bg-secondary transition duration-300 cursor-pointer"
              @click="dismiss"
            >
              Отменить
            </button>
            <button
              type="button"
              class="px-4 py-2 bg-accent text-accent-foreground rounded-full hover:bg-accent/90 transition duration-300 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="score === null || isSubmitting"
              @click="submit"
            >
              Отправить
            </button>
          </div>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<style scoped>
.nps-fade-enter-active,
.nps-fade-leave-active {
  transition: opacity 0.3s ease;
}

.nps-fade-enter-from,
.nps-fade-leave-to {
  opacity: 0;
}

.nps-scale-enter-active {
  transition: all 0.3s ease-out;
}

.nps-scale-leave-active {
  transition: all 0.2s ease-in;
}

.nps-scale-enter-from,
.nps-scale-leave-to {
  transform: scale(0.95);
  opacity: 0;
}
</style>
