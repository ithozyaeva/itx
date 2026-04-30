package models

import "time"

type AIMaterialContentType string
type AIMaterialKind string

const (
	AIMaterialContentTypePrompt AIMaterialContentType = "prompt"
	AIMaterialContentTypeLink   AIMaterialContentType = "link"
	AIMaterialContentTypeAgent  AIMaterialContentType = "agent"

	AIMaterialKindPrompt   AIMaterialKind = "prompt"
	AIMaterialKindSkill    AIMaterialKind = "skill"
	AIMaterialKindLibrary  AIMaterialKind = "library"
	AIMaterialKindTutorial AIMaterialKind = "tutorial"
	AIMaterialKindAgent    AIMaterialKind = "agent"

	AIMaterialMaxTags        = 5
	AIMaterialMaxTagLen      = 40
	AIMaterialMinTitleLen    = 3
	AIMaterialMaxTitleLen    = 120
	AIMaterialMinSummaryLen  = 30
	AIMaterialMaxSummaryLen  = 500
	AIMaterialMaxURLLen      = 2048
	AIMaterialMaxPromptBody  = 50_000
	AIMaterialMaxAgentConfig = 50_000
)

type AIMaterial struct {
	Id             int64                 `json:"id" gorm:"primaryKey"`
	AuthorId       int64                 `json:"authorId" gorm:"column:author_id;not null"`
	Author         *Member               `json:"author,omitempty" gorm:"foreignKey:AuthorId"`
	Title          string                `json:"title" gorm:"size:120;not null"`
	Summary        string                `json:"summary" gorm:"not null"`
	ContentType    AIMaterialContentType `json:"contentType" gorm:"column:content_type;size:16;not null"`
	MaterialKind   AIMaterialKind        `json:"materialKind" gorm:"column:material_kind;size:16;not null"`
	PromptBody     string                `json:"promptBody" gorm:"column:prompt_body;default:''"`
	ExternalURL    string                `json:"externalUrl" gorm:"column:external_url;size:2048;default:''"`
	AgentConfig    string                `json:"agentConfig" gorm:"column:agent_config;default:''"`
	LikesCount     int                   `json:"likesCount" gorm:"column:likes_count;default:0"`
	BookmarksCount int                   `json:"bookmarksCount" gorm:"column:bookmarks_count;default:0"`
	CommentsCount  int                   `json:"commentsCount" gorm:"column:comments_count;default:0"`
	IsHidden       bool                  `json:"isHidden" gorm:"column:is_hidden;default:false"`
	CreatedAt      time.Time             `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time             `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`

	// Заполняются сервисом, не маппятся на колонки.
	Tags       []string `json:"tags" gorm:"-:all"`
	Liked      bool     `json:"liked" gorm:"-:all"`
	Bookmarked bool     `json:"bookmarked" gorm:"-:all"`
}

func (AIMaterial) TableName() string {
	return "ai_materials"
}

type AIMaterialTag struct {
	MaterialId int64  `gorm:"primaryKey;column:material_id"`
	Tag        string `gorm:"primaryKey;size:40;column:tag"`
}

func (AIMaterialTag) TableName() string {
	return "ai_material_tags"
}

type AIMaterialLike struct {
	MaterialId int64     `gorm:"primaryKey;column:material_id"`
	MemberId   int64     `gorm:"primaryKey;column:member_id"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (AIMaterialLike) TableName() string {
	return "ai_material_likes"
}

type AIMaterialBookmark struct {
	MaterialId int64     `gorm:"primaryKey;column:material_id"`
	MemberId   int64     `gorm:"primaryKey;column:member_id"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (AIMaterialBookmark) TableName() string {
	return "ai_material_bookmarks"
}

type AIMaterialComment struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	MaterialId int64     `json:"materialId" gorm:"column:material_id;not null"`
	AuthorId   int64     `json:"authorId" gorm:"column:author_id;not null"`
	Author     *Member   `json:"author,omitempty" gorm:"foreignKey:AuthorId"`
	Body       string    `json:"body" gorm:"not null"`
	LikesCount int       `json:"likesCount" gorm:"column:likes_count;default:0"`
	IsHidden   bool      `json:"isHidden" gorm:"column:is_hidden;default:false"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`

	// Заполняется сервисом — флаг «текущий viewer лайкнул».
	Liked bool `json:"liked" gorm:"-:all"`
}

func (AIMaterialComment) TableName() string {
	return "ai_material_comments"
}

type AIMaterialCommentLike struct {
	CommentId int64     `gorm:"primaryKey;column:comment_id"`
	MemberId  int64     `gorm:"primaryKey;column:member_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (AIMaterialCommentLike) TableName() string {
	return "ai_material_comment_likes"
}

type CreateAIMaterialRequest struct {
	Title        string                `json:"title"`
	Summary      string                `json:"summary"`
	ContentType  AIMaterialContentType `json:"contentType"`
	MaterialKind AIMaterialKind        `json:"materialKind"`
	PromptBody   string                `json:"promptBody"`
	ExternalURL  string                `json:"externalUrl"`
	AgentConfig  string                `json:"agentConfig"`
	Tags         []string              `json:"tags"`
}

type UpdateAIMaterialRequest = CreateAIMaterialRequest

func IsValidAIMaterialContentType(v AIMaterialContentType) bool {
	switch v {
	case AIMaterialContentTypePrompt, AIMaterialContentTypeLink, AIMaterialContentTypeAgent:
		return true
	}
	return false
}

func IsValidAIMaterialKind(v AIMaterialKind) bool {
	switch v {
	case AIMaterialKindPrompt, AIMaterialKindSkill, AIMaterialKindLibrary, AIMaterialKindTutorial, AIMaterialKindAgent:
		return true
	}
	return false
}
