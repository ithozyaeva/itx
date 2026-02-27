import { useRouter } from 'vue-router'
import { checkAuth, isAuthenticated, isLoading, logout as logoutService } from '@/services/authService'

export function useAuth() {
  const router = useRouter()

  // Проверяем авторизацию при инициализации
  checkAuth()

  // Функция для выхода из системы
  function logout() {
    logoutService()
    router.push('/login')
  }

  return {
    isAuthenticated,
    isLoading,
    logout,
  }
}
