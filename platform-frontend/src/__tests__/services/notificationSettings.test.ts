import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      patch: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { notificationSettingsService } from '@/services/notificationSettings'

describe('notificationSettingsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('get', () => {
    it('should call GET notification-settings', async () => {
      const settings = { emailEnabled: true, pushEnabled: false }
      mockJson.mockResolvedValue(settings)
      mockApiClient.get.mockReturnValue({ json: mockJson })

      const result = await notificationSettingsService.get()

      expect(mockApiClient.get).toHaveBeenCalledWith('notification-settings')
      expect(result).toEqual(settings)
    })
  })

  describe('update', () => {
    it('should call PATCH notification-settings with settings', async () => {
      const settings = { emailEnabled: false }
      const responseData = { emailEnabled: false, pushEnabled: true }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.patch.mockReturnValue({ json: mockJson })

      const result = await notificationSettingsService.update(settings)

      expect(mockApiClient.patch).toHaveBeenCalledWith('notification-settings', { json: settings })
      expect(result).toEqual(responseData)
    })
  })
})
