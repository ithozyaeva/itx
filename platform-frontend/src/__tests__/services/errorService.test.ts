import { describe, expect, it, vi } from 'vitest'

// Mock the toast composable before importing the module under test
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({
    toast: vi.fn(),
  }),
}))

import { ErrorType, handleError } from '@/services/errorService'

function mockResponse(status: number, body?: Record<string, unknown>) {
  return {
    status,
    json: vi.fn().mockResolvedValue(body ?? {}),
  }
}

describe('errorService', () => {
  describe('handleError', () => {
    it('handles HTTPError with 401 status as authentication error', async () => {
      const error = {
        name: 'HTTPError',
        response: mockResponse(401),
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.AUTHENTICATION)
      expect(result.message).toBe('Ошибка аутентификации. Пожалуйста, войдите снова.')
      expect(result.originalError).toBe(error)
    })

    it('handles HTTPError with 403 status as authentication error', async () => {
      const error = {
        name: 'HTTPError',
        response: mockResponse(403),
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.AUTHENTICATION)
    })

    it('handles HTTPError with 400 status as validation error', async () => {
      const error = {
        name: 'HTTPError',
        response: mockResponse(400),
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.VALIDATION)
      expect(result.message).toBe('Проверьте правильность введенных данных')
    })

    it('handles HTTPError with server error message', async () => {
      const error = {
        name: 'HTTPError',
        response: mockResponse(400, { error: 'недостаточно баллов' }),
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.VALIDATION)
      expect(result.message).toBe('недостаточно баллов')
    })

    it('handles HTTPError with 500 status as server error', async () => {
      const error = {
        name: 'HTTPError',
        response: mockResponse(500),
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.SERVER)
      expect(result.message).toBe('Ошибка сервера. Пожалуйста, попробуйте позже.')
    })

    it('handles HTTPError with 502 status as server error', async () => {
      const error = {
        name: 'HTTPError',
        response: mockResponse(502),
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.SERVER)
    })

    it('handles HTTPError with unknown status as unknown error', async () => {
      const error = {
        name: 'HTTPError',
        response: mockResponse(418),
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.UNKNOWN)
      expect(result.message).toBe('Произошла неизвестная ошибка')
    })

    it('handles NetworkError', async () => {
      const error = {
        name: 'NetworkError',
      }
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.NETWORK)
      expect(result.message).toBe('Проблемы с сетевым подключением')
    })

    it('handles generic Error with message', async () => {
      const error = new Error('Something went wrong')
      const result = await handleError(error)
      expect(result.type).toBe(ErrorType.UNKNOWN)
      expect(result.message).toBe('Something went wrong')
    })

    it('handles unknown error without message', async () => {
      const error = {}
      const result = await handleError(error)
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
