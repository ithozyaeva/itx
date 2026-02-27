import type { ProfTag, TelegramUser } from './profile'

export interface ReferalLink {
  id: number
  author: TelegramUser
  company: string
  grade: Grade
  profTags: ProfTag[]
  status: ReferalLinkStatus
  vacationsCount: number
  updatedAt: string
}

export type Grade = 'senior' | 'junior' | 'middle'

type ReferalLinkStatus = 'active' | 'freezed'
