package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"
)

type ChatQuestRepository struct{}

func NewChatQuestRepository() *ChatQuestRepository {
	return &ChatQuestRepository{}
}

// CreateQuest создаёт новое задание
func (r *ChatQuestRepository) CreateQuest(quest *models.ChatQuest) error {
	return database.DB.Create(quest).Error
}

// UpdateQuest обновляет задание
func (r *ChatQuestRepository) UpdateQuest(quest *models.ChatQuest) error {
	return database.DB.Save(quest).Error
}

// DeleteQuest удаляет задание
func (r *ChatQuestRepository) DeleteQuest(id int64) error {
	return database.DB.Delete(&models.ChatQuest{}, id).Error
}

// GetActiveQuests возвращает активные задания (starts_at <= now <= ends_at, is_active)
func (r *ChatQuestRepository) GetActiveQuests() ([]models.ChatQuest, error) {
	var quests []models.ChatQuest
	now := time.Now()
	err := database.DB.Where("is_active = ? AND starts_at <= ? AND ends_at >= ?", true, now, now).
		Order("ends_at ASC").
		Find(&quests).Error
	return quests, err
}

// GetQuestsForMember возвращает активные квесты с прогрессом для участника
func (r *ChatQuestRepository) GetQuestsForMember(memberID int64) ([]models.ChatQuestWithProgress, error) {
	var result []models.ChatQuestWithProgress
	now := time.Now()
	err := database.DB.Raw(`
		SELECT cq.*,
		       COALESCE(cqp.current_count, 0) as current_count,
		       COALESCE(cqp.completed, false) as completed
		FROM chat_quests cq
		LEFT JOIN chat_quest_progress cqp ON cqp.quest_id = cq.id AND cqp.member_id = ?
		WHERE cq.is_active = true AND cq.starts_at <= ? AND cq.ends_at >= ?
		ORDER BY cqp.completed ASC NULLS FIRST, cq.ends_at ASC
	`, memberID, now, now).Scan(&result).Error
	return result, err
}

// GetAllQuests возвращает все квесты для админки с пагинацией
func (r *ChatQuestRepository) GetAllQuests(limit, offset int) ([]models.ChatQuest, int64, error) {
	var quests []models.ChatQuest
	var total int64

	database.DB.Model(&models.ChatQuest{}).Count(&total)

	err := database.DB.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&quests).Error
	return quests, total, err
}

// GetQuestByID возвращает квест по ID
func (r *ChatQuestRepository) GetQuestByID(id int64) (*models.ChatQuest, error) {
	var quest models.ChatQuest
	err := database.DB.First(&quest, id).Error
	return &quest, err
}

// IncrementProgress увеличивает прогресс на 1, возвращает новый счётчик и признак завершения
func (r *ChatQuestRepository) IncrementProgress(questID int64, memberID int64) (newCount int, targetCount int, alreadyCompleted bool, err error) {
	// Проверяем, есть ли уже прогресс
	var progress models.ChatQuestProgress
	result := database.DB.Where("quest_id = ? AND member_id = ?", questID, memberID).First(&progress)

	if result.Error != nil {
		// Создаём новый прогресс
		progress = models.ChatQuestProgress{
			QuestID:      questID,
			MemberID:     memberID,
			CurrentCount: 1,
		}
		if err = database.DB.Create(&progress).Error; err != nil {
			return
		}
		newCount = 1
	} else {
		if progress.Completed {
			alreadyCompleted = true
			newCount = progress.CurrentCount
			return
		}
		// Инкрементируем
		progress.CurrentCount++
		if err = database.DB.Save(&progress).Error; err != nil {
			return
		}
		newCount = progress.CurrentCount
	}

	// Получаем target_count из квеста
	var quest models.ChatQuest
	if err = database.DB.First(&quest, questID).Error; err != nil {
		return
	}
	targetCount = quest.TargetCount
	return
}

// MarkCompleted помечает квест выполненным для участника
func (r *ChatQuestRepository) MarkCompleted(questID int64, memberID int64) error {
	now := time.Now()
	return database.DB.Model(&models.ChatQuestProgress{}).
		Where("quest_id = ? AND member_id = ?", questID, memberID).
		Updates(map[string]interface{}{
			"completed":    true,
			"completed_at": now,
		}).Error
}

