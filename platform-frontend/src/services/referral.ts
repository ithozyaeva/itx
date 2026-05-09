import { apiClient } from './api'

export interface InviteeRow {
  id: number
  firstName: string
  lastName: string
  tg: string
  avatarUrl: string
  hasActiveSub: boolean
  joinedAt: string
}

export interface ReferralCabinet {
  code: string
  deeplink: string
  balance: number
  invitedTotal: number
  withActiveSub: number
  totalEarned: number
  recentInvitees: InviteeRow[]
}

export interface ReferrerAuthor {
  id: number
  firstName: string
  lastName: string
  tg: string
  avatarUrl: string
}

export interface ReferrerInfo {
  author: ReferrerAuthor
  seenAt: string | null
}

export const referralService = {
  async getCabinet(): Promise<ReferralCabinet> {
    return apiClient.get('members/me/referral').json<ReferralCabinet>()
  },

  async getReferrer(): Promise<ReferrerInfo | null> {
    const resp = await apiClient.get('members/me/referrer').json<{ referrer: ReferrerInfo | null }>()
    return resp.referrer
  },

  async markReferrerSeen(): Promise<void> {
    await apiClient.post('members/me/referrer/seen')
  },
}
