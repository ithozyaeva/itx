package repository

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"log"
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

	e.LoadConversionsCounts(links)

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

func (r *ReferalLinkRepository) TrackConversion(linkId int64, memberId int64) error {
	conversion := &models.ReferralConversion{
		ReferralLinkId: linkId,
		MemberId:       memberId,
	}
	return database.DB.Create(conversion).Error
}

func (r *ReferalLinkRepository) GetConversionsCount(linkId int64) (int64, error) {
	var count int64
	err := database.DB.Model(&models.ReferralConversion{}).Where("referral_link_id = ?", linkId).Count(&count).Error
	return count, err
}

func (r *ReferalLinkRepository) LoadConversionsCounts(links []models.ReferalLink) {
	if len(links) == 0 {
		return
	}

	ids := make([]int64, len(links))
	for i, link := range links {
		ids[i] = link.Id
	}

	type countResult struct {
		ReferralLinkId int64 `gorm:"column:referral_link_id"`
		Count          int64 `gorm:"column:count"`
	}

	var counts []countResult
	if err := database.DB.Model(&models.ReferralConversion{}).
		Select("referral_link_id, COUNT(*) as count").
		Where("referral_link_id IN ?", ids).
		Group("referral_link_id").
		Scan(&counts).Error; err != nil {
		log.Printf("Error loading conversion counts: %v", err)
		return
	}

	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.ReferralLinkId] = c.Count
	}

	for i := range links {
		links[i].ConversionsCount = countMap[links[i].Id]
	}
}

// ExpireLinks замораживает ссылки с истёкшим сроком действия
func (r *ReferalLinkRepository) ExpireLinks() (int64, error) {
	result := database.DB.Model(&models.ReferalLink{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status = ?", time.Now(), models.ReferalLinkActive).
		Update("status", models.ReferalLinkFreezed)
	return result.RowsAffected, result.Error
}
