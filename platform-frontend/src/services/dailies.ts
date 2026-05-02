import type { CheckInResponse, DailyTodayResponse, StreakResponse } from '@/models/dailies'
import { apiClient } from './api'

export const dailiesService = {
  async checkIn() {
    return apiClient.post('dailies/check-in').json<CheckInResponse>()
  },

  async getToday() {
    return apiClient.get('dailies/today').json<DailyTodayResponse>()
  },

  async getStreak() {
    return apiClient.get('streak/me').json<StreakResponse>()
  },
}
