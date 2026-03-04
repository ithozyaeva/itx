package models

import "time"

type ChatHighlight struct {
	Id               int64     `json:"id" gorm:"primaryKey"`
	ChatID           int64     `json:"chatId" gorm:"column:chat_id;not null"`
	MessageID        int       `json:"messageId" gorm:"column:message_id;not null"`
	AuthorTelegramID int64     `json:"authorTelegramId" gorm:"column:author_telegram_id;not null"`
	AuthorUsername   string    `json:"authorUsername" gorm:"column:author_username"`
	AuthorFirstName  string    `json:"authorFirstName" gorm:"column:author_first_name"`
	MessageText      string    `json:"messageText" gorm:"column:message_text;type:text;not null"`
	HighlightedBy    int64     `json:"highlightedBy" gorm:"column:highlighted_by;not null"`
	MemberID         *int64    `json:"memberId" gorm:"column:member_id"`
	CreatedAt        time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ChatHighlight) TableName() string {
	return "chat_highlights"
}
