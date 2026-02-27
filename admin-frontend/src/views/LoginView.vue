<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TelegramAuth from '@/components/TelegramAuth.vue'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { isAuthenticated } from '@/services/authService'

const router = useRouter()

function handleTelegramAuthSuccess() {
  if (isAuthenticated.value) {
    router.push('/dashboard')
  }
}

onMounted(() => {
  if (isAuthenticated.value) {
    router.push('/dashboard')
  }
})
</script>

<template>
  <div class="flex h-screen w-full items-center justify-center">
    <Card class="w-full max-w-md">
      <CardHeader class="space-y-1">
        <CardTitle class="text-2xl font-bold">
          Вход в систему
        </CardTitle>
        <CardDescription>
          Войдите в систему используя ваш Telegram аккаунт
        </CardDescription>
      </CardHeader>
      <CardContent>
        <TelegramAuth @auth-success="handleTelegramAuthSuccess" />
      </CardContent>
    </Card>
  </div>
</template>
