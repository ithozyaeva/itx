package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"time"
)

type ReferalLinkService struct {
	BaseService[models.ReferalLink]
	repo repository.ReferalLinkRepository
}

func NewReferalLinkService() *ReferalLinkService {
	repo := repository.NewReferalLinkRepository()
	return &ReferalLinkService{
		BaseService: NewBaseService(repo),
		repo:        *repo,
	}
}

func (s *ReferalLinkService) SearchWithMember(limit *int, offset *int, filter *repository.SearchFilter, order *repository.Order, memberId int64) (*models.RegistrySearch[models.ReferalLink], error) {
	items, count, err := s.repo.SearchWithMember(limit, offset, filter, order, memberId)
	if err != nil {
		return nil, err
	}
	return &models.RegistrySearch[models.ReferalLink]{
		Items: items,
		Total: int(count),
	}, nil
}

func (s *ReferalLinkService) AddLink(req *models.AddLinkRequest, member *models.Member) (*models.ReferalLink, error) {
	newEntity := &models.ReferalLink{
		Author:         *member,
		Company:        req.Company,
		Grade:          req.Grade,
		ProfTags:       req.ProfTags,
		Status:         models.ReferalLinkActive,
		VacationsCount: req.VacationsCount,
		ExpiresAt:      req.ExpiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.repo.Create(newEntity)
}

func (s *ReferalLinkService) UpdateLink(req *models.UpdateLinkRequest, member *models.Member) (*models.ReferalLink, error) {
	updatedEntity := &models.ReferalLink{
		Id:             req.Id,
		Author:         *member,
		Company:        req.Company,
		Grade:          req.Grade,
		ProfTags:       req.ProfTags,
		Status:         req.Status,
		VacationsCount: req.VacationsCount,
		ExpiresAt:      req.ExpiresAt,
		UpdatedAt:      time.Now(),
	}

	return s.repo.Update(updatedEntity)
}

func (s *ReferalLinkService) TrackConversion(linkId int64, memberId int64) error {
	return s.repo.TrackConversion(linkId, memberId)
}

// ExpireLinks замораживает реферальные ссылки с истёкшим сроком действия
func (s *ReferalLinkService) ExpireLinks() (int64, error) {
	return s.repo.ExpireLinks()
}
