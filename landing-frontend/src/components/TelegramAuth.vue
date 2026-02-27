<script setup lang="ts">
import type { TelegramUser } from '@/services/auth'
import { Button } from 'itx-ui-kit'
import { onMounted, ref } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'
import ConfirmServiceModal from '@/components/ConfirmServiceModal.vue'
import { useConfirmedPrivacy } from '@/composables/useUser'
import { authService } from '@/services/auth'

const props = defineProps<{ variant?: 'filled' | 'dark-filled' }>()

const emit = defineEmits<{
  (e: 'auth', user: TelegramUser, token: string): void
}>()

const isConfirmedPrivacy = useConfirmedPrivacy()
const yandexMetrika = useYandexMetrika()

const isModalOpen = ref<boolean>(false)

const btnVariant = props.variant || 'filled'

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token')

  if (token) {
    authService
      .authenticate(token)
      .then(({ user, token: authToken }) => {
        authService.setAuthHeader(authToken)
        emit('auth', user, authToken)
        window.history.replaceState({}, document.title, window.location.pathname)
      })
      .catch((error) => {
        console.error('Authentication failed:', error)
      })
  }
})

function handleClick() {
  if (isConfirmedPrivacy.value) {
    const botUrl = authService.getBotUrl()
    yandexMetrika.reachGoal('become_itx_owner_click', {
      location: props.variant || 'default',
    } as any)
    yandexMetrika.extLink(botUrl, { title: 'Стать IT-Хозяином' })
    window.open(botUrl, '_blank')
  }
  else {
    isModalOpen.value = true
  }
}
</script>

<template>
  <div class="flex items-center h-full">
    <Button
      :variant="btnVariant"
      as="a"
      class="h-full"
      @click="handleClick"
    >
      <slot>Стать IT-хозяином</slot>
    </Button>
  </div>
  <ConfirmServiceModal
    :is-open="isModalOpen"
    @close="isModalOpen = false"
  />
</template>

<style scoped>
#telegram-login {
  min-height: 40px;
}
</style>
