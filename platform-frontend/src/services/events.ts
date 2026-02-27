import type { CommunityEvent } from '@/models/event'
import { apiClient } from './api'

export const eventsService = {
  searchOld: async (limit: number, offset: number) => {
    return apiClient.get('events', { searchParams: { limit, offset, dateTo: new Date().toISOString() } }).json<{ items: CommunityEvent[], total: number }>()
  },
  searchNext: async (limit: number, offset: number) => {
    return apiClient.get('events', { searchParams: { limit, offset, dateFrom: new Date().toISOString() } }).json<{ items: CommunityEvent[], total: number }>()
  },
  applyEvent: async (eventId) => {
    return apiClient.post('events/apply', { json: { eventId } }).json<CommunityEvent>()
  },
  declineEvent: async (eventId) => {
    return apiClient.post('events/decline', { json: { eventId } }).json<CommunityEvent>()
  },
}
