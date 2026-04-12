import axios from 'axios'

export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000'

export interface TelegramUser {
  id: number
  telegramID: number
  tg: string
  firstName: string
  lastName: string
  avatarUrl?: string
}

export const authService = {
  async authenticate(token: string): Promise<{ user: TelegramUser, token: string }> {
    const response = await axios.post(`/api/auth/telegram`, { token })
    return response.data
  },

  setAuthHeader(authToken: string) {
    axios.defaults.headers.common['X-Telegram-User-Token'] = authToken
  },

  clearAuthHeader() {
    delete axios.defaults.headers.common['X-Telegram-User-Token']
  },

  getBotUrl(): string {
    return `https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_NAME}?start=from_site`
  },

  getDeepLinkUrl(): string {
    return `tg://resolve?domain=${import.meta.env.VITE_TELEGRAM_BOT_NAME}&start=from_site`
  },

  openBot(): void {
    const deepLink = this.getDeepLinkUrl()
    const webLink = this.getBotUrl()

    // Пробуем открыть через tg:// (работает без VPN если Telegram установлен)
    const start = Date.now()
    window.location.href = deepLink

    // Если через 2 секунды страница не ушла — Telegram не установлен, fallback на t.me
    setTimeout(() => {
      if (Date.now() - start < 3000) {
        window.open(webLink, '_blank')
      }
    }, 2000)
  },
}
