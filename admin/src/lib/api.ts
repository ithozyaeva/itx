import { useToast } from '@/components/ui/toast'
import router from '@/router'
import { isAuthenticated, logout, refreshToken } from '@/services/authService'
import ky from 'ky'

// Флаг для отслеживания процесса обновления токена
let isRefreshing = false
// Промис для хранения текущего процесса обновления
let refreshPromise: Promise<string | null> | null = null

const api = ky.create({
  prefixUrl: import.meta.env.SERVER_URL || 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
  hooks: {
    beforeRequest: [
      (request) => {
        const token = localStorage.getItem('token')
        if (token) {
          request.headers.set('Authorization', `${token}`)
        }
      },
    ],
    afterResponse: [
      async (request, options, response) => {
        if (response.ok)
          return response

        if (response.status === 401) {
          try {
            if (isAuthenticated.value) {
              try {
                const token = localStorage.getItem('token')
                if (!token)
                  throw new Error('Токен не найден')

                if (isRefreshing && refreshPromise) {
                  await refreshPromise
                }
                else {
                  isRefreshing = true
                  refreshPromise = refreshToken(token)
                  await refreshPromise
                  isRefreshing = false
                  refreshPromise = null
                }

                const newToken = localStorage.getItem('token')
                const newRequest = new Request(request.url, {
                  method: request.method,
                  headers: request.headers,
                  body: request.body,
                  credentials: request.credentials,
                  mode: request.mode,
                  cache: request.cache,
                })

                newRequest.headers.set('Authorization', `${newToken}`)

                return fetch(newRequest)
              }
              catch (error) {
                console.error('Не удалось обновить токен:', error)
                logout()
                router.push('/login')

                const { toast } = useToast()
                toast({
                  title: 'Сессия истекла',
                  description: 'Пожалуйста, войдите снова',
                  variant: 'destructive',
                })
              }
            }
            else {
              router.push('/login')
            }
          }
          catch (error) {
            console.error('Ошибка при обработке ответа:', error)
          }
        }

        return response
      },
    ],
  },
})

export default api
