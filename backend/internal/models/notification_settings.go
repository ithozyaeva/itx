package models

type NotificationSettings struct {
	Id             int64 `json:"id" gorm:"primaryKey"`
	MemberId       int64 `json:"memberId" gorm:"column:member_id;uniqueIndex;not null"`
	MuteAll        bool  `json:"muteAll" gorm:"column:mute_all;default:false"`
	NewEvents      bool  `json:"newEvents" gorm:"column:new_events;default:true"`
	RemindWeek     bool  `json:"remindWeek" gorm:"column:remind_week;default:true"`
	RemindDay      bool  `json:"remindDay" gorm:"column:remind_day;default:true"`
	RemindHour     bool  `json:"remindHour" gorm:"column:remind_hour;default:true"`
	EventStart     bool  `json:"eventStart" gorm:"column:event_start;default:true"`
	EventUpdates   bool  `json:"eventUpdates" gorm:"column:event_updates;default:true"`
	EventCancelled bool  `json:"eventCancelled" gorm:"column:event_cancelled;default:true"`
	DailyMorning   bool  `json:"dailyMorning" gorm:"column:daily_morning;default:true"`
	DailyEvening   bool  `json:"dailyEvening" gorm:"column:daily_evening;default:true"`
	DailyStreak    bool  `json:"dailyStreak" gorm:"column:daily_streak;default:true"`
	DailyRaffle    bool  `json:"dailyRaffle" gorm:"column:daily_raffle;default:true"`
}

func (NotificationSettings) TableName() string {
	return "notification_settings"
}
