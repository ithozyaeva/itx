import { apiClient } from './api'

export interface ChatQuestWithProgress {
  id: number
  title: string
  description: string
  questType: string
  chatId: number | null
  targetCount: number
  pointsReward: number
  startsAt: string
  endsAt: string
  isActive: boolean
  currentCount: number
  completed: boolean
}

export const chatQuestService = {
  async getActiveQuests() {
    return apiClient.get('chat-quests/active').json<ChatQuestWithProgress[]>()
  },

  async getAllQuests(filter?: string) {
    return apiClient.get('chat-quests/all', { searchParams: filter ? { filter } : {} }).json<ChatQuestWithProgress[]>()
  },
}
