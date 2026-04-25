package handler

import (
	"bufio"
	"fmt"
	"ithozyeva/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

type SSEHandler struct{}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

func (h *SSEHandler) Stream(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	hub := service.GetSSEHub()
	ch, unsubscribe := hub.Subscribe(member.Id)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer unsubscribe()

		fmt.Fprintf(w, "data: {\"type\":\"connected\"}\n\n")
		w.Flush()

		for msg := range ch {
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				log.Printf("SSE write error for member %d: %v", member.Id, err)
				return
			}
			if err = w.Flush(); err != nil {
				return
			}
		}
	})

	return nil
}

// PublishToMember отправляет событие конкретному пользователю
func PublishToMember(memberId int64, eventType string) {
	service.GetSSEHub().Publish(memberId, service.SSEEvent{Type: eventType})
}

// BroadcastEvent отправляет событие всем подключённым пользователям
func BroadcastEvent(eventType string) {
	service.GetSSEHub().Broadcast(service.SSEEvent{Type: eventType})
}
