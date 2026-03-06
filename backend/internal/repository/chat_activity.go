package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"
)

type ChatActivityRepository struct{}

func NewChatActivityRepository() *ChatActivityRepository {
	return &ChatActivityRepository{}
}

// SaveMessage сохраняет сообщение в БД
func (r *ChatActivityRepository) SaveMessage(msg *models.ChatMessage) error {
	return database.DB.Create(msg).Error
}

// GetTrackedChatIDs возвращает список активных chat_id
func (r *ChatActivityRepository) GetTrackedChatIDs() ([]int64, error) {
	var chatIDs []int64
	err := database.DB.Model(&models.TrackedChat{}).
		Where("is_active = ?", true).
		Pluck("chat_id", &chatIDs).Error
	return chatIDs, err
}

// GetTrackedChats возвращает все отслеживаемые чаты
func (r *ChatActivityRepository) GetTrackedChats() ([]models.TrackedChat, error) {
	var chats []models.TrackedChat
	err := database.DB.Where("is_active = ?", true).Order("title").Find(&chats).Error
	return chats, err
}

// GetMemberIDsByChatID возвращает member_id участников, которые писали в указанном чате
func (r *ChatActivityRepository) GetMemberIDsByChatID(chatID int64) ([]int64, error) {
	var memberIDs []int64
	err := database.DB.Raw(`
		SELECT DISTINCT cm.member_id
		FROM chat_messages cm
		WHERE cm.chat_id = ? AND cm.member_id IS NOT NULL
	`, chatID).Pluck("member_id", &memberIDs).Error
	return memberIDs, err
}

// GetMessageCountsByChat возвращает количество сообщений по чатам за период
func (r *ChatActivityRepository) GetMessageCountsByChat(from, to time.Time) ([]models.ChatMessageCount, error) {
	var counts []models.ChatMessageCount
	err := database.DB.Raw(`
		SELECT cm.chat_id, tc.title, COUNT(*) as count
		FROM chat_messages cm
		JOIN tracked_chats tc ON tc.chat_id = cm.chat_id
		WHERE cm.sent_at >= ? AND cm.sent_at < ?
		GROUP BY cm.chat_id, tc.title
		ORDER BY count DESC
	`, from, to).Scan(&counts).Error
	return counts, err
}

// GetDailyActivity возвращает количество сообщений по дням для графика
func (r *ChatActivityRepository) GetDailyActivity(chatID *int64, days int) ([]models.DailyActivity, error) {
	var activity []models.DailyActivity
	var err error

	if chatID != nil {
		err = database.DB.Raw(`
			SELECT TO_CHAR(d.day, 'YYYY-MM-DD') as date,
			       COALESCE(cnt, 0) as count
			FROM generate_series(
			  CURRENT_DATE - ?::int * INTERVAL '1 day',
			  CURRENT_DATE,
			  '1 day'::interval
			) d(day)
			LEFT JOIN (
			  SELECT DATE(sent_at) as day, COUNT(*) as cnt
			  FROM chat_messages
			  WHERE chat_id = ? AND sent_at >= CURRENT_DATE - ?::int * INTERVAL '1 day'
			  GROUP BY DATE(sent_at)
			) cm ON cm.day = d.day::date
			ORDER BY d.day
		`, days, *chatID, days).Scan(&activity).Error
	} else {
		err = database.DB.Raw(`
			SELECT TO_CHAR(d.day, 'YYYY-MM-DD') as date,
			       COALESCE(cnt, 0) as count
			FROM generate_series(
			  CURRENT_DATE - ?::int * INTERVAL '1 day',
			  CURRENT_DATE,
			  '1 day'::interval
			) d(day)
			LEFT JOIN (
			  SELECT DATE(sent_at) as day, COUNT(*) as cnt
			  FROM chat_messages
			  WHERE sent_at >= CURRENT_DATE - ?::int * INTERVAL '1 day'
			  GROUP BY DATE(sent_at)
			) cm ON cm.day = d.day::date
			ORDER BY d.day
		`, days, days).Scan(&activity).Error
	}

	return activity, err
}

// GetTopUsers возвращает топ пользователей по количеству сообщений
func (r *ChatActivityRepository) GetTopUsers(from, to time.Time, limit int) ([]models.TopUser, error) {
	var users []models.TopUser
	err := database.DB.Raw(`
		WITH user_counts AS (
			SELECT telegram_user_id, telegram_username, telegram_first_name, COUNT(*) as count
			FROM chat_messages
			WHERE sent_at >= ? AND sent_at < ?
			GROUP BY telegram_user_id, telegram_username, telegram_first_name
			ORDER BY count DESC
			LIMIT ?
		),
		user_top_chat AS (
			SELECT DISTINCT ON (uc.telegram_user_id)
				uc.telegram_user_id,
				tc.title as top_chat
			FROM user_counts uc
			JOIN chat_messages cm ON cm.telegram_user_id = uc.telegram_user_id
				AND cm.sent_at >= ? AND cm.sent_at < ?
			JOIN tracked_chats tc ON tc.chat_id = cm.chat_id
			GROUP BY uc.telegram_user_id, tc.title
			ORDER BY uc.telegram_user_id, COUNT(*) DESC
		)
		SELECT uc.telegram_user_id, uc.telegram_username, uc.telegram_first_name, uc.count,
		       COALESCE(utc.top_chat, '') as top_chat
		FROM user_counts uc
		LEFT JOIN user_top_chat utc ON utc.telegram_user_id = uc.telegram_user_id
		ORDER BY uc.count DESC
	`, from, to, limit, from, to).Scan(&users).Error
	return users, err
}

