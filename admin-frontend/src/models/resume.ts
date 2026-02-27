import type { Member } from './members'

export type WorkFormat = 'REMOTE' | 'HYBRID' | 'OFFICE' | ''

export interface Resume {
  id: number
  tgId: number
  fileName: string
  filePath: string
  workExperience?: string
  desiredPosition?: string
  workFormat?: WorkFormat
  createdAt: string
  updatedAt: string
  member?: Member
}

export interface ResumeFilter {
  workFormat?: WorkFormat
  desiredPosition?: string
  workExperience?: string
}
