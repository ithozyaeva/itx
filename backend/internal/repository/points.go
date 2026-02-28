package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type PointsRepository struct{}

func NewPointsRepository() *PointsRepository {
	return &PointsRepository{}
}

func (r *PointsRepository) AwardPoints(tx *models.PointTransaction) error {
	result := database.DB.Exec(
		`INSERT INTO point_transactions (member_id, amount, reason, source_type, source_id, description)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT (member_id, reason, source_type, source_id) DO NOTHING`,
		tx.MemberId, tx.Amount, tx.Reason, tx.SourceType, tx.SourceId, tx.Description,
	)
	return result.Error
}

func (r *PointsRepository) GetBalance(memberId int64) (int, error) {
	var balance int
	err := database.DB.Raw(
		`SELECT COALESCE(SUM(amount), 0) FROM point_transactions WHERE member_id = ?`,
		memberId,
	).Scan(&balance).Error
	return balance, err
}

func (r *PointsRepository) GetTransactions(memberId int64, limit int) ([]models.PointTransaction, error) {
	var transactions []models.PointTransaction
	err := database.DB.
		Where("member_id = ?", memberId).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

func (r *PointsRepository) GetLeaderboard(limit int) ([]models.MemberPointsBalance, error) {
	var entries []models.MemberPointsBalance
	err := database.DB.Raw(
		`SELECT pt.member_id, m.first_name, m.last_name, m.username, m.avatar_url,
		        COALESCE(SUM(pt.amount), 0) as total
		 FROM point_transactions pt
		 JOIN members m ON m.id = pt.member_id
		 GROUP BY pt.member_id, m.first_name, m.last_name, m.username, m.avatar_url
		 ORDER BY total DESC
		 LIMIT ?`,
		limit,
	).Scan(&entries).Error
	return entries, err
}

func (r *PointsRepository) GivePoints(tx *models.PointTransaction) error {
	return database.DB.Create(tx).Error
}

func (r *PointsRepository) SearchTransactions(memberId *int64, limit, offset int) ([]models.AdminPointTransaction, int64, error) {
	var items []models.AdminPointTransaction
	var total int64

	countQuery := database.DB.Table("point_transactions")
	if memberId != nil {
		countQuery = countQuery.Where("member_id = ?", *memberId)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	baseQuery := `SELECT pt.id, pt.member_id, m.first_name as member_first_name, m.last_name as member_last_name,
		        m.username as member_username, pt.amount, pt.reason, pt.source_type, pt.description, pt.created_at
		 FROM point_transactions pt
		 JOIN members m ON m.id = pt.member_id`

	var args []interface{}
	if memberId != nil {
		baseQuery += ` WHERE pt.member_id = ?`
		args = append(args, *memberId)
	}
	baseQuery += ` ORDER BY pt.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	if err := database.DB.Raw(baseQuery, args...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *PointsRepository) DeleteTransaction(id int64) error {
	return database.DB.Delete(&models.PointTransaction{}, id).Error
}

func (r *PointsRepository) CreateManualTransaction(tx *models.PointTransaction) error {
	return database.DB.Create(tx).Error
}

func (r *PointsRepository) GetPastEventsForAward(daysBack int) ([]models.Event, error) {
	var events []models.Event
	err := database.DB.
		Preload("Hosts").
		Preload("Members").
		Where("date < NOW() AND date > NOW() - INTERVAL '1 day' * ?", daysBack).
		Find(&events).Error
	return events, err
}

func (r *PointsRepository) GetMembersWithEventsInWeek(year int, week int) ([]int64, error) {
	var memberIds []int64
	err := database.DB.Raw(
		`SELECT DISTINCT em.member_id
		 FROM event_members em
		 JOIN events e ON e.id = em.event_id
		 WHERE EXTRACT(ISOYEAR FROM e.date) = ? AND EXTRACT(WEEK FROM e.date) = ?
		   AND e.date < NOW()`,
		year, week,
	).Scan(&memberIds).Error
	return memberIds, err
}

func (r *PointsRepository) GetMembersWithMonthlyEvents(year int, month int, minEvents int) ([]int64, error) {
	var memberIds []int64
	err := database.DB.Raw(
		`SELECT em.member_id
		 FROM event_members em
		 JOIN events e ON e.id = em.event_id
		 WHERE EXTRACT(YEAR FROM e.date) = ? AND EXTRACT(MONTH FROM e.date) = ?
		   AND e.date < NOW()
		 GROUP BY em.member_id
		 HAVING COUNT(DISTINCT em.event_id) >= ?`,
		year, month, minEvents,
	).Scan(&memberIds).Error
	return memberIds, err
}

func (r *PointsRepository) GetMembersWithStreak(weeks int) ([]int64, error) {
	var memberIds []int64
	err := database.DB.Raw(
		`SELECT member_id FROM (
			SELECT em.member_id,
				COUNT(DISTINCT CAST(EXTRACT(ISOYEAR FROM e.date) * 100 + EXTRACT(WEEK FROM e.date) AS INTEGER)) as active_weeks
			FROM event_members em
			JOIN events e ON e.id = em.event_id
			WHERE e.date < NOW()
			  AND e.date > NOW() - INTERVAL '1 week' * ?
			GROUP BY em.member_id
		) sub
		WHERE active_weeks >= ?`,
		weeks, weeks,
	).Scan(&memberIds).Error
	return memberIds, err
}
