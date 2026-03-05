package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type KudosRepository struct{}

func NewKudosRepository() *KudosRepository {
	return &KudosRepository{}
}

func (r *KudosRepository) Create(kudos *models.Kudos) error {
	return database.DB.Create(kudos).Error
}

func (r *KudosRepository) GetRecent(limit, offset int) ([]models.KudosPublic, int64, error) {
	var total int64
	database.DB.Model(&models.Kudos{}).Count(&total)

	var items []models.KudosPublic
	err := database.DB.Raw(`
		SELECT k.id, k.from_id, k.message, k.created_at,
			f.first_name as from_first_name, f.last_name as from_last_name, f.username as from_username, f.avatar_url as from_avatar_url,
			t.id as to_id, t.first_name as to_first_name, t.last_name as to_last_name, t.username as to_username, t.avatar_url as to_avatar_url
		FROM kudos k
		JOIN members f ON f.id = k.from_id
		JOIN members t ON t.id = k.to_id
		ORDER BY k.created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset).Scan(&items).Error

	return items, total, err
}

func (r *KudosRepository) CountTodayByFrom(fromId int64) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Kudos{}).
		Where("from_id = ? AND created_at >= CURRENT_DATE", fromId).
		Count(&count).Error
	return count, err
}

func (r *KudosRepository) GetReceivedCount(memberId int64) (int, error) {
	var count int64
	err := database.DB.Model(&models.Kudos{}).Where("to_id = ?", memberId).Count(&count).Error
	return int(count), err
}

func (r *KudosRepository) GetSentCount(memberId int64) (int, error) {
	var count int64
	err := database.DB.Model(&models.Kudos{}).Where("from_id = ?", memberId).Count(&count).Error
	return int(count), err
}
