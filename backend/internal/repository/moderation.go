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
