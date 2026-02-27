import api from '@/lib/api'

export interface DashboardStats {
  totalMembers: number
  totalMentors: number
  upcomingEvents: number
  pastEvents: number
  pendingReviews: number
  approvedReviews: number
  referralLinks: number
  resumes: number
}

export const statsService = {
  getStats: async () => {
    return api.get('stats').json<DashboardStats>()
  },
}
