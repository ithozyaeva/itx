package models

import (
	"log"
	"time"

	"ithozyeva/internal/s3resolve"

	"gorm.io/gorm"
)

const (
	ReviewOnCommunityStatusDraft    ReviewOnCommunityStatus = "DRAFT"
	ReviewOnCommunityStatusApproved ReviewOnCommunityStatus = "APPROVED"
)

type ReviewOnCommunityStatus string

type MemberRole struct {
	MemberId int64 `gorm:"primaryKey;column:member_id;not null"`
	Role     Role  `gorm:"primaryKey;size:255;not null"`
}

type Member struct {
	Id          int64        `json:"id" gorm:"primaryKey"`
	Username    string       `json:"tg" gorm:"column:username"`
	TelegramID  int64        `json:"telegramID" gorm:"column:telegram_id"`
	FirstName   string       `json:"firstName" gorm:"column:first_name"`
	LastName    string       `json:"lastName" gorm:"column:last_name"`
	Bio         string       `json:"bio" gorm:"column:bio;default:''"`
	Grade       string       `json:"grade" gorm:"column:grade;default:''"`
	Company     string       `json:"company" gorm:"column:company;default:''"`
	AvatarURL   string       `json:"avatarUrl" gorm:"column:avatar_url;default:''"`
	MemberRoles []MemberRole `json:"-" gorm:"foreignKey:MemberId;references:Id"`
	Roles       []Role       `json:"roles" gorm:"-:all"`
	Birthday    *DateOnly    `json:"birthday" gorm:"column:birthday"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	// ReferredByLinkID — реф-ссылка, по которой юзер впервые попал на платформу
	// (через /start ref_<id> в боте). Не FK: ссылка может быть удалена, но
	// атрибуция остаётся. Заполняется один раз при создании members-записи.
	ReferredByLinkID *int64 `json:"-" gorm:"column:referred_by_link_id"`
	// ReferralWelcomeSeenAt — отметка что юзер видел welcome-баннер про
	// реферрера. NULL — баннер ещё не показывали (или закрыли без ack).
	ReferralWelcomeSeenAt *time.Time `json:"-" gorm:"column:referral_welcome_seen_at"`
	// SubscriptionTier — эффективный тир подписки (EffectiveTier из subscription_users).
	// Заполняется хендлерами перед отдачей профиля. В БД не хранится.
	SubscriptionTier *SubscriptionTier `json:"subscriptionTier,omitempty" gorm:"-:all"`
}

type ReviewOnCommunity struct {
	Id       int                     `json:"id"`
	AuthorId uint                    `json:"authorId" gorm:"column:authorId"`
	Author   Member                  `json:"author" gorm:"foreignKey:authorId"`
	Text     string                  `json:"text"`
	Date     string                  `json:"date"`
	Status   ReviewOnCommunityStatus `json:"status"`
}

type ReviewOnCommunityWithAuthor struct {
	Id         int    `json:"id"`
	AuthorId   int    `json:"authorId"`
	AuthorName string `json:"authorName"`
	AuthorTg   string `json:"authorTg"`
	Text       string `json:"text"`
	Date       string `json:"date"`
}

type CreateReviewOnCommunityRequest struct {
	Text     string  `json:"text" binding:"required"`
	Date     *string `json:"date"`
	AuthorTg string  `json:"authorTg"`
}

type AddReviewOnCommunityRequest struct {
	Text string `json:"text" binding:"required"`
}

func (ReviewOnCommunity) TableName() string {
	return "reviewOnCommunity"
}

func (m *Member) GetRoleStrings() []Role {
	roles := make([]Role, len(m.MemberRoles))
	for i, r := range m.MemberRoles {
		roles[i] = r.Role
	}
	return roles
}
func (m *Member) SetRoleStrings(roleStrings []Role, memberId int64) {
	log.Default().Printf("Setting roles for member %d: %v", memberId, roleStrings)
	roles := make([]MemberRole, len(roleStrings))
	for i, r := range roleStrings {
		roles[i] = MemberRole{
			MemberId: memberId,
			Role:     r,
		}
	}

	log.Printf("Setting roles for member %d: %v", memberId, roles)

	m.MemberRoles = roles
}

func (m *Member) AfterFind(tx *gorm.DB) (err error) {
	m.Roles = m.GetRoleStrings()
	m.AvatarURL = s3resolve.ResolveS3URL(m.AvatarURL)
	return nil
}

func (m *Member) BeforeSave(tx *gorm.DB) (err error) {
	if len(m.Roles) > 0 {
		m.SetRoleStrings(m.Roles, m.Id)
	}
	return nil
}
