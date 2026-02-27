import type { CommunityEvent } from '@/models/event'
import { apiClient } from './api'

export interface EventSearchFilters {
  title?: string
  placeType?: string
}

function cleanFilters(filters?: EventSearchFilters): Record<string, string> {
  if (!filters)
    return {}
  return Object.fromEntries(
    Object.entries(filters).filter(([_, v]) => v !== undefined && v !== ''),
  )
}

export const eventsService = {
  searchOld: async (limit: number, offset: number, filters?: EventSearchFilters) => {
    return apiClient.get('events', { searchParams: { limit, offset, dateTo: new Date().toISOString(), ...cleanFilters(filters) } }).json<{ items: CommunityEvent[], total: number }>()
  },
  searchNext: async (limit: number, offset: number, filters?: EventSearchFilters) => {
    return apiClient.get('events', { searchParams: { limit, offset, dateFrom: new Date().toISOString(), ...cleanFilters(filters) } }).json<{ items: CommunityEvent[], total: number }>()
  },
  applyEvent: async (eventId: number) => {
    return apiClient.post('events/apply', { json: { eventId } }).json<CommunityEvent>()
  },
  declineEvent: async (eventId: number) => {
    return apiClient.post('events/decline', { json: { eventId } }).json<CommunityEvent>()
  },
}
