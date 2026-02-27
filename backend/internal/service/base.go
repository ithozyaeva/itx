package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type BaseService[T any] interface {
	Search(limit *int, offset *int, filter *repository.SearchFilter, order *repository.Order) (*models.RegistrySearch[T], error)
	GetById(id int64) (*T, error)
	Create(entity *T) (*T, error)
	Update(entity *T) (*T, error)
	Delete(entity *T) error
}

type baseService[T any] struct {
	repo repository.BaseRepository[T]
}

func NewBaseService[T any](repo repository.BaseRepository[T]) BaseService[T] {
	return &baseService[T]{repo: repo}
}

func (s *baseService[T]) Search(limit *int, offset *int, filter *repository.SearchFilter, order *repository.Order) (*models.RegistrySearch[T], error) {
	items, count, err := s.repo.Search(limit, offset, filter, order)
	if err != nil {
		return nil, err
	}

	return &models.RegistrySearch[T]{
		Items: items,
		Total: int(count),
	}, nil
}

func (s *baseService[T]) GetById(id int64) (*T, error) {
	return s.repo.GetById(id)
}

func (s *baseService[T]) Create(entity *T) (*T, error) {
	return s.repo.Create(entity)
}

func (s *baseService[T]) Update(entity *T) (*T, error) {
	return s.repo.Update(entity)
}

func (s *baseService[T]) Delete(entity *T) error {
	return s.repo.Delete(entity)
}
