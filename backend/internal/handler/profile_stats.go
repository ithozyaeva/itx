package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProfileStatsHandler struct {
	repo           *repository.ProfileStatsRepository
	achievementSvc *service.AchievementService
}

func NewProfileStatsHandler() *ProfileStatsHandler {
	return &ProfileStatsHandler{
		repo:           repository.NewProfileStatsRepository(),
		achievementSvc: service.NewAchievementService(),
	}
}

func (h *ProfileStatsHandler) enrichAchievements(stats *repository.ProfileStats, memberId int64) {
	earned, total, err := h.achievementSvc.GetAchievementCounts(memberId, stats.PointsBalance)
	if err != nil {
		log.Printf("enrichAchievements: failed for member %d: %v", memberId, err)
		return
	}
	stats.AchievementsEarned = earned
	stats.AchievementsTotal = total
}

func (h *ProfileStatsHandler) GetMyStats(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	stats, err := h.repo.GetStats(member.Id)
	if err != nil {
		log.Printf("get my profile stats error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки статистики профиля"})
	}
	h.enrichAchievements(stats, member.Id)
	return c.JSON(stats)
}

func (h *ProfileStatsHandler) GetMemberStats(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	stats, err := h.repo.GetStats(id)
	if err != nil {
		log.Printf("get member profile stats error (id=%d): %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки статистики профиля"})
	}
	h.enrichAchievements(stats, id)
	return c.JSON(stats)
}
