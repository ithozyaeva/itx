package models

import "time"

type Notification struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	MemberId  int64     `json:"memberId" gorm:"column:member_id"`
	Type      string    `json:"type" gorm:"default:'info'"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Read      bool      `json:"read" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}
