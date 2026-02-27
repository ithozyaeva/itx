package models

import "time"

type PlaceType string

const (
	EventOnline  PlaceType = "ONLINE"
	EventOffline PlaceType = "OFFLINE"
	EventHybrid  PlaceType = "HYBRID"
)

type RepeatPeriod string

const (
	RepeatDaily   RepeatPeriod = "DAILY"
	RepeatWeekly  RepeatPeriod = "WEEKLY"
	RepeatMonthly RepeatPeriod = "MONTHLY"
	RepeatYearly  RepeatPeriod = "YEARLY"
)

type Event struct {
	Id                       int64      `json:"id" gorm:"primaryKey"`
	Title                    string     `json:"title"`
	Description              string     `json:"description"`
	Date                     time.Time  `json:"date" time_format:"2006-01-02T15:04" time_location:"UTC"`
	Timezone                 string     `json:"timezone" gorm:"default:'UTC'"`
	PlaceType                PlaceType  `json:"placeType"`
	Place                    string     `json:"place"`
	CustomPlaceType          string     `json:"customPlaceType"`
	EventType                string     `json:"eventType"`
	Open                     bool       `json:"open"`
	VideoLink                string     `json:"videoLink" gorm:"column:video_link"`
	IsRepeating              bool       `json:"isRepeating" gorm:"default:false"`
	RepeatPeriod             *string    `json:"repeatPeriod" gorm:"column:repeat_period"`
	RepeatInterval           *int       `json:"repeatInterval" gorm:"column:repeat_interval;default:1"`
	RepeatEndDate            *time.Time `json:"repeatEndDate" gorm:"column:repeat_end_date"`
	EventTags                []EventTag `json:"eventTags" gorm:"many2many:event_event_tags;foreignKey:id;joinForeignKey:event_id;References:id;joinReferences:event_tag_id;replace:true"`
	Hosts                    []Member   `json:"hosts" gorm:"many2many:event_hosts;foreignKey:id;joinForeignKey:event_id;References:id;joinReferences:member_id;replace:true"`
	Members                  []Member   `json:"members" gorm:"many2many:event_members;foreignKey:id;joinForeignKey:event_id;References:id;joinReferences:member_id;replace:true"`
	LastRepeatingAlertSentAt *time.Time `json:"lastRepeatingAlertSentAt" gorm:"column:last_repeating_alert_sent_at"`
}

type EventTag struct {
	Id   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}
