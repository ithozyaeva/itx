package service

import (
	"errors"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"strings"
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

// SearchUpcoming делегирует в репозиторий поиск предстоящих событий,
// отсортированных по ближайшему будущему вхождению (а не по исходному date).
// Возвращает тот же формат, что и Search.
func (s *EventsService) SearchUpcoming(limit *int, offset *int, filter *repository.SearchFilter) (*models.RegistrySearch[models.Event], error) {
	items, total, err := s.repo.SearchUpcoming(limit, offset, filter)
	if err != nil {
		return nil, err
	}
	return &models.RegistrySearch[models.Event]{Items: items, Total: int(total)}, nil
}

// ResolveEventTags нормализует список тегов события: теги с Id > 0 оставляются как есть,
// а теги без Id (приходят с фронта, когда админ вводит новое имя вручную) ищутся по имени
// или создаются. Возвращает список тегов с корректными Id, пригодный для GORM many2many.
func (s *EventsService) ResolveEventTags(tags []models.EventTag) ([]models.EventTag, error) {
	if len(tags) == 0 {
		return tags, nil
	}
	resolved := make([]models.EventTag, 0, len(tags))
	seenIDs := make(map[int64]bool)
	seenNames := make(map[string]bool)
	for _, tag := range tags {
		if tag.Id > 0 {
			if seenIDs[tag.Id] {
				continue
			}
			seenIDs[tag.Id] = true
			resolved = append(resolved, tag)
			continue
		}
		name := strings.TrimSpace(tag.Name)
		if name == "" || seenNames[name] {
			continue
		}
		seenNames[name] = true
		var existing models.EventTag
		if err := database.DB.Where("name = ?", name).FirstOrCreate(&existing, models.EventTag{Name: name}).Error; err != nil {
			return nil, err
		}
		if seenIDs[existing.Id] {
			continue
		}
		seenIDs[existing.Id] = true
		resolved = append(resolved, existing)
	}
	return resolved, nil
}

func (s *EventsService) GetFutureEvents(now time.Time) ([]models.Event, error) {
	// Отсекаем прошлое на уровне БД, а не в Go. Шедулер алертов зовёт это
	// раз в минуту, и раньше каждый вызов делал SELECT * FROM events
	// (~200–400ms на 29 строках) — отсюда постоянный SLOW SQL в логах.
	events, err := s.repo.GetFutureEvents(now)
	if err != nil {
		return nil, err
	}

	// Для повторяющихся дополнительно проверяем RepeatPeriod != NULL — БД
	// этого условия не знает (в схеме period хранится как *string).
	filtered := events[:0]
	for _, ev := range events {
		if ev.IsRepeating && ev.RepeatPeriod == nil {
			continue
		}
		filtered = append(filtered, ev)
	}
	return filtered, nil
}
