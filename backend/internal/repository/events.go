package repository

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"
)

type EventRepository struct {
	BaseRepository[models.Event]
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.Event{}),
	}
}

func (e *EventRepository) Search(limit *int, offset *int, filter *SearchFilter, order *Order) ([]models.Event, int64, error) {
	var events []models.Event
	var count int64

	if err := database.DB.Model(&models.Event{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := database.DB.Model(&models.Event{}).Preload("Hosts").Preload("Members").Preload("EventTags")

	if filter != nil {
		for key, value := range *filter {
			query = query.Where(key, value)
		}
	}

	if order != nil {
		query = query.Order(fmt.Sprintf("\"%s\" %s", order.ColumnBy, order.Order))
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

func (r *EventRepository) Update(entity *models.Event) (*models.Event, error) {
	oldEvent, err := r.GetById(entity.Id)
	if err != nil {
		return nil, err
	}

	oldDateTruncated := oldEvent.Date.Truncate(time.Minute)
	newDateTruncated := entity.Date.Truncate(time.Minute)
	if !oldDateTruncated.Equal(newDateTruncated) {
		entity.LastRepeatingAlertSentAt = nil
	}

	err = database.DB.Model(&entity).Save(entity).Error

	if err != nil {
		return nil, err
	}

	database.DB.Model(&entity).Association("Members").Replace(entity.Members)
	database.DB.Model(&entity).Association("Hosts").Replace(entity.Hosts)
	database.DB.Model(&entity).Association("EventTags").Replace(entity.EventTags)

	updatedEntity, err := r.BaseRepository.GetById(entity.Id)

	if err != nil {
		return nil, err
	}

	return updatedEntity, nil
}

// GetById получает отзыв по ID с информацией о услуге
func (r *EventRepository) GetById(id int64) (*models.Event, error) {
	var event models.Event
	if err := database.DB.Preload("Hosts").Preload("Members").Preload("EventTags").First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) AddMember(eventId int, memberId int) (*models.Event, error) {
	entity, err := r.GetById(int64(eventId))
	if err != nil {
		return nil, err
	}

	err = database.DB.Model(&entity).Association("Members").Append(&models.Member{Id: int64(memberId)})

	if err != nil {
		return nil, err
	}

	entity, err = r.GetById(int64(eventId))
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *EventRepository) RemoveMember(eventId int, memberId int) (*models.Event, error) {
	entity, err := r.GetById(int64(eventId))
	if err != nil {
		return nil, err
	}

	err = database.DB.Model(&entity).Association("Members").Delete(&models.Member{Id: int64(memberId)})

	if err != nil {
		return nil, err
	}

	entity, err = r.GetById(int64(eventId))
	if err != nil {
		return nil, err
	}

	return entity, nil
}

