export interface GuildPublic {
  id: number
  name: string
  description: string
  icon: string
  color: string
  ownerId: number
  ownerFirstName: string
  ownerLastName: string
  ownerUsername: string
  ownerAvatarUrl: string
  memberCount: number
  totalPoints: number
  isMember: boolean
}

export interface GuildMemberEntry {
  memberId: number
  firstName: string
  lastName: string
  tg: string
  avatarUrl: string
  total: number
}
