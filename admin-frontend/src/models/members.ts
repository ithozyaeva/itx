export interface Member {
  id: number
  firstName: string
  lastName: string
  tg: string
  telegramID: number
  isMentor?: boolean
  roles: MemberRole[]
  birthday: string
}

export interface CreateMemberDto {
  firstName: string
  lastName: string
  tg: string
  roles: MemberRole[]
}

export interface UpdateMemberDto {
  id: number
  firstName: string
  lastName: string
  tg: string
  roles: MemberRole[]
}

export type MemberRole = 'UNSUBSCRIBER' | 'SUBSCRIBER' | 'MENTOR' | 'ADMIN' | 'EVENT_MAKER'
