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

import { achievementsService } from '@/services/achievements'

describe('achievementsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getMyAchievements', () => {
    it('should call GET achievements/me', async () => {
      const achievements = { items: [{ id: 1, name: 'First Steps' }] }
      mockJson.mockResolvedValue(achievements)

      const result = await achievementsService.getMyAchievements()

      expect(mockApiClient.get).toHaveBeenCalledWith('achievements/me')
      expect(result).toEqual(achievements)
    })
  })

  describe('getByMemberId', () => {
    it('should call GET achievements/member/:id', async () => {
      const achievements = { items: [{ id: 2, name: 'Veteran' }] }
      mockJson.mockResolvedValue(achievements)

      const result = await achievementsService.getByMemberId(7)

      expect(mockApiClient.get).toHaveBeenCalledWith('achievements/member/7')
      expect(result).toEqual(achievements)
    })
  })
})
