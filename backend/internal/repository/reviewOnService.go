package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type ReviewOnServiceRepository struct {
	BaseRepository[models.ReviewOnService]
}

func NewReviewOnServiceRepository() *ReviewOnServiceRepository {
	return &ReviewOnServiceRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.ReviewOnService{}),
	}
}

func (r *ReviewOnServiceRepository) Search(limit *int, offset *int, filter *SearchFilter, order *Order) ([]models.ReviewOnService, int64, error) {
	var reviews []models.ReviewOnService
	var count int64

	// Сначала считаем общее количество всех записей
	if err := database.DB.Model(&models.ReviewOnService{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Создаем запрос с предзагрузкой связанных данных
	query := database.DB.Model(&models.ReviewOnService{}).
		Preload("Service")

	// Применяем limit только если он передан
	if limit != nil {
		query = query.Limit(*limit)
	}

	// Применяем offset только если он передан
	if offset != nil {
		query = query.Offset(*offset)
	}

	// Выполняем запрос
	if err := query.Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}

// GetReviewsWithMentorInfo получает отзывы с информацией о менторе
func (r *ReviewOnServiceRepository) GetReviewsWithMentorInfo(limit *int, offset *int) ([]models.ReviewOnServiceWithMentor, int64, error) {
	var reviews []models.ReviewOnServiceWithMentor
	var count int64

	query := database.DB.Table("\"reviewOnService\" AS r").
		Select(`
			r.id, 
			r."serviceId", 
			s.name as service_name, 
			s."ownerId" as mentor_id,
			m.username as mentor_name,
			r.author, 
			r.text, 
			r.date
		`).
		Joins(`JOIN services s ON r."serviceId" = s.id`).
		Joins(`JOIN mentors mt ON s."ownerId" = mt.id`).
		Joins(`JOIN members m ON mt."memberId" = m.id`)

	// Сначала считаем общее количество
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Применяем пагинацию
	if limit != nil {
		query = query.Limit(*limit)
	}
	if offset != nil {
		query = query.Offset(*offset)
	}

	// Выполняем запрос
	if err := query.Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}

// GetById получает отзыв по ID с информацией о услуге
func (r *ReviewOnServiceRepository) GetById(id int64) (*models.ReviewOnService, error) {
	var review models.ReviewOnService
	if err := database.DB.Preload("Service").First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// Create создает новый отзыв
func (r *ReviewOnServiceRepository) CreateReview(review *models.ReviewOnServiceRequest) (*models.ReviewOnService, error) {
	newReview := &models.ReviewOnService{
		ServiceId: review.ServiceId,
		Author:    review.Author,
		Text:      review.Text,
		Date:      review.Date,
	}

	if err := database.DB.Create(newReview).Error; err != nil {
		return nil, err
	}

	// Загружаем связанную услугу
	if err := database.DB.Preload("Service").First(newReview, newReview.Id).Error; err != nil {
		return nil, err
	}

	return newReview, nil
}
