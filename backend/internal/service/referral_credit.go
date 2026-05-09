package service

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"
	"strings"

	"gorm.io/gorm"
)

type ReferralCreditService struct {
	repo       *repository.ReferralCreditRepository
	memberRepo *repository.MemberRepository
	settings   *AppSettingsService
}

func NewReferralCreditService() *ReferralCreditService {
	return &ReferralCreditService{
		repo:       repository.NewReferralCreditRepository(),
		memberRepo: repository.NewMemberRepository(),
		settings:   NewAppSettingsService(),
	}
}

// formatRefereeLabel — дружелюбная подпись реферала для description в credits.
// Использует имя/username если есть, иначе "#<id>". Если member не найден —
// fallback на ID. Не критично, чисто UI-косметика.
func (s *ReferralCreditService) formatRefereeLabel(refereeId int64) string {
	m, err := s.memberRepo.GetById(refereeId)
	if err != nil || m == nil {
		return fmt.Sprintf("#%d", refereeId)
	}
	if m.Username != "" {
		return "@" + m.Username
	}
	full := strings.TrimSpace(m.FirstName + " " + m.LastName)
	if full != "" {
		return full
	}
	return fmt.Sprintf("#%d", refereeId)
}

// AwardForConversion идемпотентно начисляет credits автору ссылки за
// конверсию. Сумма из app_settings.referral_conversion_credits (default 30).
// Идемпотентность: повторная конверсия по той же ссылке — no-op (защищена
// уникальным индексом в БД).
func (s *ReferralCreditService) AwardForConversion(memberId int64, linkId int64) {
	amount := s.settings.GetInt("referral_conversion_credits", 30)
	if amount <= 0 {
		return
	}
	err := s.repo.AwardIdempotent(&models.ReferralCreditTransaction{
		MemberId:    memberId,
		Amount:      amount,
		Reason:      models.CreditReasonReferalConversion,
		SourceType:  "referal_conversion",
		SourceId:    linkId,
		Description: "Конверсия по реферальной ссылке",
	})
	if err != nil {
		log.Printf("AwardForConversion error (member=%d, link=%d): %v", memberId, linkId, err)
	}
}

// AwardForCommunityReferral — награда инвайтеру за привлечение нового юзера
// в сообщество через персональный deeplink (ref_<code> в боте). Сумма —
// app_settings.community_referral_credits, default 30.
//
// Идемпотентно по (referrerId, reason='community_referral',
// source_type='community_referral', source_id=refereeId): один раз на пару
// (Алиса, Боб) на всю историю.
//
// Self-fraud: TrackConversion handler И ApplyPendingReferral отрезают
// self-referral; на уровне award проверяем третьим слоем.
func (s *ReferralCreditService) AwardForCommunityReferral(referrerId int64, refereeId int64) {
	if referrerId == refereeId || referrerId <= 0 {
		return
	}
	amount := s.settings.GetInt("community_referral_credits", 30)
	if amount <= 0 {
		return
	}
	err := s.repo.AwardIdempotent(&models.ReferralCreditTransaction{
		MemberId:    referrerId,
		Amount:      amount,
		Reason:      models.CreditReasonCommunityReferral,
		SourceType:  "community_referral",
		SourceId:    refereeId,
		Description: fmt.Sprintf("%s пришёл по персональной ссылке", s.formatRefereeLabel(refereeId)),
	})
	if err != nil {
		log.Printf("AwardForCommunityReferral error (referrer=%d, referee=%d): %v", referrerId, refereeId, err)
	}
}

