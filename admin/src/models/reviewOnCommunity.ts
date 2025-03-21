import type { Member } from './members'

export interface ReviewOnCommunity {
  id: number
  authorId: number
  author: Member
  text: string
  date: string
}

export interface ReviewOnCommunityFormData {
  id?: number
  authorId?: number
  authorTg?: string
  text: string
  date: string
}
