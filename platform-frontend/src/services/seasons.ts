import type { Season, SeasonWithLeaderboard } from '@/models/season'
import { apiClient } from './api'

export const seasonService = {
  async getAll() {
    return apiClient.get('seasons').json<Season[]>()
  },

  async getActive(limit = 20) {
    return apiClient.get('seasons/active', { searchParams: { limit } }).json<SeasonWithLeaderboard>()
  },

  async getLeaderboard(id: number, limit = 20) {
    return apiClient.get(`seasons/${id}/leaderboard`, { searchParams: { limit } }).json<SeasonWithLeaderboard>()
  },
}
