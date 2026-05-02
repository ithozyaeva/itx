package repository

import (
	"ithozyeva/database"
)

// PushTarget — минимальный набор полей юзера, нужный для адресной
// отправки Telegram-пуша.
type PushTarget struct {
	MemberId   int64  `gorm:"column:id"`
	TelegramID int64  `gorm:"column:telegram_id"`
	FirstName  string `gorm:"column:first_name"`
}

type PushTargetsRepository struct{}

func NewPushTargetsRepository() *PushTargetsRepository {
	return &PushTargetsRepository{}
}

// EligibleForDailyMorning — все юзеры, которым сегодня нужно отправить
// утренний пуш. Учитывает mute_all и индивидуальный daily_morning,
// дефолт = TRUE (для юзеров без записи в notification_settings).
func (r *PushTargetsRepository) EligibleForDailyMorning() ([]PushTarget, error) {
	return r.eligibleByFlag("daily_morning")
}

// EligibleForDailyEvening — аналогично, для вечернего пуша.
func (r *PushTargetsRepository) EligibleForDailyEvening() ([]PushTarget, error) {
	return r.eligibleByFlag("daily_evening")
}

func (r *PushTargetsRepository) eligibleByFlag(flag string) ([]PushTarget, error) {
	targets := make([]PushTarget, 0)
	// COALESCE(ns.flag, TRUE) — если у юзера нет настроек, используем
	// дефолт включенного пуша. mute_all всегда дефолтит в FALSE.
	sql := `
		SELECT m.id, m.telegram_id, m.first_name
		FROM members m
		LEFT JOIN notification_settings ns ON ns.member_id = m.id
		WHERE m.telegram_id IS NOT NULL AND m.telegram_id > 0
		  AND COALESCE(ns.mute_all, FALSE) = FALSE
		  AND COALESCE(ns.` + flag + `, TRUE) = TRUE
	`
	err := database.DB.Raw(sql).Scan(&targets).Error
	return targets, err
}
