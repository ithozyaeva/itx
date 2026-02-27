import type { TelegramUser } from './profile'

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
  videoLink: string
  isRepeating: boolean
  repeatPeriod?: string
  repeatInterval?: number
  repeatEndDate?: string
  hosts: TelegramUser[]
  members: TelegramUser[]
}

export type PlaceType = 'ONLINE' | 'OFFLINE' | 'HYBRID'
