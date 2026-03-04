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
		Preload("Assignees")

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
		Preload("Assignees").
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

func (r *TaskExchangeRepository) AddAssignee(taskId int64, memberId int64) error {
	assignee := models.TaskExchangeAssignee{
		TaskId:   taskId,
		MemberId: memberId,
	}
	return database.DB.Create(&assignee).Error
}

func (r *TaskExchangeRepository) RemoveAssignee(taskId int64, memberId int64) error {
	return database.DB.
		Where("task_id = ? AND member_id = ?", taskId, memberId).
		Delete(&models.TaskExchangeAssignee{}).Error
}

func (r *TaskExchangeRepository) GetAssigneesCount(taskId int64) (int64, error) {
	var count int64
	err := database.DB.Model(&models.TaskExchangeAssignee{}).
		Where("task_id = ?", taskId).
		Count(&count).Error
	return count, err
}

func (r *TaskExchangeRepository) IsAssignee(taskId int64, memberId int64) (bool, error) {
	var count int64
	err := database.DB.Model(&models.TaskExchangeAssignee{}).
		Where("task_id = ? AND member_id = ?", taskId, memberId).
		Count(&count).Error
	return count > 0, err
}

func (r *TaskExchangeRepository) UpdateStatus(id int64, status models.TaskExchangeStatus) error {
	return database.DB.Model(&models.TaskExchange{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *TaskExchangeRepository) Update(id int64, updates map[string]interface{}) error {
	return database.DB.Model(&models.TaskExchange{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *TaskExchangeRepository) MarkDone(id int64) (int64, error) {
	result := database.DB.Model(&models.TaskExchange{}).
		Where("id = ? AND status = ?", id, models.TaskStatusInProgress).
		Update("status", models.TaskStatusDone)
	return result.RowsAffected, result.Error
}

func (r *TaskExchangeRepository) Approve(id int64) (int64, error) {
	result := database.DB.Model(&models.TaskExchange{}).
		Where("id = ? AND status = ?", id, models.TaskStatusDone).
		Update("status", models.TaskStatusApproved)
	return result.RowsAffected, result.Error
}

func (r *TaskExchangeRepository) Reject(id int64) (int64, error) {
	result := database.DB.Model(&models.TaskExchange{}).
		Where("id = ? AND status = ?", id, models.TaskStatusDone).
		Update("status", models.TaskStatusInProgress)
	return result.RowsAffected, result.Error
}

func (r *TaskExchangeRepository) Delete(id int64) error {
	return database.DB.Delete(&models.TaskExchange{}, id).Error
}
