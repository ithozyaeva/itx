import type { AchievementsResponse } from '@/models/achievement'
import { apiClient } from './api'

export const achievementsService = {
  async getMyAchievements() {
    return apiClient.get('achievements/me').json<AchievementsResponse>()
  },

  async getByMemberId(id: number) {
    return apiClient.get(`achievements/member/${id}`).json<AchievementsResponse>()
  },
}
