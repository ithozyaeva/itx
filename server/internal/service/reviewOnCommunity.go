package service

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
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
func (s *ReviewOnCommunityService) CreateReviewOnCommunity(req *models.ReviewOnCommunityRequest) (*models.ReviewOnCommunity, error) {
	// Найти пользователя по Telegram
	member, err := repository.NewMemberRepository().GetMemberByTelegram(req.AuthorTg)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти пользователя с Telegram %s: %w", req.AuthorTg, err)
	}

	review := &models.ReviewOnCommunity{
		AuthorId: uint(member.Id),
		Text:     req.Text,
		Date:     req.Date,
	}

	// Сохранить отзыв в базе данных
	createdReview, err := s.repo.Create(review)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании отзыва: %w", err)
	}

	return createdReview, nil
}
