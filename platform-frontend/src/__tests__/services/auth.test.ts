import { beforeEach, describe, expect, it, vi } from 'vitest'

const { mockJson, mockKy } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockKy: {
      post: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('ky', () => ({
  default: mockKy,
}))

import { authService } from '@/services/auth'

describe('authService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('authenticate', () => {
    it('should call ky.post with correct args and return result', async () => {
      const response = { user: { id: 1 }, token: 'jwt-token' }
      mockJson.mockResolvedValue(response)

      const result = await authService.authenticate('tg-token-123')

      expect(mockKy.post).toHaveBeenCalledWith('/api/auth/telegram', { json: { token: 'tg-token-123' } })
      expect(result).toEqual(response)
    })
  })

  describe('clearAuthHeader', () => {
    it('should remove tg_token from localStorage', () => {
      const removeItemSpy = vi.spyOn(Storage.prototype, 'removeItem')

      authService.clearAuthHeader()

      expect(removeItemSpy).toHaveBeenCalledWith('tg_token')
      removeItemSpy.mockRestore()
    })
  })

  // logout инвалидирует токен на сервере + чистит localStorage. До этого
  // эндпоинта токен оставался валидным в БД на ~30 дней после клика «Выйти»,
  // и компрометированный tg_token продолжал работать.
  describe('logout', () => {
    beforeEach(() => {
      localStorage.clear()
    })

    it('вызывает /api/auth/telegram/logout с токеном из localStorage и чистит storage', async () => {
      localStorage.setItem('tg_token', 'tg-token-xyz')

      await authService.logout()

      expect(mockKy.post).toHaveBeenCalledWith('/api/auth/telegram/logout', {
        headers: { 'X-Telegram-User-Token': 'tg-token-xyz' },
      })
      expect(localStorage.getItem('tg_token')).toBeNull()
    })

    it('пропускает серверный вызов если токена нет', async () => {
      await authService.logout()

      expect(mockKy.post).not.toHaveBeenCalled()
      expect(localStorage.getItem('tg_token')).toBeNull()
    })

    it('best-effort: чистит storage даже если серверный logout упал', async () => {
      localStorage.setItem('tg_token', 'tg-token-zzz')
      mockKy.post.mockImplementationOnce(() => {
        throw new Error('network down')
      })

      await expect(authService.logout()).resolves.toBeUndefined()
      expect(localStorage.getItem('tg_token')).toBeNull()
    })
  })

  describe('getBotUrl', () => {
    it('should return correct bot URL', () => {
      const result = authService.getBotUrl()

      expect(result).toBe(`https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_NAME}?start=from_site`)
    })
  })
})
