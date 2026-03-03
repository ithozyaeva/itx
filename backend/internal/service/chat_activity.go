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
	trackedChatIDs map[int64]bool
	mu             sync.RWMutex
}

func NewChatActivityService() *ChatActivityService {
	s := &ChatActivityService{
		repo:           repository.NewChatActivityRepository(),
		trackedChatIDs: make(map[int64]bool),
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

	msg := &models.ChatMessage{
		ChatID:            message.Chat.ID,
		TelegramUserID:    message.From.ID,
		TelegramUsername:  message.From.UserName,
		TelegramFirstName: message.From.FirstName,
		SentAt:            time.Unix(int64(message.Date), 0),
	}

	if err := s.repo.SaveMessage(msg); err != nil {
		log.Printf("Error saving chat message: %v", err)
	}
}

// GetStats возвращает общую статистику активности
func (s *ChatActivityService) GetStats() (*models.ChatActivityStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -7)
	tomorrow := todayStart.AddDate(0, 0, 1)

	totalToday, uniqueToday, err := s.repo.GetTotalStats(todayStart, tomorrow)
	if err != nil {
		return nil, err
	}

	totalWeek, uniqueWeek, err := s.repo.GetTotalStats(weekStart, tomorrow)
	if err != nil {
		return nil, err
	}

	chatStats, err := s.repo.GetMessageCountsByChat(weekStart, tomorrow)
	if err != nil {
		return nil, err
	}

	return &models.ChatActivityStats{
		TotalMessagesToday: totalToday,
		TotalMessagesWeek:  totalWeek,
		UniqueUsersToday:   uniqueToday,
		UniqueUsersWeek:    uniqueWeek,
		ChatStats:          chatStats,
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
