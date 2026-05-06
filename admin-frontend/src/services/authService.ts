import ky from 'ky'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'

export interface TelegramAuthResponse {
  user: any
  token: string
}

export const isLoading = ref(false)
export const isAuthenticated = ref(!!localStorage.getItem('tg_token'))

// hasAdminAccess: undefined = ещё не проверяли, true/false — закэшированный
// результат проверки прав. Сбрасывается при logout/login.
let hasAdminAccess: boolean | undefined

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
    hasAdminAccess = undefined

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
  hasAdminAccess = undefined

  toast({
    title: 'Выход из системы',
    description: 'Вы успешно вышли из системы',
  })
}

/**
 * Проверяет, есть ли у пользователя хоть какие-то админ-права. Без этой
 * проверки любой Telegram-юзер с валидным токеном попадал на /dashboard,
 * хотя API на /api/admin/* всё равно отвечал 403 — UI просто пустел.
 *
 * Бэкенд возвращает [] для не-админов; считаем «нет прав». 401 → токен
 * невалиден, выкидываем как обычно. Любая иная ошибка — fail-closed.
 */
export async function ensureAdminAccess(): Promise<boolean> {
  if (hasAdminAccess !== undefined) {
    return hasAdminAccess
  }
  const tgToken = localStorage.getItem('tg_token')
  if (!tgToken) {
    hasAdminAccess = false
    return false
  }
  try {
    const permissions = await ky.get('/api/admin/me/permissions', {
      headers: { 'X-Telegram-User-Token': tgToken },
    }).json<string[]>()
    hasAdminAccess = Array.isArray(permissions) && permissions.length > 0
    return hasAdminAccess
  }
  catch {
    hasAdminAccess = false
    return false
  }
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
