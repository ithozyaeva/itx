import type { TelegramUser } from '@/models/profile'
import ky from 'ky'

export const authService = {
  async authenticate(token: string): Promise<{ user: TelegramUser, token: string }> {
    const response = await ky.post(`/api/auth/telegram`, { json: { token } })
    return await response.json()
  },

  clearAuthHeader() {
    localStorage.removeItem('tg_token')
  },

  getBotUrl(): string {
    return `https://t.me/itx_welcome_bot?start=from_site`
  },
}
