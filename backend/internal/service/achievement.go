package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type AchievementService struct {
	repo       *repository.AchievementRepository
	pointsRepo *repository.PointsRepository
}

func NewAchievementService() *AchievementService {
	return &AchievementService{
		repo:       repository.NewAchievementRepository(),
		pointsRepo: repository.NewPointsRepository(),
	}
}

var AllAchievements = []models.Achievement{
	// Events
	{Id: "first_event", Title: "Первый шаг", Description: "Посетить первое событие", Icon: "footprints", Category: models.AchievementCategoryEvents, Threshold: 1},
	{Id: "events_5", Title: "Активист", Description: "Посетить 5 событий", Icon: "flame", Category: models.AchievementCategoryEvents, Threshold: 5},
	{Id: "events_10", Title: "Завсегдатай", Description: "Посетить 10 событий", Icon: "calendar-check", Category: models.AchievementCategoryEvents, Threshold: 10},
	{Id: "events_25", Title: "Ветеран", Description: "Посетить 25 событий", Icon: "medal", Category: models.AchievementCategoryEvents, Threshold: 25},
	{Id: "events_50", Title: "Старожил", Description: "Посетить 50 событий", Icon: "history", Category: models.AchievementCategoryEvents, Threshold: 50},
	{Id: "first_host", Title: "Спикер", Description: "Провести первое событие", Icon: "mic", Category: models.AchievementCategoryEvents, Threshold: 1},
	{Id: "hosts_5", Title: "Опытный спикер", Description: "Провести 5 событий", Icon: "presentation", Category: models.AchievementCategoryEvents, Threshold: 5},
	{Id: "hosts_10", Title: "Гуру", Description: "Провести 10 событий", Icon: "graduation-cap", Category: models.AchievementCategoryEvents, Threshold: 10},
	// Points
	{Id: "points_50", Title: "Новичок", Description: "Набрать 50 баллов", Icon: "star", Category: models.AchievementCategoryPoints, Threshold: 50},
	{Id: "points_100", Title: "Сотня", Description: "Набрать 100 баллов", Icon: "trophy", Category: models.AchievementCategoryPoints, Threshold: 100},
	{Id: "points_500", Title: "Топ игрок", Description: "Набрать 500 баллов", Icon: "crown", Category: models.AchievementCategoryPoints, Threshold: 500},
	{Id: "points_1000", Title: "Легенда", Description: "Набрать 1000 баллов", Icon: "gem", Category: models.AchievementCategoryPoints, Threshold: 1000},
	// Social
	{Id: "first_review", Title: "Критик", Description: "Оставить первый отзыв", Icon: "message-square", Category: models.AchievementCategorySocial, Threshold: 1},
	{Id: "reviews_5", Title: "Эксперт мнений", Description: "Оставить 5 отзывов", Icon: "messages-square", Category: models.AchievementCategorySocial, Threshold: 5},
	{Id: "reviews_10", Title: "Обозреватель", Description: "Оставить 10 отзывов", Icon: "book-open", Category: models.AchievementCategorySocial, Threshold: 10},
	{Id: "first_referral", Title: "Амбассадор", Description: "Создать первый реферал", Icon: "share-2", Category: models.AchievementCategorySocial, Threshold: 1},
	{Id: "referrals_5", Title: "Агент влияния", Description: "Создать 5 рефералов", Icon: "users", Category: models.AchievementCategorySocial, Threshold: 5},
	{Id: "referral_convert", Title: "Вербовщик", Description: "Получить первую конверсию", Icon: "user-plus", Category: models.AchievementCategorySocial, Threshold: 1},
	// Activity
	{Id: "profile_complete", Title: "Перфекционист", Description: "Полностью заполнить профиль", Icon: "user-check", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "streak_4w", Title: "Марафонец", Description: "Серия активности 4 недели", Icon: "zap", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "resume_upload", Title: "Карьерист", Description: "Загрузить резюме", Icon: "file-text", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "first_listing", Title: "Продавец", Description: "Опубликовать первое объявление", Icon: "package", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "first_purchase", Title: "Покупатель", Description: "Оставить первую заявку на покупку", Icon: "shopping-cart", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "first_task", Title: "Инициатор", Description: "Создать первое задание", Icon: "clipboard-list", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "tasks_create_5", Title: "Менеджер", Description: "Создать 5 заданий", Icon: "list-checks", Category: models.AchievementCategoryActivity, Threshold: 5},
	{Id: "first_task_done", Title: "Исполнитель", Description: "Выполнить первое задание", Icon: "check-circle", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "tasks_done_5", Title: "Трудяга", Description: "Выполнить 5 заданий", Icon: "hard-hat", Category: models.AchievementCategoryActivity, Threshold: 5},
	{Id: "tasks_done_10", Title: "Мастер дел", Description: "Выполнить 10 заданий", Icon: "briefcase", Category: models.AchievementCategoryActivity, Threshold: 10},
	{Id: "first_quest", Title: "Квестер", Description: "Выполнить первое задание в чате", Icon: "target", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "quests_5", Title: "Охотник за квестами", Description: "Выполнить 5 заданий в чатах", Icon: "swords", Category: models.AchievementCategoryActivity, Threshold: 5},
	{Id: "chatter_of_week", Title: "Чаттерc недели", Description: "Стать самым активным участником чата за неделю", Icon: "message-circle", Category: models.AchievementCategoryActivity, Threshold: 1},
}