// GetAllQuestsForMember возвращает все квесты (активные + завершённые) с прогрессом для участника
func (r *ChatQuestRepository) GetAllQuestsForMember(memberID int64, filter string) ([]models.ChatQuestWithProgress, error) {
	var result []models.ChatQuestWithProgress
	now := time.Now()

	var query string
	switch filter {
	case "active":
		query = `
			SELECT cq.*,
			       COALESCE(cqp.current_count, 0) as current_count,
			       COALESCE(cqp.completed, false) as completed
			FROM chat_quests cq
			LEFT JOIN chat_quest_progress cqp ON cqp.quest_id = cq.id AND cqp.member_id = ?
			WHERE cq.is_active = true AND cq.starts_at <= ? AND cq.ends_at >= ?
			  AND COALESCE(cqp.completed, false) = false
			ORDER BY cq.ends_at ASC
		`
		err := database.DB.Raw(query, memberID, now, now).Scan(&result).Error
		return result, err
	case "completed":
		query = `
			SELECT cq.*,
			       COALESCE(cqp.current_count, 0) as current_count,
			       COALESCE(cqp.completed, false) as completed
			FROM chat_quests cq
			INNER JOIN chat_quest_progress cqp ON cqp.quest_id = cq.id AND cqp.member_id = ?
			WHERE cqp.completed = true
			ORDER BY cqp.completed_at DESC
		`
		err := database.DB.Raw(query, memberID).Scan(&result).Error
		return result, err
	default: // "all"
		query = `
			SELECT cq.*,
			       COALESCE(cqp.current_count, 0) as current_count,
			       COALESCE(cqp.completed, false) as completed
			FROM chat_quests cq
			LEFT JOIN chat_quest_progress cqp ON cqp.quest_id = cq.id AND cqp.member_id = ?
			WHERE cq.is_active = true AND (cq.starts_at <= ? AND cq.ends_at >= ? OR COALESCE(cqp.completed, false) = true)
			ORDER BY cqp.completed ASC NULLS FIRST, cq.ends_at ASC
		`
		err := database.DB.Raw(query, memberID, now, now).Scan(&result).Error
		return result, err
	}
}

// RecordStreakDay записывает день активности для daily_streak квеста (INSERT ON CONFLICT DO NOTHING)
func (r *ChatQuestRepository) RecordStreakDay(questID int64, memberID int64, day time.Time) error {
	return database.DB.Exec(`
		INSERT INTO chat_quest_streak_days (quest_id, member_id, day)
		VALUES (?, ?, ?)
		ON CONFLICT (quest_id, member_id, day) DO NOTHING
	`, questID, memberID, day.Format("2006-01-02")).Error
}

// GetCurrentStreak считает текущую серию последовательных дней от сегодня назад
func (r *ChatQuestRepository) GetCurrentStreak(questID int64, memberID int64) (int, error) {
	var streak int
	err := database.DB.Raw(`
		WITH days AS (
			SELECT day, ROW_NUMBER() OVER (ORDER BY day DESC) as rn
			FROM chat_quest_streak_days
			WHERE quest_id = ? AND member_id = ?
			ORDER BY day DESC
		)
		SELECT COUNT(*) FROM days
		WHERE day = CURRENT_DATE - (rn - 1) * INTERVAL '1 day'
	`, questID, memberID).Scan(&streak).Error
	return streak, err
}

// GetProgress возвращает прогресс квеста для участника
func (r *ChatQuestRepository) GetProgress(questID int64, memberID int64, progress *models.ChatQuestProgress) error {
	return database.DB.Where("quest_id = ? AND member_id = ?", questID, memberID).First(progress).Error
}

// SetProgressCount устанавливает current_count для прогресса квеста (не инкремент)
func (r *ChatQuestRepository) SetProgressCount(questID int64, memberID int64, count int) error {
	var progress models.ChatQuestProgress
	result := database.DB.Where("quest_id = ? AND member_id = ?", questID, memberID).First(&progress)

	if result.Error != nil {
		// Создаём новый прогресс
		progress = models.ChatQuestProgress{
			QuestID:      questID,
			MemberID:     memberID,
			CurrentCount: count,
		}
		return database.DB.Create(&progress).Error
	}

	progress.CurrentCount = count
	return database.DB.Save(&progress).Error
}
