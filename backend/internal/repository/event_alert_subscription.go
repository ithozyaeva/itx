package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type EventAlertSubscriptionRepository struct {
	BaseRepository[models.EventAlertSubscription]
}

func NewEventAlertSubscriptionRepository() *EventAlertSubscriptionRepository {
	return &EventAlertSubscriptionRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.EventAlertSubscription{}),
	}
}

// GetByEventAndMember получает подписку по событию и пользователю
func (r *EventAlertSubscriptionRepository) GetByEventAndMember(eventId int64, memberId int64) (*models.EventAlertSubscription, error) {
	var subscription models.EventAlertSubscription
	err := database.DB.Where("event_id = ? AND member_id = ?", eventId, memberId).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetSubscribedMembersForEvent получает всех подписанных пользователей для события
func (r *EventAlertSubscriptionRepository) GetSubscribedMembersForEvent(eventId int64) ([]models.Member, error) {
	var members []models.Member
	err := database.DB.
		Table("members").
		Joins("INNER JOIN event_alert_subscriptions ON members.id = event_alert_subscriptions.member_id").
		Where("event_alert_subscriptions.event_id = ? AND event_alert_subscriptions.status = ?", eventId, models.EventAlertStatusSubscribed).
		Where("members.telegram_id IS NOT NULL AND members.telegram_id != 0").
		Find(&members).Error
	return members, err
}

// GetPendingSubscriptionsForEvent получает все подписки со статусом PENDING для события
func (r *EventAlertSubscriptionRepository) GetPendingSubscriptionsForEvent(eventId int64) ([]models.EventAlertSubscription, error) {
	var subscriptions []models.EventAlertSubscription
	err := database.DB.
		Where("event_id = ? AND status = ?", eventId, models.EventAlertStatusPending).
		Preload("Member").
		Preload("Event").
		Find(&subscriptions).Error
	return subscriptions, err
}

// CreateOrUpdate создает или обновляет подписку
func (r *EventAlertSubscriptionRepository) CreateOrUpdate(subscription *models.EventAlertSubscription) (*models.EventAlertSubscription, error) {
	var existing models.EventAlertSubscription
	err := database.DB.Where("event_id = ? AND member_id = ?", subscription.EventId, subscription.MemberId).First(&existing).Error
	
	if err != nil {
		return r.Create(subscription)
	}
	
	existing.Status = subscription.Status
	if subscription.ReminderSentAt != nil {
		existing.ReminderSentAt = subscription.ReminderSentAt
	}
	return r.Update(&existing)
}

