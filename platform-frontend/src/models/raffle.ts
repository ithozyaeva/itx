export interface RaffleItem {
  id: number
  title: string
  description: string
  prize: string
  ticketCost: number
  maxTickets: number
  endsAt: string
  status: 'ACTIVE' | 'FINISHED'
  totalTickets: number
  myTickets: number
  winnerId: number | null
  winnerFirstName: string
  winnerLastName: string
  winnerUsername: string
  winnerAvatarUrl: string
}
