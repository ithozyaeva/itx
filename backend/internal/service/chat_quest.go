package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

		newCount, targetCount, alreadyCompleted, err := s.repo.IncrementProgress(quest.Id, *memberID)
		if err != nil {
			log.Printf("Error incrementing quest progress (quest=%d, member=%d): %v", quest.Id, *memberID, err)
			continue
		}

		if alreadyCompleted {
			continue
		}

		// Проверяем достижение цели
		if newCount >= targetCount {
			if err := s.repo.MarkCompleted(quest.Id, *memberID); err != nil {
				log.Printf("Error marking quest completed (quest=%d, member=%d): %v", quest.Id, *memberID, err)
				continue
			}

			// Начисляем баллы
			s.pointsSvc.GiveCustomPoints(
				*memberID,
				quest.PointsReward,
				models.PointReasonChatQuest,
				"chat_quest",
				quest.Id,
				quest.Title,
			)
			log.Printf("Quest completed: quest=%d, member=%d, reward=%d", quest.Id, *memberID, quest.PointsReward)
		}
	}
}

// GetActiveQuestsForMember возвращает активные квесты с прогрессом для участника
func (s *ChatQuestService) GetActiveQuestsForMember(memberID int64) ([]models.ChatQuestWithProgress, error) {
	return s.repo.GetQuestsForMember(memberID)
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
