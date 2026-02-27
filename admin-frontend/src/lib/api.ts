import ky from 'ky'
import router from '@/router'
import { isAuthenticated, logout, refreshToken } from '@/services/authService'
import { handleError } from '@/services/errorService'

// Флаг для отслеживания процесса обновления токена
let isRefreshing = false
// Промис для хранения текущего процесса обновления
let refreshPromise: Promise<string | null> | null = null

const api = ky.create({
  prefixUrl: '/api/admin/',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
  hooks: {
    beforeRequest: [
      (request) => {
        const jwtToken = localStorage.getItem('token')
        const tg_token = localStorage.getItem('tg_token')

        if (jwtToken) {
          request.headers.set('Authorization', `Bearer ${jwtToken}`)
        }

        if (tg_token) {
          request.headers.set('X-Telegram-User-Token', tg_token)
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

                newRequest.headers.set('Authorization', `Bearer ${newToken}`)

                return fetch(newRequest)
              }
              catch (error) {
                handleError(error)
                logout()
                router.push('/login')
              }
            }
            else {
              router.push('/login')
            }
          }
          catch (error) {
            handleError(error)
          }
        }

        return response
      },
    ],
  },
})

export default api
