export interface StreakMilestone {
  days: number
  reward: number
  reached: boolean
}

export interface StreakResponse {
  current: number
  longest: number
  freezesAvailable: number
  lastCheckIn: string | null
  nextThreshold: number | null
  daysToNext: number | null
  milestones: StreakMilestone[]
}

export interface CheckInResponse {
  checkInDone: boolean
  alreadyToday: boolean
  streak: StreakResponse
  raffleEntered?: boolean
}

export type DailyTaskTier = 'engagement' | 'light' | 'meaningful' | 'big'

export interface DailyTaskWithProgress {
  id: number
  code: string
  title: string
  description: string
  icon: string
  tier: DailyTaskTier
  points: number
  target: number
  triggerKey: string
  active: boolean
  createdAt: string
  progress: number
  completedAt: string | null
  awarded: boolean
}

export interface DailyCheckInState {
  done: boolean
  at: string | null
}

export interface DailyAllBonusState {
  points: number
  awarded: boolean
}

export interface DailyTodayResponse {
  checkIn: DailyCheckInState
  tasks: DailyTaskWithProgress[]
  allBonus: DailyAllBonusState
}
