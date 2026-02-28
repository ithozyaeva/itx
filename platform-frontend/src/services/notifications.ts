import { apiClient } from './api'

export interface Notification {
  id: number
  memberId: number
  type: string
  title: string
  body: string
  read: boolean
  createdAt: string
}

export const notificationService = {
  getAll: async () => {
    return apiClient.get('notifications').json<Notification[]>()
  },
  getUnreadCount: async () => {
    return apiClient.get('notifications/unread-count').json<{ count: number }>()
  },
  markAsRead: async (id: number) => {
    return apiClient.patch(`notifications/${id}/read`)
  },
  markAllAsRead: async () => {
    return apiClient.post('notifications/read-all')
  },
}
