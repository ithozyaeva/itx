export interface PointsMonth {
  month: string
  total: number
}

export interface PointsSource {
  reason: string
  total: number
}

export interface ActivityDay {
  date: string
  count: number
}

export interface ProfileStats {
  eventsAttended: number
  eventsHosted: number
  reviewsCount: number
  referralsCount: number
  kudosSent: number
  kudosReceived: number
  tasksCreated: number
  tasksDone: number
  pointsBalance: number
  memberSince: string
  pointsHistory: PointsMonth[]
  pointsBySource: PointsSource[]
  achievementsEarned: number
  achievementsTotal: number
  activityHistory: ActivityDay[]
}
