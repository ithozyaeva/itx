package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SeasonHandler struct {
	svc *service.SeasonService
}

func NewSeasonHandler() *SeasonHandler {
	return &SeasonHandler{
		svc: service.NewSeasonService(),
	}
}

func (h *SeasonHandler) GetAll(c *fiber.Ctx) error {
	seasons, err := h.svc.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(seasons)
}

func (h *SeasonHandler) GetActive(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	result, err := h.svc.GetActiveWithLeaderboard(limit)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Нет активного сезона"})
	}
	return c.JSON(result)
}

func (h *SeasonHandler) GetLeaderboard(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	result, err := h.svc.GetLeaderboard(id, limit)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Сезон не найден"})
	}
	return c.JSON(result)
}

func (h *SeasonHandler) Create(c *fiber.Ctx) error {
	season := new(models.Season)
	if err := c.BodyParser(season); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if err := h.svc.Create(season); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(season)
}

func (h *SeasonHandler) Finish(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	if err := h.svc.Finish(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
