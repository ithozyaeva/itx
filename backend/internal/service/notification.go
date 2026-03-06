package service

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"log"
)

// CreateNotification создаёт in-app уведомление для участника
func CreateNotification(memberId int64, notifType string, title string, body string) error {
	notification := models.Notification{
		MemberId: memberId,
		Type:     notifType,
		Title:    title,
		Body:     body,
	}
	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification for member %d: %v", memberId, err)
		return err
	}
	GetSSEHub().Publish(memberId, SSEEvent{Type: "notifications"})
	return nil
}
