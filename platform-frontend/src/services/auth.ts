import type { TelegramUser } from '@/models/profile'
import ky from 'ky'

export const authService = {
  async authenticate(token: string): Promise<{ user: TelegramUser, token: string }> {
    const response = await ky.post(`/api/auth/telegram`, { json: { token } })
    return await response.json()
  },

  // authenticateWebApp обменивает window.Telegram.WebApp.initData на тот же
  // tg_token, что выпускает /telegram. Бэкенд валидирует HMAC-подпись по
  // бот-токену, поэтому подделать тело со стороны браузера нельзя.
  async authenticateWebApp(initData: string): Promise<{ user: TelegramUser, token: string }> {
    const response = await ky.post(`/api/auth/telegram-webapp`, { json: { init_data: initData } })
    return await response.json()
  },

  clearAuthHeader() {
    localStorage.removeItem('tg_token')
  },

  getBotUrl(): string {
    return `https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_NAME}?start=from_site`
  },
}
