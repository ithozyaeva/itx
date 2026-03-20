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

import { casinoService } from '@/services/casino'

describe('casinoService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('coinFlip', () => {
    it('should call POST minigames/coin-flip with correct params', async () => {
      const result = { id: 1, game: 'coin-flip', won: true }
      mockJson.mockResolvedValue(result)

      const response = await casinoService.coinFlip(50, 'heads')

      expect(mockApiClient.post).toHaveBeenCalledWith('minigames/coin-flip', {
        json: { betAmount: 50, choice: 'heads' },
      })
      expect(response).toEqual(result)
    })

    it('should support tails choice', async () => {
      mockJson.mockResolvedValue({})

      await casinoService.coinFlip(100, 'tails')

      expect(mockApiClient.post).toHaveBeenCalledWith('minigames/coin-flip', {
        json: { betAmount: 100, choice: 'tails' },
      })
    })
  })

  describe('diceRoll', () => {
    it('should call POST minigames/dice-roll with target and direction', async () => {
      const result = { id: 2, game: 'dice-roll', won: false }
      mockJson.mockResolvedValue(result)

      const response = await casinoService.diceRoll(25, 50, 'over')

      expect(mockApiClient.post).toHaveBeenCalledWith('minigames/dice-roll', {
        json: { betAmount: 25, target: 50, direction: 'over' },
      })
      expect(response).toEqual(result)
    })

    it('should support under direction', async () => {
      mockJson.mockResolvedValue({})

      await casinoService.diceRoll(10, 30, 'under')

      expect(mockApiClient.post).toHaveBeenCalledWith('minigames/dice-roll', {
        json: { betAmount: 10, target: 30, direction: 'under' },
      })
    })
  })

  describe('wheelSpin', () => {
    it('should call POST minigames/wheel with betAmount', async () => {
      const result = { id: 3, game: 'wheel', multiplier: 2 }
      mockJson.mockResolvedValue(result)

      const response = await casinoService.wheelSpin(200)

      expect(mockApiClient.post).toHaveBeenCalledWith('minigames/wheel', {
        json: { betAmount: 200 },
      })
      expect(response).toEqual(result)
    })
  })

  describe('getHistory', () => {
    it('should call GET minigames/history and extract items', async () => {
      const history = [{ id: 1 }, { id: 2 }]
      mockJson.mockResolvedValue({ items: history, total: 2 })

      const result = await casinoService.getHistory()

      expect(mockApiClient.get).toHaveBeenCalledWith('minigames/history?limit=20')
      expect(result).toEqual(history)
    })

    it('should call GET minigames/history with custom limit', async () => {
      mockJson.mockResolvedValue({ items: [], total: 0 })

      await casinoService.getHistory(50)

      expect(mockApiClient.get).toHaveBeenCalledWith('minigames/history?limit=50')
    })

    it('should return empty array when items is null', async () => {
      mockJson.mockResolvedValue({ items: null, total: 0 })

      const result = await casinoService.getHistory()

      expect(result).toEqual([])
    })
  })

  describe('getStats', () => {
    it('should call GET minigames/stats', async () => {
      const stats = { balance: 100, totalBets: 5 }
      mockJson.mockResolvedValue(stats)

      const result = await casinoService.getStats()

      expect(mockApiClient.get).toHaveBeenCalledWith('minigames/stats')
      expect(result).toEqual(stats)
    })
  })
})
