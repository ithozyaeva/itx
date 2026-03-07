import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock api
const mockJson = vi.fn()
const mockBlob = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson, blob: mockBlob })),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { chatActivityService } = await import('@/services/chatActivityService')

describe('chatActivityService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getStats', () => {
    it('calls the correct endpoint', async () => {
      const mockStats = {
        totalMessagesToday: 42,
        totalMessagesWeek: 200,
        uniqueUsersToday: 10,
        uniqueUsersWeek: 30,
        totalMessagesLastWeek: 180,
        uniqueUsersLastWeek: 25,
        chatStats: [],
      }
      mockJson.mockResolvedValueOnce(mockStats)

      const result = await chatActivityService.getStats()

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/stats')
      expect(result).toEqual(mockStats)
    })
  })

  describe('getChart', () => {
    it('calls with default params when no arguments', async () => {
      mockJson.mockResolvedValueOnce([])

      await chatActivityService.getChart()

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/chart', {
        searchParams: { days: '30' },
      })
    })

    it('includes chatId when provided', async () => {
      mockJson.mockResolvedValueOnce([])

      await chatActivityService.getChart(123)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/chart', {
        searchParams: { days: '30', chat_id: '123' },
      })
    })

    it('includes userId when provided', async () => {
      mockJson.mockResolvedValueOnce([])

      await chatActivityService.getChart(undefined, 14, 456)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/chart', {
        searchParams: { days: '14', user_id: '456' },
      })
    })

    it('includes both chatId and userId when provided', async () => {
      mockJson.mockResolvedValueOnce([])

      await chatActivityService.getChart(10, 7, 20)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/chart', {
        searchParams: { days: '7', chat_id: '10', user_id: '20' },
      })
    })

    it('uses custom days parameter', async () => {
      mockJson.mockResolvedValueOnce([])

      await chatActivityService.getChart(undefined, 90)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/chart', {
        searchParams: { days: '90' },
      })
    })
  })

  describe('getTopUsers', () => {
    it('calls with default params', async () => {
      mockJson.mockResolvedValueOnce([])

      await chatActivityService.getTopUsers()

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/top-users', {
        searchParams: { days: '7', limit: '5' },
      })
    })

    it('passes custom days and limit', async () => {
      mockJson.mockResolvedValueOnce([])

      await chatActivityService.getTopUsers(30, 10)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/top-users', {
        searchParams: { days: '30', limit: '10' },
      })
    })
  })

  describe('getChats', () => {
    it('fetches tracked chats', async () => {
      const mockChats = [
        { id: 1, chatId: 100, title: 'Chat 1', chatType: 'supergroup', isActive: true },
      ]
      mockJson.mockResolvedValueOnce(mockChats)

      const result = await chatActivityService.getChats()

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/chats')
      expect(result).toEqual(mockChats)
    })
  })

  describe('getUserStats', () => {
    it('calls with userId and default days', async () => {
      const mockUserStats = {
        telegramUserId: 42,
        telegramUsername: 'user42',
        telegramFirstName: 'User',
        totalMessages: 100,
        activeChats: 3,
        avgPerDay: 3.3,
      }
      mockJson.mockResolvedValueOnce(mockUserStats)

      const result = await chatActivityService.getUserStats(42)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/user-stats', {
        searchParams: { user_id: '42', days: '30' },
      })
      expect(result).toEqual(mockUserStats)
    })

    it('passes custom days', async () => {
      mockJson.mockResolvedValueOnce({})

      await chatActivityService.getUserStats(42, 60)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/user-stats', {
        searchParams: { user_id: '42', days: '60' },
      })
    })
  })

  describe('exportCSV', () => {
    it('creates a download link with correct defaults', async () => {
      const mockBlobValue = new Blob(['csv,data'], { type: 'text/csv' })
      mockBlob.mockResolvedValueOnce(mockBlobValue)

      const mockCreateObjectURL = vi.fn(() => 'blob:http://test/csv')
      const mockRevokeObjectURL = vi.fn()
      global.URL.createObjectURL = mockCreateObjectURL
      global.URL.revokeObjectURL = mockRevokeObjectURL

      const mockClick = vi.fn()
      const mockElement = { href: '', download: '', click: mockClick } as any
      vi.spyOn(document, 'createElement').mockReturnValueOnce(mockElement)

      await chatActivityService.exportCSV()

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/export', {
        searchParams: { days: '30' },
      })
      expect(mockElement.download).toBe('chat-activity.csv')
      expect(mockClick).toHaveBeenCalled()
      expect(mockRevokeObjectURL).toHaveBeenCalledWith('blob:http://test/csv')
    })

    it('includes chatId in export request', async () => {
      const mockBlobValue = new Blob(['data'])
      mockBlob.mockResolvedValueOnce(mockBlobValue)

      global.URL.createObjectURL = vi.fn(() => 'blob:url')
      global.URL.revokeObjectURL = vi.fn()
      vi.spyOn(document, 'createElement').mockReturnValueOnce({ href: '', download: '', click: vi.fn() } as any)

      await chatActivityService.exportCSV(14, 55)

      expect(mockApi.get).toHaveBeenCalledWith('chat-activity/export', {
        searchParams: { days: '14', chat_id: '55' },
      })
    })
  })
})
