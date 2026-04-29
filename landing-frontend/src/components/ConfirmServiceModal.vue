<script setup lang="ts">
import { X as CloseIcon } from 'lucide-vue-next'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import Button from '@/components/ui/UiButton.vue'
import { useConfirmedPrivacy } from '@/composables/useUser'
import { authService } from '@/services/auth'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const isConfirmedPrivacy = useConfirmedPrivacy()
const yandexMetrika = useYandexMetrika()

function handleConfirmed() {
  isConfirmedPrivacy.value = true
  emit('close')
  yandexMetrika.reachGoal('privacy_modal_accept', {
    action: 'accept',
  } as any)
  yandexMetrika.extLink(authService.getBotUrl(), { title: 'Стать IT-Хозяином (после модалки)' })
  authService.openBot()
}

function handleClose() {
  emit('close')
}

function handleBackdropClick(event: MouseEvent) {
  if (event.target === event.currentTarget) {
    handleClose()
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center backdrop-blur-sm p-4 sm:p-6"
        @click="handleBackdropClick"
      >
        <Transition name="modal-scale">
          <div
            v-if="isOpen"
            class="bg-[#1D2723] rounded-3xl border-2 border-[#2B3D36] p-5 sm:p-9 w-full max-w-3xl relative shadow-xl max-h-[calc(100vh-2rem)] sm:max-h-[calc(100vh-3rem)] overflow-y-auto"
          >
            <button
              class="absolute right-3 top-3 cursor-pointer hover:opacity-75 transition-opacity"
              @click="handleClose"
            >
              <CloseIcon
                class="h-6 w-6 sm:h-8 sm:w-8"
              />
            </button>

            <p
              class="mb-3 pr-8 text-accent font-display uppercase font-semibold tracking-[0.04em] leading-tight text-lg sm:text-2xl"
            >
              Согласие на обработку персональных данных
            </p>

            <div class="mb-5 sm:mb-7">
              <p class="font-sans font-medium leading-[1.4] text-sm sm:text-base md:text-lg">
                Нажимая на кнопку "Принять", Вы соглашаетесь с
                <a
                  href="/privacy"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-accent underline cursor-pointer hover:text-accent/75 transition-colors"
                >
                  политикой конфиденциальности
                </a>
                и даёте согласие на обработку персональных данных согласно Федеральному закону №152-ФЗ.
              </p>
            </div>

            <div class="flex flex-col-reverse sm:flex-row justify-start gap-3 sm:gap-5">
              <Button
                variant="stroke"
                @click="handleClose"
              >
                Отклонить
              </Button>
              <Button
                variant="filled"
                @click="handleConfirmed"
              >
                Принять
              </Button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
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
