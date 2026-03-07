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

import { highlightsService } from '@/services/highlights'

describe('highlightsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getRecent', () => {
    it('should call GET highlights/recent with default limit', async () => {
      const highlights = [{ id: 1, text: 'Highlight 1' }]
      mockJson.mockResolvedValue(highlights)

      const result = await highlightsService.getRecent()

      expect(mockApiClient.get).toHaveBeenCalledWith('highlights/recent', { searchParams: { limit: 5 } })
      expect(result).toEqual(highlights)
    })

    it('should call GET highlights/recent with custom limit', async () => {
      const highlights = [{ id: 1, text: 'Highlight 1' }]
      mockJson.mockResolvedValue(highlights)

      const result = await highlightsService.getRecent(10)

      expect(mockApiClient.get).toHaveBeenCalledWith('highlights/recent', { searchParams: { limit: 10 } })
      expect(result).toEqual(highlights)
    })
  })
})
