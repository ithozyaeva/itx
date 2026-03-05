export interface Season {
  id: number
  title: string
  startDate: string
  endDate: string
  status: 'ACTIVE' | 'FINISHED'
  createdAt: string
}

export interface SeasonLeaderboardEntry {
  memberId: number
  firstName: string
  lastName: string
  tg: string
  avatarUrl: string
  total: number
  rank: number
}

export interface SeasonWithLeaderboard {
  season: Season
  leaderboard: SeasonLeaderboardEntry[]
}
