package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

// MemberRepoAdapter адаптер для репозитория участников
type MemberRepoAdapter struct {
	*repository.MemberRepository
}

// GetById реализует метод интерфейса BaseRepository
func (a *MemberRepoAdapter) GetById(id int64) (*models.Member, error) {
	return a.BaseRepository.GetById(id)
}

// Search реализует метод интерфейса BaseRepository
func (a *MemberRepoAdapter) Search(limit *int, offset *int) ([]models.Member, int64, error) {
	memberModels, count, err := a.BaseRepository.Search(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var members []models.Member
	for _, model := range memberModels {
		members = append(members, models.Member{
			Id:   model.Id,
			Tg:   model.Tg,
			Name: model.Name,
		})
	}

	return members, int64(count), nil
}

// MemberService представляет сервис для работы с участниками
type MemberService struct {
	BaseService[models.Member]
	repo *repository.MemberRepository
}

// NewMemberService создает новый экземпляр сервиса участников
func NewMemberService() *MemberService {
	repo := repository.NewMemberRepository()
	adapter := &MemberRepoAdapter{repo}

	return &MemberService{
		BaseService: NewBaseService[models.Member](adapter),
		repo:        repo,
	}
}

// Search выполняет поиск участников с пагинацией
func (s *MemberService) Search(limit *int, offset *int) (*models.RegistrySearch[models.MemberModel], error) {
	members, count, err := s.repo.Search(limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.RegistrySearch[models.MemberModel]{
		Items: members,
		Total: count,
	}, nil
}

// GetById получает участника по ID
func (s *MemberService) GetById(id int64) (*models.MemberModel, error) {
	return s.repo.GetById(id)
}
