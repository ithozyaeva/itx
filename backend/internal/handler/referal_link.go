package handler

import (
	"strings"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReferalLinkHandler struct {
	BaseHandler[models.ReferalLink]
	svc       *service.ReferalLinkService
	auditSvc  *service.AuditService
	pointsSvc *service.PointsService
}

func NewReferalLinkHandler() *ReferalLinkHandler {
	svc := service.NewReferalLinkService()
	return &ReferalLinkHandler{
		BaseHandler: *NewBaseHandler(svc),
		svc:         svc,
		auditSvc:    service.NewAuditService(),
		pointsSvc:   service.NewPointsService(),
	}
}

type ReferalLinkSearchRequest struct {
	Limit   *int    `query:"limit"`
	Offset  *int    `query:"offset"`
	Grade   *string `query:"grade"`
	Company *string `query:"company"`
	Status  *string `query:"status"`
}

func (h *ReferalLinkHandler) Search(c *fiber.Ctx) error {
	req := new(ReferalLinkSearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	filter := make(repository.SearchFilter)

	if req.Grade != nil && *req.Grade != "" {
		filter["grade = ?"] = *req.Grade
	}
	if req.Company != nil && *req.Company != "" {
		filter["company ILIKE ?"] = "%" + *req.Company + "%"
	}
	if req.Status != nil && *req.Status != "" {
		filter["status = ?"] = *req.Status
	}

	result, err := h.service.Search(req.Limit, req.Offset, &filter, nil)
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

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "referal_link", result.Id, result.Company)
	go h.pointsSvc.GiveForAction(member.Id, models.PointReasonReferalCreate, "referal_link", result.Id,
		"Создание реферальной ссылки")

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

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionUpdate, "referal_link", result.Id, result.Company)

	return c.JSON(result)
}

type TrackConversionRequest struct {
	ReferralLinkId int64 `json:"referralLinkId"`
}

func (h *ReferalLinkHandler) TrackConversion(c *fiber.Ctx) error {
	req := new(TrackConversionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	err := h.svc.TrackConversion(req.ReferralLinkId, member.Id)
	if err != nil {
		// Handle unique constraint violation (duplicate conversion) as idempotent success
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "referral_conversions_unique") {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Начисляем баллы автору реферальной ссылки за конверсию
	go func() {
		link, err := h.service.GetById(req.ReferralLinkId)
		if err == nil && link.AuthorId != 0 {
			h.pointsSvc.AwardIdempotent(link.AuthorId, models.PointReasonReferalConversion, "referal_conversion", req.ReferralLinkId,
				"Конверсия по реферальной ссылке")
		}
	}()

	return c.SendStatus(fiber.StatusOK)
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

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionDelete, "referal_link", req.Id, "")

	return c.JSON(nil)
}
