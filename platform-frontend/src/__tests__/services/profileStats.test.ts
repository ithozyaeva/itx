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

import { profileStatsService } from '@/services/profileStats'

describe('profileStatsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getMyStats', () => {
    it('should call GET profile-stats/me', async () => {
      const stats = { points: 100, eventsAttended: 5 }
      mockJson.mockResolvedValue(stats)

      const result = await profileStatsService.getMyStats()

      expect(mockApiClient.get).toHaveBeenCalledWith('profile-stats/me')
      expect(result).toEqual(stats)
    })
  })

  describe('getMemberStats', () => {
    it('should call GET profile-stats/:id', async () => {
      const stats = { points: 200, eventsAttended: 10 }
      mockJson.mockResolvedValue(stats)

      const result = await profileStatsService.getMemberStats(42)

      expect(mockApiClient.get).toHaveBeenCalledWith('profile-stats/42')
      expect(result).toEqual(stats)
    })
  })
})
