package models

import "time"

// EventAlertSubscriptionStatus статус подписки на алерты события
type EventAlertSubscriptionStatus string

const (
	// EventAlertStatusPending - пользователь еще не ответил на инициализирующий алерт
	EventAlertStatusPending EventAlertSubscriptionStatus = "PENDING"
	// EventAlertStatusSubscribed - пользователь подписан (нажал "Приду")
	EventAlertStatusSubscribed EventAlertSubscriptionStatus = "SUBSCRIBED"
	// EventAlertStatusUnsubscribed - пользователь отписан (нажал "Не приду" или не ответил на повторный)
	EventAlertStatusUnsubscribed EventAlertSubscriptionStatus = "UNSUBSCRIBED"
)

// EventAlertSubscription подписка пользователя на алерты события
type EventAlertSubscription struct {
	Id                    int64                        `json:"id" gorm:"primaryKey"`
	EventId               int64                        `json:"eventId" gorm:"column:event_id;not null;index"`
	MemberId              int64                        `json:"memberId" gorm:"column:member_id;not null;index"`
	Status                EventAlertSubscriptionStatus `json:"status" gorm:"type:varchar(50);not null;default:'PENDING'"`
	ReminderSentAt        *time.Time                   `json:"reminderSentAt" gorm:"column:reminder_sent_at"`
	CreatedAt             time.Time                    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt             time.Time                    `json:"updatedAt" gorm:"column:updated_at"`
	
	Event  Event  `json:"event,omitempty" gorm:"foreignKey:EventId;references:Id"`
	Member Member `json:"member,omitempty" gorm:"foreignKey:MemberId;references:Id"`
}

func (EventAlertSubscription) TableName() string {
	return "event_alert_subscriptions"
}

