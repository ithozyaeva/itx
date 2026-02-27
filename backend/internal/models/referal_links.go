package models

import "time"

type ReferalLink struct {
	Id             int64             `json:"id" gorm:"primaryKey"`
	AuthorId       int64             `json:"-" gorm:"column:author_id"`
	Author         Member            `json:"author" gorm:"foreignKey:author_id"`
	Company        string            `json:"company"`
	Grade          string            `json:"grade"`
	ProfTags       []ProfTag         `json:"profTags" gorm:"many2many:referal_links_tags"`
	Status         ReferalLinkStatus `json:"status"`
	VacationsCount int               `json:"vacationsCount"`
	CreatedAt      time.Time         `json:"-"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type Grade string

const (
	SeniorGrade Grade = "senior"
	JuniorGrade Grade = "junior"
	MiddleGrade Grade = "middle"
)

type ReferalLinkStatus string

const (
	ReferalLinkFreezed ReferalLinkStatus = "freezed"
	ReferalLinkActive  ReferalLinkStatus = "active"
)

type AddLinkRequest struct {
	Company        string    `json:"company"`
	Grade          string    `json:"grade"`
	ProfTags       []ProfTag `json:"profTags"`
	VacationsCount int       `json:"vacationsCount"`
}

type UpdateLinkRequest struct {
	Id             int64             `json:"id"`
	Company        string            `json:"company"`
	Grade          string            `json:"grade"`
	ProfTags       []ProfTag         `json:"profTags"`
	VacationsCount int               `json:"vacationsCount"`
	Status         ReferalLinkStatus `json:"status"`
}

type DeleteLinkRequest struct {
	Id int64 `json:"id"`
}
