import { apiClient } from './api'

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

export const referrerService = {
  async getMine(): Promise<ReferrerInfo | null> {
    const resp = await apiClient.get('members/me/referrer').json<{ referrer: ReferrerInfo | null }>()
    return resp.referrer
  },

  async markSeen(): Promise<void> {
    await apiClient.post('members/me/referrer/seen')
  },
}
