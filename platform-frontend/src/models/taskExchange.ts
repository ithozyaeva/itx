import type { TelegramUser } from '@/models/profile'

export type TaskExchangeStatus = 'OPEN' | 'IN_PROGRESS' | 'DONE' | 'APPROVED'

export interface TaskExchange {
  id: number
  title: string
  description: string
  creatorId: number
  creator: TelegramUser
  assigneeId: number | null
  assignee: TelegramUser | null
  status: TaskExchangeStatus
  createdAt: string
  updatedAt: string
}

export interface TaskExchangeSearchResponse {
  items: TaskExchange[]
  total: number
}

export interface CreateTaskExchangeRequest {
  title: string
  description: string
}
