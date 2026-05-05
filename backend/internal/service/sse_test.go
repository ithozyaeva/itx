package service

import (
	"sync"
	"testing"
	"time"
)

// TestSSEHub_PublishWhileUnsubscribe — регрессия Bug #2:
// до фикса Publish брал snapshot каналов под RLock и отпускал блокировку
// перед циклом send. Параллельный unsubscribe (под Lock) успевал
// close(ch) до отправки, и `case ch <- msg:` валил процесс с
// «panic: send on closed channel».
//
// Тест воспроизводит этот сценарий: много subscribe/unsubscribe и
// много publish одновременно. Без фикса — runtime fatal с панику в
// тестовом процессе. С фиксом — успешно завершается.
func TestSSEHub_PublishWhileUnsubscribe(t *testing.T) {
	hub := &SSEHub{
		clients: make(map[int64]map[chan string]struct{}),
	}

	const memberID int64 = 42
	const subscribers = 50
	const publishers = 8
	const duration = 200 * time.Millisecond

	// Стартуем читатели каналов: они принимают сообщения и сразу
	// отписываются по случайному моменту, имитируя disconnect клиента
	// в произвольный момент.
	var subWg sync.WaitGroup
	for i := 0; i < subscribers; i++ {
		subWg.Add(1)
		go func(idx int) {
			defer subWg.Done()
			ch, unsubscribe := hub.Subscribe(memberID)
			// Простой drain — без него переполнится буфер канала.
			done := make(chan struct{})
			go func() {
				for range ch {
				}
				close(done)
			}()
			// Случайная задержка: одни отпишутся быстро, другие позже —
			// чтобы окно гонки попадало в разные моменты publish-цикла.
			time.Sleep(time.Duration(idx%20) * time.Millisecond)
			unsubscribe()
			<-done
		}(i)
	}

	// Стартуем publisher'ы: молотим Publish в цикле всё время duration.
	stop := make(chan struct{})
	var pubWg sync.WaitGroup
	for i := 0; i < publishers; i++ {
		pubWg.Add(1)
		go func() {
			defer pubWg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					hub.Publish(memberID, SSEEvent{Type: "notifications"})
				}
			}
		}()
	}

	// Ждём субскрайберов (они закроются после своих тайм-слотов).
	subWg.Wait()
	close(stop)
	pubWg.Wait()

	// Если дошли сюда — паники не было. Доп. assert: после полного
	// unsubscribe map пуст.
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	if _, ok := hub.clients[memberID]; ok {
		t.Errorf("hub.clients[%d] остался после unsubscribe всех; expected entry deleted", memberID)
	}
}

// TestSSEHub_PublishDelivers — sanity-чек: Publish реально доставляет
// сообщения подписчикам (не превратили fix в no-op).
func TestSSEHub_PublishDelivers(t *testing.T) {
	hub := &SSEHub{
		clients: make(map[int64]map[chan string]struct{}),
	}
	const memberID int64 = 7

	ch, unsubscribe := hub.Subscribe(memberID)
	defer unsubscribe()

	hub.Publish(memberID, SSEEvent{Type: "tasks"})

	select {
	case msg := <-ch:
		if msg == "" {
			t.Fatalf("got empty message")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Publish did not deliver to subscriber within 500ms")
	}
}
