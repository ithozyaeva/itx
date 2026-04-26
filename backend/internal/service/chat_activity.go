package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatActivityService struct {
	repo           *repository.ChatActivityRepository
	memberRepo     *repository.MemberRepository
	questService   *ChatQuestService
	trackedChatIDs map[int64]bool
	memberIDCache  map[int64]int64 // telegram_user_id -> member_id
	mu             sync.RWMutex
	memberMu       sync.RWMutex
}

func NewChatActivityService() *ChatActivityService {
	s := &ChatActivityService{
		repo:           repository.NewChatActivityRepository(),
		memberRepo:     repository.NewMemberRepository(),
		questService:   NewChatQuestService(),
		trackedChatIDs: make(map[int64]bool),
		memberIDCache:  make(map[int64]int64),
	}
	s.loadTrackedChats()
	return s
}

// loadTrackedChats загружает список отслеживаемых чатов в память
func (s *ChatActivityService) loadTrackedChats() {
	chatIDs, err := s.repo.GetTrackedChatIDs()
	if err != nil {
		log.Printf("Error loading tracked chat IDs: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range chatIDs {
		s.trackedChatIDs[id] = true
	}
	log.Printf("Loaded %d tracked chats for activity monitoring", len(s.trackedChatIDs))
}

// resolveMemberID находит member_id по telegram_user_id с кешированием
func (s *ChatActivityService) resolveMemberID(telegramUserID int64) *int64 {
	// Проверяем кеш
	s.memberMu.RLock()
	if memberID, ok := s.memberIDCache[telegramUserID]; ok {
		s.memberMu.RUnlock()
		return &memberID
	}
	s.memberMu.RUnlock()

	// Ищем в БД
	member, err := s.memberRepo.GetByTelegramID(telegramUserID)
	if err != nil {
		return nil
	}

	// Кешируем
	s.memberMu.Lock()
	s.memberIDCache[telegramUserID] = member.Id
	s.memberMu.Unlock()

	return &member.Id
}

// IsTrackedChat проверяет, является ли чат отслеживаемым (группа ITX)
func (s *ChatActivityService) IsTrackedChat(chatID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.trackedChatIDs[chatID]
}

// TrackMessage проверяет, что чат отслеживается, и сохраняет сообщение
func (s *ChatActivityService) TrackMessage(message *tgbotapi.Message) {
	if message == nil || message.From == nil {
		return
	}

	s.mu.RLock()
	tracked := s.trackedChatIDs[message.Chat.ID]
	s.mu.RUnlock()

	if !tracked {
		return
	}

	// Определяем member_id
	memberID := s.resolveMemberID(message.From.ID)

	tgMessageID := message.MessageID
	msg := &models.ChatMessage{
		ChatID:            message.Chat.ID,
		TelegramUserID:    message.From.ID,
		TelegramUsername:  message.From.UserName,
		TelegramFirstName: message.From.FirstName,
		MemberID:          memberID,
		MessageText:       message.Text,
		TelegramMessageID: &tgMessageID,
		SentAt:            time.Unix(int64(message.Date), 0),
	}

	if err := s.repo.SaveMessage(msg); err != nil {
		log.Printf("Error saving chat message: %v", err)
		return
	}

	// Обрабатываем квесты
	go s.questService.ProcessMessage(message, memberID)
}

// GetRecentMessages возвращает последние N сообщений с текстом из чата
func (s *ChatActivityService) GetRecentMessages(chatID int64, limit int) ([]models.ChatMessage, error) {
	return s.repo.GetRecentMessages(chatID, limit)
}

// GetMessagesSince возвращает сообщения из чата начиная с указанного времени
func (s *ChatActivityService) GetMessagesSince(chatID int64, since time.Time) ([]models.ChatMessage, error) {
	return s.repo.GetMessagesSince(chatID, since)
}

// GetMemberIDsByChatID возвращает member_id участников конкретного чата
func (s *ChatActivityService) GetMemberIDsByChatID(chatID int64) ([]int64, error) {
	return s.repo.GetMemberIDsByChatID(chatID)
}

// GetStats возвращает общую статистику активности с данными за предыдущую неделю
func (s *ChatActivityService) GetStats() (*models.ChatActivityStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -7)
	prevWeekStart := todayStart.AddDate(0, 0, -14)
	tomorrow := todayStart.AddDate(0, 0, 1)

	totalToday, uniqueToday, err := s.repo.GetTotalStats(todayStart, tomorrow)
	if err != nil {
		return nil, err
	}

	totalWeek, uniqueWeek, err := s.repo.GetTotalStats(weekStart, tomorrow)
	if err != nil {
		return nil, err
	}

	// Статистика за предыдущую неделю для сравнения
	totalLastWeek, uniqueLastWeek, err := s.repo.GetTotalStats(prevWeekStart, weekStart)
	if err != nil {
		return nil, err
	}

	chatStats, err := s.repo.GetMessageCountsByChat(weekStart, tomorrow)
	if err != nil {
		return nil, err
	}

	return &models.ChatActivityStats{
		TotalMessagesToday:    totalToday,
		TotalMessagesWeek:     totalWeek,
		UniqueUsersToday:      uniqueToday,
		UniqueUsersWeek:       uniqueWeek,
		TotalMessagesLastWeek: totalLastWeek,
		UniqueUsersLastWeek:   uniqueLastWeek,
		ChatStats:             chatStats,
	}, nil
}

