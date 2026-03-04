package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type NotificationSettingsService struct {
	repo *repository.NotificationSettingsRepository
}

func NewNotificationSettingsService() *NotificationSettingsService {
	return &NotificationSettingsService{
		repo: repository.NewNotificationSettingsRepository(),
	}
}

func (s *NotificationSettingsService) GetByMemberId(memberId int64) (*models.NotificationSettings, error) {
	return s.repo.GetByMemberId(memberId)
}

func (s *NotificationSettingsService) GetByMemberIds(memberIds []int64) (map[int64]*models.NotificationSettings, error) {
	return s.repo.GetByMemberIds(memberIds)
}

func (s *NotificationSettingsService) Update(memberId int64, settings *models.NotificationSettings) (*models.NotificationSettings, error) {
	settings.MemberId = memberId
	return s.repo.Upsert(settings)
}
