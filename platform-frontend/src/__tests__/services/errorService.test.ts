import { describe, expect, it, vi } from 'vitest'

// Mock the toast composable before importing the module under test
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({
    toast: vi.fn(),
  }),
}))

import { ErrorType, handleError } from '@/services/errorService'

describe('errorService', () => {
  describe('handleError', () => {
    it('handles HTTPError with 401 status as authentication error', () => {
      const error = {
        name: 'HTTPError',
        response: { status: 401 },
      }
      const result = handleError(error)
      expect(result.type).toBe(ErrorType.AUTHENTICATION)
      expect(result.message).toBe('Ошибка аутентификации. Пожалуйста, войдите снова.')
      expect(result.originalError).toBe(error)
    })

    it('handles HTTPError with 403 status as authentication error', () => {
      const error = {
        name: 'HTTPError',
        response: { status: 403 },
      }
      const result = handleError(error)
      expect(result.type).toBe(ErrorType.AUTHENTICATION)
    })

    it('handles HTTPError with 400 status as validation error', () => {
      const error = {
        name: 'HTTPError',
        response: { status: 400 },
      }
      const result = handleError(error)
      expect(result.type).toBe(ErrorType.VALIDATION)
      expect(result.message).toBe('Проверьте правильность введенных данных')
    })

    it('handles HTTPError with 500 status as server error', () => {
      const error = {
        name: 'HTTPError',
        response: { status: 500 },
      }
      const result = handleError(error)
      expect(result.type).toBe(ErrorType.SERVER)
      expect(result.message).toBe('Ошибка сервера. Пожалуйста, попробуйте позже.')
    })

    it('handles HTTPError with 502 status as server error', () => {
      const error = {
        name: 'HTTPError',
        response: { status: 502 },
      }
      const result = handleError(error)
      expect(result.type).toBe(ErrorType.SERVER)
    })

    it('handles HTTPError with unknown status as unknown error', () => {
      const error = {
        name: 'HTTPError',
        response: { status: 418 },
      }
      const result = handleError(error)
      expect(result.type).toBe(ErrorType.UNKNOWN)
      expect(result.message).toBe('Произошла неизвестная ошибка')
    })

    it('handles NetworkError', () => {
      const error = {
        name: 'NetworkError',
      }
      const result = handleError(error)
      expect(result.type).toBe(ErrorType.NETWORK)
      expect(result.message).toBe('Проблемы с сетевым подключением')
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
      expect(result.message).toBe('Произошла неизвестная ошибка')
    })
  })

  describe('ErrorType enum', () => {
    it('has expected values', () => {
      expect(ErrorType.NETWORK).toBe('network')
      expect(ErrorType.VALIDATION).toBe('validation')
      expect(ErrorType.AUTHENTICATION).toBe('authentication')
      expect(ErrorType.AUTHORIZATION).toBe('authorization')
      expect(ErrorType.SERVER).toBe('server')
      expect(ErrorType.UNKNOWN).toBe('unknown')
    })
  })
})
