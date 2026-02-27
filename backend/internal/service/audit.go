package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"
	"time"
)

type AuditService struct {
	repo *repository.AuditLogRepository
}

func NewAuditService() *AuditService {
	return &AuditService{
		repo: repository.NewAuditLogRepository(),
	}
}

func (s *AuditService) Log(actorId int64, actorName string, actorType models.ActorType, action models.AuditAction, entityType string, entityId int64, entityName string) {
	entry := &models.AuditLog{
		ActorId:    actorId,
		ActorName:  actorName,
		ActorType:  actorType,
		Action:     action,
		EntityType: entityType,
		EntityId:   entityId,
		EntityName: entityName,
		CreatedAt:  time.Now(),
	}

	if _, err := s.repo.Create(entry); err != nil {
		log.Printf("Failed to create audit log: %v", err)
	}
}

func (s *AuditService) Search(limit *int, offset *int, filter *repository.SearchFilter) (*models.RegistrySearch[models.AuditLog], error) {
	items, count, err := s.repo.Search(limit, offset, filter, &repository.Order{
		ColumnBy: "created_at",
		Order:    "DESC",
	})
	if err != nil {
		return nil, err
	}

	return &models.RegistrySearch[models.AuditLog]{
		Items: items,
		Total: int(count),
	}, nil
}
