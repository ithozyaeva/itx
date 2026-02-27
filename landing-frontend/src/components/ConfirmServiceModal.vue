<script setup lang="ts">
import { Button, CloseIcon, Typography } from 'itx-ui-kit'
import { useYandexMetrika } from 'yandex-metrika-vue3'
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
  const botUrl = authService.getBotUrl()
  yandexMetrika.reachGoal('privacy_modal_accept', {
    action: 'accept',
  } as any)
  yandexMetrika.extLink(botUrl, { title: 'Стать IT-Хозяином (после модалки)' })
  window.open(botUrl, '_blank')
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
  <Transition name="modal-fade">
    <div
      v-if="isOpen"
      class="fixed h-screen inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center backdrop-blur-sm"
      @click="handleBackdropClick"
    >
      <Transition name="modal-scale">
        <div
          v-if="isOpen"
          class="bg-[#1D2723] rounded-3xl border-2 border-[#2B3D36] p-9 w-full mx-5 max-w-3xl  relative shadow-xl"
        >
          <button
            class="absolute right-3 top-3 cursor-pointer hover:opacity-75 transition-opacity"
            @click="handleClose"
          >
            <CloseIcon
              class="h-8 w-8"
            />
          </button>

          <Typography
            variant="h3"
            as="p"
            class="mb-3 text-accent"
          >
            Согласие на обработку персональных данных
          </Typography>

          <div class="mb-7">
            <Typography
              as="p"
              variant="body-l"
            >
              Нажимая на кнопку "Принять", Вы соглашаетесь с
              <Typography
                as="a"
                variant="body-l"
                href="/privacy"
                target="_blank"
                class="text-accent underline cursor-pointer hover:text-accent/75 transition-colors"
              >
                политикой конфиденциальности
              </Typography>
              и даёте согласие на обработку персональных данных согласно Федеральному закону №152-ФЗ.
            </Typography>
          </div>

          <div class="flex flex-col-reverse sm:flex-row justify-start gap-5">
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
