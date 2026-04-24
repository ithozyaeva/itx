package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
)

// MemberRepoAdapter адаптер для репозитория участников
type MemberRepoAdapter struct {
	*repository.MemberRepository
}

// MemberService представляет сервис для работы с участниками
type MemberService struct {
	BaseService[models.Member]
	repo             *repository.MemberRepository
	mentorRepo       *repository.MentorRepository
	subscriptionRepo *repository.SubscriptionRepository
}

// NewMemberService создает новый экземпляр сервиса участников
func NewMemberService() *MemberService {
	repo := repository.NewMemberRepository()
	adapter := &MemberRepoAdapter{repo}

	return &MemberService{
		BaseService:      NewBaseService[models.Member](adapter),
		repo:             repo,
		mentorRepo:       repository.NewMentorRepository(),
		subscriptionRepo: repository.NewSubscriptionRepository(),
	}
}

// GetEffectiveTier возвращает эффективный тир подписки пользователя
// (ManualTierID в приоритете, иначе ResolvedTierID). nil если подписки нет.
// Используется хендлерами профиля, чтобы отдать фронту актуальный уровень
// без чтения roles, которые с тирами не синхронизируются.
func (s *MemberService) GetEffectiveTier(telegramID int64) *models.SubscriptionTier {
	user, err := s.subscriptionRepo.GetUser(telegramID)
	if err != nil {
		return nil
	}
	tierID := user.EffectiveTierID()
	if tierID == nil {
		return nil
	}
	tier, err := s.subscriptionRepo.GetTier(*tierID)
	if err != nil {
		return nil
	}
	return tier
}

func (s *MemberService) GetTodayBirthdays() ([]string, error) {
	return s.repo.GetTodayBirthdays()
}

func (s *MemberService) GetMentor(memberId int64) (*models.MentorModel, error) {
	mentorDb, err := s.mentorRepo.GetByMemberID(memberId)

	if err != nil {
		return nil, err
	}

	mentor := mentorDb.ToModel()

	return &mentor, nil
}

func (s *MemberService) GetPermissions(memberId int64) ([]models.Permission, error) {
	return s.repo.GetMemberPermissions(memberId)
}

func (s *MemberService) GetAllPermissions() ([]models.Permission, error) {
	return s.repo.GetAllPermissions()
}

func (s *MemberService) GetByTelegramID(telegramID int64) (*models.Member, error) {
	return s.repo.GetByTelegramID(telegramID)
}

func (s *MemberService) GetByUsername(username string) (*models.Member, error) {
	return s.repo.GetMemberByTelegram(username)
}

func (s *MemberService) GetSubscribedMembersWithTelegram() ([]models.Member, error) {
	return s.repo.GetSubscribedMembersWithTelegram()
}

// IsAdminByTelegramID checks if a user with the given Telegram ID has the ADMIN role.
func (s *MemberService) IsAdminByTelegramID(telegramID int64) bool {
	member, err := s.repo.GetByTelegramID(telegramID)
	if err != nil {
		return false
	}
	return utils.HasRole(member.Roles, models.MemberRoleAdmin)
}

// GetAdminTelegramIDs returns Telegram IDs of all members with the ADMIN role.
func (s *MemberService) GetAdminTelegramIDs() []int64 {
	members, err := s.repo.GetMembersByRole(models.MemberRoleAdmin)
	if err != nil {
		return nil
	}
	var ids []int64
	for _, m := range members {
		if m.TelegramID != 0 {
			ids = append(ids, m.TelegramID)
		}
	}
	return ids
}
