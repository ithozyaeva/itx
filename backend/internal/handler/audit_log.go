package handler

import (
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuditLogHandler struct {
	svc *service.AuditService
}

func NewAuditLogHandler() *AuditLogHandler {
	return &AuditLogHandler{
		svc: service.NewAuditService(),
	}
}

type AuditLogSearchRequest struct {
	Limit      *int    `query:"limit"`
	Offset     *int    `query:"offset"`
	ActorType  *string `query:"actorType"`
	Action     *string `query:"action"`
	EntityType *string `query:"entityType"`
}

func (h *AuditLogHandler) Search(c *fiber.Ctx) error {
	req := new(AuditLogSearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	filter := make(repository.SearchFilter)

	if req.ActorType != nil && *req.ActorType != "" {
		filter["actor_type = ?"] = *req.ActorType
	}
	if req.Action != nil && *req.Action != "" {
		filter["action = ?"] = *req.Action
	}
	if req.EntityType != nil && *req.EntityType != "" {
		filter["entity_type = ?"] = *req.EntityType
	}

	result, err := h.svc.Search(req.Limit, req.Offset, &filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
