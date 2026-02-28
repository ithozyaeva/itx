export interface PointTransaction {
  id: number
  memberId: number
  amount: number
  reason: string
  sourceType: string
  sourceId: number
  description: string
  createdAt: string
}

export interface PointsSummary {
  balance: number
  transactions: PointTransaction[]
}

export interface LeaderboardEntry {
  memberId: number
  firstName: string
  lastName: string
  tg: string
  avatarUrl: string
  total: number
}
