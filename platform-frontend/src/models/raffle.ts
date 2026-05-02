export type RaffleKind = 'manual' | 'daily'
export type RaffleEntryRule = 'purchase' | 'auto_check_in'

export interface RaffleItem {
  id: number
  title: string
  description: string
  prize: string
  ticketCost: number
  maxTickets: number
  endsAt: string
  status: 'ACTIVE' | 'FINISHED'
  kind?: RaffleKind
  entryRule?: RaffleEntryRule
  dayKey?: string | null
  totalTickets: number
  myTickets: number
  winnerId: number | null
  winnerFirstName: string
  winnerLastName: string
  winnerUsername: string
  winnerAvatarUrl: string
}