// achievementReasonMap maps achievement IDs to the PointReason they track.
var achievementReasonMap = map[string]models.PointReason{
	"first_event":      models.PointReasonEventAttend,
	"events_5":         models.PointReasonEventAttend,
	"events_10":        models.PointReasonEventAttend,
	"events_25":        models.PointReasonEventAttend,
	"events_50":        models.PointReasonEventAttend,
	"first_host":       models.PointReasonEventHost,
	"hosts_5":          models.PointReasonEventHost,
	"hosts_10":         models.PointReasonEventHost,
	"first_review":     models.PointReasonReviewCommunity,
	"reviews_5":        models.PointReasonReviewCommunity,
	"reviews_10":       models.PointReasonReviewCommunity,
	"first_referral":   models.PointReasonReferalCreate,
	"referrals_5":      models.PointReasonReferalCreate,
	"referral_convert": models.PointReasonReferalConversion,
	"profile_complete": models.PointReasonProfileComplete,
	"streak_4w":        models.PointReasonStreak4Weeks,
	"resume_upload":    models.PointReasonResumeUpload,
	"first_listing":    models.PointReasonMarketplaceCreate,
	"first_purchase":   models.PointReasonMarketplaceBuy,
	"first_task":       models.PointReasonTaskCreate,
	"tasks_create_5":   models.PointReasonTaskCreate,
	"first_task_done":  models.PointReasonTaskExecute,
	"tasks_done_5":     models.PointReasonTaskExecute,
	"tasks_done_10":    models.PointReasonTaskExecute,
	"first_quest":      models.PointReasonChatQuest,
	"quests_5":         models.PointReasonChatQuest,
	"chatter_of_week":  models.PointReasonChatterOfWeek,
}

func (s *AchievementService) GetUserAchievements(memberId int64) (*models.AchievementsResponse, error) {
	reasonCounts, err := s.repo.GetReasonCounts(memberId)
	if err != nil {
		return nil, err
	}

	balance, err := s.pointsRepo.GetBalance(memberId)
	if err != nil {
		return nil, err
	}

	return s.buildAchievements(reasonCounts, balance), nil
}

// GetAchievementCounts returns only the unlocked/total counts using a pre-computed balance.
func (s *AchievementService) GetAchievementCounts(memberId int64, balance int) (earned int, total int, err error) {
	reasonCounts, err := s.repo.GetReasonCounts(memberId)
	if err != nil {
		return 0, 0, err
	}

	resp := s.buildAchievements(reasonCounts, balance)
	return resp.UnlockedCount, resp.TotalCount, nil
}

func (s *AchievementService) buildAchievements(reasonCounts map[models.PointReason]int, balance int) *models.AchievementsResponse {
	// Also count reviews on services
	reviewCount := reasonCounts[models.PointReasonReviewCommunity] + reasonCounts[models.PointReasonReviewService]

	var items []models.UserAchievement
	unlockedCount := 0

	for _, a := range AllAchievements {
		ua := models.UserAchievement{
			Achievement: a,
		}

		var progress int
		switch a.Id {
		case "points_50", "points_100", "points_500", "points_1000":
			progress = balance
		case "first_review", "reviews_5", "reviews_10":
			progress = reviewCount
		default:
			reason, ok := achievementReasonMap[a.Id]
			if ok {
				progress = reasonCounts[reason]
			}
		}

		ua.Progress = progress
		if progress >= a.Threshold {
			ua.Unlocked = true
			unlockedCount++
		}

		items = append(items, ua)
	}

	return &models.AchievementsResponse{
		Items:         items,
		TotalCount:    len(items),
		UnlockedCount: unlockedCount,
	}
}
