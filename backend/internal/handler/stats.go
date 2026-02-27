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
