package repository

import (
	"errors"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModerationRepository struct{}

func NewModerationRepository() *ModerationRepository {
	return &ModerationRepository{}
}

// LogAction записывает модерационное действие.
func (r *ModerationRepository) LogAction(action *models.ModerationAction) error {
	if action.Meta == "" {
		action.Meta = "{}"
	}
	return database.DB.Create(action).Error
}

// MessagesForCleanup возвращает telegram_message_id сообщений пользователя в
// чате за период since (>=). Только записи с непустым telegram_message_id —
// старые сообщения без него Telegram API удалить не позволит.
func (r *ModerationRepository) MessagesForCleanup(chatID, userID int64, since time.Time) ([]int, error) {
	var ids []int
	err := database.DB.Model(&models.ChatMessage{}).
		Where("chat_id = ? AND telegram_user_id = ? AND sent_at >= ? AND telegram_message_id IS NOT NULL",
			chatID, userID, since).
		Order("sent_at ASC").
		Pluck("telegram_message_id", &ids).Error
	return ids, err
}

// DeleteCleanedMessages удаляет записи из chat_messages по списку telegram_message_id
// (после успешного удаления в Telegram). Возвращает кол-во удалённых строк.
func (r *ModerationRepository) DeleteCleanedMessages(chatID int64, messageIDs []int) (int64, error) {
	if len(messageIDs) == 0 {
		return 0, nil
	}
	res := database.DB.Where("chat_id = ? AND telegram_message_id IN ?", chatID, messageIDs).
		Delete(&models.ChatMessage{})
	return res.RowsAffected, res.Error
}

// FindOpenVoteban — открытое голосование на (chat_id, target_user_id), если есть.
func (r *ModerationRepository) FindOpenVoteban(chatID, targetUserID int64) (*models.Voteban, error) {
	var v models.Voteban
	err := database.DB.Where("chat_id = ? AND target_user_id = ? AND status = ?",
		chatID, targetUserID, models.VotebanStatusOpen).
		First(&v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

// CreateVoteban сохраняет новое голосование.
func (r *ModerationRepository) CreateVoteban(v *models.Voteban) error {
	return database.DB.Create(v).Error
}

// GetVoteban читает запись по id.
func (r *ModerationRepository) GetVoteban(id int64) (*models.Voteban, error) {
	var v models.Voteban
	if err := database.DB.First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// UpsertVote ставит/обновляет голос. Возвращает true, если запись создалась
// впервые (для логов / UX), и финальное значение голоса.
func (r *ModerationRepository) UpsertVote(votebanID, voterID int64, vote int16) error {
	now := time.Now()
	v := models.VotebanVote{
		VotebanID:   votebanID,
		VoterUserID: voterID,
		Vote:        vote,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "voteban_id"}, {Name: "voter_user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"vote":       vote,
			"updated_at": now,
		}),
	}).Create(&v).Error
}

// GetVote возвращает текущий голос пользователя или nil, если не голосовал.
func (r *ModerationRepository) GetVote(votebanID, voterID int64) (*int16, error) {
	var v models.VotebanVote
	err := database.DB.Where("voteban_id = ? AND voter_user_id = ?", votebanID, voterID).
		First(&v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &v.Vote, nil
}

// CountVotes считает «за» и «против» по voteban_id.
func (r *ModerationRepository) CountVotes(votebanID int64) (models.VotebanTally, error) {
	var tally models.VotebanTally
	type row struct {
		Vote  int16
		Count int
	}
	var rows []row
	err := database.DB.Model(&models.VotebanVote{}).
		Select("vote, COUNT(*) as count").
		Where("voteban_id = ?", votebanID).
		Group("vote").
		Scan(&rows).Error
	if err != nil {
		return tally, err
	}
	for _, r := range rows {
		switch r.Vote {
		case models.VotebanVoteFor:
			tally.For = r.Count
		case models.VotebanVoteAgainst:
			tally.Against = r.Count
		}
	}
	return tally, nil
}

// FinalizeVoteban переводит голосование в финальный статус.
func (r *ModerationRepository) FinalizeVoteban(id int64, status string) error {
	now := time.Now()
	return database.DB.Model(&models.Voteban{}).
		Where("id = ? AND status = ?", id, models.VotebanStatusOpen).
		Updates(map[string]interface{}{
			"status":       status,
			"finalized_at": now,
		}).Error
}

// ListExpiredOpenVotebans возвращает голосования с истёкшим окном для finalize.
func (r *ModerationRepository) ListExpiredOpenVotebans(now time.Time) ([]models.Voteban, error) {
	var list []models.Voteban
	err := database.DB.Where("status = ? AND expires_at <= ?", models.VotebanStatusOpen, now).
		Find(&list).Error
	return list, err
}

// LatestVotebanCreatedInChat возвращает время последнего созданного voteban в
// чате (любой статус), или nil если голосований ещё не было. Используется
// для cooldown «не чаще одного /voteban в чате раз в N минут».
func (r *ModerationRepository) LatestVotebanCreatedInChat(chatID int64) (*time.Time, error) {
	var t time.Time
	err := database.DB.Model(&models.Voteban{}).
		Select("created_at").
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(1).
		Scan(&t).Error
	if err != nil {
		return nil, err
	}
	if t.IsZero() {
		return nil, nil
	}
	return &t, nil
}

// LatestVotebanCreatedByInitiator возвращает время последнего voteban в чате,
// запущенного конкретным юзером.
func (r *ModerationRepository) LatestVotebanCreatedByInitiator(chatID, initiatorID int64) (*time.Time, error) {
	var t time.Time
	err := database.DB.Model(&models.Voteban{}).
		Select("created_at").
		Where("chat_id = ? AND initiator_user_id = ?", chatID, initiatorID).
		Order("created_at DESC").
		Limit(1).
		Scan(&t).Error
	if err != nil {
		return nil, err
	}
	if t.IsZero() {
		return nil, nil
	}
	return &t, nil
}

// ListExpiredUnnotifiedActions возвращает санкции, у которых срок истёк, а
// уведомление в чат(ы) ещё не отправлено. Watcher шлёт алерт и помечает запись.
func (r *ModerationRepository) ListExpiredUnnotifiedActions(now time.Time) ([]models.ModerationAction, error) {
	var list []models.ModerationAction
	err := database.DB.
		Where("action IN ? AND expires_at IS NOT NULL AND expires_at <= ? AND expired_notified_at IS NULL",
			models.ModerationActionsWithExpiry, now).
		Order("expires_at ASC").
		Limit(100).
		Find(&list).Error
	return list, err
}

// MarkActionExpiredNotified ставит expired_notified_at = NOW() для записи.
func (r *ModerationRepository) MarkActionExpiredNotified(id int64) error {
	now := time.Now()
	return database.DB.Model(&models.ModerationAction{}).
		Where("id = ? AND expired_notified_at IS NULL", id).
		Update("expired_notified_at", now).Error
}

// --- Global bans ---

// UpsertGlobalBan создаёт/обновляет запись (по PK user_id).
func (r *ModerationRepository) UpsertGlobalBan(b *models.GlobalBan) error {
	now := time.Now()
	b.UpdatedAt = now
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	return database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"banned_by", "reason", "expires_at", "updated_at",
		}),
	}).Create(b).Error
}

