package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type ChatHighlightRepository struct {
	BaseRepository[models.ChatHighlight]
}

func NewChatHighlightRepository() *ChatHighlightRepository {
	return &ChatHighlightRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.ChatHighlight{}),
	}
}

// GetRecent returns the most recent highlights
func (r *ChatHighlightRepository) GetRecent(limit int) ([]models.ChatHighlight, error) {
	var highlights []models.ChatHighlight
	err := database.DB.Order("created_at DESC").Limit(limit).Find(&highlights).Error
	if err != nil {
		return nil, err
	}
	return highlights, nil
}

// SearchHighlights returns highlights with pagination
func (r *ChatHighlightRepository) SearchHighlights(limit, offset int) ([]models.ChatHighlight, int64, error) {
	var highlights []models.ChatHighlight
	var count int64

	err := database.DB.Model(&models.ChatHighlight{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = database.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&highlights).Error
	if err != nil {
		return nil, 0, err
	}

	return highlights, count, nil
}
