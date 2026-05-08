package service

import (
	"context"
	"log"
	"strings"
	"time"

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
	referalRepo      *repository.ReferalLinkRepository
	referalSvc       *ReferalLinkService
	creditsSvc       *ReferralCreditService
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
		referalRepo:      repository.NewReferalLinkRepository(),
		referalSvc:       NewReferalLinkService(),
		creditsSvc:       NewReferralCreditService(),
	}
}

// ApplyPendingReferral вычитывает Redis-pending-attribution для члена
// (positioned заранее ботом /start ref_<id>) и фиксирует её в БД +
// дёргает конверсию + начисляет credits автору ссылки.
//
// Идемпотентно по нескольким уровням:
//   - SetReferredByLinkID игнорится, если у юзера уже стоит referred_by_link_id
//   - TrackConversion защищён уникальным индексом (link_id, member_id)
//   - AwardForConversion идемпотентен по (member, reason, source_type, source_id)
//
// Anti-self-fraud: если автор ссылки совпадает с текущим членом — игнорим.
//
// Вызывать ОДИН РАЗ при создании или re-login существующего члена в
// auth-handler'е. Ошибки логируются, но не пробрасываются: атрибуция —
// best-effort, не должна блокировать auth flow.
func (s *MemberService) ApplyPendingReferral(ctx context.Context, member *models.Member, pending *PendingReferralService) {
	if member == nil || member.TelegramID == 0 || pending == nil {
		return
	}
	if member.ReferredByLinkID != nil {
		return // already attributed
	}
	linkID, err := pending.GetAndDelete(ctx, member.TelegramID)
	if err != nil {
		log.Printf("ApplyPendingReferral: redis read failed (member=%d, tg=%d): %v", member.Id, member.TelegramID, err)
		return
	}
	if linkID == 0 {
		return
	}
	link, err := s.referalRepo.GetById(linkID)
	if err != nil {
		log.Printf("ApplyPendingReferral: link %d not found for member %d", linkID, member.Id)
		return
	}
	if link.AuthorId == member.Id {
		log.Printf("ApplyPendingReferral: ignoring self-referral member=%d link=%d", member.Id, linkID)
		return
	}
	written, err := s.repo.SetReferredByLinkID(member.Id, linkID)
	if err != nil {
		log.Printf("ApplyPendingReferral: SetReferredByLinkID failed (member=%d, link=%d): %v", member.Id, linkID, err)
		return
	}
	if !written {
		return // race: другой запрос успел зафиксировать первым
	}
	member.ReferredByLinkID = &linkID

	// Конверсия + награда. TrackConversion может вернуть duplicate-key —
	// это OK (старая конверсия от той же пары), просто логируем и идём
	// дальше: AwardForConversion идемпотентен.
	if err := s.referalSvc.TrackConversion(linkID, member.Id); err != nil {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "duplicate key") && !strings.Contains(errMsg, "referral_conversions_unique") {
			log.Printf("ApplyPendingReferral: TrackConversion failed (link=%d, member=%d): %v", linkID, member.Id, err)
		}
	}
	s.creditsSvc.AwardForConversion(link.AuthorId, linkID)
	log.Printf("ApplyPendingReferral: attributed member=%d → link=%d (author=%d)", member.Id, linkID, link.AuthorId)
}

// GetReferrer возвращает данные реферрера юзера для welcome-баннера на
// фронте. Вернёт nil, если у юзера нет referred_by_link_id или ссылка/автор
// удалены. Поле SeenAt — для подсказки фронту, показывать ли модалку.
func (s *MemberService) GetReferrer(member *models.Member) (*ReferrerInfo, error) {
	if member == nil || member.ReferredByLinkID == nil {
		return nil, nil
	}
	link, err := s.referalRepo.GetById(*member.ReferredByLinkID)
	if err != nil {
		// Ссылка удалена — атрибуция остаётся, но баннер показать нечего.
		return nil, nil
	}
	return &ReferrerInfo{
		Author: link.Author,
		SeenAt: member.ReferralWelcomeSeenAt,
	}, nil
}

// MarkReferralWelcomeSeen — фронт отметил, что юзер закрыл welcome-баннер.
// После этого GetReferrer всё ещё возвращает данные, но фронт по полю SeenAt
// решает не показывать модалку повторно.
func (s *MemberService) MarkReferralWelcomeSeen(memberID int64) error {
	return s.repo.SetReferralWelcomeSeenAt(memberID)
}

// ReferrerInfo — payload для GET /platform/members/me/referrer.
type ReferrerInfo struct {
	Author models.Member `json:"author"`
	SeenAt *time.Time    `json:"seenAt"`
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

// ClaimUsername освобождает username у других участников, чтобы переданный
// keepID мог им владеть без коллизии с UNIQUE-индексом. См. ReleaseUsername.
func (s *MemberService) ClaimUsername(username string, keepID int64) error {
	return s.repo.ReleaseUsername(username, keepID)
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