// GetGlobalBan возвращает запись или (nil, nil) если её нет.
func (r *ModerationRepository) GetGlobalBan(userID int64) (*models.GlobalBan, error) {
	var b models.GlobalBan
	err := database.DB.Where("user_id = ?", userID).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// DeleteGlobalBan удаляет запись (используется при /globalunban).
func (r *ModerationRepository) DeleteGlobalBan(userID int64) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.GlobalBan{}).Error
}

// ListActiveGlobalBans возвращает все активные баны (не истёкшие).
func (r *ModerationRepository) ListActiveGlobalBans(now time.Time) ([]models.GlobalBan, error) {
	var list []models.GlobalBan
	err := database.DB.Where("expires_at IS NULL OR expires_at > ?", now).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

// --- Admin UI listing с обогащением ---

// ModerationActionView — строка для админки: action + display-поля
// (username/first_name + chat_title), полученные join'ами.
type ModerationActionView struct {
	models.ModerationAction
	TargetUsername  string `json:"targetUsername" gorm:"column:target_username"`
	TargetFirstName string `json:"targetFirstName" gorm:"column:target_first_name"`
	ChatTitle       string `json:"chatTitle" gorm:"column:chat_title"`
}

// VotebanView — voteban + раскладка голосов.
type VotebanView struct {
	models.Voteban
	VotesFor     int `json:"votesFor" gorm:"column:votes_for"`
	VotesAgainst int `json:"votesAgainst" gorm:"column:votes_against"`
}

// listEnrichedActions возвращает список с username/first_name из последнего
// сообщения юзера в чате (если есть) и chat_title из tracked_chats или
// subscription_chats. Сортировка по created_at DESC.
func (r *ModerationRepository) listEnrichedActions(where string, args ...interface{}) ([]ModerationActionView, error) {
	var rows []ModerationActionView
	q := `
		SELECT
			a.*,
			COALESCE((
				SELECT cm.telegram_username FROM chat_messages cm
				WHERE cm.telegram_user_id = a.target_user_id
				  AND cm.telegram_username <> ''
				ORDER BY cm.sent_at DESC LIMIT 1
			), '') AS target_username,
			COALESCE((
				SELECT cm.telegram_first_name FROM chat_messages cm
				WHERE cm.telegram_user_id = a.target_user_id
				  AND cm.telegram_first_name <> ''
				ORDER BY cm.sent_at DESC LIMIT 1
			), '') AS target_first_name,
			COALESCE((
				SELECT title FROM tracked_chats WHERE chat_id = a.chat_id LIMIT 1
			), (
				SELECT title FROM subscription_chats WHERE id = a.chat_id LIMIT 1
			), '') AS chat_title
		FROM bot_moderation_actions a
		WHERE ` + where + `
		ORDER BY a.created_at DESC
		LIMIT 200
	`
	err := database.DB.Raw(q, args...).Scan(&rows).Error
	return rows, err
}

// ListActiveSanctions возвращает действующие санкции (ban/mute/voteban_*),
// у которых expires_at IS NULL OR > NOW(). Список типов задаётся строками,
// а не константами, чтобы покрыть и legacy `voteban_mute`, и будущий
// `voteban_kick` (#295) без жёсткого зависания на других PR.
func (r *ModerationRepository) ListActiveSanctions() ([]ModerationActionView, error) {
	return r.listEnrichedActions(
		"a.action IN ('ban','mute','voteban_mute','voteban_kick') AND (a.expires_at IS NULL OR a.expires_at > NOW())",
	)
}

// ListRecentActions — лог последних 200 модерационных действий (любой тип).
func (r *ModerationRepository) ListRecentActions() ([]ModerationActionView, error) {
	return r.listEnrichedActions("1=1")
}

// GetActionByID — действие по id с обогащением.
func (r *ModerationRepository) GetActionByID(id int64) (*ModerationActionView, error) {
	rows, err := r.listEnrichedActions("a.id = ?", id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

// ListOpenVotebansEnriched — открытые голосования с раскладкой голосов.
func (r *ModerationRepository) ListOpenVotebansEnriched() ([]VotebanView, error) {
	var rows []VotebanView
	err := database.DB.Raw(`
		SELECT
			v.*,
			COALESCE((SELECT COUNT(*) FROM bot_voteban_votes WHERE voteban_id = v.id AND vote = 1), 0)  AS votes_for,
			COALESCE((SELECT COUNT(*) FROM bot_voteban_votes WHERE voteban_id = v.id AND vote = -1), 0) AS votes_against
		FROM bot_votebans v
		WHERE v.status = ?
		ORDER BY v.created_at DESC
	`, models.VotebanStatusOpen).Scan(&rows).Error
	return rows, err
}

// KnownChatIDs возвращает уникальные chat_id из subscription_chats и активных
// tracked_chats — тот набор, по которому проходим при глобальном бане/анбане.
func (r *ModerationRepository) KnownChatIDs() ([]int64, error) {
	var ids []int64
	err := database.DB.Raw(`
		SELECT DISTINCT chat_id FROM (
			SELECT id AS chat_id FROM subscription_chats
			UNION
			SELECT chat_id FROM tracked_chats WHERE is_active = TRUE
		) c
	`).Pluck("chat_id", &ids).Error
	return ids, err
}
