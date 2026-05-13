// Auth-логика, портированная из landing-frontend/src/services/auth.ts +
// landing-frontend/src/composables/useUser.ts. Без vue/vueuse — голый TS.

export interface TelegramUser {
  id: number
  telegramID: number
  tg: string
  firstName: string
  lastName: string
  avatarUrl?: string
}

const TG_USER_KEY = 'tg_user:v2'
const TG_TOKEN_KEY = 'tg_token'
const PRIVACY_KEY = 'confirmed_privacy'
const TG_USER_TTL_MS = 60 * 60 * 1000

const BOT_NAME = (import.meta.env.PUBLIC_TELEGRAM_BOT_NAME as string) || 'it_hozyaeva_bot'

export function getBotUrl(): string {
  return `https://t.me/${BOT_NAME}?start=from_site`
}

export function getDeepLinkUrl(): string {
  return `tg://resolve?domain=${BOT_NAME}&start=from_site`
}

export function getToken(): string | null {
  return localStorage.getItem(TG_TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TG_TOKEN_KEY, token)
}

export function getUser(): TelegramUser | null {
  try {
    const raw = localStorage.getItem(TG_USER_KEY)
    if (!raw)
      return null
    const parsed = JSON.parse(raw) as { data: TelegramUser, savedAt: number } | null
    if (!parsed || typeof parsed.savedAt !== 'number' || Date.now() - parsed.savedAt > TG_USER_TTL_MS)
      return null
    return parsed.data
  }
  catch {
    return null
  }
}

export function setUser(user: TelegramUser): void {
  localStorage.setItem(TG_USER_KEY, JSON.stringify({ data: user, savedAt: Date.now() }))
}

export function isPrivacyConfirmed(): boolean {
  try {
    return localStorage.getItem(PRIVACY_KEY) === 'true'
  }
  catch {
    return false
  }
}

export function setPrivacyConfirmed(): void {
  localStorage.setItem(PRIVACY_KEY, 'true')
}

export function openTelegramBot(): void {
  const start = Date.now()
  window.location.href = getDeepLinkUrl()

  setTimeout(() => {
    if (Date.now() - start < 3000)
      window.open(getBotUrl(), '_blank')
  }, 2000)
}

export async function authenticate(token: string): Promise<{ user: TelegramUser, token: string } | null> {
  try {
    const res = await fetch('/api/auth/telegram', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token }),
    })
    if (!res.ok)
      return null
    return await res.json()
  }
  catch (e) {
    console.error('Authentication failed:', e)
    return null
  }
}

export function reachGoal(name: string, params?: Record<string, unknown>): void {
  const ym = (window as any).ym
  const id = import.meta.env.PUBLIC_YANDEX_METRIKA_ID
  if (typeof ym === 'function' && id)
    ym(id, 'reachGoal', name, params)
}

export function extLink(url: string, params?: Record<string, unknown>): void {
  const ym = (window as any).ym
  const id = import.meta.env.PUBLIC_YANDEX_METRIKA_ID
  if (typeof ym === 'function' && id)
    ym(id, 'extLink', url, params)
}
