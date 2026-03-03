package handler

import (
	"fmt"
	"ithozyeva/internal/service"
	"strconv"
	"strings"

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

	// Фильтр по пользователю
	if uid := c.Query("user_id"); uid != "" {
		if parsed, err := strconv.ParseInt(uid, 10, 64); err == nil {
			activity, err := h.service.GetDailyActivityByUser(parsed, days)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки графика"})
			}
			return c.JSON(activity)
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

// GetUserStats возвращает статистику конкретного пользователя
func (h *ChatActivityHandler) GetUserStats(c *fiber.Ctx) error {
	uid := c.Query("user_id")
	if uid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Необходимо указать user_id"})
	}

	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный user_id"})
	}

	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 90 {
			days = parsed
		}
	}

	stats, err := h.service.GetUserStats(userID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки статистики пользователя"})
	}
	return c.JSON(stats)
}

// ExportCSV возвращает данные активности в формате CSV
func (h *ChatActivityHandler) ExportCSV(c *fiber.Ctx) error {
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

	rows, err := h.service.GetMessagesForExport(chatID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка экспорта"})
	}

	// Формируем CSV
	var sb strings.Builder
	sb.WriteString("Дата,Чат,Пользователь,Сообщений\n")
	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%d\n",
			row.Date,
			escapeCsvField(row.ChatTitle),
			escapeCsvField(row.TelegramUsername),
			row.MessageCount,
		))
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=chat-activity.csv")
	return c.SendString(sb.String())
}

func escapeCsvField(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}
