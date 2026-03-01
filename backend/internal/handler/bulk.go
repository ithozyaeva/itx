package handler

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type BulkHandler struct {
	auditSvc  *service.AuditService
	pointsSvc *service.PointsService
}

func NewBulkHandler() *BulkHandler {
	return &BulkHandler{
		auditSvc:  service.NewAuditService(),
		pointsSvc: service.NewPointsService(),
	}
}

type BulkIdsRequest struct {
	Ids []int64 `json:"ids"`
}

func (h *BulkHandler) BulkDeleteEvents(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.Event{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	go func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "event", id, fmt.Sprintf("bulk delete #%d", id))
		}
	}()

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteMentors(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.MentorDbShortModel{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	go func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "mentor", id, fmt.Sprintf("bulk delete #%d", id))
		}
	}()

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteMembers(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.Member{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	go func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "member", id, fmt.Sprintf("bulk delete #%d", id))
		}
	}()

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteReviews(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.ReviewOnCommunity{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	go func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "review_on_community", id, fmt.Sprintf("bulk delete #%d", id))
		}
	}()

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkApproveReviews(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Model(&models.ReviewOnCommunity{}).Where("id IN ?", req.Ids).Update("status", "APPROVED").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var reviews []models.ReviewOnCommunity
	database.DB.Where("id IN ?", req.Ids).Find(&reviews)

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	go func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionApprove, "review_on_community", id, fmt.Sprintf("bulk approve #%d", id))
		}
		for _, review := range reviews {
			h.pointsSvc.AwardIdempotent(int64(review.AuthorId), models.PointReasonReviewCommunity, "review_community", int64(review.Id), "Отзыв о сообществе")
		}
	}()

	return c.JSON(fiber.Map{"approved": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteMentorsReviews(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.ReviewOnService{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	go func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "review_on_service", id, fmt.Sprintf("bulk delete #%d", id))
		}
	}()

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}
