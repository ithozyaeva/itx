package models

import (
	"ithozyeva/internal/s3resolve"
	"time"

	"gorm.io/gorm"
)

type SeasonStatus string

const (
	SeasonStatusActive   SeasonStatus = "ACTIVE"
	SeasonStatusFinished SeasonStatus = "FINISHED"
)

type Season struct {
	Id        int64        `json:"id" gorm:"primaryKey"`
	Title     string       `json:"title" gorm:"column:title;not null"`
	StartDate time.Time    `json:"startDate" gorm:"column:start_date;not null"`
	EndDate   time.Time    `json:"endDate" gorm:"column:end_date;not null"`
	Status    SeasonStatus `json:"status" gorm:"column:status;default:'ACTIVE'"`
	CreatedAt time.Time    `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

type SeasonLeaderboardEntry struct {
	MemberId  int64  `json:"memberId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"tg"`
	AvatarURL string `json:"avatarUrl"`
	Total     int    `json:"total"`
	Rank      int    `json:"rank"`
}

func (e *SeasonLeaderboardEntry) AfterFind(tx *gorm.DB) (err error) {
	e.AvatarURL = s3resolve.ResolveS3URL(e.AvatarURL)
	return nil
}

type SeasonWithLeaderboard struct {
	Season      Season                   `json:"season"`
	Leaderboard []SeasonLeaderboardEntry `json:"leaderboard"`
}
