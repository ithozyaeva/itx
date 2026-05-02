package models

import (
	"time"

	"ithozyeva/internal/s3resolve"

	"gorm.io/gorm"
)

type PointReason string

const (
	PointReasonEventAttend      PointReason = "event_attend"
	PointReasonEventHost        PointReason = "event_host"
	PointReasonReviewCommunity  PointReason = "review_community"
	PointReasonReviewService    PointReason = "review_service"
	PointReasonResumeUpload     PointReason = "resume_upload"
	PointReasonReferalCreate    PointReason = "referal_create"
	PointReasonReferalConversion PointReason = "referal_conversion"
	PointReasonProfileComplete  PointReason = "profile_complete"
	PointReasonWeeklyActivity   PointReason = "weekly_activity"
	PointReasonMonthlyActive    PointReason = "monthly_active"
	PointReasonStreak4Weeks     PointReason = "streak_4weeks"
	PointReasonAdminManual      PointReason = "admin_manual"
	PointReasonTaskCreate          PointReason = "task_create"
	PointReasonTaskExecute         PointReason = "task_execute"
	PointReasonMarketplaceCreate   PointReason = "marketplace_create"
	PointReasonMarketplaceBuy      PointReason = "marketplace_buy"
	PointReasonChatQuest           PointReason = "chat_quest"
	PointReasonChatterOfWeek       PointReason = "chatter_of_week"
	PointReasonKudosReceived       PointReason = "kudos_received"
	PointReasonRaffleSpend         PointReason = "raffle_spend"
	PointReasonCasinoBet           PointReason = "casino_bet"
	PointReasonCasinoWin           PointReason = "casino_win"
	PointReasonDailyCheckIn        PointReason = "daily_checkin"
	PointReasonDailyStreak3        PointReason = "daily_streak_3"
	PointReasonDailyStreak7        PointReason = "daily_streak_7"
	PointReasonDailyStreak14       PointReason = "daily_streak_14"
	PointReasonDailyStreak30       PointReason = "daily_streak_30"
	PointReasonDailyTaskComplete   PointReason = "daily_task_complete"
	PointReasonDailyAllTasksBonus  PointReason = "daily_all_tasks_bonus"
	PointReasonChallengeComplete   PointReason = "challenge_complete"
	PointReasonDailyRaffleWin      PointReason = "daily_raffle_win"
)

var PointValues = map[PointReason]int{
	PointReasonEventAttend:      10,
	PointReasonEventHost:        25,
	PointReasonReviewCommunity:  15,
	PointReasonReviewService:    10,
	PointReasonResumeUpload:     10,
	PointReasonReferalCreate:    5,
	PointReasonReferalConversion: 30,
	PointReasonProfileComplete:  20,
	PointReasonWeeklyActivity:   5,
	PointReasonMonthlyActive:    30,
	PointReasonStreak4Weeks:     50,
	PointReasonTaskCreate:          15,
	PointReasonTaskExecute:         25,
	PointReasonMarketplaceCreate:   15,
	PointReasonMarketplaceBuy:      10,
	PointReasonChatterOfWeek:       15,
	PointReasonKudosReceived:       5,
	PointReasonDailyCheckIn:        5,
	PointReasonDailyStreak3:        15,
	PointReasonDailyStreak7:        50,
	PointReasonDailyStreak14:       150,
	PointReasonDailyStreak30:       500,
	PointReasonDailyAllTasksBonus:  50,
	PointReasonDailyRaffleWin:      100,
}

type PointTransaction struct {
	Id          int64       `json:"id" gorm:"primaryKey"`
	MemberId    int64       `json:"memberId" gorm:"column:member_id;not null"`
	Amount      int         `json:"amount" gorm:"not null"`
	Reason      PointReason `json:"reason" gorm:"column:reason;size:50;not null"`
	SourceType  string      `json:"sourceType" gorm:"column:source_type;size:50;not null"`
	SourceId    int64       `json:"sourceId" gorm:"column:source_id;not null;default:0"`
	Description string      `json:"description" gorm:"column:description;default:''"`
	CreatedAt   time.Time   `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

type MemberPointsBalance struct {
	MemberId  int64  `json:"memberId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"tg"`
	AvatarURL string `json:"avatarUrl"`
	Total     int    `json:"total"`
}

func (m *MemberPointsBalance) AfterFind(tx *gorm.DB) (err error) {
	m.AvatarURL = s3resolve.ResolveS3URL(m.AvatarURL)
	return nil
}

type PointsSummary struct {
	Balance      int                `json:"balance"`
	Transactions []PointTransaction `json:"transactions"`
}

type AdminPointTransaction struct {
	Id              int64       `json:"id"`
	MemberId        int64       `json:"memberId"`
	MemberFirstName string      `json:"memberFirstName"`
	MemberLastName  string      `json:"memberLastName"`
	MemberUsername  string      `json:"memberUsername"`
	Amount          int         `json:"amount"`
	Reason          PointReason `json:"reason"`
	SourceType      string      `json:"sourceType"`
	Description     string      `json:"description"`
	CreatedAt       time.Time   `json:"createdAt"`
}

type AdminAwardRequest struct {
	MemberId    int64  `json:"memberId"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}
