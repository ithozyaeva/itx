import api from '@/lib/api'

export interface ChatQuest {
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
  createdAt: string
}

export interface ChatQuestCreateRequest {
  title: string
  description: string
  questType: string
  chatId: number | null
  targetCount: number
  pointsReward: number
  startsAt: string
  endsAt: string
  isActive: boolean
}

export interface ChatQuestsResponse {
  items: ChatQuest[]
  total: number
}

export const chatQuestService = {
  async getAll(limit = 20, offset = 0) {
    return api.get('chat-quests/', {
      searchParams: { limit: String(limit), offset: String(offset) },
    }).json<ChatQuestsResponse>()
  },
  async create(quest: ChatQuestCreateRequest) {
    return api.post('chat-quests/', { json: quest }).json<ChatQuest>()
  },
  async update(id: number, quest: Partial<ChatQuestCreateRequest>) {
    return api.put(`chat-quests/${id}`, { json: quest }).json<ChatQuest>()
  },
  async remove(id: number) {
    await api.delete(`chat-quests/${id}`)
  },
}
