package models

import "time"

type TrackedChat struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	ChatID    int64     `json:"chatId" gorm:"column:chat_id;uniqueIndex"`
	Title     string    `json:"title"`
	ChatType  string    `json:"chatType" gorm:"column:chat_type;default:supergroup"`
	IsActive  bool      `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (TrackedChat) TableName() string {
	return "tracked_chats"
}

type ChatMessage struct {
	Id                int64     `json:"id" gorm:"primaryKey"`
	ChatID            int64     `json:"chatId" gorm:"column:chat_id"`
	TelegramUserID    int64     `json:"telegramUserId" gorm:"column:telegram_user_id"`
	TelegramUsername  string    `json:"telegramUsername" gorm:"column:telegram_username"`
	TelegramFirstName string    `json:"telegramFirstName" gorm:"column:telegram_first_name"`
	SentAt            time.Time `json:"sentAt" gorm:"column:sent_at"`
	CreatedAt         time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}

// ChatMessageCount — количество сообщений по чату
type ChatMessageCount struct {
	ChatID int64  `json:"chatId"`
	Title  string `json:"title"`
	Count  int64  `json:"count"`
}

// DailyActivity — сообщения по дням
type DailyActivity struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// TopUser — топ пользователь по сообщениям
type TopUser struct {
	TelegramUserID    int64  `json:"telegramUserId"`
	TelegramUsername  string `json:"telegramUsername"`
	TelegramFirstName string `json:"telegramFirstName"`
	Count             int64  `json:"count"`
	TopChat           string `json:"topChat"`
}

// ChatActivityStats — общая статистика активности
type ChatActivityStats struct {
	TotalMessagesToday int64             `json:"totalMessagesToday"`
	TotalMessagesWeek  int64             `json:"totalMessagesWeek"`
	UniqueUsersToday   int64             `json:"uniqueUsersToday"`
	UniqueUsersWeek    int64             `json:"uniqueUsersWeek"`
	ChatStats          []ChatMessageCount `json:"chatStats"`
}
