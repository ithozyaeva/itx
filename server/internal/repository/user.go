package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type UserRepository interface {
	GetUserByLogin(login string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
