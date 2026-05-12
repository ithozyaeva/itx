import type { TelegramUser } from '@/models/profile'
import ky from 'ky'
import { getTelegramWebApp } from '@/composables/useTelegramWebApp'
import { useToken } from '@/composables/useToken'
import { useUser } from '@/composables/useUser'
import { handleError } from './errorService'

const localStorageUser = useUser()
const localStorageToken = useToken()

let isRefreshing = false
let refreshPromise: Promise<{ token: string, user: TelegramUser }> | null = null

function doMiniAppReauth(): Promise<{ token: string, user: TelegramUser }> {
  const tg = getTelegramWebApp()
  if (!tg || !tg.initData) {
    return Promise.reject(new Error('Not in miniapp'))
  }
  return ky.post('/api/auth/telegram-webapp', {
    json: { init_data: tg.initData },
  }).json<{ token: string, user: TelegramUser }>()
}

function doRefresh(): Promise<{ token: string, user: TelegramUser }> {
  // Refresh опирается ТОЛЬКО на текущий X-Telegram-User-Token: знание
  // Telegram-ID не должно давать продлить чужую сессию. Внутри miniapp
  // на провал отката есть второй путь — initData всегда свежий и подписан
  // bot-token'ом, поэтому Telegram-клиент может восстановить сессию без
  // редиректа на лендинг и без тоста «Unauthorized» (раньше юзер видел его
  // на полсекунды между неудачным refresh и тихим re-auth в App.vue).
  const token = localStorageToken.value
  if (!token) {
    return doMiniAppReauth()
  }
  return ky.post('/api/auth/telegram/refresh', {
    headers: {
      'X-Telegram-User-Token': token,
    },
  }).json<{ token: string, user: TelegramUser }>().catch((refreshErr) => {
    const tg = getTelegramWebApp()
    if (tg && tg.initData) {
      return doMiniAppReauth().catch(() => {
        throw refreshErr
      })
    }
    throw refreshErr
  })
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

      // Backend гейтит премиум-эндпоинты для UNSUBSCRIBER через 403
      // {"error":"subscription_required","redirect":"/"}.
      // Уводим юзера на главную, где есть teaser «Открой полный доступ» —
      // лобовое приземление на /tariffs из любого экрана воспринимается как
      // принудительная продажа, особенно если пришли по ссылке из бота.
      if (response.status === 403) {
        try {
          const body = await response.clone().json() as { error?: string, redirect?: string }
          if (body?.error === 'subscription_required') {
            const router = (await import('@/router')).default
            const target = body.redirect || '/'
            if (router.currentRoute.value.path !== target) {
              router.push(target)
            }
          }
        }
        catch {
          // Ignored — body не JSON, передадим ошибку дальше как есть.
        }
      }

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
          // Refresh не помог — токен в localStorage не подходит ни одному
          // живому token в БД (например, после массовой инвалидации в #325
          // или просто протухший после 30 дней). Чистим оба ключа: иначе
          // App.vue при следующем заходе снова видит tg_token, снова дёргает
          // /me → 401 → редирект → бесконечный цикл «платформа на полсекунды
          // → лендос». Лендинг тоже сразу увидит «не залогинен» и покажет
          // TelegramAuth вместо «Перейти в платформу».
          localStorageToken.value = null
          localStorageUser.value = null
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
    if (isRefreshing)
      return
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
