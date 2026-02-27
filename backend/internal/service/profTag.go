package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type ProfTagService = BaseService[models.ProfTag]

func NewProfTagService() ProfTagService {
	repo := repository.NewProfTagRepository()
	return NewBaseService(repo)
}
