package models

import (
	"ithozyeva/internal/s3resolve"
	"time"

	"gorm.io/gorm"
)

type Kudos struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	FromId     int64     `json:"fromId" gorm:"column:from_id;not null"`
	ToId       int64     `json:"toId" gorm:"column:to_id;not null"`
	Message    string    `json:"message" gorm:"column:message;not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	From       Member    `json:"from" gorm:"foreignKey:FromId;references:Id"`
	To         Member    `json:"to" gorm:"foreignKey:ToId;references:Id"`
}

type KudosPublic struct {
	Id            int64     `json:"id"`
	FromId        int64     `json:"fromId"`
	FromFirstName string    `json:"fromFirstName"`
	FromLastName  string    `json:"fromLastName"`
	FromUsername  string    `json:"fromUsername"`
	FromAvatarURL string    `json:"fromAvatarUrl"`
	ToId          int64     `json:"toId"`
	ToFirstName   string    `json:"toFirstName"`
	ToLastName    string    `json:"toLastName"`
	ToUsername    string    `json:"toUsername"`
	ToAvatarURL   string    `json:"toAvatarUrl"`
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (k *KudosPublic) AfterFind(tx *gorm.DB) (err error) {
	k.FromAvatarURL = s3resolve.ResolveS3URL(k.FromAvatarURL)
	k.ToAvatarURL = s3resolve.ResolveS3URL(k.ToAvatarURL)
	return nil
}

type CreateKudosRequest struct {
	ToId    int64  `json:"toId"`
	Message string `json:"message"`
}
