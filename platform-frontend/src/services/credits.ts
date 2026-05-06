import { apiClient } from './api'

export interface CreditTransaction {
  id: number
  memberId: number
  amount: number
  reason: string
  sourceType: string
  sourceId: number
  description: string
  createdAt: string
}

export interface CreditsSummary {
  balance: number
  transactions: CreditTransaction[]
}

export interface PurchaseResult {
  tier_id: number
  tier_slug: string
  tier_name: string
  tier_level: number
  price_credits: number
  balance_left: number
  expires_at: string
}

export const creditsService = {
  async getMine() {
    return apiClient.get('credits/me').json<CreditsSummary>()
  },

  async purchaseTier(slug: string) {
    return apiClient.post('subscriptions/purchase', {
      json: { tier_slug: slug },
    }).json<PurchaseResult>()
  },
}
