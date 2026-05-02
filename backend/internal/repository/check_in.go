package repository

import (
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type CheckInRepository struct{}

func NewCheckInRepository() *CheckInRepository {
	return &CheckInRepository{}
}

// Insert tries to record a check-in for (memberId, day). Returns inserted=true
// only when a new row was created (idempotent against double-click).
func (r *CheckInRepository) Insert(memberId int64, day time.Time) (inserted bool, err error) {
	res := database.DB.Exec(
		`INSERT INTO daily_check_ins (member_id, day) VALUES (?, ?) ON CONFLICT (member_id, day) DO NOTHING`,
		memberId, day,
	)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *CheckInRepository) HasCheckedIn(memberId int64, day time.Time) (bool, error) {
	var count int64
	err := database.DB.Model(&models.DailyCheckIn{}).
		Where("member_id = ? AND day = ?", memberId, day).
		Count(&count).Error
	return count > 0, err
}

func (r *CheckInRepository) Get(memberId int64, day time.Time) (*models.DailyCheckIn, error) {
	var c models.DailyCheckIn
	err := database.DB.Where("member_id = ? AND day = ?", memberId, day).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CheckedInMembers — set member_id, сделавших check-in за указанный день.
// Используется в batch-сценариях (например, вечерний пуш) вместо N
// отдельных HasCheckedIn-запросов.
func (r *CheckInRepository) CheckedInMembers(day time.Time) (map[int64]bool, error) {
	var ids []int64
	err := database.DB.Model(&models.DailyCheckIn{}).
		Where("day = ?", day).
		Pluck("member_id", &ids).Error
	if err != nil {
		return nil, err
	}
	out := make(map[int64]bool, len(ids))
	for _, id := range ids {
		out[id] = true
	}
	return out, nil
}
