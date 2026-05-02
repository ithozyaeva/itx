package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"gorm.io/gorm"
)

// TelegramSender — функция отправки прямого сообщения юзеру.
// Подменяется из main.go реальной реализацией bot.SendDirectMessage,
// чтобы избежать цикла импортов service ↔ bot.
type TelegramSender func(chatID int64, text string)

var telegramSenderRef atomic.Pointer[TelegramSender]

// SetTelegramSender регистрирует реализацию отправителя пушей.
// Если sender не зарегистрирован — все push-функции тихо no-op.
func SetTelegramSender(send TelegramSender) {
	telegramSenderRef.Store(&send)
}

func sendTelegram(chatID int64, text string) {
	p := telegramSenderRef.Load()
	if p == nil || *p == nil {
		return
	}
	(*p)(chatID, text)
}

// pushTarget описывает цель пуша: chat-id и текущие notification_settings.
type pushTarget struct {
	chatID   int64
	settings models.NotificationSettings
}

// resolvePushTarget находит telegram-id и settings по member_id.
// Если settings отсутствует — возвращает дефолтную (всё включено).
func resolvePushTarget(memberId int64) (*pushTarget, error) {
	var member models.Member
	if err := database.DB.Select("id, telegram_id").First(&member, memberId).Error; err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}
	if member.TelegramID == 0 {
		return nil, errors.New("у пользователя нет telegram_id")
	}

	var ns models.NotificationSettings
	if err := database.DB.Where("member_id = ?", memberId).First(&ns).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ns = models.NotificationSettings{
				MemberId:     memberId,
				DailyMorning: true,
				DailyEvening: true,
				DailyStreak:  true,
				DailyRaffle:  true,
			}
		} else {
			return nil, fmt.Errorf("get notification_settings: %w", err)
		}
	}

	return &pushTarget{
		chatID:   member.TelegramID,
		settings: ns,
	}, nil
}

// PushStreakThreshold отправляет уведомление при пересечении порога стрика.
// Идемпотентность не на уровне пуша (Telegram сам не дедуплицирует) — поэтому
// вызывайте один раз на CrossedThreshold.
func PushStreakThreshold(memberId int64, days, reward int) {
	t, err := resolvePushTarget(memberId)
	if err != nil {
		log.Printf("push streak target (member=%d): %v", memberId, err)
		return
	}
	if t.settings.MuteAll || !t.settings.DailyStreak {
		return
	}
	text := fmt.Sprintf("🔥 <b>%d дней подряд!</b>\nТы заработал бонус +%d баллов за стрик. Так держать!", days, reward)
	sendTelegram(t.chatID, text)
}

// PushDailyRaffleWin уведомляет победителя ежедневного розыгрыша.
func PushDailyRaffleWin(memberId int64, prize string) {
	t, err := resolvePushTarget(memberId)
	if err != nil {
		log.Printf("push raffle win target (member=%d): %v", memberId, err)
		return
	}
	if t.settings.MuteAll || !t.settings.DailyRaffle {
		return
	}
	text := fmt.Sprintf("🎉 <b>Ты выиграл ежедневный розыгрыш!</b>\nПриз: %s начислен на твой счёт.", prize)
	sendTelegram(t.chatID, text)
}

// pushBatchPace — пауза между Telegram-сообщениями в массовой рассылке,
// чтобы не упереться в rate-limit (Telegram: ~30 msg/sec, у нас лимит
// мягче). 100мс ≈ 10 msg/sec, безопасно даже для 1000+ подписчиков
// (полная рассылка <2 минут).
const pushBatchPace = 100 * time.Millisecond

// in-memory anti-double-send. Сбрасывается при рестарте, что приемлемо
// для once-per-day рассылок (риск повторной отправки в 1 день только при
// деплое сразу после успешной рассылки — приемлемая редкость).
var (
	lastMorningPushDayUnix atomic.Int64
	lastEveningPushDayUnix atomic.Int64
)

