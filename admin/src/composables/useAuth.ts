import type { LoginCredentials } from '@/services/authService'
import { checkAuth, isAuthenticated, isLoading, login as loginService, logout as logoutService } from '@/services/authService'
import { useRouter } from 'vue-router'

export function useAuth() {
  const router = useRouter()

  // Проверяем авторизацию при инициализации
  checkAuth()

  // Функция для входа в систему
  async function login(credentials: LoginCredentials) {
    const token = await loginService(credentials)

    if (token) {
      // Перенаправляем на дашборд после успешного входа
      router.push('/dashboard')
    }

    return token
  }

  // Функция для выхода из системы
  function logout() {
    logoutService()
    router.push('/login')
  }

  return {
    isAuthenticated,
    isLoading,
    login,
    logout,
  }
}
