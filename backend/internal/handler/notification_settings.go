package handler

import (
	"log"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type NotificationSettingsHandler struct {
	svc *service.NotificationSettingsService
}

type updateNotificationSettingsRequest struct {
	MuteAll        bool `json:"muteAll"`
	NewEvents      bool `json:"newEvents"`
	RemindWeek     bool `json:"remindWeek"`
	RemindDay      bool `json:"remindDay"`
	RemindHour     bool `json:"remindHour"`
	EventStart     bool `json:"eventStart"`
	EventUpdates   bool `json:"eventUpdates"`
	EventCancelled bool `json:"eventCancelled"`
	DailyMorning   bool `json:"dailyMorning"`
	DailyEvening   bool `json:"dailyEvening"`
	DailyStreak    bool `json:"dailyStreak"`
	DailyRaffle    bool `json:"dailyRaffle"`
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
		log.Printf("get notification settings error (member=%d): %v", memberId, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки настроек уведомлений"})
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
		MuteAll:        req.MuteAll,
		NewEvents:      req.NewEvents,
		RemindWeek:     req.RemindWeek,
		RemindDay:      req.RemindDay,
		RemindHour:     req.RemindHour,
		EventStart:     req.EventStart,
		EventUpdates:   req.EventUpdates,
		EventCancelled: req.EventCancelled,
		DailyMorning:   req.DailyMorning,
		DailyEvening:   req.DailyEvening,
		DailyStreak:    req.DailyStreak,
		DailyRaffle:    req.DailyRaffle,
	}

	result, err := h.svc.Update(memberId, settings)
	if err != nil {
		log.Printf("update notification settings error (member=%d): %v", memberId, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления настроек уведомлений"})
	}

	return c.JSON(result)
}
