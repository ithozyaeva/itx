export type ChallengeKind = 'weekly' | 'monthly'

export interface ChallengeWithProgress {
  instanceId: number
  templateId: number
  code: string
  title: string
  description: string
  icon: string
  kind: ChallengeKind
  metricKey: string
  target: number
  rewardPoints: number
  achievementCode: string | null
  startsAt: string
  endsAt: string
  progress: number
  completedAt: string | null
  awarded: boolean
}

export interface ChallengesResponse {
  weekly: ChallengeWithProgress[]
  monthly: ChallengeWithProgress[]
}
