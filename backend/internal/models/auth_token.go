package models

import (
	"time"
)

type AuthToken struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	TelegramID int64     `json:"telegram_id" gorm:"uniqueIndex"`
	ExpiredAt  time.Time `json:"expired_at"`
	Token      string    `json:"token"`
}
