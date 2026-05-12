package service

import (
	"log"
	"runtime/debug"
)

// SafeGo запускает fn в фоновой goroutine, гарантируя что panic не уронит
// весь процесс. До появления хелпера в коде было 15+ голых `go func() {...}()`
// в хендлерах и боте; один nil-deref в любой из них валил backend целиком
// (Go runtime: панику в goroutine без recover нельзя поймать снаружи).
//
// label — короткое описание для лога, чтобы по «safego panic [auth-role-sync]»
// можно было быстро найти место.
//
// Используйте вместо bare go-launch'а для любых fire-and-forget задач:
// нотификации, побочные эффекты на points/credits, обновления subscription,
// SSE-publish и пр.
//
// Уже-обёрнутые goroutines (gamification_hook.go, mentor repo) можно по
// желанию переписать на SafeGo, но не обязательно — поведение идентичное.
func SafeGo(label string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("safego panic [%s]: %v\n%s", label, r, debug.Stack())
			}
		}()
		fn()
	}()
}
