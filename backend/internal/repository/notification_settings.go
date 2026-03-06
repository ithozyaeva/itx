package repository

import (
	"errors"

	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

type NotificationSettingsRepository struct {
	BaseRepository[models.NotificationSettings]
}

func NewNotificationSettingsRepository() *NotificationSettingsRepository {
	return &NotificationSettingsRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.NotificationSettings{}),
	}
}

// GetByMemberId returns notification settings for a member, creating defaults if none exist
func (r *NotificationSettingsRepository) GetByMemberId(memberId int64) (*models.NotificationSettings, error) {
	var settings models.NotificationSettings
	err := database.DB.Where("member_id = ?", memberId).First(&settings).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// Create default settings
		settings = models.NotificationSettings{
			MemberId:       memberId,
			NewEvents:      true,
			RemindWeek:     true,
			RemindDay:      true,
			RemindHour:     true,
			EventStart:     true,
			EventUpdates:   true,
			EventCancelled: true,
		}
		if err := database.DB.Create(&settings).Error; err != nil {
			return nil, err
		}
		return &settings, nil
	}
	return &settings, nil
}

// GetByMemberIds returns notification settings for multiple members (batch)
func (r *NotificationSettingsRepository) GetByMemberIds(memberIds []int64) (map[int64]*models.NotificationSettings, error) {
	var settingsList []models.NotificationSettings
	err := database.DB.Where("member_id IN ?", memberIds).Find(&settingsList).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int64]*models.NotificationSettings)
	for i := range settingsList {
		result[settingsList[i].MemberId] = &settingsList[i]
	}
	return result, nil
}

// Upsert creates or updates notification settings
func (r *NotificationSettingsRepository) Upsert(settings *models.NotificationSettings) (*models.NotificationSettings, error) {
	var existing models.NotificationSettings
	err := database.DB.Where("member_id = ?", settings.MemberId).First(&existing).Error

	if err != nil {
		if err := database.DB.Create(settings).Error; err != nil {
			return nil, err
		}
		return settings, nil
	}

	existing.MuteAll = settings.MuteAll
	existing.NewEvents = settings.NewEvents
	existing.RemindWeek = settings.RemindWeek
	existing.RemindDay = settings.RemindDay
	existing.RemindHour = settings.RemindHour
	existing.EventStart = settings.EventStart
	existing.EventUpdates = settings.EventUpdates
	existing.EventCancelled = settings.EventCancelled

	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}
