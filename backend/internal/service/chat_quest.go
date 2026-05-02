package service

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendTelegramDMFunc — функция для отправки Telegram DM, устанавливается извне (из bot пакета) для избежания circular import
var SendTelegramDMFunc func(chatID int64, text string)

type ChatQuestService struct {
	repo       *repository.ChatQuestRepository
	pointsSvc  *PointsService
	memberRepo repository.MemberRepositoryInterface
}

func NewChatQuestService() *ChatQuestService {
	return &ChatQuestService{
		repo:       repository.NewChatQuestRepository(),
		pointsSvc:  NewPointsService(),
		memberRepo: repository.NewMemberRepository(),
	}
}

// ProcessMessage обрабатывает сообщение для всех активных квестов
func (s *ChatQuestService) ProcessMessage(message *tgbotapi.Message, memberID *int64) {
	if memberID == nil {
		return
	}

	// Параллельно с чат-квестами — трекаем дейлик «3 сообщения в чате».
	TrackDailyTrigger(*memberID, "chat_message", 1)

	quests, err := s.repo.GetActiveQuests()
	if err != nil {
		log.Printf("Error getting active quests: %v", err)
		return
	}

	for _, quest := range quests {
		// Проверяем условие chat_id (NULL = любой чат)
		if quest.ChatID != nil && *quest.ChatID != message.Chat.ID {
			continue
		}

		switch quest.QuestType {
		case models.QuestTypeDailyStreak:
			s.processDailyStreak(quest, *memberID)
		default: // message_count
			s.processMessageCount(quest, *memberID)
		}
	}
}

// processMessageCount обрабатывает квест типа message_count
func (s *ChatQuestService) processMessageCount(quest models.ChatQuest, memberID int64) {
	newCount, targetCount, alreadyCompleted, err := s.repo.IncrementProgress(quest.Id, memberID)
	if err != nil {
		log.Printf("Error incrementing quest progress (quest=%d, member=%d): %v", quest.Id, memberID, err)
		return
	}

	if alreadyCompleted {
		return
	}

	if newCount >= targetCount {
		s.completeQuest(quest, memberID)
	}
}

// processDailyStreak обрабатывает квест типа daily_streak
func (s *ChatQuestService) processDailyStreak(quest models.ChatQuest, memberID int64) {
	today := time.Now().Truncate(24 * time.Hour)

	// Записываем день активности
	if err := s.repo.RecordStreakDay(quest.Id, memberID, today); err != nil {
		log.Printf("Error recording streak day (quest=%d, member=%d): %v", quest.Id, memberID, err)
		return
	}

	// Считаем текущую серию
	streak, err := s.repo.GetCurrentStreak(quest.Id, memberID)
	if err != nil {
		log.Printf("Error getting current streak (quest=%d, member=%d): %v", quest.Id, memberID, err)
		return
	}

	// Обновляем прогресс
	if err := s.repo.SetProgressCount(quest.Id, memberID, streak); err != nil {
		log.Printf("Error setting progress count (quest=%d, member=%d): %v", quest.Id, memberID, err)
		return
	}

	// Проверяем достижение цели
	if streak >= quest.TargetCount {
		// Проверяем, не завершён ли уже
		var progress models.ChatQuestProgress
		if err := s.repo.GetProgress(quest.Id, memberID, &progress); err == nil && progress.Completed {
			return
		}
		s.completeQuest(quest, memberID)
	}
}

// completeQuest завершает квест, начисляет баллы и отправляет уведомления
func (s *ChatQuestService) completeQuest(quest models.ChatQuest, memberID int64) {
	if err := s.repo.MarkCompleted(quest.Id, memberID); err != nil {
		log.Printf("Error marking quest completed (quest=%d, member=%d): %v", quest.Id, memberID, err)
		return
	}

	// Начисляем баллы
	s.pointsSvc.GiveCustomPoints(
		memberID,
		quest.PointsReward,
		models.PointReasonChatQuest,
		"chat_quest",
		quest.Id,
		quest.Title,
	)
	log.Printf("Quest completed: quest=%d, member=%d, reward=%d", quest.Id, memberID, quest.PointsReward)

	// In-app уведомление
	go func() {
		title := "Задание выполнено!"
		body := fmt.Sprintf("Вы выполнили задание «%s» и получили %d баллов!", quest.Title, quest.PointsReward)
		if err := CreateNotification(memberID, "quest_complete", title, body); err != nil {
			log.Printf("Error creating quest completion notification (quest=%d, member=%d): %v", quest.Id, memberID, err)
		}
	}()

	// Telegram DM уведомление
	go func() {
		if SendTelegramDMFunc == nil {
			return
		}
		// Проверяем MuteAll
		notifSettingsRepo := repository.NewNotificationSettingsRepository()
		if ns, err := notifSettingsRepo.GetByMemberId(memberID); err == nil && ns.MuteAll {
			return
		}
		member, err := s.memberRepo.GetById(memberID)
		if err != nil {
			log.Printf("Error getting member for telegram notification (member=%d): %v", memberID, err)
			return
		}
		if member.TelegramID == 0 {
			return
		}
		text := fmt.Sprintf("🎉 <b>Задание выполнено!</b>\n\nВы выполнили задание «%s» и получили <b>%d баллов</b>!", quest.Title, quest.PointsReward)
		SendTelegramDMFunc(member.TelegramID, text)
	}()
}

// GetActiveQuestsForMember возвращает активные квесты с прогрессом для участника
func (s *ChatQuestService) GetActiveQuestsForMember(memberID int64) ([]models.ChatQuestWithProgress, error) {
	return s.repo.GetQuestsForMember(memberID)
}

// GetAllQuestsForMember возвращает все квесты с прогрессом для участника с фильтром
func (s *ChatQuestService) GetAllQuestsForMember(memberID int64, filter string) ([]models.ChatQuestWithProgress, error) {
	return s.repo.GetAllQuestsForMember(memberID, filter)
}

// CRUD для админки

func (s *ChatQuestService) CreateQuest(quest *models.ChatQuest) error {
	return s.repo.CreateQuest(quest)
}

func (s *ChatQuestService) UpdateQuest(quest *models.ChatQuest) error {
	return s.repo.UpdateQuest(quest)
}

func (s *ChatQuestService) DeleteQuest(id int64) error {
	return s.repo.DeleteQuest(id)
}

func (s *ChatQuestService) GetAllQuests(limit, offset int) ([]models.ChatQuest, int64, error) {
	return s.repo.GetAllQuests(limit, offset)
}

func (s *ChatQuestService) GetQuestByID(id int64) (*models.ChatQuest, error) {
	return s.repo.GetQuestByID(id)
}
