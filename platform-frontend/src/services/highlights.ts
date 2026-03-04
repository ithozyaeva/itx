import type { ChatHighlight } from '@/models/highlight'
import { apiClient } from './api'

export const highlightsService = {
  async getRecent(limit = 5): Promise<ChatHighlight[]> {
    const response = await apiClient.get('highlights/recent', {
      searchParams: { limit },
    })
    return response.json<ChatHighlight[]>()
  },
}
