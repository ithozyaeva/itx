import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
      patch: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { notificationService } from '@/services/notifications'

describe('notificationService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('should call GET notifications', async () => {
      const notifications = [{ id: 1, title: 'Test' }]
      mockJson.mockResolvedValue(notifications)

      const result = await notificationService.getAll()

      expect(mockApiClient.get).toHaveBeenCalledWith('notifications')
      expect(result).toEqual(notifications)
    })
  })

  describe('getUnreadCount', () => {
    it('should call GET notifications/unread-count', async () => {
      const countData = { count: 5 }
      mockJson.mockResolvedValue(countData)

      const result = await notificationService.getUnreadCount()

      expect(mockApiClient.get).toHaveBeenCalledWith('notifications/unread-count')
      expect(result).toEqual(countData)
    })
  })

  describe('markAsRead', () => {
    it('should call PATCH notifications/:id/read', async () => {
      await notificationService.markAsRead(42)

      expect(mockApiClient.patch).toHaveBeenCalledWith('notifications/42/read')
    })
  })

  describe('markAllAsRead', () => {
    it('should call POST notifications/read-all', async () => {
      await notificationService.markAllAsRead()

      expect(mockApiClient.post).toHaveBeenCalledWith('notifications/read-all')
    })
  })
})
