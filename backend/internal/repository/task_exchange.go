package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type TaskExchangeRepository struct{}

func NewTaskExchangeRepository() *TaskExchangeRepository {
	return &TaskExchangeRepository{}
}

func (r *TaskExchangeRepository) Search(status *string, limit, offset int) ([]models.TaskExchange, int64, error) {
	var items []models.TaskExchange
	var total int64

	countQuery := database.DB.Model(&models.TaskExchange{})
	if status != nil && *status != "" {
		countQuery = countQuery.Where("status = ?", *status)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := database.DB.
		Preload("Creator").
		Preload("Assignee")

	if status != nil && *status != "" {
		findQuery = findQuery.Where("status = ?", *status)
	}

	err := findQuery.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

func (r *TaskExchangeRepository) GetById(id int64) (*models.TaskExchange, error) {
	var task models.TaskExchange
	err := database.DB.
		Preload("Creator").
		Preload("Assignee").
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskExchangeRepository) Create(task *models.TaskExchange) (*models.TaskExchange, error) {
	if err := database.DB.Create(task).Error; err != nil {
		return nil, err
	}
	return r.GetById(task.Id)
}

func (r *TaskExchangeRepository) Assign(id int64, assigneeId int64) error {
	return database.DB.Model(&models.TaskExchange{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"assignee_id": assigneeId,
			"status":      models.TaskStatusInProgress,
		}).Error
}

func (r *TaskExchangeRepository) Unassign(id int64) error {
	return database.DB.Model(&models.TaskExchange{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"assignee_id": nil,
			"status":      models.TaskStatusOpen,
		}).Error
}

func (r *TaskExchangeRepository) MarkDone(id int64) error {
	return database.DB.Model(&models.TaskExchange{}).
		Where("id = ?", id).
		Update("status", models.TaskStatusDone).Error
}

func (r *TaskExchangeRepository) Approve(id int64) error {
	return database.DB.Model(&models.TaskExchange{}).
		Where("id = ?", id).
		Update("status", models.TaskStatusApproved).Error
}

func (r *TaskExchangeRepository) Reject(id int64) error {
	return database.DB.Model(&models.TaskExchange{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"assignee_id": nil,
			"status":      models.TaskStatusOpen,
		}).Error
}

func (r *TaskExchangeRepository) Delete(id int64) error {
	return database.DB.Delete(&models.TaskExchange{}, id).Error
}
