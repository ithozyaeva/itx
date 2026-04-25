package service

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"strings"
)

type FeedbackService struct {
	repo *repository.FeedbackRepository
}

func NewFeedbackService() *FeedbackService {
	return &FeedbackService{
		repo: repository.NewFeedbackRepository(),
	}
}

const maxFeedbackCommentLen = 2000

func (s *FeedbackService) Create(member *models.Member, req models.CreateFeedbackRequest) (*models.Feedback, error) {
	if req.Score < 0 || req.Score > 10 {
		return nil, fmt.Errorf("score must be between 0 and 10")
	}

	if member != nil {
		count, err := s.repo.CountTodayByMember(member.Id)
		if err != nil {
			return nil, err
		}
		if count >= 1 {
			return nil, fmt.Errorf("можно отправить не более одного отзыва в сутки")
		}
	}

	feedback := &models.Feedback{
		Score: req.Score,
	}

	if req.Comment != nil {
		trimmed := strings.TrimSpace(*req.Comment)
		if len(trimmed) > maxFeedbackCommentLen {
			return nil, fmt.Errorf("комментарий слишком длинный (макс. %d символов)", maxFeedbackCommentLen)
		}
		if trimmed != "" {
			feedback.Comment = &trimmed
		}
	}

	if member != nil {
		feedback.UserId = &member.Id
	}

	if err := s.repo.Create(feedback); err != nil {
		return nil, err
	}

	return feedback, nil
}

func (s *FeedbackService) List(limit, offset int) ([]models.FeedbackPublic, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(limit, offset)
}
