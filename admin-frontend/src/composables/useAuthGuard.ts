import type { RouteLocationNormalized } from 'vue-router'
import { useAuth } from './useAuth'

export function useAuthGuard() {
  const { isAuthenticated } = useAuth()

  const guardRoute = (to: RouteLocationNormalized) => {
    // Проверяем, требует ли маршрут аутентификации
    if (to.meta.requiresAuth && !isAuthenticated.value) {
      // Если требуется аутентификация, но пользователь не авторизован,
      // перенаправляем на страницу входа
      return { path: '/login', query: { redirect: to.fullPath } }
    }

    // Если пользователь авторизован и пытается перейти на страницу входа,
    // перенаправляем на дашборд
    if (to.path === '/login' && isAuthenticated.value) {
      return { path: '/dashboard' }
    }

    // В остальных случаях разрешаем переход
    return true
  }

  return {
    guardRoute,
  }
}
