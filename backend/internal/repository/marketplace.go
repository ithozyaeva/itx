package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type MarketplaceRepository struct{}

func NewMarketplaceRepository() *MarketplaceRepository {
	return &MarketplaceRepository{}
}

func (r *MarketplaceRepository) Search(status *string, limit, offset int) ([]models.MarketplaceItem, int64, error) {
	var items []models.MarketplaceItem
	var total int64

	countQuery := database.DB.Model(&models.MarketplaceItem{})
	if status != nil && *status != "" {
		countQuery = countQuery.Where("status = ?", *status)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := database.DB.
		Preload("Seller").
		Preload("Buyer")

	if status != nil && *status != "" {
		findQuery = findQuery.Where("status = ?", *status)
	}

	err := findQuery.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

func (r *MarketplaceRepository) GetById(id int64) (*models.MarketplaceItem, error) {
	var item models.MarketplaceItem
	err := database.DB.
		Preload("Seller").
		Preload("Buyer").
		First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketplaceRepository) Create(item *models.MarketplaceItem) (*models.MarketplaceItem, error) {
	if err := database.DB.Create(item).Error; err != nil {
		return nil, err
	}
	return r.GetById(item.Id)
}

func (r *MarketplaceRepository) Update(id int64, updates map[string]interface{}) error {
	return database.DB.Model(&models.MarketplaceItem{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MarketplaceRepository) RequestPurchase(id int64, buyerId int64) (int64, error) {
	result := database.DB.Model(&models.MarketplaceItem{}).
		Where("id = ? AND status = ?", id, models.MarketplaceStatusActive).
		Updates(map[string]interface{}{
			"buyer_id": buyerId,
			"status":   models.MarketplaceStatusReserved,
		})
	return result.RowsAffected, result.Error
}

func (r *MarketplaceRepository) CancelPurchase(id int64) (int64, error) {
	result := database.DB.Model(&models.MarketplaceItem{}).
		Where("id = ? AND status = ?", id, models.MarketplaceStatusReserved).
		Updates(map[string]interface{}{
			"buyer_id": nil,
			"status":   models.MarketplaceStatusActive,
		})
	return result.RowsAffected, result.Error
}

func (r *MarketplaceRepository) MarkSold(id int64) (int64, error) {
	result := database.DB.Model(&models.MarketplaceItem{}).
		Where("id = ? AND status = ?", id, models.MarketplaceStatusReserved).
		Update("status", models.MarketplaceStatusSold)
	return result.RowsAffected, result.Error
}

func (r *MarketplaceRepository) Delete(id int64) error {
	return database.DB.Delete(&models.MarketplaceItem{}, id).Error
}
