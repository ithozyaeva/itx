package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type SeasonService struct {
	repo *repository.SeasonRepository
}

func NewSeasonService() *SeasonService {
	return &SeasonService{
		repo: repository.NewSeasonRepository(),
	}
}

func (s *SeasonService) GetActive() (*models.Season, error) {
	return s.repo.GetActive()
}

func (s *SeasonService) GetAll() ([]models.Season, error) {
	return s.repo.GetAll()
}

func (s *SeasonService) Create(season *models.Season) error {
	return s.repo.Create(season)
}

func (s *SeasonService) Update(season *models.Season) error {
	return s.repo.Update(season)
}

func (s *SeasonService) Finish(id int64) error {
	return s.repo.FinishSeason(id)
}

func (s *SeasonService) GetLeaderboard(seasonId int64, limit int) (*models.SeasonWithLeaderboard, error) {
	season, err := s.repo.GetById(seasonId)
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	leaderboard, err := s.repo.GetLeaderboard(seasonId, limit)
	if err != nil {
		return nil, err
	}

	return &models.SeasonWithLeaderboard{
		Season:      *season,
		Leaderboard: leaderboard,
	}, nil
}

func (s *SeasonService) GetActiveWithLeaderboard(limit int) (*models.SeasonWithLeaderboard, error) {
	season, err := s.repo.GetActive()
	if err != nil {
		return nil, err
	}
	return s.GetLeaderboard(season.Id, limit)
}
