package handler

import (
	"log"
	"strconv"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type FeedbackHandler struct {
	svc *service.FeedbackService
}

func NewFeedbackHandler() *FeedbackHandler {
	return &FeedbackHandler{
		svc: service.NewFeedbackService(),
	}
}

func (h *FeedbackHandler) Create(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	req := new(models.CreateFeedbackRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	feedback, err := h.svc.Create(member, *req)
	if err != nil {
		log.Printf("create feedback error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(feedback)
}

func (h *FeedbackHandler) AdminList(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	items, total, err := h.svc.List(limit, offset)
	if err != nil {
		log.Printf("admin list feedback error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки отзывов"})
	}

	return c.JSON(fiber.Map{
		"items": items,
		"total": total,
	})
}
