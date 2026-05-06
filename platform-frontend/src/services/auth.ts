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

  // logout инвалидирует tg_token серверно (POST /api/auth/telegram/logout),
  // затем чистит локалсторадж. До добавления серверного эндпоинта клик
  // «Выйти» только удалял токен из localStorage, но в auth_tokens он жил
  // до natural expiry (~30 дней) — украденный токен оставался валидным.
  // best-effort: ошибки сети не блокируют локальный logout, иначе залогинить
  // обратно нельзя без восстановления связи.
  async logout(): Promise<void> {
    const token = localStorage.getItem('tg_token')
    if (token) {
      try {
        await ky.post('/api/auth/telegram/logout', {
          headers: { 'X-Telegram-User-Token': token },
        })
      }
      catch {
        // Игнорируем — серверный logout best-effort, локальный всегда успешный.
      }
    }
    localStorage.removeItem('tg_token')
  },

  // clearAuthHeader — синхронная версия для аварийных путей (apiClient после
  // неудачного refresh: токен уже точно мёртв, серверный вызов смысла не имеет).
  clearAuthHeader() {
    localStorage.removeItem('tg_token')
  },

  getBotUrl(): string {
    return `https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_NAME}?start=from_site`
  },
}
