package service

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"

	"ithozyeva/database"
	"ithozyeva/internal/models"

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
func PushDailyRaffleWin(memberId int64, raffleId int64, prize string) {
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
	_ = raffleId
}
