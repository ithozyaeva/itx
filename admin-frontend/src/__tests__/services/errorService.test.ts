import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockToast = vi.fn()
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({ toast: mockToast }),
}))

const { handleError, ErrorType } = await import('@/services/errorService')

describe('errorService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  describe('handleError', () => {
    it('handles HTTPError 401 as authentication error', () => {
      const error = { name: 'HTTPError', response: { status: 401 } }

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.AUTHENTICATION)
      expect(result.message).toContain('аутентификации')
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Ошибка',
        description: result.message,
        variant: 'destructive',
      })
    })

    it('handles HTTPError 403 as authentication error', () => {
      const error = { name: 'HTTPError', response: { status: 403 } }

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.AUTHENTICATION)
    })

    it('handles HTTPError 400 as validation error', () => {
      const error = { name: 'HTTPError', response: { status: 400 } }

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.VALIDATION)
      expect(result.message).toContain('данных')
    })

    it('handles HTTPError 500+ as server error', () => {
      const error = { name: 'HTTPError', response: { status: 500 } }

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.SERVER)
      expect(result.message).toContain('сервера')
    })

    it('handles HTTPError 502 as server error', () => {
      const error = { name: 'HTTPError', response: { status: 502 } }

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.SERVER)
    })

    it('handles HTTPError with unknown status as unknown error', () => {
      const error = { name: 'HTTPError', response: { status: 418 } }

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.UNKNOWN)
    })

    it('handles NetworkError', () => {
      const error = { name: 'NetworkError' }

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.NETWORK)
      expect(result.message).toContain('сетевым')
    })

    it('handles generic Error with message', () => {
      const error = new Error('Something went wrong')

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.UNKNOWN)
      expect(result.message).toBe('Something went wrong')
    })

    it('handles unknown error without message', () => {
      const error = {}

      const result = handleError(error)

      expect(result.type).toBe(ErrorType.UNKNOWN)
      expect(result.message).toContain('неизвестная')
    })

    it('preserves originalError', () => {
      const error = new Error('test')

      const result = handleError(error)

      expect(result.originalError).toBe(error)
    })

    it('logs error to console', () => {
      handleError(new Error('test'))

      expect(console.error).toHaveBeenCalledWith('[App Error]', expect.any(Object))
    })

    it('always shows a toast notification', () => {
      handleError(new Error('test'))

      expect(mockToast).toHaveBeenCalledTimes(1)
      expect(mockToast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Ошибка',
          variant: 'destructive',
        }),
      )
    })
  })
})
