import type { Member } from './members'

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
  maxParticipants: number
  hosts: Member[]
  eventTags: EventTag[]
}

export interface EventTag {
  id: number
  name: string
}

export type PlaceType = 'ONLINE' | 'OFFLINE' | 'HYBRID'
