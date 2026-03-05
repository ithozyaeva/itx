package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type KudosHandler struct {
	svc *service.KudosService
}

func NewKudosHandler() *KudosHandler {
	return &KudosHandler{
		svc: service.NewKudosService(),
	}
}

func (h *KudosHandler) Send(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	req := new(models.CreateKudosRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	kudos, err := h.svc.Send(member.Id, req.ToId, req.Message)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(kudos)
}

func (h *KudosHandler) GetRecent(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	items, total, err := h.svc.GetRecent(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"items": items,
		"total": total,
	})
}
