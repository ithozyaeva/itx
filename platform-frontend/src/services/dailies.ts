import type { CheckInResponse, DailyTodayResponse, StreakResponse } from '@/models/dailies'
import type { RaffleItem } from '@/models/raffle'
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

  async getDailyRaffle() {
    return apiClient.get('raffles/daily/today').json<RaffleItem | { raffle: null }>()
  },
}
