package repository

import (
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type DailyTaskRepository struct{}

func NewDailyTaskRepository() *DailyTaskRepository {
	return &DailyTaskRepository{}
}

func (r *DailyTaskRepository) GetActiveByTrigger(triggerKey string) ([]models.DailyTask, error) {
	var tasks []models.DailyTask
	err := database.DB.Where("active = TRUE AND trigger_key = ?", triggerKey).Find(&tasks).Error
	return tasks, err
}

func (r *DailyTaskRepository) GetActive() ([]models.DailyTask, error) {
	var tasks []models.DailyTask
	err := database.DB.Where("active = TRUE").Order("tier, points").Find(&tasks).Error
	return tasks, err
}

func (r *DailyTaskRepository) GetByIds(ids []int64) ([]models.DailyTask, error) {
	var tasks []models.DailyTask
	if len(ids) == 0 {
		return tasks, nil
	}
	err := database.DB.Where("id IN ?", ids).Find(&tasks).Error
	return tasks, err
}

func (r *DailyTaskRepository) GetSetByDay(day time.Time) (*models.DailyTaskSet, error) {
	var set models.DailyTaskSet
	err := database.DB.Where("day = ?", day).First(&set).Error
	if err != nil {
		return nil, err
	}
	return &set, nil
}

func (r *DailyTaskRepository) UpsertSet(set *models.DailyTaskSet) error {
	return database.DB.Exec(
		`INSERT INTO daily_task_sets (day, task_ids) VALUES (?, ?) ON CONFLICT (day) DO NOTHING`,
		set.Day, set.TaskIds,
	).Error
}

func (r *DailyTaskRepository) GetMyProgress(memberId int64, day time.Time) ([]models.DailyTaskProgress, error) {
	var progress []models.DailyTaskProgress
	err := database.DB.Where("member_id = ? AND day = ?", memberId, day).Find(&progress).Error
	return progress, err
}
