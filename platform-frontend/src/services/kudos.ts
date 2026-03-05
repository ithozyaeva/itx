import type { KudosItem } from '@/models/kudos'
import { apiClient } from './api'

export const kudosService = {
  async getRecent(limit = 20, offset = 0) {
    return apiClient.get('kudos', { searchParams: { limit, offset } }).json<{ items: KudosItem[], total: number }>()
  },

  async send(toId: number, message: string) {
    return apiClient.post('kudos', { json: { toId, message } }).json<KudosItem>()
  },
}
