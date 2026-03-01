package service

import (
	"errors"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"time"
)

var ErrParticipantLimitReached = errors.New("достигнут лимит участников")

type EventsService struct {
	BaseService[models.Event]
	repo repository.EventRepository
}

func NewEventsService() *EventsService {
	repo := repository.NewEventRepository()
	return &EventsService{
		BaseService: NewBaseService(repo),
		repo:        *repo,
	}
}

func (s *EventsService) AddMember(eventId int, memberId int) (*models.Event, error) {
	event, err := s.repo.GetById(int64(eventId))
	if err != nil {
		return nil, err
	}

	if event.MaxParticipants > 0 && len(event.Members) >= event.MaxParticipants {
		return nil, ErrParticipantLimitReached
	}

	return s.repo.AddMember(eventId, memberId)
}

func (s *EventsService) RemoveMember(eventId int, memberId int) (*models.Event, error) {
	return s.repo.RemoveMember(eventId, memberId)
}

func (s *EventsService) GetUpcomingEvents(limit int) ([]models.Event, error) {
	return s.repo.GetUpcoming(limit)
}

func (s *EventsService) GetFutureEvents(now time.Time) ([]models.Event, error) {
	allEvents, _, err := s.repo.Search(nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var futureEvents []models.Event
	for _, event := range allEvents {
		if event.IsRepeating && event.RepeatPeriod != nil {
			if event.RepeatEndDate != nil && now.After(*event.RepeatEndDate) {
				continue
			}
			futureEvents = append(futureEvents, event)
		} else {
			if event.Date.After(now) || event.Date.Equal(now) {
				futureEvents = append(futureEvents, event)
			}
		}
	}

	return futureEvents, nil
}
