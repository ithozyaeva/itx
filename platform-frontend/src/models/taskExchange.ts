import type { TelegramUser } from '@/models/profile'

export type TaskExchangeStatus = 'OPEN' | 'IN_PROGRESS' | 'DONE' | 'APPROVED'

export interface TaskExchange {
  id: number
  title: string
  description: string
  creatorId: number
  creator: TelegramUser
  maxAssignees: number
  assignees: TelegramUser[]
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
  maxAssignees: number
}

export interface UpdateTaskExchangeRequest {
  title?: string
  description?: string
  maxAssignees?: number
}
