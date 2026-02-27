import type { TelegramUser } from './auth'
import axios from 'axios'

export interface CommunityEvent {
  id: number
  title: string
  description: string
  date: string
  timezone: string
  placeType: PlaceType
  place: string
  customPlaceType: string
  eventType: string
  open: boolean
  isRepeating: boolean
  repeatPeriod?: string
  repeatInterval?: number
  repeatEndDate?: string
  hosts: TelegramUser[]
  eventTags: { id: number, name: string }[]
  videoLink: string
}

export type PlaceType = 'ONLINE' | 'OFFLINE' | 'HYBRID'
export const eventService = {
  getOld: async () => {
    const response: { items: CommunityEvent[], total: number } = (await axios.get('/api/events/old')).data
    return response.items
  },
  getNext: async () => {
    const response: { items: CommunityEvent[], total: number } = (await axios.get('/api/events/next')).data
    return response.items
  },
  getICS: (eventId: number) => {
    try {
      window.location.href = `/api/events/ics?eventId=${eventId}`
    }
    catch (error) {
      console.error('Failed to download ICS file:', error)
    }
  },
}
