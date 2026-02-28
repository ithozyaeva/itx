package handler

import (
	"strconv"

	"ithozyeva/database"
	"ithozyeva/internal/models"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) GetMy(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	var notifications []models.Notification
	if err := database.DB.Where("member_id = ?", member.Id).
		Order("created_at DESC").
		Limit(50).
		Find(&notifications).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(notifications)
}

func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	var count int64
	if err := database.DB.Model(&models.Notification{}).
		Where("member_id = ? AND read = false", member.Id).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"count": count})
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	result := database.DB.Model(&models.Notification{}).
		Where("id = ? AND member_id = ?", id, member.Id).
		Update("read", true)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	if err := database.DB.Model(&models.Notification{}).
		Where("member_id = ? AND read = false", member.Id).
		Update("read", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// CreateNotification is a helper to create notifications from other handlers
func CreateNotification(memberId int64, notifType string, title string, body string) error {
	notification := models.Notification{
		MemberId: memberId,
		Type:     notifType,
		Title:    title,
		Body:     body,
	}
	return database.DB.Create(&notification).Error
}

// CreateNotificationsForEventMembers creates notifications for all members of an event
func CreateNotificationsForEventMembers(eventId int64, notifType string, title string, body string) {
	var members []models.Member
	database.DB.Raw(`
		SELECT m.* FROM members m
		JOIN event_members em ON em.member_id = m.id
		WHERE em.event_id = ?
	`, eventId).Scan(&members)

	for _, member := range members {
		CreateNotification(member.Id, notifType, title, body)
	}
}
