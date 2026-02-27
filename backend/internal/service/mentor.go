package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
)

// MentorServiceInterface интерфейс для сервиса менторов
type MentorServiceInterface interface {
	BaseService[models.MentorDbShortModel]
	GetByIdFull(id int64) (*models.MentorModel, error)
	GetServices(id int64) ([]models.Service, error)
	AddReviewToService(review *models.ReviewOnService) (*models.ReviewOnService, error)
	CreateWithRelations(request *models.MentorDbModel) (*models.MentorModel, error)
	UpdateWithRelations(request *models.MentorDbModel) (*models.MentorModel, error)
	GetAllWithRelations(limit *int, offset *int) (*models.RegistrySearch[models.MentorModel], error)
}

// MentorService реализует интерфейс MentorServiceInterface
type MentorService struct {
	BaseService[models.MentorDbShortModel]
	repo       *repository.MentorRepository
	memberRepo *repository.MemberRepository
}

// NewMentorService создает новый экземпляр сервиса менторов
func NewMentorService() *MentorService {
	repo := repository.NewMentorRepository()
	return &MentorService{
		BaseService: NewBaseService[models.MentorDbShortModel](repo),
		repo:        repo,
		memberRepo:  repository.NewMemberRepository(),
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

	mentor := mentorDb.ToModel()

	return &mentor, nil
}

// AddReviewToService добавляет отзыв к услуге ментора
func (s *MentorService) AddReviewToService(review *models.ReviewOnService) (*models.ReviewOnService, error) {
	return s.repo.AddReviewToService(review)
}

// CreateWithRelations создает нового ментора со всеми связанными сущностями
func (s *MentorService) CreateWithRelations(request *models.MentorDbModel) (*models.MentorModel, error) {
	mentorDb, err := s.repo.CreateWithRelations(request)
	if err != nil {
		return nil, err
	}

	connectedMember, err := s.memberRepo.GetById(request.MemberId)

	if err != nil {
		return nil, err
	}

	connectedMember.Roles = append(connectedMember.Roles, models.MemberRoleMentor)

	_, err = s.memberRepo.Update(connectedMember)

	if err != nil {
		return nil, err
	}

	// Преобразуем в MentorModel
	mentor := mentorDb.ToModel()

	return &mentor, nil
}

// UpdateWithRelations обновляет ментора со всеми связанными сущностями
func (s *MentorService) UpdateWithRelations(request *models.MentorDbModel) (*models.MentorModel, error) {
	mentorDb, err := s.repo.UpdateWithRelations(request)
	if err != nil {
		return nil, err
	}

	// Преобразуем в MentorModel
	mentor := mentorDb.ToModel()

	return &mentor, nil
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
		mentor := mentorDb.ToModel()
		mentors = append(mentors, mentor)
	}

	return &models.RegistrySearch[models.MentorModel]{
		Items: mentors,
		Total: int(total),
	}, nil
}

func (s *MentorService) GetByMemberID(memberId int64) (*models.MentorDbModel, error) {
	return s.repo.GetByMemberID(memberId)
}

func (s *MentorService) Delete(mentor *models.MentorDbShortModel) error {
	connectedMember, err := s.memberRepo.GetById(mentor.MemberId)

	if err != nil {
		return err
	}

	connectedMember.Roles = utils.RemoveRole(connectedMember.Roles, models.MemberRoleMentor)

	_, err = s.memberRepo.Update(connectedMember)

	if err != nil {
		return err
	}

	s.repo.BaseRepository.Delete(mentor)

	return nil
}
