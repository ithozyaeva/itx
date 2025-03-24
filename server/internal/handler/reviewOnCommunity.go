package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReviewOnCommunityHandler struct {
	BaseHandler[models.ReviewOnCommunity]
	svc *service.ReviewOnCommunityService
}

func NewReviewOnCommunityHandler() *ReviewOnCommunityHandler {
	svc := service.NewReviewOnCommunityService()
	return &ReviewOnCommunityHandler{
		BaseHandler: *NewBaseHandler(svc),
		svc:         svc,
	}
}

func (h *ReviewOnCommunityHandler) GetAllWithAuthor(c *fiber.Ctx) error {
	req := new(models.SearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.GetAllWithAuthor(req.Limit, req.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *ReviewOnCommunityHandler) AddReview(c *fiber.Ctx) error {
	review := new(models.ReviewOnCommunityRequest)
	if err := c.BodyParser(review); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.CreateReviewOnCommunity(review)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
