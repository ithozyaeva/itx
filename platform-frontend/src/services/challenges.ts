import type { ChallengesResponse } from '@/models/challenges'
import { apiClient } from './api'

export const challengesService = {
  async getMine() {
    return apiClient.get('challenges').json<ChallengesResponse>()
  },
}