// GetActivityChart возвращает данные для графика активности
func (s *ChatActivityService) GetActivityChart(chatID *int64, days int) ([]models.DailyActivity, error) {
	return s.repo.GetDailyActivity(chatID, days)
}

// GetTopUsers возвращает топ пользователей
func (s *ChatActivityService) GetTopUsers(days int, limit int) ([]models.TopUser, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	from := todayStart.AddDate(0, 0, -days)
	to := todayStart.AddDate(0, 0, 1)
	return s.repo.GetTopUsers(from, to, limit)
}

// GetTrackedChats возвращает список отслеживаемых чатов
func (s *ChatActivityService) GetTrackedChats() ([]models.TrackedChat, error) {
	return s.repo.GetTrackedChats()
}

// AddTrackedChat включает отслеживание активности для чата и прогревает кеш.
// Используется, когда бота только что добавили в группу.
func (s *ChatActivityService) AddTrackedChat(chatID int64, title string, chatType string) error {
	if err := s.repo.UpsertTrackedChat(chatID, title, chatType); err != nil {
		return err
	}
	s.mu.Lock()
	s.trackedChatIDs[chatID] = true
	s.mu.Unlock()
	return nil
}

// RemoveTrackedChat снимает чат с отслеживания (при удалении бота из группы).
// История сообщений сохраняется.
func (s *ChatActivityService) RemoveTrackedChat(chatID int64) error {
	if err := s.repo.DeactivateTrackedChat(chatID); err != nil {
		return err
	}
	s.mu.Lock()
	delete(s.trackedChatIDs, chatID)
	s.mu.Unlock()
	return nil
}

// GetUserStats возвращает статистику конкретного пользователя
func (s *ChatActivityService) GetUserStats(userID int64, days int) (*models.UserStats, error) {
	return s.repo.GetUserStats(userID, days)
}

// GetDailyActivityByUser возвращает активность пользователя по дням
func (s *ChatActivityService) GetDailyActivityByUser(userID int64, days int) ([]models.DailyActivity, error) {
	return s.repo.GetDailyActivityByUser(userID, days)
}

// GetMessagesForExport возвращает данные для CSV экспорта
func (s *ChatActivityService) GetMessagesForExport(chatID *int64, days int) ([]models.ExportRow, error) {
	return s.repo.GetMessagesForExport(chatID, days)
}

// CountUserMessagesInChatSince — сколько сообщений у юзера в чате за период.
func (s *ChatActivityService) CountUserMessagesInChatSince(chatID, userID int64, since time.Time) (int64, error) {
	return s.repo.CountUserMessagesInChatSince(chatID, userID, since)
}

// LookupUserIDByUsername ищет telegram_user_id по @username среди сообщений
// в этом чате. Используется командами модерации, когда есть только @username
// (например, /unban @user). 0 — не найден.
func (s *ChatActivityService) LookupUserIDByUsername(chatID int64, username string) (int64, error) {
	return s.repo.LookupUserIDByUsername(chatID, username)
}

// CleanupOldMessages удаляет сообщения старше retentionDays дней
func (s *ChatActivityService) CleanupOldMessages(retentionDays int) {
	beforeDate := time.Now().AddDate(0, 0, -retentionDays)
	count, err := s.repo.DeleteOldMessages(beforeDate)
	if err != nil {
		log.Printf("Error cleaning up old chat messages: %v", err)
		return
	}
	if count > 0 {
		log.Printf("Cleaned up %d old chat messages (older than %d days)", count, retentionDays)
	}
}
