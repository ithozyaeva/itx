package models

import "time"

type SubscriptionTier struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Slug  string `json:"slug" gorm:"uniqueIndex;size:50"`
	Name  string `json:"name" gorm:"size:100"`
	Level int    `json:"level" gorm:"uniqueIndex"`
}

func (SubscriptionTier) TableName() string { return "subscription_tiers" }

type SubscriptionChat struct {
	ID              int64  `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Title           string `json:"title" gorm:"size:255"`
	ChatType        string `json:"chat_type" gorm:"size:50;default:supergroup"`
	AnchorForTierID *uint  `json:"anchor_for_tier_id"`
}

func (SubscriptionChat) TableName() string { return "subscription_chats" }

type SubscriptionTierChat struct {
	TierID uint  `json:"tier_id" gorm:"primaryKey"`
	ChatID int64 `json:"chat_id" gorm:"primaryKey"`
}

func (SubscriptionTierChat) TableName() string { return "subscription_tier_chats" }

type SubscriptionUser struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Username       *string    `json:"username" gorm:"size:255"`
	FullName       string     `json:"full_name" gorm:"size:255"`
	ResolvedTierID *uint      `json:"resolved_tier_id"`
	ManualTierID   *uint      `json:"manual_tier_id"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	LastCheckAt    *time.Time `json:"last_check_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (SubscriptionUser) TableName() string { return "subscription_users" }

// EffectiveTierID returns ManualTierID if set, otherwise ResolvedTierID.
func (u *SubscriptionUser) EffectiveTierID() *uint {
	if u.ManualTierID != nil {
		return u.ManualTierID
	}
	return u.ResolvedTierID
}

type SubscriptionUserChatAccess struct {
	UserID    int64      `json:"user_id" gorm:"primaryKey"`
	ChatID    int64      `json:"chat_id" gorm:"primaryKey"`
	GrantedAt time.Time  `json:"granted_at" gorm:"default:now()"`
	RevokedAt *time.Time `json:"revoked_at"`
}

func (SubscriptionUserChatAccess) TableName() string { return "subscription_user_chat_access" }

type SubscriptionAuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action" gorm:"size:50"`
	Details   string    `json:"details" gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time `json:"created_at"`
}

func (SubscriptionAuditLog) TableName() string { return "subscription_audit_logs" }
