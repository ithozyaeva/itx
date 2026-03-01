package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type AchievementRepository struct{}

func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{}
}

func (r *AchievementRepository) GetReasonCounts(memberId int64) (map[models.PointReason]int, error) {
	type reasonCount struct {
		Reason models.PointReason
		Count  int
	}

	var results []reasonCount
	err := database.DB.Raw(
		`SELECT reason, COUNT(*) as count FROM point_transactions WHERE member_id = ? GROUP BY reason`,
		memberId,
	).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[models.PointReason]int)
	for _, rc := range results {
		counts[rc.Reason] = rc.Count
	}
	return counts, nil
}
