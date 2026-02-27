package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReviewOnServiceHandler struct {
	BaseHandler[models.ReviewOnService]
	svc *service.ReviewOnServiceService
}

func NewReviewOnServiceHandler() *ReviewOnServiceHandler {
	svc := service.NewReviewOnServiceService()
	return &ReviewOnServiceHandler{
		BaseHandler: *NewBaseHandler[models.ReviewOnService](svc),
		svc:         svc,
	}
}

// Search выполняет поиск отзывов с пагинацией
func (h *ReviewOnServiceHandler) Search(c *fiber.Ctx) error {
	req := new(models.SearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.Search(req.Limit, req.Offset, nil, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// GetReviewsWithMentorInfo получает отзывы с информацией о менторе
func (h *ReviewOnServiceHandler) GetReviewsWithMentorInfo(c *fiber.Ctx) error {
	req := new(models.SearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.GetReviewsWithMentorInfo(req.Limit, req.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// CreateReview создает новый отзыв
func (h *ReviewOnServiceHandler) CreateReview(c *fiber.Ctx) error {
	request := new(models.ReviewOnServiceRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.Create(&models.ReviewOnService{
		Text:      request.Text,
		Author:    request.Author,
		ServiceId: request.ServiceId,
		Date:      request.Date,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetById получает отзыв по ID
func (h *ReviewOnServiceHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	result, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Отзыв не найден"})
	}

	return c.JSON(result)
}
