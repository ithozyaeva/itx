package service

import (
	"encoding/json"
	"log"
	"sync"
)

// SSEEvent представляет событие для отправки клиентам
type SSEEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// SSEHub управляет подключениями SSE клиентов
type SSEHub struct {
	mu      sync.RWMutex
	clients map[int64]map[chan string]struct{} // memberId -> set of channels
}

var sseHub = &SSEHub{
	clients: make(map[int64]map[chan string]struct{}),
}

// GetSSEHub возвращает глобальный SSE хаб
func GetSSEHub() *SSEHub {
	return sseHub
}

// Subscribe добавляет клиента для memberId, возвращает канал и функцию отписки
func (h *SSEHub) Subscribe(memberId int64) (chan string, func()) {
	ch := make(chan string, 16)
	h.mu.Lock()
	if h.clients[memberId] == nil {
		h.clients[memberId] = make(map[chan string]struct{})
	}
	h.clients[memberId][ch] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		delete(h.clients[memberId], ch)
		if len(h.clients[memberId]) == 0 {
			delete(h.clients, memberId)
		}
		h.mu.Unlock()
		close(ch)
	}

	return ch, unsubscribe
}

// Publish отправляет событие конкретному пользователю
func (h *SSEHub) Publish(memberId int64, event SSEEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("SSE marshal error: %v", err)
		return
	}
	msg := string(data)

	h.mu.RLock()
	channels := h.clients[memberId]
	h.mu.RUnlock()

	for ch := range channels {
		select {
		case ch <- msg:
		default:
			// канал переполнен — пропускаем
		}
	}
}

// Broadcast отправляет событие всем подключённым пользователям
func (h *SSEHub) Broadcast(event SSEEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("SSE marshal error: %v", err)
		return
	}
	msg := string(data)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, channels := range h.clients {
		for ch := range channels {
			select {
			case ch <- msg:
			default:
			}
		}
	}
}
