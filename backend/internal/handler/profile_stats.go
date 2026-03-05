package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProfileStatsHandler struct {
	repo *repository.ProfileStatsRepository
}

func NewProfileStatsHandler() *ProfileStatsHandler {
	return &ProfileStatsHandler{
		repo: repository.NewProfileStatsRepository(),
	}
}

func (h *ProfileStatsHandler) GetMyStats(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	stats, err := h.repo.GetStats(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

func (h *ProfileStatsHandler) GetMemberStats(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	stats, err := h.repo.GetStats(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}
