package models

import (
	"ithozyeva/internal/s3resolve"
	"time"

	"gorm.io/gorm"
)

type Guild struct {
	Id          int64         `json:"id" gorm:"primaryKey"`
	Name        string        `json:"name" gorm:"column:name;not null;uniqueIndex"`
	Description string        `json:"description" gorm:"column:description;default:''"`
	Icon        string        `json:"icon" gorm:"column:icon;default:'users'"`
	Color       string        `json:"color" gorm:"column:color;default:'#6366f1'"`
	OwnerId     int64         `json:"ownerId" gorm:"column:owner_id;not null"`
	CreatedAt   time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	Members     []GuildMember `json:"members,omitempty" gorm:"foreignKey:GuildId"`
}

type GuildMember struct {
	GuildId  int64     `json:"guildId" gorm:"primaryKey;column:guild_id"`
	MemberId int64     `json:"memberId" gorm:"primaryKey;column:member_id"`
	JoinedAt time.Time `json:"joinedAt" gorm:"column:joined_at;autoCreateTime"`
}

type GuildPublic struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Icon            string `json:"icon"`
	Color           string `json:"color"`
	OwnerId         int64  `json:"ownerId"`
	OwnerFirstName  string `json:"ownerFirstName"`
	OwnerLastName   string `json:"ownerLastName"`
	OwnerUsername   string `json:"ownerUsername"`
	OwnerAvatarURL  string `json:"ownerAvatarUrl"`
	MemberCount     int    `json:"memberCount"`
	TotalPoints     int    `json:"totalPoints"`
	IsMember        bool   `json:"isMember"`
}

func (g *GuildPublic) AfterFind(tx *gorm.DB) (err error) {
	g.OwnerAvatarURL = s3resolve.ResolveS3URL(g.OwnerAvatarURL)
	return nil
}

type CreateGuildRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}
