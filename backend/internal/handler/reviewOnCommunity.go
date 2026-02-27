package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

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
	review := new(models.AddReviewOnCommunityRequest)
	if err := c.BodyParser(review); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	author := c.Locals("member").(*models.Member)

	err := h.svc.CreateReviewOnCommunity(&models.CreateReviewOnCommunityRequest{
		AuthorTg: author.Username,
		Text:     review.Text,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *ReviewOnCommunityHandler) CreateReview(c *fiber.Ctx) error {
	review := new(models.CreateReviewOnCommunityRequest)
	if err := c.BodyParser(review); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	err := h.svc.CreateReviewOnCommunity(review)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)

}

func (h *ReviewOnCommunityHandler) GetApproved(c *fiber.Ctx) error {
	result, err := h.svc.GetApproved()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *ReviewOnCommunityHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.Approve(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
