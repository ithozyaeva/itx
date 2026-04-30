package models

import "time"

// CommentEntityType — поддерживаемые типы родительских сущностей.
// Расширяется по мере появления новых разделов с обсуждениями.
type CommentEntityType string

const (
	CommentEntityAIMaterial CommentEntityType = "ai_material"
	CommentEntityEvent      CommentEntityType = "event"

	CommentMinLen = 1
	CommentMaxLen = 4_000
)

func IsValidCommentEntityType(t CommentEntityType) bool {
	switch t {
	case CommentEntityAIMaterial, CommentEntityEvent:
		return true
	}
	return false
}

// Comment — единая сущность для комментариев ко всем разделам платформы.
// Polymorphic-связь через (entity_type, entity_id) — никакого FK, целостность
// поддерживается на уровне сервиса (visibility-checker для каждого
// entity_type) и триггером comments_count_recalc для денормализации
// счётчика на parent-таблице.
type Comment struct {
	Id         int64             `json:"id" gorm:"primaryKey"`
	EntityType CommentEntityType `json:"entityType" gorm:"column:entity_type;size:32;not null"`
	EntityId   int64             `json:"entityId" gorm:"column:entity_id;not null"`
	AuthorId   int64             `json:"authorId" gorm:"column:author_id;not null"`
	Author     *Member           `json:"author,omitempty" gorm:"foreignKey:AuthorId"`
	Body       string            `json:"body" gorm:"not null"`
	LikesCount int               `json:"likesCount" gorm:"column:likes_count;default:0"`
	IsHidden   bool              `json:"isHidden" gorm:"column:is_hidden;default:false"`
	CreatedAt  time.Time         `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time         `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`

	// Заполняется сервисом — флаг «текущий viewer лайкнул».
	Liked bool `json:"liked" gorm:"-:all"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentLike struct {
	CommentId int64     `gorm:"primaryKey;column:comment_id"`
	MemberId  int64     `gorm:"primaryKey;column:member_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}

type CreateCommentRequest struct {
	Body string `json:"body"`
}
