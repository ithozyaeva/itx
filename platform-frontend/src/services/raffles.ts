import type { RaffleItem } from '@/models/raffle'
import { apiClient } from './api'

export const raffleService = {
  async getAll() {
    return apiClient.get('raffles').json<RaffleItem[]>()
  },

  async buyTickets(id: number, count = 1) {
    return apiClient.post(`raffles/${id}/buy`, { json: { count } }).json()
  },
}
