import api from '@/lib/api'

export interface ChatActivityStats {
  totalMessagesToday: number
  totalMessagesWeek: number
  uniqueUsersToday: number
  uniqueUsersWeek: number
  totalMessagesLastWeek: number
  uniqueUsersLastWeek: number
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

export interface UserStats {
  telegramUserId: number
  telegramUsername: string
  telegramFirstName: string
  totalMessages: number
  activeChats: number
  avgPerDay: number
}

export const chatActivityService = {
  getStats: async () => {
    return api.get('chat-activity/stats').json<ChatActivityStats>()
  },
  getChart: async (chatId?: number, days = 30, userId?: number) => {
    const searchParams: Record<string, string> = { days: String(days) }
    if (chatId)
      searchParams.chat_id = String(chatId)
    if (userId)
      searchParams.user_id = String(userId)
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
  getUserStats: async (userId: number, days = 30) => {
    return api.get('chat-activity/user-stats', {
      searchParams: { user_id: String(userId), days: String(days) },
    }).json<UserStats>()
  },
  exportCSV: async (days = 30, chatId?: number) => {
    const searchParams: Record<string, string> = { days: String(days) }
    if (chatId)
      searchParams.chat_id = String(chatId)
    const response = await api.get('chat-activity/export', { searchParams })
    const blob = await response.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'chat-activity.csv'
    a.click()
    URL.revokeObjectURL(url)
  },
}
