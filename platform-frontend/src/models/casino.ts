export interface CasinoBetResult {
  id: number
  game: string
  betAmount: number
  betChoice: string
  result: string
  multiplier: number
  payout: number
  profit: number
  balance: number
  createdAt: string
}

export interface CasinoFeedItem {
  id: number
  memberFirstName: string
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

export interface CasinoStats {
  balance: number
  totalBets: number
  totalWagered: number
  totalWon: number
  totalProfit: number
}
