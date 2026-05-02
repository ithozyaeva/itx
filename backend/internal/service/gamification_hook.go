package service

import (
	"log"
	"sync"
)

// gamificationHooks — синглтон-точка интеграции дейликов и челленджей в
// существующие хендлеры. Все вызовы выполняются в фоне и под защитой
// recover, чтобы сбой геймификации не уронил основное действие.
var (
	gamificationOnce sync.Once
	dailyTaskSvc     *DailyTaskService
	challengeSvc     *ChallengeService
)

func ensureGamificationHooks() {
	gamificationOnce.Do(func() {
		dailyTaskSvc = NewDailyTaskService()
		challengeSvc = NewChallengeService()
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

// TrackChallengeMetric инкрементирует прогресс по активным челлендж-инстансам
// с указанным metric_key. Аналогичен TrackDailyTrigger.
//
// Используйте после действий, влияющих на еженедельные/ежемесячные челленджи
// (events_attended, kudos_received, comments_posted и т.д.).
func TrackChallengeMetric(memberId int64, metricKey string, n int) {
	ensureGamificationHooks()
	if memberId == 0 || metricKey == "" {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("challenge-metric panic (member=%d, metric=%s): %v", memberId, metricKey, r)
			}
		}()
		challengeSvc.IncrementProgress(memberId, metricKey, n)
	}()
}
