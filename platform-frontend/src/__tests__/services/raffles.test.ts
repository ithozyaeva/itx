import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { raffleService } from '@/services/raffles'

describe('raffleService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('should call GET raffles', async () => {
      const raffles = [{ id: 1, name: 'Raffle 1' }]
      mockJson.mockResolvedValue(raffles)

      const result = await raffleService.getAll()

      expect(mockApiClient.get).toHaveBeenCalledWith('raffles')
      expect(result).toEqual(raffles)
    })
  })

  describe('buyTickets', () => {
    it('should call POST raffles/:id/buy with default count', async () => {
      mockJson.mockResolvedValue({})

      await raffleService.buyTickets(10)

      expect(mockApiClient.post).toHaveBeenCalledWith('raffles/10/buy', { json: { count: 1 } })
    })

    it('should call POST raffles/:id/buy with custom count', async () => {
      mockJson.mockResolvedValue({})

      await raffleService.buyTickets(10, 5)

      expect(mockApiClient.post).toHaveBeenCalledWith('raffles/10/buy', { json: { count: 5 } })
    })
  })
})
