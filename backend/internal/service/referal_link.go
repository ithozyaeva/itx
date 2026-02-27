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

func (s *ReferalLinkService) AddLink(req *models.AddLinkRequest, member *models.Member) (*models.ReferalLink, error) {
	newEntity := &models.ReferalLink{
		Author:         *member,
		Company:        req.Company,
		Grade:          req.Grade,
		ProfTags:       req.ProfTags,
		Status:         models.ReferalLinkActive,
		VacationsCount: req.VacationsCount,
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
		UpdatedAt:      time.Now(),
	}

	return s.repo.Update(updatedEntity)
}
