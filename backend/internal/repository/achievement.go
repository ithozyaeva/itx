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

// GetExplicitGrants — все коды ачивок, выданных пользователю явно через
// achievement_grants (например, по completion челленджа). Дополняет
// reason-based достижения из AllAchievements.
func (r *AchievementRepository) GetExplicitGrants(memberId int64) (map[string]bool, error) {
	type row struct {
		Code string
	}
	var rows []row
	err := database.DB.Raw(
		`SELECT code FROM achievement_grants WHERE member_id = ?`,
		memberId,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(rows))
	for _, r := range rows {
		out[r.Code] = true
	}
	return out, nil
}

// GrantExplicit идемпотентно выдаёт явную ачивку — используется
// сервисом челленджей при completion. ON CONFLICT — идемпотентно.
func (r *AchievementRepository) GrantExplicit(memberId int64, code string) error {
	return database.DB.Exec(
		`INSERT INTO achievement_grants (member_id, code) VALUES (?, ?) ON CONFLICT (member_id, code) DO NOTHING`,
		memberId, code,
	).Error
}
