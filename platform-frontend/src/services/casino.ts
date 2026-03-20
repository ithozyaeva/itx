import type { CasinoBetResult, CasinoFeedItem, CasinoStats } from '@/models/casino'
import { apiClient } from './api'

export const casinoService = {
  async coinFlip(betAmount: number, choice: 'heads' | 'tails') {
    return apiClient.post('minigames/coin-flip', { json: { betAmount, choice } }).json<CasinoBetResult>()
  },

  async diceRoll(betAmount: number, target: number, direction: 'over' | 'under') {
    return apiClient.post('minigames/dice-roll', { json: { betAmount, target, direction } }).json<CasinoBetResult>()
  },

  async wheelSpin(betAmount: number) {
    return apiClient.post('minigames/wheel', { json: { betAmount } }).json<CasinoBetResult>()
  },

  async getHistory(limit = 20) {
    const res = await apiClient.get(`minigames/history?limit=${limit}`).json<{ items: CasinoBetResult[], total: number }>()
    return res.items ?? []
  },

  async getFeed(limit = 20) {
    const res = await apiClient.get(`minigames/feed?limit=${limit}`).json<{ items: CasinoFeedItem[] }>()
    return res.items ?? []
  },

  async getStats() {
    return apiClient.get('minigames/stats').json<CasinoStats>()
  },
}
