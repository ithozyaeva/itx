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

export interface MonthlyStats {
  month: string
  count: number
}

export interface ChartStats {
  memberGrowth: MonthlyStats[]
  eventAttendance: MonthlyStats[]
}

export const statsService = {
  getStats: async () => {
    return api.get('stats').json<DashboardStats>()
  },
  getChartStats: async () => {
    return api.get('stats/charts').json<ChartStats>()
  },
}
