package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type ChatHighlightService struct {
	repo *repository.ChatHighlightRepository
}

func NewChatHighlightService() *ChatHighlightService {
	return &ChatHighlightService{
		repo: repository.NewChatHighlightRepository(),
	}
}

func (s *ChatHighlightService) Create(highlight *models.ChatHighlight) (*models.ChatHighlight, error) {
	return s.repo.Create(highlight)
}

func (s *ChatHighlightService) GetRecent(limit int) ([]models.ChatHighlight, error) {
	return s.repo.GetRecent(limit)
}

func (s *ChatHighlightService) Search(limit, offset int) ([]models.ChatHighlight, int64, error) {
	return s.repo.SearchHighlights(limit, offset)
}
