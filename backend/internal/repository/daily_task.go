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

// Admin CRUD ---------------------------------------------------------------

func (r *DailyTaskRepository) GetAllAdmin() ([]models.DailyTask, error) {
	var tasks []models.DailyTask
	err := database.DB.Order("active DESC, tier, points").Find(&tasks).Error
	return tasks, err
}

func (r *DailyTaskRepository) GetById(id int64) (*models.DailyTask, error) {
	var t models.DailyTask
	if err := database.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *DailyTaskRepository) Create(t *models.DailyTask) error {
	return database.DB.Create(t).Error
}

func (r *DailyTaskRepository) Update(t *models.DailyTask) error {
	return database.DB.Save(t).Error
}

func (r *DailyTaskRepository) Delete(id int64) error {
	return database.DB.Delete(&models.DailyTask{}, id).Error
}

// GetRecentSets — последние N МСК-дней с их составом, свежие сверху.
// Используется в админке для аудита: что выпадало в какой день.
func (r *DailyTaskRepository) GetRecentSets(limit int) ([]models.DailyTaskSet, error) {
	var sets []models.DailyTaskSet
	err := database.DB.Order("day DESC").Limit(limit).Find(&sets).Error
	return sets, err
}

// AwardedCountsForDay — количество awarded задач каждого юзера за day.
// Один SQL вместо N. Используется в batch-сценариях (вечерний пуш).
// taskIds ограничивает выборку сегодняшним сетом, чтобы случайно не
// захватить устаревшие/удалённые задачи.
func (r *DailyTaskRepository) AwardedCountsForDay(day time.Time, taskIds []int64) (map[int64]int, error) {
	if len(taskIds) == 0 {
		return map[int64]int{}, nil
	}
	type row struct {
		MemberId int64 `gorm:"column:member_id"`
		Count    int
	}
	var rows []row
	err := database.DB.Raw(
		`SELECT member_id, COUNT(*) AS count
		 FROM daily_task_progress
		 WHERE day = ? AND awarded = TRUE AND task_id IN ?
		 GROUP BY member_id`,
		day, taskIds,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[int64]int, len(rows))
	for _, r := range rows {
		out[r.MemberId] = r.Count
	}
	return out, nil
}
