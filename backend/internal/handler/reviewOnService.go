package handler

import (
	"log"
	"strconv"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReviewOnServiceHandler struct {
	BaseHandler[models.ReviewOnService]
	svc       *service.ReviewOnServiceService
	auditSvc  *service.AuditService
	pointsSvc *service.PointsService
}

func NewReviewOnServiceHandler() *ReviewOnServiceHandler {
	svc := service.NewReviewOnServiceService()
	return &ReviewOnServiceHandler{
		BaseHandler: *NewBaseHandler[models.ReviewOnService](svc),
		svc:         svc,
		auditSvc:    service.NewAuditService(),
		pointsSvc:   service.NewPointsService(),
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
		log.Printf("search service reviews error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка поиска отзывов"})
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
		log.Printf("get reviews with mentor info error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки отзывов"})
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
		log.Printf("create service review error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка создания отзыва"})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "review_on_service", int64(result.Id), result.Author)

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Delete переопределяет базовый метод Delete для добавления аудит-лога
func (h *ReviewOnServiceHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Отзыв не найден"})
	}

	if err := h.svc.Delete(entity); err != nil {
		log.Printf("delete service review error (id=%d): %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления отзыва"})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionDelete, "review_on_service", id, entity.Author)

	return c.SendStatus(fiber.StatusNoContent)
}

// Approve одобряет отзыв на услугу
func (h *ReviewOnServiceHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	result, err := h.svc.Approve(id)
	if err != nil {
		log.Printf("approve service review error (id=%d): %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка одобрения отзыва"})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionApprove, "review_on_service", id, result.Author)

	if result.AuthorMemberId != nil {
		go h.pointsSvc.AwardIdempotent(*result.AuthorMemberId, models.PointReasonReviewService, "review_service", int64(result.Id), "Отзыв на услугу ментора")
		go CreateNotification(*result.AuthorMemberId, "review_approved", "Отзыв одобрен", "Ваш отзыв на услугу ментора был одобрен")
	}

	return c.JSON(result)
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
