import ky from 'ky'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import { handleError } from '@/services/errorService'

export interface LoginCredentials {
  login: string
  password: string
}

export interface TokenResponse {
  token: string
}

export interface TelegramAuthResponse {
  user: any
  token: string
}

export const isLoading = ref(false)
export const isAuthenticated = ref(!!localStorage.getItem('token') || !!localStorage.getItem('tg_token'))
const { toast } = useToast()

/**
 * Авторизация пользователя
 *
 * @param credentials - Учетные данные пользователя
 */
export async function login(credentials: LoginCredentials): Promise<string | null> {
  try {
    isLoading.value = true
    const response = await ky.post('/api/auth/login', {
      json: credentials,
    }).json<TokenResponse>()

    localStorage.setItem('token', response.token)
    isAuthenticated.value = true

    toast({
      title: 'Успешный вход',
      description: 'Вы успешно вошли в систему',
    })

    return response.token
  }
  catch (error) {
    let errorMessage = 'Произошла ошибка при входе'

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
 * Авторизация через Telegram
 *
 * @param token - Telegram токен
 */
export async function loginWithTelegram(token: string): Promise<TelegramAuthResponse | null> {
  try {
    isLoading.value = true
    const response = await ky.post('/api/auth/telegram', {
      json: { token },
    }).json<TelegramAuthResponse>()

    // Store both tokens if needed
    localStorage.setItem('token', response.token)
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
 * Обновление JWT токена
 *
 * @param token - Текущий JWT токен
 */
export async function refreshToken(token: string): Promise<string | null> {
  try {
    const response = await ky.post('/api/auth/refresh', {
      json: { token },
      hooks: {
        beforeRequest: [
          (request) => {
            request.headers.delete('Authorization')
          },
        ],
      },
    }).json<TokenResponse>()

    localStorage.setItem('token', response.token)
    isAuthenticated.value = true

    return response.token
  }
  catch (error) {
    handleError(error)
    return null
  }
}

/**
 * Проверка валидности токена
 *
 * @param token - JWT токен
 */
export function isTokenValid(token: string): boolean {
  if (!token)
    return false

  try {
    const payload = JSON.parse(atob(token.split('.')[1]))

    const expirationTime = payload.exp * 1000
    return Date.now() < expirationTime
  }
  catch (error) {
    console.error('Ошибка при проверке токена:', error)
    return false
  }
}

/**
 * Выход из системы
 */
export function logout(): void {
  localStorage.removeItem('token')
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
  const jwtToken = localStorage.getItem('token')
  const tg_token = localStorage.getItem('tg_token')

  // Check either token type
  if ((jwtToken && isTokenValid(jwtToken)) || tg_token) {
    isAuthenticated.value = true
    return true
  }

  isAuthenticated.value = false
  return false
}