// SendDailyMorningPush массово рассылает утренний пуш с составом сегодняшних
// дейликов и временем до ежедневного розыгрыша. Идемпотентен в пределах
// одного процесса: повторный вызов в тот же МСК-день — no-op.
func SendDailyMorningPush() {
	day := utils.MSKToday()
	dayUnix := day.Unix()
	if !lastMorningPushDayUnix.CompareAndSwap(0, dayUnix) &&
		lastMorningPushDayUnix.Load() == dayUnix {
		return
	}
	// Если предыдущий день — обновим. CAS выше обработал только
	// случай 0→day; для перехода day(N-1)→day(N) делаем явный Store.
	lastMorningPushDayUnix.Store(dayUnix)

	ensureGamificationHooks()
	taskRepo := repository.NewDailyTaskRepository()
	set, err := taskRepo.GetSetByDay(day)
	if err != nil || set == nil {
		log.Printf("morning push: today's daily set not ready: %v", err)
		return
	}
	tasks, err := taskRepo.GetByIds([]int64(set.TaskIds))
	if err != nil || len(tasks) == 0 {
		log.Printf("morning push: daily tasks lookup failed: %v", err)
		return
	}

	// Сохраняем порядок set.TaskIds.
	tasksById := make(map[int64]models.DailyTask, len(tasks))
	for _, t := range tasks {
		tasksById[t.Id] = t
	}
	titles := make([]string, 0, len(set.TaskIds))
	totalPoints := 0
	for _, id := range set.TaskIds {
		if t, ok := tasksById[id]; ok {
			titles = append(titles, fmt.Sprintf("• %s (+%d)", t.Title, t.Points))
			totalPoints += t.Points
		}
	}
	bonus := models.PointValues[models.PointReasonDailyAllTasksBonus]

	hoursToRaffle := int(time.Until(utils.MSKEndOfDay(day)).Hours())
	if hoursToRaffle < 0 {
		hoursToRaffle = 0
	}

	body := strings.Join([]string{
		"<b>Доброе утро! Сегодняшние дейлики:</b>",
		"",
		strings.Join(titles, "\n"),
		"",
		fmt.Sprintf("Заработай до %d баллов + %d бонус за все 5.", totalPoints, bonus),
		fmt.Sprintf("До ежедневного розыгрыша на 100 баллов — ~%d ч.", hoursToRaffle),
	}, "\n")

	pushRepo := repository.NewPushTargetsRepository()
	targets, err := pushRepo.EligibleForDailyMorning()
	if err != nil {
		log.Printf("morning push targets: %v", err)
		return
	}
	go broadcastPush(targets, body)
}

// SendDailyEveningPush — точечный пуш только тем, у кого сегодня сделан
// check-in, но awarded < 5: «час до розыгрыша, у тебя X/5 дейликов».
// Тех, кто не делал check-in, не теребим — иначе пуш превращается в спам.
func SendDailyEveningPush() {
	day := utils.MSKToday()
	dayUnix := day.Unix()
	if !lastEveningPushDayUnix.CompareAndSwap(0, dayUnix) &&
		lastEveningPushDayUnix.Load() == dayUnix {
		return
	}
	lastEveningPushDayUnix.Store(dayUnix)

	ensureGamificationHooks()
	pushRepo := repository.NewPushTargetsRepository()
	targets, err := pushRepo.EligibleForDailyEvening()
	if err != nil {
		log.Printf("evening push targets: %v", err)
		return
	}
	if len(targets) == 0 {
		return
	}

	// Подгрузим набор задач сегодня одной выборкой.
	taskRepo := repository.NewDailyTaskRepository()
	set, err := taskRepo.GetSetByDay(day)
	if err != nil || set == nil {
		return
	}
	totalTasks := len(set.TaskIds)

	// Batch: один SQL для check-ins + один для awarded-counts. Раньше
	// каждый юзер генерировал 2 SQL — на 1000 подписчиков это 2000
	// запросов подряд. Теперь — 2.
	checkInRepo := repository.NewCheckInRepository()
	checkedIn, err := checkInRepo.CheckedInMembers(day)
	if err != nil {
		log.Printf("evening push checked-in members: %v", err)
		return
	}
	awardedCounts, err := taskRepo.AwardedCountsForDay(day, []int64(set.TaskIds))
	if err != nil {
		log.Printf("evening push awarded counts: %v", err)
		return
	}

	go func() {
		for _, t := range targets {
			if !checkedIn[t.MemberId] {
				continue
			}
			awarded := awardedCounts[t.MemberId]
			if awarded >= totalTasks {
				continue
			}
			text := fmt.Sprintf(
				"⏰ Час до ежедневного розыгрыша!\nТы выполнил <b>%d/%d</b> дейликов сегодня — успей ещё.",
				awarded, totalTasks,
			)
			sendTelegram(t.TelegramID, text)
			time.Sleep(pushBatchPace)
		}
	}()
}

// broadcastPush массово отправляет одинаковый текст всем targets,
// соблюдая pushBatchPace, чтобы не словить rate-limit.
func broadcastPush(targets []repository.PushTarget, text string) {
	for _, t := range targets {
		sendTelegram(t.TelegramID, text)
		time.Sleep(pushBatchPace)
	}
}
