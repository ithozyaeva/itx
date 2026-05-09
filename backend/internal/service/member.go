package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"gorm.io/gorm"
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
	creditsRepo      *repository.ReferralCreditRepository
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
		creditsRepo:      repository.NewReferralCreditRepository(),
		creditsSvc:       NewReferralCreditService(),
	}
}

// EnsureReferralCode гарантирует наличие персонального реф-кода у юзера.
// Если кода ещё нет — генерирует и сохраняет атомарно через AssignReferralCode.
// Идемпотентно: повторный вызов возвращает существующий код.
func (s *MemberService) EnsureReferralCode(member *models.Member) (string, error) {
	if member == nil || member.Id == 0 {
		return "", errors.New("member required")
	}
	if member.ReferralCode != nil && *member.ReferralCode != "" {
		return *member.ReferralCode, nil
	}
	code, _, err := s.repo.AssignReferralCode(member.Id)
	if err != nil {
		return "", err
	}
	member.ReferralCode = &code
	return code, nil
}

// BackfillReferralCodes — startup-job: проходит batch'ами по members без кода
// и генерирует им. Без блокировки старта приложения: запускать в горутине.
//
// Безопасно при повторном запуске и при concurrent run на нескольких инстансах
// (AssignReferralCode идемпотентен через WHERE referral_code IS NULL).
//
// limit — размер батча (default 200), max — max членов за один проход
// (default 100k = практический no-limit). Возвращает количество обработанных.
func (s *MemberService) BackfillReferralCodes(limit, max int) (int, error) {
	if limit <= 0 {
		limit = 200
	}
	if max <= 0 {
		max = 100000
	}
	processed := 0
	for processed < max {
		ids, err := s.repo.MembersWithoutReferralCode(limit)
		if err != nil {
			return processed, err
		}
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			if _, _, err := s.repo.AssignReferralCode(id); err != nil {
				log.Printf("BackfillReferralCodes: AssignReferralCode failed for member %d: %v", id, err)
				continue
			}
			processed++
		}
		// Если последний батч короче limit — больше юзеров без кода нет.
		if len(ids) < limit {
			break
		}
	}
	return processed, nil
}

// ApplyPendingReferral — Боб только что авторизовался; если у него
// в Redis-pending есть referrer_member_id (от бот-deeplink ref_<code>),
// фиксируем атрибуцию + дёргаем конверсию + начисляем credits инвайтеру.
//
// Многоуровневая идемпотентность:
//   - GetAndDelete атомарно (Redis), повторный вызов получит 0
//   - SetReferredByMemberID с WHERE IS NULL (БД), параллельный auth не перетрёт
//   - AwardForConversion идемпотентен по уникальному индексу (member, reason, source)
//
// Anti-self-fraud: бот уже отрезал self-referral, но повторно проверяем по
// member.Id == referrerMemberID (defense-in-depth, цепочка может быть отравлена
// прямой манипуляцией Redis).
//
// Порядок «эффекты → флаг»: TrackConversion + Award идут ДО SetReferredByMemberID.
// Если что-то падает посередине, флаг IS NULL остаётся, при следующем auth
// повторяем (idempotent). Иначе крэш на Award оставлял бы атрибуцию без credits.
func (s *MemberService) ApplyPendingReferral(ctx context.Context, member *models.Member, pending *PendingReferralService) {
	if member == nil || member.TelegramID == 0 || pending == nil {
		return
	}
	if member.ReferredByMemberID != nil {
		return // already attributed
	}
	referrerID, err := pending.GetAndDelete(ctx, member.TelegramID)
	if err != nil {
		log.Printf("ApplyPendingReferral: redis read failed (member=%d, tg=%d): %v", member.Id, member.TelegramID, err)
		return
	}
	if referrerID == 0 {
		return
	}
	if referrerID == member.Id {
		log.Printf("ApplyPendingReferral: ignoring self-referral member=%d", member.Id)
		return
	}
	// Подтверждаем что инвайтер существует. Если был удалён, атрибуцию не делаем.
	referrer, err := s.repo.GetById(referrerID)
	if err != nil || referrer == nil {
		log.Printf("ApplyPendingReferral: referrer %d not found", referrerID)
		return
	}

	// Конверсия в community-trекe фиксируется через credits-награду напрямую.
	// referral_conversions у нас остаётся для вакансий, не смешиваем.
	// Идемпотентность: уникальный индекс на (member, reason, source_type, source_id),
	// source_id=member.Id (Боб), уйдёт ровно раз за пару (Алиса, Боб).
	s.creditsSvc.AwardForCommunityReferral(referrerID, member.Id)

	written, err := s.repo.SetReferredByMemberID(member.Id, referrerID)
	if err != nil {
		log.Printf("ApplyPendingReferral: SetReferredByMemberID failed (member=%d, referrer=%d): %v", member.Id, referrerID, err)
		return
	}
	if !written {
		// Параллельный auth уже зафиксировал — Award был no-op'ом по индексу,
		// состояние корректно.
		return
	}
	member.ReferredByMemberID = &referrerID
	log.Printf("ApplyPendingReferral: attributed member=%d → referrer=%d", member.Id, referrerID)
}

