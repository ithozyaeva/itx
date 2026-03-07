package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

type CasinoRepository struct{}

func NewCasinoRepository() *CasinoRepository {
	return &CasinoRepository{}
}

func (r *CasinoRepository) PlaceBet(memberId int64, bet *models.CasinoBet, won bool) (int, error) {
	var balance int

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Check balance
		if err := tx.Raw(
			`SELECT COALESCE(SUM(amount), 0) FROM point_transactions WHERE member_id = ?`,
			memberId,
		).Scan(&balance).Error; err != nil {
			return err
		}

		if balance < bet.BetAmount {
			return gorm.ErrInvalidData
		}

		// Debit bet amount
		debitTx := &models.PointTransaction{
			MemberId:    memberId,
			Amount:      -bet.BetAmount,
			Reason:      models.PointReasonCasinoBet,
			SourceType:  "casino",
			SourceId:    0,
			Description: "Ставка: " + bet.Game,
		}
		if err := tx.Create(debitTx).Error; err != nil {
			return err
		}

		// Credit winnings if won
		if won && bet.Payout > 0 {
			creditTx := &models.PointTransaction{
				MemberId:    memberId,
				Amount:      bet.Payout,
				Reason:      models.PointReasonCasinoWin,
				SourceType:  "casino",
				SourceId:    0,
				Description: "Выигрыш: " + bet.Game,
			}
			if err := tx.Create(creditTx).Error; err != nil {
				return err
			}
		}

		// Create casino bet record
		if err := tx.Create(bet).Error; err != nil {
			return err
		}

		// Recalculate balance
		if err := tx.Raw(
			`SELECT COALESCE(SUM(amount), 0) FROM point_transactions WHERE member_id = ?`,
			memberId,
		).Scan(&balance).Error; err != nil {
			return err
		}

		return nil
	})

	return balance, err
}

func (r *CasinoRepository) GetHistory(memberId int64, limit, offset int) ([]models.CasinoHistoryItem, int64, error) {
	items := make([]models.CasinoHistoryItem, 0)
	var total int64

	if err := database.DB.Model(&models.CasinoBet{}).Where("member_id = ?", memberId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Raw(`
		SELECT id, game, bet_amount, bet_choice, result, multiplier, payout, profit, created_at
		FROM casino_bets
		WHERE member_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, memberId, limit, offset).Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *CasinoRepository) GetStats(memberId int64) (*models.CasinoStats, error) {
	stats := &models.CasinoStats{}
	err := database.DB.Raw(`
		SELECT
			COUNT(*) as total_bets,
			COALESCE(SUM(bet_amount), 0) as total_wagered,
			COALESCE(SUM(payout), 0) as total_payout,
			COALESCE(MAX(payout), 0) as biggest_win
		FROM casino_bets
		WHERE member_id = ?
	`, memberId).Scan(stats).Error
	return stats, err
}

func (r *CasinoRepository) GetAdminStats() (*models.CasinoAdminStats, error) {
	stats := &models.CasinoAdminStats{}
	err := database.DB.Raw(`
		SELECT
			COUNT(*) as total_bets,
			COALESCE(SUM(bet_amount), 0) as total_wagered,
			COALESCE(SUM(payout), 0) as total_payout,
			COALESCE(SUM(bet_amount), 0) - COALESCE(SUM(payout), 0) as house_profit,
			COUNT(DISTINCT member_id) as unique_players
		FROM casino_bets
	`).Scan(stats).Error
	return stats, err
}

func (r *CasinoRepository) SearchBets(username *string, game *string, limit, offset int) ([]models.CasinoAdminBet, int64, error) {
	items := make([]models.CasinoAdminBet, 0)
	var total int64

	countQuery := database.DB.Table("casino_bets cb").Joins("JOIN members m ON m.id = cb.member_id")
	if username != nil {
		countQuery = countQuery.Where("m.username ILIKE ?", "%"+*username+"%")
	}
	if game != nil {
		countQuery = countQuery.Where("cb.game = ?", *game)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	baseQuery := `SELECT cb.id, cb.member_id, m.first_name as member_first_name, m.last_name as member_last_name,
		m.username as member_username, cb.game, cb.bet_amount, cb.bet_choice, cb.result, cb.multiplier, cb.payout, cb.profit, cb.created_at
		FROM casino_bets cb
		JOIN members m ON m.id = cb.member_id`

	var args []interface{}
	var where []string
	if username != nil {
		where = append(where, "m.username ILIKE ?")
		args = append(args, "%"+*username+"%")
	}
	if game != nil {
		where = append(where, "cb.game = ?")
		args = append(args, *game)
	}
	if len(where) > 0 {
		baseQuery += " WHERE "
		for i, w := range where {
			if i > 0 {
				baseQuery += " AND "
			}
			baseQuery += w
		}
	}
	baseQuery += ` ORDER BY cb.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	if err := database.DB.Raw(baseQuery, args...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
