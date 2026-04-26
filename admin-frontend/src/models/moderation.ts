export type ModerationActionType
  = | 'ban'
    | 'unban'
    | 'mute'
    | 'unmute'
    | 'cleanup'
    | 'voteban_mute'
    | 'voteban_kick'
    | 'globalban'
    | 'globalunban'

export interface ModerationAction {
  id: number
  chatId: number
  targetUserId: number
  actorUserId: number
  action: ModerationActionType
  reason: string | null
  durationSeconds: number | null
  expiresAt: string | null
  meta: string
  createdAt: string
  /** Обогащение из join'ов */
  targetUsername: string
  targetFirstName: string
  chatTitle: string
}

export interface GlobalBan {
  userId: number
  bannedBy: number
  reason: string | null
  expiresAt: string | null
  createdAt: string
  updatedAt: string
}

export interface VotebanRow {
  id: number
  chatId: number
  chatTitle: string
  targetUserId: number
  targetUsername: string
  targetFirstName: string
  initiatorUserId: number
  triggerMessageId: number | null
  pollMessageId: number
  requiredVotes: number
  muteSeconds: number
  expiresAt: string
  status: 'open' | 'passed' | 'failed' | 'cancelled'
  finalizedAt: string | null
  createdAt: string
  votesFor: number
  votesAgainst: number
}
