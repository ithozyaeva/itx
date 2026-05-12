package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/utils"
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
	// «Сегодня» считаем по МСК-полночи, а не по CURRENT_DATE (UTC session TZ):
	// иначе дневной лимит сбрасывается в 03:00 MSK, а не в 00:00 как ожидает
	// юзер. Симметрично с birthday-checker и daily_task (см. utils/msktime.go).
	since := utils.MSKToday()
	err := database.DB.Model(&models.Feedback{}).
		Where("user_id = ? AND created_at >= ?", memberId, since).
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
