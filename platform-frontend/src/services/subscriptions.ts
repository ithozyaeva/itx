import { apiClient } from './api'

export interface PublicTier {
  id: number
  slug: string
  name: string
  level: number
  price: number
  boosty_url: string
  description: string
  features: string[]
}

export const subscriptionsService = {
  async getPublicTiers(): Promise<PublicTier[]> {
    const response = await apiClient.get('subscriptions/tiers')
    const data = await response.json<{ items: PublicTier[] }>()
    return data.items
  },
}
