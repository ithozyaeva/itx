import type { TelegramUser } from '@/models/profile'
import { useToken } from '@/composables/useToken'
import { useUser } from '@/composables/useUser'
import ky from 'ky'

const localStorageUser = useUser()
const localStorageToken = useToken()

export const apiClient = ky.create({
  prefixUrl: '/api/platform/',
  retry: 1,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = localStorage.getItem('tg_token')
        if (token) {
          request.headers.set('X-Telegram-User-Token', `${token}`)
        }
      },
    ],
    afterResponse: [(request, _, response) => {
      if (response.ok)
        return response

      if (response.status === 401) {
        try {
          if (!localStorageUser.value) {
            throw new Error('No user')
          }
          const userId = localStorageUser.value.telegramID
          if (!userId) {
            throw new Error('No user id')
          }
          ky.post('/api/auth/telegram/refresh', {
            json: {
              token: btoa(userId?.toString()),
            },
          }).json<{ token: string, user: TelegramUser }>().then(({ token, user }) => {
            localStorageToken.value = token
            localStorageUser.value = user

            const newRequest = new Request(request.url, {
              method: request.method,
              headers: request.headers,
              body: request.body,
              credentials: request.credentials,
              mode: request.mode,
              cache: request.cache,
            })

            newRequest.headers.set('X-Telegram-User-Token', `${token}`)

            return fetch(newRequest)
          })
        }
        catch (e) {
          console.error(e)
        }
      }

      return response
    }],
  },
})
