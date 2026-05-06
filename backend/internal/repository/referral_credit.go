package repository

import (
	"errors"
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

// ErrInsufficientCredits — баланс пользователя меньше суммы списания.
// Хендлер маппит её в HTTP 402 Payment Required, чтобы UI мог отличить
// «не хватило кредитов» от прочих ошибок.
var ErrInsufficientCredits = errors.New("insufficient referral credits")

type ReferralCreditRepository struct{}

func NewReferralCreditRepository() *ReferralCreditRepository {
	return &ReferralCreditRepository{}
}

func (r *ReferralCreditRepository) GetBalance(memberId int64) (int, error) {
	var balance int
	err := database.DB.Raw(
		`SELECT COALESCE(SUM(amount), 0) FROM referral_credit_transactions WHERE member_id = ?`,
		memberId,
	).Scan(&balance).Error
	return balance, err
}

func (r *ReferralCreditRepository) GetTransactions(memberId int64, limit int) ([]models.ReferralCreditTransaction, error) {
	var transactions []models.ReferralCreditTransaction
	err := database.DB.
		Where("member_id = ?", memberId).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

// Award — обычный INSERT без идемпотентности. Используется для admin_manual
// и subscription_purchase (списание), где каждый вызов = одна уникальная
// транзакция, повторений быть не должно.
func (r *ReferralCreditRepository) Award(tx *models.ReferralCreditTransaction) error {
	return database.DB.Create(tx).Error
}

// AwardIdempotent — INSERT с защитой от дублирования (повторный вызов с
// той же тройкой member/reason/source = no-op). Реализован через NOT EXISTS
// (как в points), чтобы не зависеть от частичного уникального индекса
// в синтаксисе ON CONFLICT.
func (r *ReferralCreditRepository) AwardIdempotent(tx *models.ReferralCreditTransaction) error {
	return database.DB.Exec(
		`INSERT INTO referral_credit_transactions (member_id, amount, reason, source_type, source_id, description)
		 SELECT ?, ?, ?, ?, ?, ?
		 WHERE NOT EXISTS (
		     SELECT 1 FROM referral_credit_transactions
		     WHERE member_id = ? AND reason = ? AND source_type = ? AND source_id = ?
		 )`,
		tx.MemberId, tx.Amount, tx.Reason, tx.SourceType, tx.SourceId, tx.Description,
		tx.MemberId, tx.Reason, tx.SourceType, tx.SourceId,
	).Error
}

// Spend атомарно списывает credits внутри переданной транзакции и
// возвращает ID созданной transaction-записи (для idempotency-key
// последующих наград, например AwardForReferralPurchase).
//
// Использует SELECT … FOR UPDATE по строкам member_id, чтобы между
// проверкой баланса и INSERT'ом отрицательного списания не вклинилась
// параллельная покупка с тем же балансом.
//
// Возвращает ErrInsufficientCredits, если баланс < amount, без записи
// транзакции — вся внешняя БД-транзакция должна откатиться.
func (r *ReferralCreditRepository) Spend(
	db *gorm.DB,
	memberId int64,
	amount int,
	reason models.ReferralCreditReason,
	sourceType string,
	sourceId int64,
	description string,
) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("spend amount must be positive")
	}

	var balance int
	if err := db.Raw(
		`SELECT COALESCE(SUM(amount), 0) FROM referral_credit_transactions
		 WHERE member_id = ? FOR UPDATE`,
		memberId,
	).Scan(&balance).Error; err != nil {
		return 0, err
	}
	if balance < amount {
		return 0, ErrInsufficientCredits
	}

	tx := &models.ReferralCreditTransaction{
		MemberId:    memberId,
		Amount:      -amount,
		Reason:      reason,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
	}
	if err := db.Create(tx).Error; err != nil {
		return 0, err
	}
	return tx.Id, nil
}

// AwardTx — INSERT внутри переданной транзакции (без идемпотентности).
// Используется в PurchaseTierWithCredits, чтобы списание и любые «бонус»-
// начисления попали в одну атомарную единицу с выдачей подписки.
func (r *ReferralCreditRepository) AwardTx(db *gorm.DB, tx *models.ReferralCreditTransaction) error {
	return db.Create(tx).Error
}

func (r *ReferralCreditRepository) SearchTransactions(username *string, limit, offset int) ([]models.AdminCreditTransaction, int64, error) {
	items := make([]models.AdminCreditTransaction, 0)
	var total int64

	countQuery := database.DB.Table("referral_credit_transactions rct").
		Joins("JOIN members m ON m.id = rct.member_id")
	if username != nil {
		countQuery = countQuery.Where("m.username ILIKE ?", "%"+*username+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	baseQuery := `SELECT rct.id, rct.member_id, m.first_name as member_first_name, m.last_name as member_last_name,
	        m.username as member_username, rct.amount, rct.reason, rct.source_type, rct.description, rct.created_at
	 FROM referral_credit_transactions rct
	 JOIN members m ON m.id = rct.member_id`

	var args []interface{}
	if username != nil {
		baseQuery += ` WHERE m.username ILIKE ?`
		args = append(args, "%"+*username+"%")
	}
	baseQuery += ` ORDER BY rct.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	if err := database.DB.Raw(baseQuery, args...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
