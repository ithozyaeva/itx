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

import { seasonService } from '@/services/seasons'

describe('seasonService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('should call GET seasons', async () => {
      const seasons = [{ id: 1, name: 'Season 1' }]
      mockJson.mockResolvedValue(seasons)

      const result = await seasonService.getAll()

      expect(mockApiClient.get).toHaveBeenCalledWith('seasons')
      expect(result).toEqual(seasons)
    })
  })

  describe('getActive', () => {
    it('should call GET seasons/active with default limit', async () => {
      const season = { id: 1, leaderboard: [] }
      mockJson.mockResolvedValue(season)

      const result = await seasonService.getActive()

      expect(mockApiClient.get).toHaveBeenCalledWith('seasons/active', { searchParams: { limit: 20 } })
      expect(result).toEqual(season)
    })

    it('should call GET seasons/active with custom limit', async () => {
      const season = { id: 1, leaderboard: [] }
      mockJson.mockResolvedValue(season)

      const result = await seasonService.getActive(50)

      expect(mockApiClient.get).toHaveBeenCalledWith('seasons/active', { searchParams: { limit: 50 } })
      expect(result).toEqual(season)
    })
  })

  describe('getLeaderboard', () => {
    it('should call GET seasons/:id/leaderboard with default limit', async () => {
      const leaderboard = { id: 3, leaderboard: [{ memberId: 1, points: 100 }] }
      mockJson.mockResolvedValue(leaderboard)

      const result = await seasonService.getLeaderboard(3)

      expect(mockApiClient.get).toHaveBeenCalledWith('seasons/3/leaderboard', { searchParams: { limit: 20 } })
      expect(result).toEqual(leaderboard)
    })

    it('should call GET seasons/:id/leaderboard with custom limit', async () => {
      const leaderboard = { id: 5, leaderboard: [] }
      mockJson.mockResolvedValue(leaderboard)

      const result = await seasonService.getLeaderboard(5, 10)

      expect(mockApiClient.get).toHaveBeenCalledWith('seasons/5/leaderboard', { searchParams: { limit: 10 } })
      expect(result).toEqual(leaderboard)
    })
  })
})
