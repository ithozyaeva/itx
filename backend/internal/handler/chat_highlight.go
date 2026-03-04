package handler

import (
	"strconv"

	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ChatHighlightHandler struct {
	svc *service.ChatHighlightService
}

func NewChatHighlightHandler() *ChatHighlightHandler {
	return &ChatHighlightHandler{
		svc: service.NewChatHighlightService(),
	}
}

func (h *ChatHighlightHandler) GetRecent(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	highlights, err := h.svc.GetRecent(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(highlights)
}

func (h *ChatHighlightHandler) Search(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	highlights, total, err := h.svc.Search(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"items": highlights,
		"total": total,
	})
}
