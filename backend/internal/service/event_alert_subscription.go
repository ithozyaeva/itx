package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type EventAlertSubscriptionService struct {
	repo *repository.EventAlertSubscriptionRepository
}

func NewEventAlertSubscriptionService() *EventAlertSubscriptionService {
	return &EventAlertSubscriptionService{
		repo: repository.NewEventAlertSubscriptionRepository(),
	}
}

// CreateSubscription создает подписку на алерты события
func (s *EventAlertSubscriptionService) CreateSubscription(eventId int64, memberId int64) (*models.EventAlertSubscription, error) {
	subscription := &models.EventAlertSubscription{
		EventId:  eventId,
		MemberId: memberId,
		Status:   models.EventAlertStatusPending,
	}
	return s.repo.Create(subscription)
}

// UpdateSubscriptionStatus обновляет статус подписки
func (s *EventAlertSubscriptionService) UpdateSubscriptionStatus(eventId int64, memberId int64, status models.EventAlertSubscriptionStatus) (*models.EventAlertSubscription, error) {
	subscription, err := s.repo.GetByEventAndMember(eventId, memberId)
	if err != nil {
		return nil, err
	}
	subscription.Status = status
	return s.repo.Update(subscription)
}

// GetSubscribedMembersForEvent получает всех подписанных пользователей для события
func (s *EventAlertSubscriptionService) GetSubscribedMembersForEvent(eventId int64) ([]models.Member, error) {
	return s.repo.GetSubscribedMembersForEvent(eventId)
}

// GetPendingSubscriptionsForEvent получает все подписки со статусом PENDING для события
func (s *EventAlertSubscriptionService) GetPendingSubscriptionsForEvent(eventId int64) ([]models.EventAlertSubscription, error) {
	return s.repo.GetPendingSubscriptionsForEvent(eventId)
}

// CreateOrUpdate создает или обновляет подписку
func (s *EventAlertSubscriptionService) CreateOrUpdate(subscription *models.EventAlertSubscription) (*models.EventAlertSubscription, error) {
	return s.repo.CreateOrUpdate(subscription)
}

