package repository

import (
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	Search(limit *int, offset *int) ([]T, int64, error)
	GetById(id int64) (*T, error)
	Create(entity *T) (*T, error)
	Update(entity *T) (*T, error)
	Delete(entity *T) error
}

type baseRepository[T any] struct {
	db    *gorm.DB
	model *T
}

func NewBaseRepository[T any](db *gorm.DB, model *T) BaseRepository[T] {
	return &baseRepository[T]{db: db, model: model}
}

// Реализация методов
func (r *baseRepository[T]) Search(limit *int, offset *int) ([]T, int64, error) {
	var entities []T
	var count int64

	// Сначала считаем общее количество всех записей
	if err := r.db.Model(r.model).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Создаем запрос
	query := r.db.Model(r.model)

	// Применяем limit только если он передан
	if limit != nil {
		query = query.Limit(*limit)
	}

	// Применяем offset только если он передан
	if offset != nil {
		query = query.Offset(*offset)
	}

	// Выполняем запрос
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
