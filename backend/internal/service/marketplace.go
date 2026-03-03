package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
)

type MarketplaceService struct {
	repo *repository.MarketplaceRepository
}

func NewMarketplaceService() *MarketplaceService {
	return &MarketplaceService{
		repo: repository.NewMarketplaceRepository(),
	}
}

func (s *MarketplaceService) Search(status *string, limit, offset int) ([]models.MarketplaceItem, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.Search(status, limit, offset)
}

func (s *MarketplaceService) GetById(id int64) (*models.MarketplaceItem, error) {
	return s.repo.GetById(id)
}

func (s *MarketplaceService) Create(req *models.CreateMarketplaceItemRequest, sellerId int64, imageFileName string, imageContent []byte, imageContentType string) (*models.MarketplaceItem, error) {
	if req.Title == "" {
		return nil, errors.New("название обязательно")
	}

	item := &models.MarketplaceItem{
		Title:           req.Title,
		Description:     req.Description,
		Price:           req.Price,
		City:            req.City,
		CanShip:         req.CanShip,
		Condition:       req.Condition,
		Defects:         req.Defects,
		PackageContents: req.PackageContents,
		ContactTelegram: req.ContactTelegram,
		ContactEmail:    req.ContactEmail,
		ContactPhone:    req.ContactPhone,
		SellerId:        sellerId,
		Status:          models.MarketplaceStatusActive,
	}

	if item.Condition == "" {
		item.Condition = models.MarketplaceConditionUsed
	}

	if len(imageContent) > 0 && imageFileName != "" {
		s3Client, err := utils.NewS3Client()
		if err != nil {
			log.Printf("marketplace: s3 client error: %v", err)
		} else {
			ext := filepath.Ext(imageFileName)
			if ext == "" {
				ext = ".jpg"
			}
			key := fmt.Sprintf("marketplace/%d/%s%s", sellerId, uuid.New().String(), ext)

			if err := s3Client.UploadPublic(context.Background(), key, imageContent, imageContentType); err != nil {
				log.Printf("marketplace: s3 upload error: %v", err)
			} else {
				item.ImagePath = s3Client.GetPublicURL(key)
			}
		}
	}

	return s.repo.Create(item)
}

func (s *MarketplaceService) Update(id int64, req *models.CreateMarketplaceItemRequest, memberId int64, isAdmin bool) (*models.MarketplaceItem, error) {
	item, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("объявление не найдено")
	}

	if !isAdmin && item.SellerId != memberId {
		return nil, errors.New("только автор может редактировать объявление")
	}

	updates := map[string]interface{}{
		"title":            req.Title,
		"description":      req.Description,
		"price":            req.Price,
		"city":             req.City,
		"can_ship":         req.CanShip,
		"condition":        req.Condition,
		"defects":          req.Defects,
		"package_contents": req.PackageContents,
		"contact_telegram": req.ContactTelegram,
		"contact_email":    req.ContactEmail,
		"contact_phone":    req.ContactPhone,
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	return s.repo.GetById(id)
}

func (s *MarketplaceService) RequestPurchase(id int64, buyerId int64) (*models.MarketplaceItem, error) {
	item, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("объявление не найдено")
	}

	if item.SellerId == buyerId {
		return nil, errors.New("нельзя купить своё объявление")
	}

	rows, err := s.repo.RequestPurchase(id, buyerId)
	if err != nil {
		return nil, errors.New("не удалось забронировать")
	}
	if rows == 0 {
		return nil, errors.New("объявление недоступно для покупки")
	}

	return s.repo.GetById(id)
}

func (s *MarketplaceService) CancelPurchase(id int64, memberId int64, isAdmin bool) (*models.MarketplaceItem, error) {
	item, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("объявление не найдено")
	}

	if !isAdmin && item.SellerId != memberId && (item.BuyerId == nil || *item.BuyerId != memberId) {
		return nil, errors.New("нет прав для отмены заявки")
	}

	rows, err := s.repo.CancelPurchase(id)
	if err != nil {
		return nil, errors.New("не удалось отменить бронь")
	}
	if rows == 0 {
		return nil, errors.New("объявление не забронировано")
	}

	return s.repo.GetById(id)
}

func (s *MarketplaceService) MarkSold(id int64, memberId int64, isAdmin bool) (*models.MarketplaceItem, error) {
	item, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("объявление не найдено")
	}

	if !isAdmin && item.SellerId != memberId {
		return nil, errors.New("только продавец может подтвердить продажу")
	}

	rows, err := s.repo.MarkSold(id)
	if err != nil {
		return nil, errors.New("не удалось подтвердить продажу")
	}
	if rows == 0 {
		return nil, errors.New("объявление не забронировано")
	}

	return s.repo.GetById(id)
}

func (s *MarketplaceService) Delete(id int64, memberId int64, isAdmin bool) error {
	item, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("объявление не найдено")
	}

	if isAdmin {
		return s.repo.Delete(id)
	}

	if item.SellerId != memberId {
		return errors.New("только автор может удалить объявление")
	}

	if item.Status != models.MarketplaceStatusActive {
		return errors.New("можно удалить только активные объявления")
	}

	return s.repo.Delete(id)
}

func ParseCondition(val string) models.MarketplaceItemCondition {
	upper := models.MarketplaceItemCondition(strings.ToUpper(val))
	if upper == models.MarketplaceConditionNew || upper == models.MarketplaceConditionUsed {
		return upper
	}
	return models.MarketplaceConditionUsed
}
