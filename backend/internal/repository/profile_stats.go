package repository

import (
	"ithozyeva/database"
)

type ProfileStatsRepository struct{}

func NewProfileStatsRepository() *ProfileStatsRepository {
	return &ProfileStatsRepository{}
}

type ProfileStats struct {
	EventsAttended    int              `json:"eventsAttended"`
	EventsHosted      int              `json:"eventsHosted"`
	ReviewsCount      int              `json:"reviewsCount"`
	ReferralsCount    int              `json:"referralsCount"`
	KudosSent         int              `json:"kudosSent"`
	KudosReceived     int              `json:"kudosReceived"`
	TasksCreated      int              `json:"tasksCreated"`
	TasksDone         int              `json:"tasksDone"`
	PointsBalance     int              `json:"pointsBalance"`
	MemberSince       string           `json:"memberSince"`
	PointsHistory     []PointsMonth    `json:"pointsHistory"`
	PointsBySource    []PointsSource   `json:"pointsBySource"`
	AchievementsEarned int             `json:"achievementsEarned"`
	AchievementsTotal  int             `json:"achievementsTotal"`
	ActivityHistory   []ActivityDay    `json:"activityHistory"`
}

type PointsMonth struct {
	Month string `json:"month"`
	Total int    `json:"total"`
}

type PointsSource struct {
	Reason string `json:"reason"`
	Total  int    `json:"total"`
}

type ActivityDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func (r *ProfileStatsRepository) GetStats(memberId int64) (*ProfileStats, error) {
	stats := &ProfileStats{}

	// Events attended
	if err := database.DB.Raw(`SELECT COUNT(*) FROM event_members WHERE member_id = ?`, memberId).Scan(&stats.EventsAttended).Error; err != nil {
		return nil, err
	}

	// Events hosted
	if err := database.DB.Raw(`SELECT COUNT(*) FROM event_hosts WHERE member_id = ?`, memberId).Scan(&stats.EventsHosted).Error; err != nil {
		return nil, err
	}

	// Reviews
	if err := database.DB.Raw(`SELECT COUNT(*) FROM "reviewOnCommunity" WHERE "authorId" = ?`, memberId).Scan(&stats.ReviewsCount).Error; err != nil {
		return nil, err
	}

	// Referrals
	if err := database.DB.Raw(`SELECT COUNT(*) FROM referal_links WHERE author_id = ?`, memberId).Scan(&stats.ReferralsCount).Error; err != nil {
		return nil, err
	}

	// Kudos
	if err := database.DB.Raw(`SELECT COUNT(*) FROM kudos WHERE from_id = ?`, memberId).Scan(&stats.KudosSent).Error; err != nil {
		return nil, err
	}
	if err := database.DB.Raw(`SELECT COUNT(*) FROM kudos WHERE to_id = ?`, memberId).Scan(&stats.KudosReceived).Error; err != nil {
		return nil, err
	}

	// Tasks
	if err := database.DB.Raw(`SELECT COUNT(*) FROM task_exchanges WHERE creator_id = ?`, memberId).Scan(&stats.TasksCreated).Error; err != nil {
		return nil, err
	}
	if err := database.DB.Raw(`
		SELECT COUNT(*) FROM task_exchange_assignees tea
		JOIN task_exchanges te ON te.id = tea.task_id
		WHERE tea.member_id = ? AND te.status = 'APPROVED'
	`, memberId).Scan(&stats.TasksDone).Error; err != nil {
		return nil, err
	}

	// Balance
	if err := database.DB.Raw(`SELECT COALESCE(SUM(amount), 0) FROM point_transactions WHERE member_id = ?`, memberId).Scan(&stats.PointsBalance).Error; err != nil {
		return nil, err
	}

	// Member since
	if err := database.DB.Raw(`SELECT TO_CHAR(created_at, 'YYYY-MM-DD') FROM members WHERE id = ?`, memberId).Scan(&stats.MemberSince).Error; err != nil {
		return nil, err
	}

	// Points history (last 6 months)
	var history []PointsMonth
	if err := database.DB.Raw(`
		SELECT TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM') as month,
			SUM(amount) as total
		FROM point_transactions
		WHERE member_id = ? AND created_at >= NOW() - INTERVAL '6 months'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY month ASC
	`, memberId).Scan(&history).Error; err != nil {
		return nil, err
	}
	stats.PointsHistory = history

	// Points by source
	var sources []PointsSource
	if err := database.DB.Raw(`
		SELECT reason, SUM(amount) as total
		FROM point_transactions
		WHERE member_id = ?
		GROUP BY reason
		ORDER BY total DESC
	`, memberId).Scan(&sources).Error; err != nil {
		return nil, err
	}
	stats.PointsBySource = sources

	// Activity history (last 12 weeks)
	var activity []ActivityDay
	if err := database.DB.Raw(`
		SELECT TO_CHAR(DATE(created_at), 'YYYY-MM-DD') as date, COUNT(*) as count
		FROM point_transactions
		WHERE member_id = ? AND created_at >= NOW() - INTERVAL '12 weeks'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, memberId).Scan(&activity).Error; err != nil {
		return nil, err
	}
	stats.ActivityHistory = activity

	return stats, nil
}
