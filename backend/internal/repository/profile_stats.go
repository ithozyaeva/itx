package repository

import (
	"ithozyeva/database"
)

type ProfileStatsRepository struct{}

func NewProfileStatsRepository() *ProfileStatsRepository {
	return &ProfileStatsRepository{}
}

type ProfileStats struct {
	EventsAttended int            `json:"eventsAttended"`
	EventsHosted   int            `json:"eventsHosted"`
	ReviewsCount   int            `json:"reviewsCount"`
	ReferralsCount int            `json:"referralsCount"`
	KudosSent      int            `json:"kudosSent"`
	KudosReceived  int            `json:"kudosReceived"`
	TasksCreated   int            `json:"tasksCreated"`
	TasksDone      int            `json:"tasksDone"`
	PointsBalance  int            `json:"pointsBalance"`
	MemberSince    string         `json:"memberSince"`
	PointsHistory  []PointsMonth  `json:"pointsHistory"`
}

type PointsMonth struct {
	Month string `json:"month"`
	Total int    `json:"total"`
}

func (r *ProfileStatsRepository) GetStats(memberId int64) (*ProfileStats, error) {
	stats := &ProfileStats{}

	// Events attended
	database.DB.Raw(`SELECT COUNT(*) FROM event_members WHERE member_id = ?`, memberId).Scan(&stats.EventsAttended)

	// Events hosted
	database.DB.Raw(`SELECT COUNT(*) FROM event_hosts WHERE member_id = ?`, memberId).Scan(&stats.EventsHosted)

	// Reviews
	database.DB.Raw(`SELECT COUNT(*) FROM "reviewOnCommunity" WHERE "authorId" = ?`, memberId).Scan(&stats.ReviewsCount)

	// Referrals
	database.DB.Raw(`SELECT COUNT(*) FROM referal_links WHERE member_id = ?`, memberId).Scan(&stats.ReferralsCount)

	// Kudos
	database.DB.Raw(`SELECT COUNT(*) FROM kudos WHERE from_id = ?`, memberId).Scan(&stats.KudosSent)
	database.DB.Raw(`SELECT COUNT(*) FROM kudos WHERE to_id = ?`, memberId).Scan(&stats.KudosReceived)

	// Tasks
	database.DB.Raw(`SELECT COUNT(*) FROM task_exchanges WHERE creator_id = ?`, memberId).Scan(&stats.TasksCreated)
	database.DB.Raw(`
		SELECT COUNT(*) FROM task_exchange_assignees tea
		JOIN task_exchanges te ON te.id = tea.task_exchange_id
		WHERE tea.member_id = ? AND te.status = 'APPROVED'
	`, memberId).Scan(&stats.TasksDone)

	// Balance
	database.DB.Raw(`SELECT COALESCE(SUM(amount), 0) FROM point_transactions WHERE member_id = ?`, memberId).Scan(&stats.PointsBalance)

	// Member since
	database.DB.Raw(`SELECT TO_CHAR(created_at, 'YYYY-MM-DD') FROM members WHERE id = ?`, memberId).Scan(&stats.MemberSince)

	// Points history (last 6 months)
	var history []PointsMonth
	database.DB.Raw(`
		SELECT TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM') as month,
			SUM(amount) as total
		FROM point_transactions
		WHERE member_id = ? AND created_at >= NOW() - INTERVAL '6 months'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY month ASC
	`, memberId).Scan(&history)
	stats.PointsHistory = history

	return stats, nil
}
