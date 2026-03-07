import type { CasinoBetResult, CasinoStats } from '@/models/casino'
import { apiClient } from './api'

export const casinoService = {
  async coinFlip(betAmount: number, choice: 'heads' | 'tails') {
    return apiClient.post('casino/coin-flip', { json: { betAmount, choice } }).json<CasinoBetResult>()
  },

  async diceRoll(betAmount: number, target: number, direction: 'over' | 'under') {
    return apiClient.post('casino/dice-roll', { json: { betAmount, target, direction } }).json<CasinoBetResult>()
  },

  async wheelSpin(betAmount: number) {
    return apiClient.post('casino/wheel', { json: { betAmount } }).json<CasinoBetResult>()
  },

  async getHistory(limit = 20) {
    const res = await apiClient.get(`casino/history?limit=${limit}`).json<{ items: CasinoBetResult[], total: number }>()
    return res.items ?? []
  },

  async getStats() {
    return apiClient.get('casino/stats').json<CasinoStats>()
  },
}
