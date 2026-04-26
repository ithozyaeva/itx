package models

import "time"

const (
	ModerationActionBan         = "ban"
	ModerationActionUnban       = "unban"
	ModerationActionMute        = "mute"
	ModerationActionUnmute      = "unmute"
	ModerationActionCleanup     = "cleanup"
	ModerationActionVotebanMute = "voteban_mute" // legacy (T1)
	ModerationActionVotebanKick = "voteban_kick"
	ModerationActionGlobalBan   = "globalban"
	ModerationActionGlobalUnban = "globalunban"
)

// ModerationActionsWithExpiry — действия, для которых имеет смысл слать
// алерт «срок истёк» в чат(ы). Глобальный бан обрабатывается отдельно
// (chat_id=0, нужно слать в каждый из затронутых чатов из meta).
var ModerationActionsWithExpiry = []string{
	ModerationActionBan,
	ModerationActionMute,
	ModerationActionVotebanMute,
	ModerationActionVotebanKick,
}

// ModerationAction — журнал модерационных действий бота.
type ModerationAction struct {
	Id                int64      `json:"id" gorm:"primaryKey"`
	ChatID            int64      `json:"chatId" gorm:"column:chat_id"`
	TargetUserID      int64      `json:"targetUserId" gorm:"column:target_user_id"`
	ActorUserID       int64      `json:"actorUserId" gorm:"column:actor_user_id"`
	Action            string     `json:"action" gorm:"column:action"`
	Reason            *string    `json:"reason" gorm:"column:reason"`
	DurationSeconds   *int       `json:"durationSeconds" gorm:"column:duration_seconds"`
	ExpiresAt         *time.Time `json:"expiresAt" gorm:"column:expires_at"`
	Meta              string     `json:"meta" gorm:"column:meta;type:jsonb;default:'{}'"`
	ExpiredNotifiedAt *time.Time `json:"expiredNotifiedAt" gorm:"column:expired_notified_at"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"column:created_at"`
}

func (ModerationAction) TableName() string {
	return "bot_moderation_actions"
}

const (
	VotebanStatusOpen      = "open"
	VotebanStatusPassed    = "passed"
	VotebanStatusFailed    = "failed"
	VotebanStatusCancelled = "cancelled"
)

const (
	VotebanVoteFor     int16 = 1
	VotebanVoteAgainst int16 = -1
)

// Voteban — голосование за временный мут пользователя в групповом чате.
type Voteban struct {
	Id               int64      `json:"id" gorm:"primaryKey"`
	ChatID           int64      `json:"chatId" gorm:"column:chat_id"`
	ChatTitle        string     `json:"chatTitle" gorm:"column:chat_title"`
	TargetUserID     int64      `json:"targetUserId" gorm:"column:target_user_id"`
	TargetUsername   string     `json:"targetUsername" gorm:"column:target_username"`
	TargetFirstName  string     `json:"targetFirstName" gorm:"column:target_first_name"`
	InitiatorUserID  int64      `json:"initiatorUserId" gorm:"column:initiator_user_id"`
	TriggerMessageID *int       `json:"triggerMessageId" gorm:"column:trigger_message_id"`
	PollMessageID    int        `json:"pollMessageId" gorm:"column:poll_message_id"`
	RequiredVotes    int        `json:"requiredVotes" gorm:"column:required_votes"`
	MuteSeconds      int        `json:"muteSeconds" gorm:"column:mute_seconds"`
	ExpiresAt        time.Time  `json:"expiresAt" gorm:"column:expires_at"`
	Status           string     `json:"status" gorm:"column:status"`
	FinalizedAt      *time.Time `json:"finalizedAt" gorm:"column:finalized_at"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"column:created_at"`
}

func (Voteban) TableName() string {
	return "bot_votebans"
}

// VotebanVote — голос конкретного юзера в рамках одного voteban.
type VotebanVote struct {
	VotebanID   int64     `json:"votebanId" gorm:"column:voteban_id;primaryKey"`
	VoterUserID int64     `json:"voterUserId" gorm:"column:voter_user_id;primaryKey"`
	Vote        int16     `json:"vote" gorm:"column:vote"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (VotebanVote) TableName() string {
	return "bot_voteban_votes"
}

// VotebanTally — агрегированные счётчики голосов.
type VotebanTally struct {
	For     int `json:"for"`
	Against int `json:"against"`
}

// GlobalBan — запись о глобальной блокировке пользователя. ExpiresAt=nil
// означает permanent. Снимается через /globalunban.
type GlobalBan struct {
	UserID    int64      `json:"userId" gorm:"column:user_id;primaryKey"`
	BannedBy  int64      `json:"bannedBy" gorm:"column:banned_by"`
	Reason    *string    `json:"reason" gorm:"column:reason"`
	ExpiresAt *time.Time `json:"expiresAt" gorm:"column:expires_at"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at"`
}

func (GlobalBan) TableName() string {
	return "bot_global_bans"
}

// IsActive возвращает true, если бан ещё действует.
func (g GlobalBan) IsActive(now time.Time) bool {
	if g.ExpiresAt == nil {
		return true
	}
	return g.ExpiresAt.After(now)
}
