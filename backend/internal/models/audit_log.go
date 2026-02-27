package models

import "time"

type ActorType string

const (
	ActorTypeAdmin    ActorType = "admin"
	ActorTypePlatform ActorType = "platform"
)

type AuditAction string

const (
	AuditActionCreate  AuditAction = "create"
	AuditActionUpdate  AuditAction = "update"
	AuditActionDelete  AuditAction = "delete"
	AuditActionApprove AuditAction = "approve"
)

type AuditLog struct {
	Id         int64      `json:"id" gorm:"primaryKey"`
	ActorId    int64      `json:"actorId" gorm:"column:actor_id"`
	ActorName  string     `json:"actorName" gorm:"column:actor_name"`
	ActorType  ActorType  `json:"actorType" gorm:"column:actor_type"`
	Action     AuditAction `json:"action" gorm:"column:action"`
	EntityType string     `json:"entityType" gorm:"column:entity_type"`
	EntityId   int64      `json:"entityId" gorm:"column:entity_id"`
	EntityName string     `json:"entityName" gorm:"column:entity_name"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"column:created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
