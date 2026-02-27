package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type EventTagRepository = BaseRepository[models.EventTag]

func NewEventTagRepository() EventTagRepository {
	return NewBaseRepository(database.DB, &models.EventTag{})
}