import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { pointsService } from '@/services/points'

describe('pointsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getMyPoints', () => {
    it('should call GET points/me and return summary', async () => {
      const summary = { total: 100, history: [] }
      mockJson.mockResolvedValue(summary)

      const result = await pointsService.getMyPoints()

      expect(mockApiClient.get).toHaveBeenCalledWith('points/me')
      expect(result).toEqual(summary)
    })
  })

  describe('getLeaderboard', () => {
    it('should call GET points/leaderboard with default limit', async () => {
      const leaderboard = { items: [{ id: 1, points: 200 }] }
      mockJson.mockResolvedValue(leaderboard)

      const result = await pointsService.getLeaderboard()

      expect(mockApiClient.get).toHaveBeenCalledWith('points/leaderboard', { searchParams: { limit: 20 } })
      expect(result).toEqual(leaderboard)
    })

    it('should call GET points/leaderboard with custom limit', async () => {
      const leaderboard = { items: [] }
      mockJson.mockResolvedValue(leaderboard)

      const result = await pointsService.getLeaderboard(50)

      expect(mockApiClient.get).toHaveBeenCalledWith('points/leaderboard', { searchParams: { limit: 50 } })
      expect(result).toEqual(leaderboard)
    })
  })
})