// GetTotalStats возвращает общую статистику
func (r *ChatActivityRepository) GetTotalStats(from, to time.Time) (totalMessages int64, uniqueUsers int64, err error) {
	err = database.DB.Raw(`
		SELECT COUNT(*) FROM chat_messages WHERE sent_at >= ? AND sent_at < ?
	`, from, to).Scan(&totalMessages).Error
	if err != nil {
		return
	}

	err = database.DB.Raw(`
		SELECT COUNT(DISTINCT telegram_user_id) FROM chat_messages WHERE sent_at >= ? AND sent_at < ?
	`, from, to).Scan(&uniqueUsers).Error
	return
}

// GetUserStats возвращает статистику конкретного пользователя
func (r *ChatActivityRepository) GetUserStats(userID int64, days int) (*models.UserStats, error) {
	var stats models.UserStats
	err := database.DB.Raw(`
		SELECT
			telegram_user_id,
			MAX(telegram_username) as telegram_username,
			MAX(telegram_first_name) as telegram_first_name,
			COUNT(*) as total_messages,
			COUNT(DISTINCT chat_id) as active_chats,
			ROUND(COUNT(*)::numeric / GREATEST(?, 1), 1) as avg_per_day
		FROM chat_messages
		WHERE telegram_user_id = ? AND sent_at >= CURRENT_DATE - ?::int * INTERVAL '1 day'
		GROUP BY telegram_user_id
	`, days, userID, days).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetDailyActivityByUser возвращает активность конкретного пользователя по дням
func (r *ChatActivityRepository) GetDailyActivityByUser(userID int64, days int) ([]models.DailyActivity, error) {
	var activity []models.DailyActivity
	err := database.DB.Raw(`
		SELECT TO_CHAR(d.day, 'YYYY-MM-DD') as date,
		       COALESCE(cnt, 0) as count
		FROM generate_series(
		  CURRENT_DATE - ?::int * INTERVAL '1 day',
		  CURRENT_DATE,
		  '1 day'::interval
		) d(day)
		LEFT JOIN (
		  SELECT DATE(sent_at) as day, COUNT(*) as cnt
		  FROM chat_messages
		  WHERE telegram_user_id = ? AND sent_at >= CURRENT_DATE - ?::int * INTERVAL '1 day'
		  GROUP BY DATE(sent_at)
		) cm ON cm.day = d.day::date
		ORDER BY d.day
	`, days, userID, days).Scan(&activity).Error
	return activity, err
}

// GetMessagesForExport возвращает агрегированные данные для CSV экспорта
func (r *ChatActivityRepository) GetMessagesForExport(chatID *int64, days int) ([]models.ExportRow, error) {
	var rows []models.ExportRow
	var err error

	if chatID != nil {
		err = database.DB.Raw(`
			SELECT TO_CHAR(DATE(cm.sent_at), 'YYYY-MM-DD') as date,
			       tc.title as chat_title,
			       COALESCE(cm.telegram_username, cm.telegram_first_name) as telegram_username,
			       COUNT(*) as message_count
			FROM chat_messages cm
			JOIN tracked_chats tc ON tc.chat_id = cm.chat_id
			WHERE cm.chat_id = ? AND cm.sent_at >= CURRENT_DATE - ?::int * INTERVAL '1 day'
			GROUP BY DATE(cm.sent_at), tc.title, COALESCE(cm.telegram_username, cm.telegram_first_name)
			ORDER BY date DESC, message_count DESC
		`, *chatID, days).Scan(&rows).Error
	} else {
		err = database.DB.Raw(`
			SELECT TO_CHAR(DATE(cm.sent_at), 'YYYY-MM-DD') as date,
			       tc.title as chat_title,
			       COALESCE(cm.telegram_username, cm.telegram_first_name) as telegram_username,
			       COUNT(*) as message_count
			FROM chat_messages cm
			JOIN tracked_chats tc ON tc.chat_id = cm.chat_id
			WHERE cm.sent_at >= CURRENT_DATE - ?::int * INTERVAL '1 day'
			GROUP BY DATE(cm.sent_at), tc.title, COALESCE(cm.telegram_username, cm.telegram_first_name)
			ORDER BY date DESC, message_count DESC
		`, days).Scan(&rows).Error
	}

	return rows, err
}

// DeleteOldMessages удаляет сообщения старше указанной даты
func (r *ChatActivityRepository) DeleteOldMessages(beforeDate time.Time) (int64, error) {
	result := database.DB.Where("sent_at < ?", beforeDate).Delete(&models.ChatMessage{})
	return result.RowsAffected, result.Error
}
