package service

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"time"
)

type ReviewOnCommunityService struct {
	BaseService[models.ReviewOnCommunity]
	repo *repository.ReviewOnCommunityRepository
}

func NewReviewOnCommunityService() *ReviewOnCommunityService {
	repo := repository.NewReviewOnCommunityRepository()
	return &ReviewOnCommunityService{
		BaseService: NewBaseService(repo),
		repo:        repo,
	}
}

func (s *ReviewOnCommunityService) GetAllWithAuthor(limit *int, offset *int) (*models.RegistrySearch[models.ReviewOnCommunity], error) {
	reviews, count, err := s.repo.GetAllWithAuthor(limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.RegistrySearch[models.ReviewOnCommunity]{
		Items: reviews,
		Total: int(count),
	}, nil
}

// CreateReviewOnCommunity создает новый отзыв о сообществе
func (s *ReviewOnCommunityService) CreateReviewOnCommunity(req *models.CreateReviewOnCommunityRequest) error {
	// Найти пользователя по Telegram
	member, err := repository.NewMemberRepository().GetMemberByTelegram(req.AuthorTg)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя с Telegram %s: %w", req.AuthorTg, err)
	}

	date := time.Now().Format("2006-01-02")
	if req.Date != nil && *req.Date != "" {
		date = *req.Date
	}

	review := &models.ReviewOnCommunity{
		AuthorId: uint(member.Id),
		Text:     req.Text,
		Date:     date,
		Status:   models.ReviewOnCommunityStatusDraft,
	}

	_, err = s.repo.Create(review)
	if err != nil {
		return fmt.Errorf("ошибка при создании отзыва: %w", err)
	}

	return nil
}

func (s *ReviewOnCommunityService) GetApproved() (*[]models.ReviewOnCommunity, error) {
	reviews, err := s.repo.GetApproved(nil)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *ReviewOnCommunityService) Approve(id int64) (*models.ReviewOnCommunity, error) {
	existedReview, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	if existedReview.Status == models.ReviewOnCommunityStatusApproved {
		return existedReview, nil
	}

	existedReview.Status = models.ReviewOnCommunityStatusApproved

	approvedReview, err := s.repo.Update(existedReview)
	if err != nil {
		return nil, err
	}

	return approvedReview, nil
}
