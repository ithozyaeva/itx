package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

type AuthTokenRepository struct {
	db *gorm.DB
}

func NewAuthTokenRepository() *AuthTokenRepository {
	return &AuthTokenRepository{db: database.DB}
}

func (r *AuthTokenRepository) Create(token *models.AuthToken) (*models.AuthToken, error) {
	if err := r.db.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (r *AuthTokenRepository) GetByTelegramID(telegramID int64) (*models.AuthToken, error) {
	var user models.AuthToken
	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthTokenRepository) Update(token *models.AuthToken) (*models.AuthToken, error) {
	if err := r.db.Save(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (r *AuthTokenRepository) Delete(token *models.AuthToken) error {
	return r.db.Delete(token).Error
}

func (r *AuthTokenRepository) GetById(id int64) (*models.AuthToken, error) {
	var user models.AuthToken
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthTokenRepository) Search(limit *int, offset *int) ([]models.AuthToken, int64, error) {
	var users []models.AuthToken
	var count int64

	if err := r.db.Model(&models.AuthToken{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&models.AuthToken{})
	if limit != nil {
		query = query.Limit(*limit)
	}
	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *AuthTokenRepository) GetByToken(token string) (*models.AuthToken, error) {
	var user models.AuthToken
	err := r.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
