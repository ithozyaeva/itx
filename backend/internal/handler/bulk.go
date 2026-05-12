package handler

import (
	"fmt"
	"log"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
		log.Printf("bulk delete events error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка массового удаления событий"})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	service.SafeGo("audit bulk-delete events", func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "event", id, fmt.Sprintf("bulk delete #%d", id))
		}
	})

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteMentors(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.MentorDbShortModel{}).Error; err != nil {
		log.Printf("bulk delete mentors error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка массового удаления менторов"})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	service.SafeGo("audit bulk-delete mentors", func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "mentor", id, fmt.Sprintf("bulk delete #%d", id))
		}
	})

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteMembers(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.Member{}).Error; err != nil {
		log.Printf("bulk delete members error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка массового удаления участников"})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	service.SafeGo("audit bulk-delete members", func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "member", id, fmt.Sprintf("bulk delete #%d", id))
		}
	})

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteReviews(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.ReviewOnCommunity{}).Error; err != nil {
		log.Printf("bulk delete reviews error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка массового удаления отзывов"})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	service.SafeGo("audit bulk-delete reviews", func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "review_on_community", id, fmt.Sprintf("bulk delete #%d", id))
		}
	})

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}

func (h *BulkHandler) BulkApproveReviews(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	var reviews []models.ReviewOnCommunity
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ReviewOnCommunity{}).
			Where("id IN ? AND status != ?", req.Ids, "APPROVED").
			Update("status", "APPROVED").Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", req.Ids).Find(&reviews).Error
	})
	if err != nil {
		log.Printf("BulkApproveReviews error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка при одобрении отзывов"})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	service.SafeGo("audit bulk-approve community reviews", func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionApprove, "review_on_community", id, fmt.Sprintf("bulk approve #%d", id))
		}
		for _, review := range reviews {
			h.pointsSvc.AwardIdempotent(int64(review.AuthorId), models.PointReasonReviewCommunity, "review_community", int64(review.Id), "Отзыв о сообществе")
			CreateNotification(int64(review.AuthorId), "review_approved", "Отзыв одобрен", "Ваш отзыв о сообществе был одобрен")
		}
	})

	return c.JSON(fiber.Map{"approved": len(req.Ids)})
}

func (h *BulkHandler) BulkApproveServiceReviews(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	var reviews []models.ReviewOnService
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ReviewOnService{}).
			Where("id IN ? AND status != ?", req.Ids, "APPROVED").
			Update("status", "APPROVED").Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", req.Ids).Find(&reviews).Error
	})
	if err != nil {
		log.Printf("BulkApproveServiceReviews error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка при одобрении отзывов"})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	service.SafeGo("audit bulk-approve service reviews", func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionApprove, "review_on_service", id, fmt.Sprintf("bulk approve #%d", id))
		}
		for _, review := range reviews {
			if review.AuthorMemberId != nil {
				h.pointsSvc.AwardIdempotent(*review.AuthorMemberId, models.PointReasonReviewService, "review_service", int64(review.Id), "Отзыв на услугу ментора")
				CreateNotification(*review.AuthorMemberId, "review_approved", "Отзыв одобрен", "Ваш отзыв на услугу ментора был одобрен")
			}
		}
	})

	return c.JSON(fiber.Map{"approved": len(req.Ids)})
}

func (h *BulkHandler) BulkDeleteMentorsReviews(c *fiber.Ctx) error {
	req := new(BulkIdsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Where("id IN ?", req.Ids).Delete(&models.ReviewOnService{}).Error; err != nil {
		log.Printf("bulk delete mentor reviews error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка массового удаления отзывов"})
	}

	actorId, actorName, actorType := getActorId(c), getActorName(c), getActorType(c)
	service.SafeGo("audit bulk-delete mentor reviews", func() {
		for _, id := range req.Ids {
			h.auditSvc.Log(actorId, actorName, actorType, models.AuditActionDelete, "review_on_service", id, fmt.Sprintf("bulk delete #%d", id))
		}
	})

	return c.JSON(fiber.Map{"deleted": len(req.Ids)})
}
