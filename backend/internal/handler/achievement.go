package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AchievementHandler struct {
	svc *service.AchievementService
}

func NewAchievementHandler() *AchievementHandler {
	return &AchievementHandler{
		svc: service.NewAchievementService(),
	}
}

func (h *AchievementHandler) GetMyAchievements(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	resp, err := h.svc.GetUserAchievements(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить достижения"})
	}

	return c.JSON(resp)
}

func (h *AchievementHandler) GetMemberAchievements(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	resp, err := h.svc.GetUserAchievements(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить достижения"})
	}

	return c.JSON(resp)
}
