import type { NotificationSettings } from '@/models/profile'
import { apiClient } from './api'

export const notificationSettingsService = {
  async get(): Promise<NotificationSettings> {
    const response = await apiClient.get('notification-settings')
    return response.json<NotificationSettings>()
  },

  async update(settings: Partial<NotificationSettings>): Promise<NotificationSettings> {
    const response = await apiClient.patch('notification-settings', { json: settings })
    return response.json<NotificationSettings>()
  },
}
