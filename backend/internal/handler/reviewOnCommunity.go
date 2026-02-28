package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReviewOnCommunityHandler struct {
	BaseHandler[models.ReviewOnCommunity]
	svc      *service.ReviewOnCommunityService
	auditSvc *service.AuditService
}

func NewReviewOnCommunityHandler() *ReviewOnCommunityHandler {
	svc := service.NewReviewOnCommunityService()
	return &ReviewOnCommunityHandler{
		BaseHandler: *NewBaseHandler(svc),
		svc:         svc,
		auditSvc:    service.NewAuditService(),
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

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "review_on_community", 0, review.AuthorTg)

	return c.SendStatus(fiber.StatusOK)

}

func (h *ReviewOnCommunityHandler) GetMyReviews(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	reviews, err := h.svc.GetByAuthorId(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reviews)
}

func (h *ReviewOnCommunityHandler) UpdateMyReview(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	review, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Отзыв не найден"})
	}
	if int64(review.AuthorId) != member.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Нет доступа"})
	}

	req := new(models.AddReviewOnCommunityRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	review.Text = req.Text
	result, err := h.svc.Update(review)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *ReviewOnCommunityHandler) DeleteMyReview(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	review, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Отзыв не найден"})
	}
	if int64(review.AuthorId) != member.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Нет доступа"})
	}

	if err := h.svc.Delete(review); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
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

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionApprove, "review_on_community", int64(id), result.Text)

	return c.JSON(result)
}

// Delete переопределяет базовый метод Delete для добавления аудит-лога
func (h *ReviewOnCommunityHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity, err := h.service.GetById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Отзыв не найден"})
	}

	if err := h.service.Delete(entity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionDelete, "review_on_community", int64(id), entity.Text)

	return c.SendStatus(fiber.StatusNoContent)
}
