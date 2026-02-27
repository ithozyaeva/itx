import ky from 'ky'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'

export interface TelegramAuthResponse {
  user: any
  token: string
}

export const isLoading = ref(false)
export const isAuthenticated = ref(!!localStorage.getItem('tg_token'))

/**
 * Авторизация через Telegram
 *
 * @param token - Telegram токен
 */
export async function loginWithTelegram(token: string): Promise<TelegramAuthResponse | null> {
  const { toast } = useToast()
  try {
    isLoading.value = true
    const response = await ky.post('/api/auth/telegram', {
      json: { token },
    }).json<TelegramAuthResponse>()

    localStorage.setItem('tg_token', token)
    isAuthenticated.value = true

    toast({
      title: 'Успешный вход',
      description: 'Вы успешно вошли через Telegram',
    })

    return response
  }
  catch (error) {
    let errorMessage = 'Произошла ошибка при входе через Telegram'

    if (error instanceof Error) {
      errorMessage = error.message
    }

    toast({
      title: 'Ошибка входа',
      description: errorMessage,
      variant: 'destructive',
    })

    return null
  }
  finally {
    isLoading.value = false
  }
}

/**
 * Выход из системы
 */
export function logout(): void {
  const { toast } = useToast()
  localStorage.removeItem('tg_token')
  isAuthenticated.value = false

  toast({
    title: 'Выход из системы',
    description: 'Вы успешно вышли из системы',
  })
}

/**
 * Проверка авторизации
 */
export function checkAuth(): boolean {
  const tgToken = localStorage.getItem('tg_token')

  if (tgToken) {
    isAuthenticated.value = true
    return true
  }

  isAuthenticated.value = false
  return false
}
