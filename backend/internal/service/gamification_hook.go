package service

import (
	"log"
	"sync"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

// gamificationHooks — синглтон-точка интеграции дейликов и челленджей в
// существующие хендлеры. Все вызовы выполняются в фоне и под защитой
// recover, чтобы сбой геймификации не уронил основное действие.
var (
	gamificationOnce sync.Once
	dailyTaskSvc     *DailyTaskService
	challengeSvc     *ChallengeService
	dailyRaffleSvc   *DailyRaffleService
	raffleRepo       *repository.RaffleRepository
)

func ensureGamificationHooks() {
	gamificationOnce.Do(func() {
		dailyTaskSvc = NewDailyTaskService()
		challengeSvc = NewChallengeService()
		dailyRaffleSvc = NewDailyRaffleService()
		raffleRepo = repository.NewRaffleRepository()
	})
}

// TrackDailyTrigger инкрементирует прогресс дейликов с указанным trigger_key.
// Безопасен для вызова из любого хендлера: не блокирует, не паникует.
//
// Вызывайте сразу после успешного основного действия (создание комментария,
// отправка kudos и т.п.). Идемпотентность достигается на уровне
// daily_task_progress (UNIQUE на member+day+task).
//
// Перед запуском goroutine consult'имся с triggerCache: если сегодняшний
// набор не содержит этот trigger_key — выходим без планирования SQL.
// Это сильно снижает нагрузку на view-эндпоинтах, которые вызываются часто.
func TrackDailyTrigger(memberId int64, triggerKey string, n int) {
	ensureGamificationHooks()
	if memberId == 0 || triggerKey == "" {
		return
	}
	if hasTrigger, cacheValid := todayHasTrigger(triggerKey); cacheValid && !hasTrigger {
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

// AwardRaffleTicket идемпотентно выдаёт один билет в сегодняшний daily-раффл
// за конкретную активность. Повторный вызов с тем же sourceType/sourceId не
// плодит билеты (UNIQUE-индекс uniq_raffle_ticket_source).
//
// Безопасен для вызова из любого хендлера: не блокирует, не паникует.
// Если daily-раффла нет (cron ещё не создал) или он не активен — no-op.
func AwardRaffleTicket(memberId int64, sourceType string, sourceId int64) {
	ensureGamificationHooks()
	if memberId == 0 || sourceType == "" {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("award-raffle-ticket panic (member=%d, source=%s/%d): %v",
					memberId, sourceType, sourceId, r)
			}
		}()
		raffle, err := dailyRaffleSvc.EnsureTodayRaffle()
		if err != nil {
			log.Printf("award-raffle-ticket ensure today (member=%d): %v", memberId, err)
			return
		}
		if raffle == nil || raffle.Status != models.RaffleStatusActive {
			return
		}
		awarded, err := raffleRepo.AwardTicketTx(database.DB, raffle.Id, memberId, sourceType, sourceId)
		if err != nil {
			log.Printf("award-raffle-ticket insert (member=%d, source=%s/%d): %v",
				memberId, sourceType, sourceId, err)
			return
		}
		if awarded {
			GetSSEHub().Publish(memberId, SSEEvent{Type: "raffles"})
		}
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
