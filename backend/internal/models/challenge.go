package models

import "time"

const (
	ChallengeKindWeekly  = "weekly"
	ChallengeKindMonthly = "monthly"
)

// ChallengeTemplate — шаблон челленджа (наполняется seed-миграцией, управляется админом)
type ChallengeTemplate struct {
	Id              int64     `json:"id" gorm:"primaryKey"`
	Code            string    `json:"code" gorm:"column:code;size:64;uniqueIndex;not null"`
	Title           string    `json:"title" gorm:"column:title;size:255;not null"`
	Description     string    `json:"description" gorm:"column:description;type:text;default:''"`
	Icon            string    `json:"icon" gorm:"column:icon;size:64;not null;default:'trophy'"`
	Kind            string    `json:"kind" gorm:"column:kind;size:16;not null"`
	MetricKey       string    `json:"metricKey" gorm:"column:metric_key;size:64;not null"`
	Target          int       `json:"target" gorm:"column:target;not null"`
	RewardPoints    int       `json:"rewardPoints" gorm:"column:reward_points;not null"`
	AchievementCode *string   `json:"achievementCode" gorm:"column:achievement_code;size:64"`
	Active          bool      `json:"active" gorm:"column:active;default:true"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (ChallengeTemplate) TableName() string {
	return "challenge_templates"
}

// ChallengeInstance — конкретный запуск шаблона на конкретный период
type ChallengeInstance struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	TemplateId int64     `json:"templateId" gorm:"column:template_id;not null"`
	Kind       string    `json:"kind" gorm:"column:kind;size:16;not null"`
	StartsAt   time.Time `json:"startsAt" gorm:"column:starts_at;not null"`
	EndsAt     time.Time `json:"endsAt" gorm:"column:ends_at;not null"`
	PeriodKey  string    `json:"periodKey" gorm:"column:period_key;size:32;not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (ChallengeInstance) TableName() string {
	return "challenge_instances"
}

// ChallengeProgress — прогресс конкретного юзера по конкретному инстансу
type ChallengeProgress struct {
	Id          int64      `json:"id" gorm:"primaryKey"`
	InstanceId  int64      `json:"instanceId" gorm:"column:instance_id;not null"`
	MemberId    int64      `json:"memberId" gorm:"column:member_id;not null"`
	Progress    int        `json:"progress" gorm:"column:progress;not null;default:0"`
	CompletedAt *time.Time `json:"completedAt" gorm:"column:completed_at"`
	Awarded     bool       `json:"awarded" gorm:"column:awarded;default:false"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (ChallengeProgress) TableName() string {
	return "challenge_progress"
}

// ChallengeWithProgress — DTO для UI: инстанс + шаблон + мой прогресс
type ChallengeWithProgress struct {
	InstanceId      int64      `json:"instanceId"`
	TemplateId      int64      `json:"templateId"`
	Code            string     `json:"code"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Icon            string     `json:"icon"`
	Kind            string     `json:"kind"`
	MetricKey       string     `json:"metricKey"`
	Target          int        `json:"target"`
	RewardPoints    int        `json:"rewardPoints"`
	AchievementCode *string    `json:"achievementCode"`
	StartsAt        time.Time  `json:"startsAt"`
	EndsAt          time.Time  `json:"endsAt"`
	Progress        int        `json:"progress"`
	CompletedAt     *time.Time `json:"completedAt"`
	Awarded         bool       `json:"awarded"`
}

// ChallengesResponse — ответ GET /challenges
type ChallengesResponse struct {
	Weekly  []ChallengeWithProgress `json:"weekly"`
	Monthly []ChallengeWithProgress `json:"monthly"`
}
