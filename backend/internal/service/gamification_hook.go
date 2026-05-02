package service

import (
	"log"
	"sync"
)

// gamificationHooks — синглтон-точка интеграции дейликов в существующие
// хендлеры. Все вызовы выполняются в фоне и под защитой recover, чтобы
// сбой геймификации не уронил основное действие пользователя.
var (
	gamificationOnce sync.Once
	dailyTaskSvc     *DailyTaskService
)

func ensureGamificationHooks() {
	gamificationOnce.Do(func() {
		dailyTaskSvc = NewDailyTaskService()
	})
}

// TrackDailyTrigger инкрементирует прогресс дейликов с указанным trigger_key.
// Безопасен для вызова из любого хендлера: не блокирует, не паникует.
//
// Вызывайте сразу после успешного основного действия (создание комментария,
// отправка kudos и т.п.). Идемпотентность достигается на уровне
// daily_task_progress (UNIQUE на member+day+task).
func TrackDailyTrigger(memberId int64, triggerKey string, n int) {
	ensureGamificationHooks()
	if memberId == 0 || triggerKey == "" {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("daily-trigger panic (member=%d, trigger=%s): %v", memberId, triggerKey, r)
			}
		}()
		dailyTaskSvc.IncrementProgress(memberId, triggerKey, n)
	}()
}