// ReferrerInfo — payload для GET /platform/members/me/referrer (welcome-баннер).
// Узкий DTO с публичными полями: не отдаём models.Member чтобы не утекли
// telegramID/bio/birthday/roles (см. PR #342, тот же класс PII-утечек).
type ReferrerInfo struct {
	Author ReferrerAuthor `json:"author"`
	SeenAt *time.Time     `json:"seenAt"`
}

// ReferrerAuthor — публичные поля инвайтера для welcome-баннера.
type ReferrerAuthor struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Tg        string `json:"tg"`
	AvatarURL string `json:"avatarUrl"`
}

// GetReferrer — данные инвайтера для welcome-баннера. Возвращает nil
// если у юзера нет инвайтера или баннер уже видели (фронт не запрашивает
// в этом случае автора, экономим БД).
func (s *MemberService) GetReferrer(member *models.Member) (*ReferrerInfo, error) {
	if member == nil || member.ReferredByMemberID == nil {
		return nil, nil
	}
	if member.ReferralWelcomeSeenAt != nil {
		return nil, nil
	}
	referrer, err := s.repo.GetById(*member.ReferredByMemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // инвайтер удалён — баннер показывать нечего
		}
		return nil, err
	}
	return &ReferrerInfo{
		Author: ReferrerAuthor{
			Id:        referrer.Id,
			FirstName: referrer.FirstName,
			LastName:  referrer.LastName,
			Tg:        referrer.Username,
			AvatarURL: referrer.AvatarURL,
		},
		SeenAt: member.ReferralWelcomeSeenAt,
	}, nil
}

// MarkReferralWelcomeSeen — фронт отметил, что юзер закрыл welcome-баннер.
// First-write-wins на уровне БД (см. SetReferralWelcomeSeenAt).
func (s *MemberService) MarkReferralWelcomeSeen(memberID int64) error {
	return s.repo.SetReferralWelcomeSeenAt(memberID)
}

// ReferralCabinet — комплексный payload для страницы /referral.
type ReferralCabinet struct {
	Code           string                       `json:"code"`
	Deeplink       string                       `json:"deeplink"`
	Balance        int                          `json:"balance"`
	InvitedTotal   int64                        `json:"invitedTotal"`
	WithActiveSub  int64                        `json:"withActiveSub"`
	TotalEarned    int                          `json:"totalEarned"`
	RecentInvitees []repository.InviteeRow      `json:"recentInvitees"`
}

// GetReferralCabinet — собирает данные для страницы рефкабинета: персональный
// код+deeplink, баланс кредитов, статистика, список последних приглашённых.
//
// Если у юзера ещё нет кода — генерирует на лету (lazy-create при первом
// открытии). После первого визита код уже есть и dедик-запрос пропускается.
func (s *MemberService) GetReferralCabinet(member *models.Member, botUsername string) (*ReferralCabinet, error) {
	code, err := s.EnsureReferralCode(member)
	if err != nil {
		return nil, err
	}
	stats, err := s.repo.GetReferralStats(member.Id)
	if err != nil {
		return nil, err
	}
	balance, err := s.creditsRepo.GetBalance(member.Id)
	if err != nil {
		return nil, err
	}
	totalEarned, err := s.creditsRepo.GetTotalEarned(member.Id)
	if err != nil {
		return nil, err
	}
	invitees, err := s.repo.GetInvitees(member.Id, 20)
	if err != nil {
		return nil, err
	}
	deeplink := ""
	if botUsername != "" {
		deeplink = "https://t.me/" + botUsername + "?start=ref_" + code
	}
	return &ReferralCabinet{
		Code:           code,
		Deeplink:       deeplink,
		Balance:        balance,
		InvitedTotal:   stats.InvitedTotal,
		WithActiveSub:  stats.WithActiveSub,
		TotalEarned:    totalEarned,
		RecentInvitees: invitees,
	}, nil
}

// FindReferrerForMember — единая точка поиска инвайтера для credits-наград
// (first/recurring purchase). Приоритет:
//  1. members.referred_by_member_id (community-программа, PR #350)
//  2. Latest referral_conversions JOIN referal_links (legacy, по вакансии)
//
// Возвращает 0 если ни одного источника не нашлось. Используется
// awardReferralRewardsFor в subscription_service.
func (s *MemberService) FindReferrerForMember(memberID int64) int64 {
	m, err := s.repo.GetById(memberID)
	if err == nil && m != nil && m.ReferredByMemberID != nil {
		return *m.ReferredByMemberID
	}
	id, err := s.referalRepo.GetReferrerForMember(memberID)
	if err == nil && id > 0 {
		return id
	}
	return 0
}

// HelperFunctions — вспомогательные re-exports для других пакетов.
// strings.Contains для тестов на duplicate-key error в TrackConversion.
var _ = strings.Contains

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

// GetByReferralCode — lookup для бот-deeplink-обработки.
func (s *MemberService) GetByReferralCode(code string) (*models.Member, error) {
	return s.repo.GetByReferralCode(code)
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
