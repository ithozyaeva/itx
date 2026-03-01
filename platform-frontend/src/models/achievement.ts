export interface Achievement {
  id: string
  title: string
  description: string
  icon: string
  category: AchievementCategory
  threshold: number
}

export type AchievementCategory = 'events' | 'points' | 'social' | 'activity'

export interface UserAchievement extends Achievement {
  unlocked: boolean
  progress: number
}

export interface AchievementsResponse {
  items: UserAchievement[]
  totalCount: number
  unlockedCount: number
}
