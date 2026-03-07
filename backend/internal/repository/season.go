package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type SeasonRepository struct{}

func NewSeasonRepository() *SeasonRepository {
	return &SeasonRepository{}
}

func (r *SeasonRepository) GetActive() (*models.Season, error) {
	var season models.Season
	err := database.DB.Where("status = ?", models.SeasonStatusActive).First(&season).Error
	if err != nil {
		return nil, err
	}
	return &season, nil
}

func (r *SeasonRepository) GetAll() ([]models.Season, error) {
	var seasons []models.Season
	err := database.DB.Order("start_date DESC").Find(&seasons).Error
	return seasons, err
}

func (r *SeasonRepository) GetById(id int64) (*models.Season, error) {
	var season models.Season
	err := database.DB.First(&season, id).Error
	if err != nil {
		return nil, err
	}
	return &season, nil
}

func (r *SeasonRepository) Create(season *models.Season) error {
	return database.DB.Create(season).Error
}

func (r *SeasonRepository) Update(season *models.Season) error {
	return database.DB.Save(season).Error
}

func (r *SeasonRepository) GetLeaderboard(seasonId int64, limit int) ([]models.SeasonLeaderboardEntry, error) {
	var season models.Season
	if err := database.DB.First(&season, seasonId).Error; err != nil {
		return nil, err
	}

	entries := make([]models.SeasonLeaderboardEntry, 0)
	err := database.DB.Raw(`
		SELECT m.id as member_id, m.first_name, m.last_name, m.username, m.avatar_url,
			COALESCE(SUM(pt.amount), 0) as total,
			ROW_NUMBER() OVER (ORDER BY COALESCE(SUM(pt.amount), 0) DESC) as rank
		FROM members m
		LEFT JOIN point_transactions pt ON pt.member_id = m.id
			AND pt.created_at >= ? AND pt.created_at <= ?
		GROUP BY m.id, m.first_name, m.last_name, m.username, m.avatar_url
		HAVING COALESCE(SUM(pt.amount), 0) > 0
		ORDER BY total DESC
		LIMIT ?
	`, season.StartDate, season.EndDate, limit).Scan(&entries).Error

	return entries, err
}

func (r *SeasonRepository) FinishSeason(id int64) error {
	return database.DB.Model(&models.Season{}).Where("id = ?", id).Update("status", models.SeasonStatusFinished).Error
}
