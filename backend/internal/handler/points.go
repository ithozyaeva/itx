package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PointsHandler struct {
	svc *service.PointsService
}

func NewPointsHandler() *PointsHandler {
	return &PointsHandler{
		svc: service.NewPointsService(),
	}
}

func (h *PointsHandler) GetMyPoints(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	summary, err := h.svc.GetMyPoints(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить баллы"})
	}

	return c.JSON(summary)
}

func (h *PointsHandler) GetLeaderboard(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	entries, err := h.svc.GetLeaderboard(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить рейтинг"})
	}

	return c.JSON(fiber.Map{"items": entries})
}
