package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type EventTagService = BaseService[models.EventTag]

func NewEventTagService() EventTagService {
	repo := repository.NewEventTagRepository()
	return NewBaseService(repo)
}