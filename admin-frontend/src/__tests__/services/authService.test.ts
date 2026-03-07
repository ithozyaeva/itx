import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock useToast
const mockToast = vi.fn()
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({ toast: mockToast }),
}))

// Mock ky
const mockJson = vi.fn()
const mockKy = {
  post: vi.fn(() => ({ json: mockJson })),
}
vi.mock('ky', () => ({ default: mockKy }))

const { isAuthenticated, isLoading, loginWithTelegram, logout, checkAuth } = await import('@/services/authService')

describe('authService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    isAuthenticated.value = false
    isLoading.value = false
  })

  describe('loginWithTelegram', () => {
    it('sends token to API and stores it', async () => {
      const mockResponse = { user: { id: 1, name: 'Test' }, token: 'jwt-token' }
      mockJson.mockResolvedValueOnce(mockResponse)

      const result = await loginWithTelegram('tg-auth-token')

      expect(mockKy.post).toHaveBeenCalledWith('/api/auth/telegram', {
        json: { token: 'tg-auth-token' },
      })
      expect(result).toEqual(mockResponse)
      expect(localStorage.getItem('tg_token')).toBe('tg-auth-token')
      expect(isAuthenticated.value).toBe(true)
    })

    it('shows success toast on login', async () => {
      mockJson.mockResolvedValueOnce({ user: {}, token: 'x' })

      await loginWithTelegram('token')

      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешный вход',
        description: 'Вы успешно вошли через Telegram',
      })
    })

    it('returns null on error', async () => {
      mockJson.mockRejectedValueOnce(new Error('Auth failed'))

      const result = await loginWithTelegram('bad-token')

      expect(result).toBeNull()
      expect(isAuthenticated.value).toBe(false)
    })

    it('shows error toast on failure', async () => {
      mockJson.mockRejectedValueOnce(new Error('Auth failed'))

      await loginWithTelegram('bad-token')

      expect(mockToast).toHaveBeenCalledWith({
        title: 'Ошибка входа',
        description: 'Auth failed',
        variant: 'destructive',
      })
    })

    it('shows generic error message for non-Error objects', async () => {
      mockJson.mockRejectedValueOnce('string error')

      await loginWithTelegram('bad-token')

      expect(mockToast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Ошибка входа',
          description: expect.stringContaining('ошибка'),
        }),
      )
    })

    it('sets and resets isLoading', async () => {
      let loadingDuringRequest = false
      mockJson.mockImplementationOnce(() => {
        loadingDuringRequest = isLoading.value
        return Promise.resolve({ user: {}, token: 'x' })
      })

      await loginWithTelegram('token')

      expect(loadingDuringRequest).toBe(true)
      expect(isLoading.value).toBe(false)
    })

    it('resets isLoading even on error', async () => {
      mockJson.mockRejectedValueOnce(new Error('fail'))

      await loginWithTelegram('token')

      expect(isLoading.value).toBe(false)
    })
  })

  describe('logout', () => {
    it('removes token from localStorage', () => {
      localStorage.setItem('tg_token', 'some-token')
      isAuthenticated.value = true

      logout()

      expect(localStorage.getItem('tg_token')).toBeNull()
      expect(isAuthenticated.value).toBe(false)
    })

    it('shows logout toast', () => {
      logout()

      expect(mockToast).toHaveBeenCalledWith({
        title: 'Выход из системы',
        description: 'Вы успешно вышли из системы',
      })
    })
  })

  describe('checkAuth', () => {
    it('returns true when token exists', () => {
      localStorage.setItem('tg_token', 'valid-token')

      const result = checkAuth()

      expect(result).toBe(true)
      expect(isAuthenticated.value).toBe(true)
    })

    it('returns false when no token', () => {
      const result = checkAuth()

      expect(result).toBe(false)
      expect(isAuthenticated.value).toBe(false)
    })

    it('sets isAuthenticated to false when token is removed', () => {
      isAuthenticated.value = true

      const result = checkAuth()

      expect(result).toBe(false)
      expect(isAuthenticated.value).toBe(false)
    })
  })
})
