export interface PointsMonth {
  month: string
  total: number
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
}
