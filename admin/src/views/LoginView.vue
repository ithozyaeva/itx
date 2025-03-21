<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '@/composables/useAuth'
import { requiredRule, useFormValidation } from '@/composables/useFormValidation'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const { isLoading, isAuthenticated, login } = useAuth()
const validationRules = {
  login: [requiredRule],
  password: [requiredRule],
}

const { values, errors, touched, validate, handleBlur, isValid } = useFormValidation({
  login: '',
  password: '',
}, validationRules)

const router = useRouter()

/** Сабмиттер формы. */
async function handleSubmit(e: Event) {
  e.preventDefault()

  // Валидируем форму перед отправкой
  if (!validate()) {
    return
  }

  try {
    await login({
      login: values.value.login,
      password: values.value.password,
    })
  }
  catch (error) {
    // Ошибка уже обрабатывается в composable useAuth
    console.error('Login error:', error)
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
          Введите ваши учетные данные для доступа к панели администратора
        </CardDescription>
      </CardHeader>
      <form @submit="handleSubmit">
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="login">Логин</Label>
            <Input
              id="login"
              v-model="values.login"
              type="text"
              placeholder="Введите ваш логин"
              :class="{ 'border-red-500': touched.login && errors.login }"
              @blur="handleBlur('login')"
            />
            <p v-if="touched.login && errors.login" class="text-sm text-red-500">
              {{ errors.login }}
            </p>
          </div>
          <div class="space-y-2">
            <Label for="password">Пароль</Label>
            <Input
              id="password"
              v-model="values.password"
              type="password"
              placeholder="Введите ваш пароль"
              :class="{ 'border-red-500': touched.password && errors.password }"
              @blur="handleBlur('password')"
            />
            <p v-if="touched.password && errors.password" class="text-sm text-red-500">
              {{ errors.password }}
            </p>
          </div>
        </CardContent>
        <CardFooter>
          <Button
            type="submit"
            class="w-full"
            :disabled="isLoading || !isValid"
          >
            {{ isLoading ? 'Вход...' : 'Войти' }}
          </Button>
        </CardFooter>
      </form>
    </Card>
  </div>
</template>