// AwardForFirstPurchase — крупная единоразовая выплата инвайтеру за то,
// что реферал впервые активировал подписку (любым способом: Boosty-anchor
// или credits). share из app_settings.referral_first_purchase_share
// (default 0.5 = 50%).
//
// Идемпотентность: уникальный индекс на (referrer, reason='referral_purchase_first',
// source_type='referral_first_paid', source_id=refereeMemberID). Повторный
// вызов из любого потока — no-op.
func (s *ReferralCreditService) AwardForFirstPurchase(referrerId int64, refereeId int64, priceCents int) {
	if referrerId == refereeId || referrerId <= 0 {
		return // defense-in-depth: TrackConversion handler уже не даёт self-conversion, но если когда-нибудь просочится — не платим
	}
	share := s.settings.GetFloat("referral_first_purchase_share", 0.5)
	amount := int(float64(priceCents) / 100.0 * share)
	if amount <= 0 {
		return
	}
	err := s.repo.AwardIdempotent(&models.ReferralCreditTransaction{
		MemberId:    referrerId,
		Amount:      amount,
		Reason:      models.CreditReasonReferralPurchaseFirst,
		SourceType:  "referral_first_paid",
		SourceId:    refereeId,
		Description: fmt.Sprintf("Реферал %s впервые оформил подписку", s.formatRefereeLabel(refereeId)),
	})
	if err != nil {
		log.Printf("AwardForFirstPurchase error (referrer=%d, referee=%d): %v", referrerId, refereeId, err)
	}
}

// AwardForRecurringPurchase — ежемесячная выплата инвайтеру за активного
// реферала. periodKey — строка YYYY-MM, идёт в source_type, чтобы
// уникальный индекс на (member_id, reason, source_type, source_id)
// гарантировал «не более одной выплаты в месяц на пару (referrer, referee)».
//
// Дёргается из PeriodicCheck для каждого активного юзера; благодаря
// идемпотентности безопасно вызывать на каждом тикере.
func (s *ReferralCreditService) AwardForRecurringPurchase(referrerId int64, refereeId int64, priceCents int, periodKey string) {
	if referrerId == refereeId || referrerId <= 0 {
		return
	}
	share := s.settings.GetFloat("referral_purchase_share", 0.2)
	amount := int(float64(priceCents) / 100.0 * share)
	if amount <= 0 {
		return
	}
	err := s.repo.AwardIdempotent(&models.ReferralCreditTransaction{
		MemberId:    referrerId,
		Amount:      amount,
		Reason:      models.CreditReasonReferralPurchaseRecurring,
		SourceType:  "ref_paid:" + periodKey,
		SourceId:    refereeId,
		Description: fmt.Sprintf("Реферал %s активен в %s", s.formatRefereeLabel(refereeId), periodKey),
	})
	if err != nil {
		log.Printf("AwardForRecurringPurchase error (referrer=%d, referee=%d, period=%s): %v",
			referrerId, refereeId, periodKey, err)
	}
}

// AdminAward — ручная выдача credits из админ-панели.
// Положительная сумма — обычный INSERT. Отрицательная — списание через
// Spend (FOR UPDATE + проверка баланса), чтобы не загнать юзера в минус
// и держать инвариант: баланс >= 0 для всех путей, кроме явных корректировок.
func (s *ReferralCreditService) AdminAward(memberId int64, amount int, description string) error {
	if amount >= 0 {
		return s.repo.Award(&models.ReferralCreditTransaction{
			MemberId:    memberId,
			Amount:      amount,
			Reason:      models.CreditReasonAdminManual,
			SourceType:  "admin",
			SourceId:    0,
			Description: description,
		})
	}
	return database.DB.Transaction(func(tx *gorm.DB) error {
		_, err := s.repo.Spend(tx, memberId, -amount,
			models.CreditReasonAdminManual, "admin", 0, description)
		return err
	})
}

func (s *ReferralCreditService) GetBalance(memberId int64) (int, error) {
	return s.repo.GetBalance(memberId)
}

func (s *ReferralCreditService) GetSummary(memberId int64) (*models.ReferralCreditSummary, error) {
	balance, err := s.repo.GetBalance(memberId)
	if err != nil {
		return nil, err
	}
	transactions, err := s.repo.GetTransactions(memberId, 50)
	if err != nil {
		return nil, err
	}
	return &models.ReferralCreditSummary{
		Balance:      balance,
		Transactions: transactions,
	}, nil
}

func (s *ReferralCreditService) SearchTransactions(username *string, limit, offset int) ([]models.AdminCreditTransaction, int64, error) {
	return s.repo.SearchTransactions(username, limit, offset)
}
