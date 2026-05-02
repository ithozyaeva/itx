package models

import "time"

// DailyCheckIn — запись о ежедневном check-in (1 на юзера на МСК-день)
type DailyCheckIn struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	MemberId  int64     `json:"memberId" gorm:"column:member_id;not null;uniqueIndex:uniq_check_in_member_day,priority:1"`
	Day       time.Time `json:"day" gorm:"column:day;type:date;not null;uniqueIndex:uniq_check_in_member_day,priority:2"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (DailyCheckIn) TableName() string {
	return "daily_check_ins"
}

// MemberStreak — сводка по стрику юзера
type MemberStreak struct {
	MemberId          int64      `json:"memberId" gorm:"column:member_id;primaryKey"`
	CurrentStreak     int        `json:"currentStreak" gorm:"column:current_streak;not null;default:0"`
	LongestStreak     int        `json:"longestStreak" gorm:"column:longest_streak;not null;default:0"`
	LastCheckInDate   *time.Time `json:"lastCheckInDate" gorm:"column:last_check_in_date;type:date"`
	FreezesAvailable  int        `json:"freezesAvailable" gorm:"column:freezes_available;not null;default:1"`
	FreezeWeekYear    *int       `json:"freezeWeekYear" gorm:"column:freeze_week_year"`
	FreezeWeekNum     *int       `json:"freezeWeekNum" gorm:"column:freeze_week_num"`
	UpdatedAt         time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (MemberStreak) TableName() string {
	return "member_streaks"
}

// StreakMilestone — описывает порог стрика для UI
type StreakMilestone struct {
	Days    int  `json:"days"`
	Reward  int  `json:"reward"`
	Reached bool `json:"reached"`
}

// StreakResponse — ответ GET /streak/me
type StreakResponse struct {
	Current          int               `json:"current"`
	Longest          int               `json:"longest"`
	FreezesAvailable int               `json:"freezesAvailable"`
	LastCheckIn      *time.Time        `json:"lastCheckIn"`
	NextThreshold    *int              `json:"nextThreshold"`
	DaysToNext       *int              `json:"daysToNext"`
	Milestones       []StreakMilestone `json:"milestones"`
}

// CheckInResponse — ответ POST /dailies/check-in
type CheckInResponse struct {
	CheckInDone   bool           `json:"checkInDone"`
	AlreadyToday  bool           `json:"alreadyToday"`
	Streak        StreakResponse `json:"streak"`
	RaffleEntered bool           `json:"raffleEntered"`
}
