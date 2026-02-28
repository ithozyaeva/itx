package handler

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"github.com/gofiber/fiber/v2"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

type DashboardStats struct {
	TotalMembers   int64 `json:"totalMembers"`
	TotalMentors   int64 `json:"totalMentors"`
	UpcomingEvents int64 `json:"upcomingEvents"`
	PastEvents     int64 `json:"pastEvents"`
	PendingReviews int64 `json:"pendingReviews"`
	ApprovedReviews int64 `json:"approvedReviews"`
	ReferralLinks  int64 `json:"referralLinks"`
	Resumes        int64 `json:"resumes"`
}

func (h *StatsHandler) GetStats(c *fiber.Ctx) error {
	var stats DashboardStats

	database.DB.Model(&models.Member{}).Count(&stats.TotalMembers)
	database.DB.Model(&models.MentorDbShortModel{}).Count(&stats.TotalMentors)
	database.DB.Model(&models.Event{}).Where("date >= CURRENT_TIMESTAMP").Count(&stats.UpcomingEvents)
	database.DB.Model(&models.Event{}).Where("date < CURRENT_TIMESTAMP").Count(&stats.PastEvents)
	database.DB.Model(&models.ReviewOnCommunity{}).Where("status = ?", "PENDING").Count(&stats.PendingReviews)
	database.DB.Model(&models.ReviewOnCommunity{}).Where("status = ?", "APPROVED").Count(&stats.ApprovedReviews)
	database.DB.Model(&models.ReferalLink{}).Count(&stats.ReferralLinks)
	database.DB.Model(&models.Resume{}).Count(&stats.Resumes)

	return c.JSON(stats)
}

type MonthlyStats struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type ChartStats struct {
	MemberGrowth    []MonthlyStats `json:"memberGrowth"`
	EventAttendance []MonthlyStats `json:"eventAttendance"`
}

func (h *StatsHandler) GetChartStats(c *fiber.Ctx) error {
	var chartStats ChartStats

	// Event attendance by month (last 12 months)
	if err := database.DB.Raw(`
		SELECT TO_CHAR(date_trunc('month', e.date), 'YYYY-MM') as month,
		       COUNT(DISTINCT em.member_id) as count
		FROM events e
		JOIN event_members em ON em.event_id = e.id
		WHERE e.date >= NOW() - INTERVAL '12 months'
		GROUP BY date_trunc('month', e.date)
		ORDER BY date_trunc('month', e.date)
	`).Scan(&chartStats.EventAttendance).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Total members count (cumulative snapshot is not possible without created_at,
	// so we return the total count per month as a flat line for now)
	var totalMembers int64
	if err := database.DB.Model(&models.Member{}).Count(&totalMembers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate last 12 months with total count
	if err := database.DB.Raw(`
		SELECT TO_CHAR(d.month, 'YYYY-MM') as month, CAST(? AS bigint) as count
		FROM generate_series(
		  date_trunc('month', NOW() - INTERVAL '11 months'),
		  date_trunc('month', NOW()),
		  '1 month'::interval
		) d(month)
		ORDER BY d.month
	`, totalMembers).Scan(&chartStats.MemberGrowth).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(chartStats)
}
