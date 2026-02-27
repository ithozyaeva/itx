package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	Search(limit *int, offset *int, filter *SearchFilter, order *Order) ([]T, int64, error)
	GetById(id int64) (*T, error)
	Create(entity *T) (*T, error)
	Update(entity *T) (*T, error)
	Delete(entity *T) error
}

type baseRepository[T any] struct {
	db    *gorm.DB
	model *T
}

type SearchFilter = map[string]interface{}

type Order struct {
	ColumnBy string
	Order    string
}

func NewBaseRepository[T any](db *gorm.DB, model *T) BaseRepository[T] {
	return &baseRepository[T]{db: db, model: model}
}

// Реализация методов
func (r *baseRepository[T]) Search(limit *int, offset *int, filter *SearchFilter, order *Order) ([]T, int64, error) {
	var entities []T
	var count int64

	query := r.db.Model(r.model)

	if filter != nil {
		for key, value := range *filter {
			query = query.Where(key, value)
		}
	}

	if order != nil {
		query = query.Order(fmt.Sprintf("\"%s\" %s", order.ColumnBy, order.Order))
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, count, nil
}

func (r *baseRepository[T]) GetById(id int64) (*T, error) {
	entity := new(T)
	if err := r.db.First(entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *baseRepository[T]) Create(entity *T) (*T, error) {
	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *baseRepository[T]) Update(entity *T) (*T, error) {
	if err := r.db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *baseRepository[T]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}
