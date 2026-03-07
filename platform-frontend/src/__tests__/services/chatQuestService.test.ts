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

import { chatQuestService } from '@/services/chatQuestService'

describe('chatQuestService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getActiveQuests', () => {
    it('should call GET chat-quests/active', async () => {
      const quests = [{ id: 1, title: 'Daily quest' }]
      mockJson.mockResolvedValue(quests)

      const result = await chatQuestService.getActiveQuests()

      expect(mockApiClient.get).toHaveBeenCalledWith('chat-quests/active')
      expect(result).toEqual(quests)
    })
  })

  describe('getAllQuests', () => {
    it('should call GET chat-quests/all without filter', async () => {
      const quests = [{ id: 1 }, { id: 2 }]
      mockJson.mockResolvedValue(quests)

      const result = await chatQuestService.getAllQuests()

      expect(mockApiClient.get).toHaveBeenCalledWith('chat-quests/all', { searchParams: {} })
      expect(result).toEqual(quests)
    })

    it('should call GET chat-quests/all with filter', async () => {
      const quests = [{ id: 1, completed: true }]
      mockJson.mockResolvedValue(quests)

      const result = await chatQuestService.getAllQuests('completed')

      expect(mockApiClient.get).toHaveBeenCalledWith('chat-quests/all', { searchParams: { filter: 'completed' } })
      expect(result).toEqual(quests)
    })
  })
})
