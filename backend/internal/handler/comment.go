package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func respondCommentErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrCommentNotFound),
		errors.Is(err, service.ErrEntityNotFound),
		errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, service.ErrCommentForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
}

// ListForEntity создаёт хендлер для GET /<entity>/:id/comments. Замыкает
// тип сущности — фронт не должен передавать его явно, путь говорит сам
// за себя.
func (h *CommentHandler) ListForEntity(entityType models.CommentEntityType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		member, err := getMember(c)
		if err != nil {
			return err
		}
		entityID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
		}
		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		items, total, err := h.svc.List(entityType, entityID, member, hasAdminRole(member), limit, offset)
		if err != nil {
			return respondCommentErr(c, err)
		}
		return c.JSON(fiber.Map{"items": items, "total": total})
	}
}

func (h *CommentHandler) CreateForEntity(entityType models.CommentEntityType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		member, err := getMember(c)
		if err != nil {
			return err
		}
		entityID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
		}
		var req models.CreateCommentRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
		}
		created, err := h.svc.Create(entityType, entityID, member, req.Body)
		if err != nil {
			return respondCommentErr(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}
}

func (h *CommentHandler) Update(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	var req models.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	updated, err := h.svc.Update(id, member.Id, req.Body, hasAdminRole(member))
	if err != nil {
		return respondCommentErr(c, err)
	}
	return c.JSON(updated)
}

func (h *CommentHandler) Delete(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	if err := h.svc.Delete(id, member.Id, hasAdminRole(member)); err != nil {
		return respondCommentErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CommentHandler) ToggleLike(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	liked, count, err := h.svc.ToggleLike(id, member)
	if err != nil {
		return respondCommentErr(c, err)
	}
	return c.JSON(fiber.Map{"liked": liked, "likesCount": count})
}

func (h *CommentHandler) SetHidden(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	type body struct {
		Hidden bool `json:"hidden"`
	}
	var b body
	if err := c.BodyParser(&b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	if err := h.svc.SetHidden(id, b.Hidden, hasAdminRole(member)); err != nil {
		return respondCommentErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
