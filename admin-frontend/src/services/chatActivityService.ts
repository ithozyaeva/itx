import api from '@/lib/api'

export interface ChatActivityStats {
  totalMessagesToday: number
  totalMessagesWeek: number
  uniqueUsersToday: number
  uniqueUsersWeek: number
  chatStats: ChatMessageCount[]
}

export interface ChatMessageCount {
  chatId: number
  title: string
  count: number
}

export interface DailyActivity {
  date: string
  count: number
}

export interface TopUser {
  telegramUserId: number
  telegramUsername: string
  telegramFirstName: string
  count: number
  topChat: string
}

export interface TrackedChat {
  id: number
  chatId: number
  title: string
  chatType: string
  isActive: boolean
}

export const chatActivityService = {
  getStats: async () => {
    return api.get('chat-activity/stats').json<ChatActivityStats>()
  },
  getChart: async (chatId?: number, days = 30) => {
    const searchParams: Record<string, string> = { days: String(days) }
    if (chatId)
      searchParams.chat_id = String(chatId)
    return api.get('chat-activity/chart', { searchParams }).json<DailyActivity[]>()
  },
  getTopUsers: async (days = 7, limit = 5) => {
    return api.get('chat-activity/top-users', {
      searchParams: { days: String(days), limit: String(limit) },
    }).json<TopUser[]>()
  },
  getChats: async () => {
    return api.get('chat-activity/chats').json<TrackedChat[]>()
  },
}
