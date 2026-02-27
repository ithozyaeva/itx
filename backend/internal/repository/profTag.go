package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type ProfTagRepository = BaseRepository[models.ProfTag]

func NewProfTagRepository() ProfTagRepository {
	return NewBaseRepository(database.DB, &models.ProfTag{})
}
