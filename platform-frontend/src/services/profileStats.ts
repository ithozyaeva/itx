import type { ProfileStats } from '@/models/profileStats'
import { apiClient } from './api'

export const profileStatsService = {
  async getMyStats() {
    return apiClient.get('profile-stats/me').json<ProfileStats>()
  },

  async getMemberStats(id: number) {
    return apiClient.get(`profile-stats/${id}`).json<ProfileStats>()
  },
}
