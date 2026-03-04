package models

import "time"

const (
	QuestTypeMessageCount = "message_count"
	QuestTypeDailyStreak  = "daily_streak"
)

type ChatQuest struct {
	Id           int64     `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"column:title;size:255;not null"`
	Description  string    `json:"description" gorm:"column:description;type:text"`
	QuestType    string    `json:"questType" gorm:"column:quest_type;size:50;not null;default:message_count"`
	ChatID       *int64    `json:"chatId" gorm:"column:chat_id"`
	TargetCount  int       `json:"targetCount" gorm:"column:target_count;not null"`
	PointsReward int       `json:"pointsReward" gorm:"column:points_reward;not null;default:10"`
	StartsAt     time.Time `json:"startsAt" gorm:"column:starts_at;not null"`
	EndsAt       time.Time `json:"endsAt" gorm:"column:ends_at;not null"`
	IsActive     bool      `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ChatQuest) TableName() string {
	return "chat_quests"
}

type ChatQuestProgress struct {
	Id           int64      `json:"id" gorm:"primaryKey"`
	QuestID      int64      `json:"questId" gorm:"column:quest_id;not null"`
	MemberID     int64      `json:"memberId" gorm:"column:member_id;not null"`
	CurrentCount int        `json:"currentCount" gorm:"column:current_count;not null;default:0"`
	Completed    bool       `json:"completed" gorm:"column:completed;default:false"`
	CompletedAt  *time.Time `json:"completedAt" gorm:"column:completed_at"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at"`
}

func (ChatQuestProgress) TableName() string {
	return "chat_quest_progress"
}

// ChatQuestWithProgress — квест с прогрессом для API ответа
type ChatQuestWithProgress struct {
	ChatQuest
	CurrentCount int  `json:"currentCount"`
	Completed    bool `json:"completed"`
}

// ChatQuestStreakDay — запись дня активности для daily_streak квестов
type ChatQuestStreakDay struct {
	Id       int64     `json:"id" gorm:"primaryKey"`
	QuestID  int64     `json:"questId" gorm:"column:quest_id;not null"`
	MemberID int64     `json:"memberId" gorm:"column:member_id;not null"`
	Day      time.Time `json:"day" gorm:"column:day;type:date;not null"`
}

func (ChatQuestStreakDay) TableName() string {
	return "chat_quest_streak_days"
}
