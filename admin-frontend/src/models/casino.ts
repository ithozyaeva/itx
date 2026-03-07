export interface CasinoBet {
  id: number
  memberId: number
  memberFirstName: string
  memberLastName: string
  memberUsername: string
  game: string
  betAmount: number
  betChoice: string
  result: string
  multiplier: number
  payout: number
  profit: number
  createdAt: string
}

export interface GameStats {
  game: string
  totalBets: number
  totalWagered: number
  totalPayout: number
  houseProfit: number
}

export interface CasinoAdminStats {
  totalBets: number
  totalWagered: number
  totalPayout: number
  houseProfit: number
  uniquePlayers: number
  gameStats: GameStats[]
}
