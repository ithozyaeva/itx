package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type FeedbackRepository struct{}

func NewFeedbackRepository() *FeedbackRepository {
	return &FeedbackRepository{}
}

func (r *FeedbackRepository) Create(feedback *models.Feedback) error {
	return database.DB.Create(feedback).Error
}

func (r *FeedbackRepository) CountTodayByMember(memberId int64) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Feedback{}).
		Where("user_id = ? AND created_at >= CURRENT_DATE", memberId).
		Count(&count).Error
	return count, err
}

func (r *FeedbackRepository) List(limit, offset int) ([]models.FeedbackPublic, int64, error) {
	var total int64
	if err := database.DB.Model(&models.Feedback{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]models.FeedbackPublic, 0)
	err := database.DB.Raw(`
		SELECT f.id, f.user_id, f.score, f.comment, f.created_at,
			m.first_name AS user_first_name, m.last_name AS user_last_name, m.username AS user_username
		FROM feedbacks f
		LEFT JOIN members m ON m.id = f.user_id
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset).Scan(&items).Error

	return items, total, err
}
