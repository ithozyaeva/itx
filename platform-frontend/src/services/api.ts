import type { TelegramUser } from '@/models/profile'
import ky from 'ky'
import { useToken } from '@/composables/useToken'
import { useUser } from '@/composables/useUser'
import { handleError } from './errorService'

const localStorageUser = useUser()
const localStorageToken = useToken()

let isRefreshing = false
let refreshPromise: Promise<{ token: string, user: TelegramUser }> | null = null

function doRefresh(): Promise<{ token: string, user: TelegramUser }> {
  if (!localStorageUser.value) {
    return Promise.reject(new Error('No user'))
  }
  const userId = localStorageUser.value.telegramID
  if (!userId) {
    return Promise.reject(new Error('No user id'))
  }
  return ky.post('/api/auth/telegram/refresh', {
    json: {
      token: btoa(userId.toString()),
    },
  }).json<{ token: string, user: TelegramUser }>()
}

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
    afterResponse: [async (request, _, response) => {
      if (response.ok)
        return response

      if (response.status === 401) {
        try {
          if (!isRefreshing) {
            isRefreshing = true
            refreshPromise = doRefresh().finally(() => {
              isRefreshing = false
            })
          }

          const { token, user } = await refreshPromise!

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
        }
        catch (e) {
          handleError(e)
        }
      }

      return response
    }],
  },
})

let refreshInterval: ReturnType<typeof setInterval> | null = null

export function startProactiveRefresh() {
  stopProactiveRefresh()
  refreshInterval = setInterval(async () => {
    try {
      const { token, user } = await doRefresh()
      localStorageToken.value = token
      localStorageUser.value = user
    }
    catch {
      // silent fail — reactive refresh on 401 will handle it
    }
  }, 20 * 60 * 1000) // every 20 minutes
}

export function stopProactiveRefresh() {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
}
