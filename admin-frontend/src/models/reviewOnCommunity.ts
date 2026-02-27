import type { Member } from './members'

export interface ReviewOnCommunity {
  id: number
  authorId: number
  author: Member
  text: string
  date: string
  status: ReviewStatus
}

export interface ReviewOnCommunityFormData {
  id?: number
  authorId?: number
  authorTg?: string
  text: string
  date: string
}

export type ReviewStatus = 'DRAFT' | 'APPROVED'
