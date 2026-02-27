import axios from 'axios'

export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000'

export interface TelegramUser {
  id: number
  telegramID: number
  tg: string
  firstName: string
  lastName: string
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
}
