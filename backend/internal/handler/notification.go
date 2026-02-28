package handler

import (
	"log"
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

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Уведомление не найдено"})
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

// GetEventMemberIds returns member IDs for an event (call before deleting the event)
func GetEventMemberIds(eventId int64) []int64 {
	var memberIds []int64
	if err := database.DB.Raw(`
		SELECT member_id FROM event_members WHERE event_id = ?
	`, eventId).Scan(&memberIds).Error; err != nil {
		log.Printf("Error getting event member IDs for event %d: %v", eventId, err)
	}
	return memberIds
}

// CreateNotificationsForMembers creates notifications for a list of member IDs
func CreateNotificationsForMembers(memberIds []int64, notifType string, title string, body string) {
	for _, memberId := range memberIds {
		if err := CreateNotification(memberId, notifType, title, body); err != nil {
			log.Printf("Error creating notification for member %d: %v", memberId, err)
		}
	}
}
