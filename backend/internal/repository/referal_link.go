package repository

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"
)

type ReferalLinkRepository struct {
	BaseRepository[models.ReferalLink]
}

func NewReferalLinkRepository() *ReferalLinkRepository {
	return &ReferalLinkRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.ReferalLink{}),
	}
}

func (e *ReferalLinkRepository) Search(limit *int, offset *int, filter *SearchFilter, order *Order) ([]models.ReferalLink, int64, error) {
	var links []models.ReferalLink
	var count int64

	if err := database.DB.Model(&models.ReferalLink{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := database.DB.Model(&models.ReferalLink{}).Preload("Author").Preload("ProfTags")

	if filter != nil {
		for key, value := range *filter {
			query = query.Where(key, value)
		}
	}

	if order != nil {
		query = query.Order(fmt.Sprintf("\"%s\" %s", order.ColumnBy, order.Order))
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&links).Error; err != nil {
		return nil, 0, err
	}

	return links, count, nil
}

func (r *ReferalLinkRepository) Update(entity *models.ReferalLink) (*models.ReferalLink, error) {
	err := database.DB.Model(&entity).Save(entity).Error

	if err != nil {
		return nil, err
	}

	database.DB.Model(&entity).Association("ProfTags").Replace(entity.ProfTags)

	updatedEntity, err := r.GetById(entity.Id)

	if err != nil {
		return nil, err
	}

	return updatedEntity, nil
}

func (r *ReferalLinkRepository) GetById(id int64) (*models.ReferalLink, error) {
	var event models.ReferalLink
	if err := database.DB.Preload("Author").Preload("ProfTags").First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// ExpireLinks замораживает ссылки с истёкшим сроком действия
func (r *ReferalLinkRepository) ExpireLinks() (int64, error) {
	result := database.DB.Model(&models.ReferalLink{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status = ?", time.Now(), models.ReferalLinkActive).
		Update("status", models.ReferalLinkFreezed)
	return result.RowsAffected, result.Error
}
