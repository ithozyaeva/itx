package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type ReviewOnCommunityRepository struct {
	BaseRepository[models.ReviewOnCommunity]
}

func NewReviewOnCommunityRepository() *ReviewOnCommunityRepository {
	return &ReviewOnCommunityRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.ReviewOnCommunity{}),
	}
}
func (r *ReviewOnCommunityRepository) GetAllWithAuthor(limit *int, offset *int) ([]models.ReviewOnCommunity, int64, error) {
	var reviews []models.ReviewOnCommunity
	var count int64

	query := database.DB.Model(&models.ReviewOnCommunity{}).Preload("Author")

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Применяем limit только если он передан
	if limit != nil {
		query = query.Limit(*limit)
	}

	// Применяем offset только если он передан
	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}
