package models

import (
	"time"

	"github.com/lib/pq"
)

const (
	DailyTaskTierEngagement  = "engagement"
	DailyTaskTierLight       = "light"
	DailyTaskTierMeaningful  = "meaningful"
	DailyTaskTierBig         = "big"
)

// DailyTask — шаблон ежедневного задания (управляется админом, наполняется seed-миграцией)
type DailyTask struct {
	Id          int64     `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"column:code;size:64;uniqueIndex;not null"`
	Title       string    `json:"title" gorm:"column:title;size:255;not null"`
	Description string    `json:"description" gorm:"column:description;type:text;default:''"`
	Icon        string    `json:"icon" gorm:"column:icon;size:64;not null;default:'circle'"`
	Tier        string    `json:"tier" gorm:"column:tier;size:16;not null"`
	Points      int       `json:"points" gorm:"column:points;not null"`
	Target      int       `json:"target" gorm:"column:target;not null;default:1"`
	TriggerKey  string    `json:"triggerKey" gorm:"column:trigger_key;size:64;not null"`
	Active      bool      `json:"active" gorm:"column:active;default:true"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (DailyTask) TableName() string {
	return "daily_tasks"
}

// DailyTaskSet — общий на всех юзеров набор задач, выбранных на конкретный МСК-день
type DailyTaskSet struct {
	Day       time.Time     `json:"day" gorm:"column:day;type:date;primaryKey"`
	TaskIds   pq.Int64Array `json:"taskIds" gorm:"column:task_ids;type:bigint[];not null"`
	CreatedAt time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (DailyTaskSet) TableName() string {
	return "daily_task_sets"
}

// DailyTaskProgress — прогресс конкретного юзера по задаче конкретного дня
type DailyTaskProgress struct {
	Id           int64      `json:"id" gorm:"primaryKey"`
	MemberId     int64      `json:"memberId" gorm:"column:member_id;not null;index:idx_dtp_member_day,priority:1"`
	Day          time.Time  `json:"day" gorm:"column:day;type:date;not null;index:idx_dtp_member_day,priority:2"`
	TaskId       int64      `json:"taskId" gorm:"column:task_id;not null"`
	Progress     int        `json:"progress" gorm:"column:progress;not null;default:0"`
	CompletedAt  *time.Time `json:"completedAt" gorm:"column:completed_at"`
	Awarded      bool       `json:"awarded" gorm:"column:awarded;default:false"`
	BonusAwarded bool       `json:"bonusAwarded" gorm:"column:bonus_awarded;default:false"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (DailyTaskProgress) TableName() string {
	return "daily_task_progress"
}

// DailyTaskWithProgress — DTO для UI: задача + прогресс конкретного юзера на сегодня
type DailyTaskWithProgress struct {
	DailyTask
	Progress    int        `json:"progress"`
	CompletedAt *time.Time `json:"completedAt"`
	Awarded     bool       `json:"awarded"`
}

// DailyTodayResponse — полный ответ GET /dailies/today
type DailyTodayResponse struct {
	CheckIn  DailyCheckInState       `json:"checkIn"`
	Tasks    []DailyTaskWithProgress `json:"tasks"`
	AllBonus DailyAllBonusState      `json:"allBonus"`
}

type DailyCheckInState struct {
	Done bool       `json:"done"`
	At   *time.Time `json:"at"`
}

type DailyAllBonusState struct {
	Points  int  `json:"points"`
	Awarded bool `json:"awarded"`
}
