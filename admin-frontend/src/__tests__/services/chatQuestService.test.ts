import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock api
const mockJson = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson })),
  post: vi.fn(() => ({ json: mockJson })),
  put: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { chatQuestService } = await import('@/services/chatQuestService')

describe('chatQuestService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('fetches quests with default pagination', async () => {
      const mockResponse = {
        items: [{ id: 1, title: 'Quest 1' }],
        total: 1,
      }
      mockJson.mockResolvedValueOnce(mockResponse)

      const result = await chatQuestService.getAll()

      expect(mockApi.get).toHaveBeenCalledWith('chat-quests/', {
        searchParams: { limit: '20', offset: '0' },
      })
      expect(result).toEqual(mockResponse)
    })

    it('passes custom limit and offset', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await chatQuestService.getAll(50, 10)

      expect(mockApi.get).toHaveBeenCalledWith('chat-quests/', {
        searchParams: { limit: '50', offset: '10' },
      })
    })
  })

  describe('create', () => {
    it('creates a quest with correct data', async () => {
      const questData = {
        title: 'New Quest',
        description: 'Description',
        questType: 'message_count',
        chatId: 123,
        targetCount: 50,
        pointsReward: 100,
        startsAt: '2026-03-01T00:00:00Z',
        endsAt: '2026-04-01T00:00:00Z',
        isActive: true,
      }
      const createdQuest = { id: 1, ...questData, createdAt: '2026-03-01T00:00:00Z' }
      mockJson.mockResolvedValueOnce(createdQuest)

      const result = await chatQuestService.create(questData)

      expect(mockApi.post).toHaveBeenCalledWith('chat-quests/', { json: questData })
      expect(result).toEqual(createdQuest)
    })

    it('handles null chatId', async () => {
      const questData = {
        title: 'Global Quest',
        description: 'For all chats',
        questType: 'message_count',
        chatId: null,
        targetCount: 100,
        pointsReward: 200,
        startsAt: '2026-03-01T00:00:00Z',
        endsAt: '2026-04-01T00:00:00Z',
        isActive: true,
      }
      mockJson.mockResolvedValueOnce({ id: 2, ...questData })

      await chatQuestService.create(questData)

      expect(mockApi.post).toHaveBeenCalledWith('chat-quests/', { json: questData })
    })
  })

  describe('update', () => {
    it('updates a quest with partial data', async () => {
      const updateData = { title: 'Updated Title', isActive: false }
      const updatedQuest = { id: 1, ...updateData }
      mockJson.mockResolvedValueOnce(updatedQuest)

      const result = await chatQuestService.update(1, updateData)

      expect(mockApi.put).toHaveBeenCalledWith('chat-quests/1', { json: updateData })
      expect(result).toEqual(updatedQuest)
    })
  })

  describe('remove', () => {
    it('deletes a quest by id', async () => {
      mockApi.delete.mockResolvedValueOnce(undefined)

      await chatQuestService.remove(5)

      expect(mockApi.delete).toHaveBeenCalledWith('chat-quests/5')
    })
  })
})
