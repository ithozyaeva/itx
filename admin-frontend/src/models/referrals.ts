import type { Member } from './members'

export interface ProfTag {
  id: number
  title: string
}

export interface ReferralLink {
  id: number
  authorId: number
  author: Member
  company: string
  grade: string
  profTags: ProfTag[]
  status: 'active' | 'freezed'
  vacationsCount: number
  expiresAt?: string
  conversionsCount: number
  updatedAt: string
}
