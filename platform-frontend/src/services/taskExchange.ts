import type { CreateTaskExchangeRequest, TaskExchange, TaskExchangeSearchResponse } from '@/models/taskExchange'
import { apiClient } from './api'

export const taskExchangeService = {
  async getAll(params?: { status?: string, limit?: number, offset?: number }) {
    const searchParams = new URLSearchParams()
    if (params?.status)
      searchParams.set('status', params.status)
    if (params?.limit)
      searchParams.set('limit', params.limit.toString())
    if (params?.offset)
      searchParams.set('offset', params.offset.toString())

    return apiClient.get('tasks', { searchParams }).json<TaskExchangeSearchResponse>()
  },

  async getById(id: number) {
    return apiClient.get(`tasks/${id}`).json<TaskExchange>()
  },

  async create(data: CreateTaskExchangeRequest) {
    return apiClient.post('tasks', { json: data }).json<TaskExchange>()
  },

  async assign(id: number) {
    return apiClient.post(`tasks/${id}/assign`).json<TaskExchange>()
  },

  async unassign(id: number) {
    return apiClient.post(`tasks/${id}/unassign`).json<TaskExchange>()
  },

  async markDone(id: number) {
    return apiClient.post(`tasks/${id}/done`).json<TaskExchange>()
  },

  async approve(id: number) {
    return apiClient.post(`tasks/${id}/approve`).json<TaskExchange>()
  },

  async reject(id: number) {
    return apiClient.post(`tasks/${id}/reject`).json<TaskExchange>()
  },

  async remove(id: number) {
    return apiClient.delete(`tasks/${id}`)
  },
}
