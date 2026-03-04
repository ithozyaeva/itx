package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type NotificationSettingsHandler struct {
	svc *service.NotificationSettingsService
}

type updateNotificationSettingsRequest struct {
	NewEvents      bool `json:"newEvents"`
	RemindWeek     bool `json:"remindWeek"`
	RemindDay      bool `json:"remindDay"`
	RemindHour     bool `json:"remindHour"`
	EventStart     bool `json:"eventStart"`
	EventUpdates   bool `json:"eventUpdates"`
	EventCancelled bool `json:"eventCancelled"`
}

func NewNotificationSettingsHandler() *NotificationSettingsHandler {
	return &NotificationSettingsHandler{
		svc: service.NewNotificationSettingsService(),
	}
}

func (h *NotificationSettingsHandler) GetMy(c *fiber.Ctx) error {
	memberId := getActorId(c)
	if memberId == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Не авторизован"})
	}

	settings, err := h.svc.GetByMemberId(memberId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(settings)
}

func (h *NotificationSettingsHandler) UpdateMy(c *fiber.Ctx) error {
	memberId := getActorId(c)
	if memberId == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Не авторизован"})
	}

	req := new(updateNotificationSettingsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	settings := &models.NotificationSettings{
		NewEvents:      req.NewEvents,
		RemindWeek:     req.RemindWeek,
		RemindDay:      req.RemindDay,
		RemindHour:     req.RemindHour,
		EventStart:     req.EventStart,
		EventUpdates:   req.EventUpdates,
		EventCancelled: req.EventCancelled,
	}

	result, err := h.svc.Update(memberId, settings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
