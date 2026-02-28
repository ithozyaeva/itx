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
	ExpiresAt        *time.Time        `json:"expiresAt,omitempty" gorm:"column:expires_at"`
	ConversionsCount int64             `json:"conversionsCount" gorm:"-"`
	CreatedAt        time.Time         `json:"-"`
	UpdatedAt        time.Time         `json:"updatedAt"`
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
	Company        string     `json:"company"`
	Grade          string     `json:"grade"`
	ProfTags       []ProfTag  `json:"profTags"`
	VacationsCount int        `json:"vacationsCount"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
}

type UpdateLinkRequest struct {
	Id             int64             `json:"id"`
	Company        string            `json:"company"`
	Grade          string            `json:"grade"`
	ProfTags       []ProfTag         `json:"profTags"`
	VacationsCount int               `json:"vacationsCount"`
	Status         ReferalLinkStatus `json:"status"`
	ExpiresAt      *time.Time        `json:"expiresAt,omitempty"`
}

type DeleteLinkRequest struct {
	Id int64 `json:"id"`
}

type ReferralConversion struct {
	Id             int64     `json:"id" gorm:"primaryKey"`
	ReferralLinkId int64     `json:"referralLinkId" gorm:"column:referral_link_id"`
	MemberId       int64     `json:"memberId" gorm:"column:member_id"`
	ConvertedAt    time.Time `json:"convertedAt" gorm:"column:converted_at;default:CURRENT_TIMESTAMP"`
}

func (ReferralConversion) TableName() string {
	return "referral_conversions"
}
