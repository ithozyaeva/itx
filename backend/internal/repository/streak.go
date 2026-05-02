package repository

import (
	"errors"

	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

type StreakRepository struct{}

func NewStreakRepository() *StreakRepository {
	return &StreakRepository{}
}

func (r *StreakRepository) Get(memberId int64) (*models.MemberStreak, error) {
	var s models.MemberStreak
	err := database.DB.Where("member_id = ?", memberId).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.MemberStreak{MemberId: memberId, FreezesAvailable: 1}, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StreakRepository) Save(s *models.MemberStreak) error {
	return database.DB.Save(s).Error
}
