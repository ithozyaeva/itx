package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReferalLinkHandler struct {
	BaseHandler[models.ReferalLink]
	svc *service.ReferalLinkService
}

func NewReferalLinkHandler() *ReferalLinkHandler {
	svc := service.NewReferalLinkService()
	return &ReferalLinkHandler{
		BaseHandler: *NewBaseHandler(svc),
		svc:         svc,
	}
}

func (h *ReferalLinkHandler) Search(c *fiber.Ctx) error {
	req := new(models.SearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.service.Search(req.Limit, req.Offset, nil, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *ReferalLinkHandler) AddLink(c *fiber.Ctx) error {
	req := new(models.AddLinkRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	if member == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Произошла ошибка при получении пользователя"})
	}

	result, err := h.svc.AddLink(req, member)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *ReferalLinkHandler) UpdateLink(c *fiber.Ctx) error {
	req := new(models.UpdateLinkRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	existedLink, err := h.service.GetById(req.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if member.Id != existedLink.Author.Id {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Нельзя изменять чужие реферальные ссылки"})
	}

	result, err := h.svc.UpdateLink(req, member)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *ReferalLinkHandler) DeleteLink(c *fiber.Ctx) error {
	req := new(models.DeleteLinkRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	existedLink, err := h.service.GetById(req.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if member.Id != existedLink.Author.Id {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Нельзя изменять чужие реферальные ссылки"})
	}

	err = h.svc.Delete(&models.ReferalLink{Id: req.Id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(nil)
}
