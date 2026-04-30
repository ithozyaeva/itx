package handler

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"
)

type AIMaterialHandler struct {
	svc *service.AIMaterialService
}

func NewAIMaterialHandler() *AIMaterialHandler {
	return &AIMaterialHandler{svc: service.NewAIMaterialService()}
}

// respondAIMaterialErr мапит sentinel-ошибки сервиса на HTTP-коды.
// Всё остальное (валидация и пр.) — 400.
func respondAIMaterialErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrAIMaterialNotFound),
		errors.Is(err, service.ErrAIMaterialCommentNotFound),
		errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, service.ErrAIMaterialForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
}

func (h *AIMaterialHandler) Search(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	filter := repository.AIMaterialFilter{
		Kind:     c.Query("kind"),
		// Теги хранятся в lowercase (см. service.normalizeTags), поэтому
		// фильтр должен быть case-insensitive — нормализуем на входе.
		Tag:      strings.ToLower(strings.TrimSpace(c.Query("tag"))),
		Query:    c.Query("q"),
		Sort:     c.Query("sort"),
		ViewerID: member.Id,
		Limit:    limit,
		Offset:   offset,
	}
	if c.Query("bookmarked") == "true" {
		filter.Bookmarked = true
	}
	if c.Query("mine") == "true" {
		filter.AuthorID = member.Id
		filter.IncludeHidden = true
	}

	if filter.Kind != "" && !models.IsValidAIMaterialKind(models.AIMaterialKind(filter.Kind)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректная категория"})
	}

	items, total, err := h.svc.Search(filter)
	if err != nil {
		log.Printf("AIMaterial search error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить материалы"})
	}

	return c.JSON(fiber.Map{"items": items, "total": total})
}

func (h *AIMaterialHandler) GetByID(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	item, err := h.svc.GetByID(id, member.Id, hasAdminRole(member))
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.JSON(item)
}

func (h *AIMaterialHandler) Create(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	var req models.CreateAIMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	item, err := h.svc.Create(&req, member.Id)
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *AIMaterialHandler) Update(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	var req models.UpdateAIMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	item, err := h.svc.Update(id, &req, member.Id, hasAdminRole(member))
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.JSON(item)
}

func (h *AIMaterialHandler) Delete(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.svc.Delete(id, member.Id, hasAdminRole(member)); err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// SetHidden — admin-only мягкое скрытие материала из листинга.
func (h *AIMaterialHandler) SetHidden(c *fiber.Ctx) error {
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
		return respondAIMaterialErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AIMaterialHandler) ToggleLike(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	liked, count, err := h.svc.ToggleLike(id, member.Id, hasAdminRole(member))
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.JSON(fiber.Map{"liked": liked, "likesCount": count})
}

func (h *AIMaterialHandler) ToggleBookmark(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	bookmarked, count, err := h.svc.ToggleBookmark(id, member.Id, hasAdminRole(member))
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.JSON(fiber.Map{"bookmarked": bookmarked, "bookmarksCount": count})
}

func (h *AIMaterialHandler) ListComments(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	items, err := h.svc.ListComments(id, member.Id, hasAdminRole(member))
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.JSON(fiber.Map{"items": items})
}

type aiMaterialCommentBody struct {
	Body string `json:"body"`
}

func (h *AIMaterialHandler) CreateComment(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	var body aiMaterialCommentBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	created, err := h.svc.CreateComment(id, member.Id, body.Body, hasAdminRole(member))
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *AIMaterialHandler) UpdateComment(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	var body aiMaterialCommentBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	updated, err := h.svc.UpdateComment(id, member.Id, body.Body, hasAdminRole(member))
	if err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.JSON(updated)
}

func (h *AIMaterialHandler) DeleteComment(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	if err := h.svc.DeleteComment(id, member.Id, hasAdminRole(member)); err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AIMaterialHandler) SetCommentHidden(c *fiber.Ctx) error {
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
	if err := h.svc.SetCommentHidden(id, b.Hidden, hasAdminRole(member)); err != nil {
		return respondAIMaterialErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AIMaterialHandler) TopTags(c *fiber.Ctx) error {
	q := c.Query("q")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	tags, err := h.svc.TopTags(q, limit)
	if err != nil {
		log.Printf("AIMaterial top tags error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить теги"})
	}
	return c.JSON(fiber.Map{"tags": tags})
}

func hasAdminRole(member *models.Member) bool {
	for _, role := range member.Roles {
		if role == models.MemberRoleAdmin {
			return true
		}
	}
	return false
}
