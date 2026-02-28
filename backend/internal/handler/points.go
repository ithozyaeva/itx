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

func (h *PointsHandler) AdminSearch(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")
	username := c.Query("username")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 {
		limit = 20
	}

	var usernamePtr *string
	if username != "" {
		usernamePtr = &username
	}

	items, total, err := h.svc.SearchTransactions(usernamePtr, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить транзакции"})
	}

	return c.JSON(fiber.Map{"items": items, "total": total})
}

func (h *PointsHandler) AdminAward(c *fiber.Ctx) error {
	var req models.AdminAwardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	if req.MemberId <= 0 || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "memberId и amount обязательны и должны быть > 0"})
	}

	if err := h.svc.AdminAwardPoints(req.MemberId, req.Amount, req.Description); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось начислить баллы"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *PointsHandler) AdminDelete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.svc.DeleteTransaction(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось удалить транзакцию"})
	}

	return c.JSON(fiber.Map{"success": true})
}
