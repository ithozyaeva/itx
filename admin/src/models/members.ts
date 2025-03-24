export interface Member {
  id: number
  name: string
  tg: string
  isMentor: boolean
}

export interface CreateMemberDto {
  name: string
  tg: string
}

export interface UpdateMemberDto {
  id: number
  name: string
  tg: string
}
