import type { GuildMemberEntry, GuildPublic } from '@/models/guild'
import { apiClient } from './api'

export const guildService = {
  async getAll() {
    return apiClient.get('guilds').json<GuildPublic[]>()
  },

  async create(data: { name: string, description: string, icon: string, color: string }) {
    return apiClient.post('guilds', { json: data }).json<GuildPublic>()
  },

  async update(id: number, data: { name: string, description: string, icon: string, color: string }) {
    return apiClient.put(`guilds/${id}`, { json: data }).json()
  },

  async join(id: number) {
    return apiClient.post(`guilds/${id}/join`).json()
  },

  async leave(id: number) {
    return apiClient.post(`guilds/${id}/leave`).json()
  },

  async remove(id: number) {
    await apiClient.delete(`guilds/${id}`)
  },

  async getMembers(id: number) {
    return apiClient.get(`guilds/${id}/members`).json<GuildMemberEntry[]>()
  },
}
