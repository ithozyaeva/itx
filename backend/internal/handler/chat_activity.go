package handler

import (
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ChatActivityHandler struct {
	service *service.ChatActivityService
}

func NewChatActivityHandler() *ChatActivityHandler {
	return &ChatActivityHandler{
		service: service.NewChatActivityService(),
	}
}

// GetStats возвращает общую статистику активности чатов
func (h *ChatActivityHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.GetStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки статистики"})
	}
	return c.JSON(stats)
}

// GetChart возвращает данные для графика активности
func (h *ChatActivityHandler) GetChart(c *fiber.Ctx) error {
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 90 {
			days = parsed
		}
	}

	var chatID *int64
	if cid := c.Query("chat_id"); cid != "" {
		if parsed, err := strconv.ParseInt(cid, 10, 64); err == nil {
			chatID = &parsed
		}
	}

	activity, err := h.service.GetActivityChart(chatID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки графика"})
	}
	return c.JSON(activity)
}

// GetTopUsers возвращает топ пользователей по активности
func (h *ChatActivityHandler) GetTopUsers(c *fiber.Ctx) error {
	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 90 {
			days = parsed
		}
	}

	limit := 5
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	users, err := h.service.GetTopUsers(days, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки топ пользователей"})
	}
	return c.JSON(users)
}

// GetChats возвращает список отслеживаемых чатов
func (h *ChatActivityHandler) GetChats(c *fiber.Ctx) error {
	chats, err := h.service.GetTrackedChats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки чатов"})
	}
	return c.JSON(chats)
}
