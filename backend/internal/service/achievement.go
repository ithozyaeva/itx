package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type AchievementService struct {
	repo     *repository.AchievementRepository
	pointsRepo *repository.PointsRepository
}

func NewAchievementService() *AchievementService {
	return &AchievementService{
		repo:     repository.NewAchievementRepository(),
		pointsRepo: repository.NewPointsRepository(),
	}
}

var AllAchievements = []models.Achievement{
	// Events
	{Id: "first_event", Title: "Первый шаг", Description: "Посетить первое событие", Icon: "footprints", Category: models.AchievementCategoryEvents, Threshold: 1},
	{Id: "events_5", Title: "Активист", Description: "Посетить 5 событий", Icon: "flame", Category: models.AchievementCategoryEvents, Threshold: 5},
	{Id: "events_10", Title: "Завсегдатай", Description: "Посетить 10 событий", Icon: "calendar-check", Category: models.AchievementCategoryEvents, Threshold: 10},
	{Id: "events_25", Title: "Ветеран", Description: "Посетить 25 событий", Icon: "medal", Category: models.AchievementCategoryEvents, Threshold: 25},
	{Id: "first_host", Title: "Спикер", Description: "Провести первое событие", Icon: "mic", Category: models.AchievementCategoryEvents, Threshold: 1},
	{Id: "hosts_5", Title: "Опытный спикер", Description: "Провести 5 событий", Icon: "presentation", Category: models.AchievementCategoryEvents, Threshold: 5},
	// Points
	{Id: "points_50", Title: "Новичок", Description: "Набрать 50 баллов", Icon: "star", Category: models.AchievementCategoryPoints, Threshold: 50},
	{Id: "points_100", Title: "Сотня", Description: "Набрать 100 баллов", Icon: "trophy", Category: models.AchievementCategoryPoints, Threshold: 100},
	{Id: "points_500", Title: "Топ игрок", Description: "Набрать 500 баллов", Icon: "crown", Category: models.AchievementCategoryPoints, Threshold: 500},
	// Social
	{Id: "first_review", Title: "Критик", Description: "Оставить первый отзыв", Icon: "message-square", Category: models.AchievementCategorySocial, Threshold: 1},
	{Id: "reviews_5", Title: "Эксперт мнений", Description: "Оставить 5 отзывов", Icon: "messages-square", Category: models.AchievementCategorySocial, Threshold: 5},
	{Id: "first_referral", Title: "Амбассадор", Description: "Создать первый реферал", Icon: "share-2", Category: models.AchievementCategorySocial, Threshold: 1},
	{Id: "referral_convert", Title: "Вербовщик", Description: "Получить первую конверсию", Icon: "user-plus", Category: models.AchievementCategorySocial, Threshold: 1},
	// Activity
	{Id: "profile_complete", Title: "Перфекционист", Description: "Полностью заполнить профиль", Icon: "user-check", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "streak_4w", Title: "Марафонец", Description: "Серия активности 4 недели", Icon: "zap", Category: models.AchievementCategoryActivity, Threshold: 1},
	{Id: "resume_upload", Title: "Карьерист", Description: "Загрузить резюме", Icon: "file-text", Category: models.AchievementCategoryActivity, Threshold: 1},
}

// achievementReasonMap maps achievement IDs to the PointReason they track.
var achievementReasonMap = map[string]models.PointReason{
	"first_event":      models.PointReasonEventAttend,
	"events_5":         models.PointReasonEventAttend,
	"events_10":        models.PointReasonEventAttend,
	"events_25":        models.PointReasonEventAttend,
	"first_host":       models.PointReasonEventHost,
	"hosts_5":          models.PointReasonEventHost,
	"first_review":     models.PointReasonReviewCommunity,
	"reviews_5":        models.PointReasonReviewCommunity,
	"first_referral":   models.PointReasonReferalCreate,
	"referral_convert": models.PointReasonReferalConversion,
	"profile_complete": models.PointReasonProfileComplete,
	"streak_4w":        models.PointReasonStreak4Weeks,
	"resume_upload":    models.PointReasonResumeUpload,
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
		case "points_50", "points_100", "points_500":
			progress = balance
		case "first_review", "reviews_5":
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
	}, nil
}
