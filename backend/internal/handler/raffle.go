package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RaffleHandler struct {
	svc *service.RaffleService
}

func NewRaffleHandler() *RaffleHandler {
	return &RaffleHandler{
		svc: service.NewRaffleService(),
	}
}

func (h *RaffleHandler) GetAll(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	items, err := h.svc.GetAll(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *RaffleHandler) GetAllAdmin(c *fiber.Ctx) error {
	items, err := h.svc.GetAllAdmin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *RaffleHandler) BuyTickets(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	req := new(models.BuyTicketRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := h.svc.BuyTickets(id, member.Id, req.Count); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	BroadcastEvent("raffles")
	return c.JSON(fiber.Map{"ok": true})
}

func (h *RaffleHandler) Create(c *fiber.Ctx) error {
	raffle := new(models.Raffle)
	if err := c.BodyParser(raffle); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if err := h.svc.Create(raffle); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(raffle)
}

func (h *RaffleHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
