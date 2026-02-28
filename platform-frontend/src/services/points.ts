import type { LeaderboardEntry, PointsSummary } from '@/models/points'
import { apiClient } from './api'

export const pointsService = {
  async getMyPoints() {
    return apiClient.get('points/me').json<PointsSummary>()
  },

  async getLeaderboard(limit = 20) {
    return apiClient.get('points/leaderboard', { searchParams: { limit } }).json<{ items: LeaderboardEntry[] }>()
  },
}
