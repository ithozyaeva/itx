package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type AuditLogRepository struct {
	BaseRepository[models.AuditLog]
}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.AuditLog{}),
	}
}
