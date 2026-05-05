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
	member, err := getMember(c)
	if err != nil {
		return err
	}

	var notifications []models.Notification
	if err := database.DB.Where("member_id = ?", member.Id).
		Order("created_at DESC").
		Limit(50).
		Find(&notifications).Error; err != nil {
		log.Printf("get notifications error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки уведомлений"})
	}

	return c.JSON(notifications)
}

func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	var count int64
	if err := database.DB.Model(&models.Notification{}).
		Where("member_id = ? AND read = false", member.Id).
		Count(&count).Error; err != nil {
		log.Printf("get unread count error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка получения количества уведомлений"})
	}

	return c.JSON(fiber.Map{"count": count})
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	result := database.DB.Model(&models.Notification{}).
		Where("id = ? AND member_id = ?", id, member.Id).
		Update("read", true)

	if result.Error != nil {
		log.Printf("mark notification as read error (id=%d, member=%d): %v", id, member.Id, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления уведомления"})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Уведомление не найдено"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	if err := database.DB.Model(&models.Notification{}).
		Where("member_id = ? AND read = false", member.Id).
		Update("read", true).Error; err != nil {
		log.Printf("mark all notifications as read error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления уведомлений"})
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
	if err := database.DB.Create(&notification).Error; err != nil {
		return err
	}
	PublishToMember(memberId, "notifications")
	return nil
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

// CreateNotificationsForMembers батчем создаёт уведомления для списка
// memberIds (один INSERT на CreateInBatches вместо N×INSERT) и затем
// шлёт SSE-publish каждому. Для большого ивента (200+ участников) это
// сокращает время с N round-trip'ов до Postgres до одного-двух батчей.
func CreateNotificationsForMembers(memberIds []int64, notifType string, title string, body string) {
	if len(memberIds) == 0 {
		return
	}
	rows := make([]models.Notification, len(memberIds))
	for i, id := range memberIds {
		rows[i] = models.Notification{
			MemberId: id,
			Type:     notifType,
			Title:    title,
			Body:     body,
		}
	}
	if err := database.DB.CreateInBatches(rows, 500).Error; err != nil {
		log.Printf("CreateNotificationsForMembers batch insert (%d rows): %v", len(rows), err)
		return
	}
	for _, id := range memberIds {
		PublishToMember(id, "notifications")
	}
}
