package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

// MentorServiceInterface интерфейс для сервиса менторов
type MentorServiceInterface interface {
	BaseService[models.MentorDbShortModel]
	GetByIdFull(id int64) (*models.MentorModel, error)
	GetServices(id int64) ([]models.Service, error)
	FindByTag(tagId int64, limit int, offset int) (*models.RegistrySearch[models.MentorModel], error)
	AddReviewToService(review *models.ReviewOnService) (*models.ReviewOnService, error)
	CreateWithRelations(request *models.MentorCreateUpdateRequest) (*models.MentorModel, error)
	UpdateWithRelations(request *models.MentorCreateUpdateRequest) (*models.MentorModel, error)
	GetAllWithRelations(limit *int, offset *int) (*models.RegistrySearch[models.MentorModel], error)
}

// MentorService реализует интерфейс MentorServiceInterface
type MentorService struct {
	BaseService[models.MentorDbShortModel]
	repo *repository.MentorRepository
}

// NewMentorService создает новый экземпляр сервиса менторов
func NewMentorService() *MentorService {
	repo := repository.NewMentorRepository()
	return &MentorService{
		BaseService: NewBaseService[models.MentorDbShortModel](repo),
		repo:        repo,
	}
}

func (s *MentorService) GetServices(id int64) ([]models.Service, error) {
	services, err := s.repo.GetServices(id)

	if err != nil {
		return nil, err
	}

	return services, nil
}

// GetByIdFull получает ментора по ID с полной информацией и преобразует в MentorModel
func (s *MentorService) GetByIdFull(id int64) (*models.MentorModel, error) {
	mentorDb, err := s.repo.GetByIdFull(id)
	if err != nil {
		return nil, err
	}

	// Преобразуем в MentorModel
	mentor := &models.MentorModel{
		Id:         mentorDb.Id,
		Tg:         mentorDb.Member.Tg,
		Name:       mentorDb.Member.Name,
		Occupation: mentorDb.Occupation,
		Experience: mentorDb.Experience,
		Order:      mentorDb.Order,
		ProfTags:   mentorDb.ProfTags,
		Contacts:   mentorDb.Contacts,
		Services:   mentorDb.Services,
		MemberId:   int(mentorDb.MemberId),
	}

	return mentor, nil
}

// FindByTag находит менторов по тегу и преобразует в MentorModel
func (s *MentorService) FindByTag(tagId int64, limit int, offset int) (*models.RegistrySearch[models.MentorModel], error) {
	mentorsDb, count, err := s.repo.FindByTag(tagId, limit, offset)
	if err != nil {
		return nil, err
	}

	// Преобразуем в []MentorModel
	mentors := make([]models.MentorModel, len(mentorsDb))
	for i, mentorDb := range mentorsDb {
		mentors[i] = models.MentorModel{
			Id:         mentorDb.Id,
			Tg:         mentorDb.Member.Tg,
			Name:       mentorDb.Member.Name,
			Occupation: mentorDb.Occupation,
			Experience: mentorDb.Experience,
			Order:      mentorDb.Order,
			ProfTags:   mentorDb.ProfTags,
			Contacts:   mentorDb.Contacts,
			Services:   mentorDb.Services,
			MemberId:   int(mentorDb.MemberId),
		}
	}

	return &models.RegistrySearch[models.MentorModel]{
		Items: mentors,
		Total: int(count),
	}, nil
}

// AddReviewToService добавляет отзыв к услуге ментора
func (s *MentorService) AddReviewToService(review *models.ReviewOnService) (*models.ReviewOnService, error) {
	return s.repo.AddReviewToService(review)
}

// CreateWithRelations создает нового ментора со всеми связанными сущностями
func (s *MentorService) CreateWithRelations(request *models.MentorCreateUpdateRequest) (*models.MentorModel, error) {
	mentorDb, err := s.repo.CreateWithRelations(request)
	if err != nil {
		return nil, err
	}

	// Преобразуем в MentorModel
	mentor := &models.MentorModel{
		Id:         mentorDb.Id,
		Tg:         mentorDb.Member.Tg,
		Name:       mentorDb.Member.Name,
		Occupation: mentorDb.Occupation,
		Experience: mentorDb.Experience,
		Order:      mentorDb.Order,
		ProfTags:   mentorDb.ProfTags,
		Contacts:   mentorDb.Contacts,
		Services:   mentorDb.Services,
		MemberId:   int(mentorDb.MemberId),
	}

	return mentor, nil
}

// UpdateWithRelations обновляет ментора со всеми связанными сущностями
func (s *MentorService) UpdateWithRelations(request *models.MentorCreateUpdateRequest) (*models.MentorModel, error) {
	mentorDb, err := s.repo.UpdateWithRelations(request)
	if err != nil {
		return nil, err
	}

	// Преобразуем в MentorModel
	mentor := &models.MentorModel{
		Id:         mentorDb.Id,
		Tg:         mentorDb.Member.Tg,
		Name:       mentorDb.Member.Name,
		Occupation: mentorDb.Occupation,
		Experience: mentorDb.Experience,
		Order:      mentorDb.Order,
		ProfTags:   mentorDb.ProfTags,
		Contacts:   mentorDb.Contacts,
		Services:   mentorDb.Services,
		MemberId:   int(mentorDb.MemberId),
	}

	return mentor, nil
}

// GetAllWithRelations получает всех менторов с полной информацией
func (s *MentorService) GetAllWithRelations(limit *int, offset *int) (*models.RegistrySearch[models.MentorModel], error) {
	mentorsDb, total, err := s.repo.GetAllWithRelations(limit, offset)
	if err != nil {
		return nil, err
	}

	// Преобразуем модели базы данных в модели представления
	var mentors []models.MentorModel
	for _, mentorDb := range mentorsDb {
		mentor := models.MentorModel{
			Id:         mentorDb.Id,
			Name:       mentorDb.Member.Name,
			Tg:         mentorDb.Member.Tg,
			Occupation: mentorDb.Occupation,
			Experience: mentorDb.Experience,
			Order:      mentorDb.Order,
			ProfTags:   mentorDb.ProfTags,
			Contacts:   mentorDb.Contacts,
			Services:   mentorDb.Services,
			MemberId:   int(mentorDb.MemberId),
		}
		mentors = append(mentors, mentor)
	}

	return &models.RegistrySearch[models.MentorModel]{
		Items: mentors,
		Total: int(total),
	}, nil
}
