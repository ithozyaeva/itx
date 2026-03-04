package models

import "time"

type TaskExchangeStatus string

const (
	TaskStatusOpen       TaskExchangeStatus = "OPEN"
	TaskStatusInProgress TaskExchangeStatus = "IN_PROGRESS"
	TaskStatusDone       TaskExchangeStatus = "DONE"
	TaskStatusApproved   TaskExchangeStatus = "APPROVED"
)

type TaskExchange struct {
	Id           int64              `json:"id" gorm:"primaryKey"`
	Title        string             `json:"title" gorm:"column:title;size:255;not null"`
	Description  string             `json:"description" gorm:"column:description;default:''"`
	CreatorId    int64              `json:"creatorId" gorm:"column:creator_id;not null"`
	Creator      Member             `json:"creator" gorm:"foreignKey:CreatorId"`
	MaxAssignees int                `json:"maxAssignees" gorm:"column:max_assignees;default:1"`
	Assignees    []Member           `json:"assignees" gorm:"many2many:task_exchange_assignees;joinForeignKey:TaskId;joinReferences:MemberId"`
	Status       TaskExchangeStatus `json:"status" gorm:"column:status;size:20;default:'OPEN'"`
	CreatedAt    time.Time          `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time          `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (TaskExchange) TableName() string {
	return "task_exchanges"
}

type TaskExchangeAssignee struct {
	TaskId    int64     `json:"taskId" gorm:"column:task_id;primaryKey"`
	MemberId  int64     `json:"memberId" gorm:"column:member_id;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (TaskExchangeAssignee) TableName() string {
	return "task_exchange_assignees"
}

type CreateTaskExchangeRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MaxAssignees int    `json:"maxAssignees"`
}

type UpdateTaskExchangeRequest struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	MaxAssignees *int    `json:"maxAssignees"`
}
