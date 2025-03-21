package models

type Member struct {
	Id   int64  `json:"id"`
	Tg   string `json:"tg"`
	Name string `json:"name"`
}

type ReviewOnCommunity struct {
	Id       int    `json:"id"`
	AuthorId uint   `json:"authorId" gorm:"column:authorId"`
	Author   Member `json:"author" gorm:"foreignKey:authorId"`
	Text     string `json:"text"`
	Date     string `json:"date"`
}

type ReviewOnCommunityWithAuthor struct {
	Id         int    `json:"id"`
	AuthorId   int    `json:"authorId"`
	AuthorName string `json:"authorName"`
	AuthorTg   string `json:"authorTg"`
	Text       string `json:"text"`
	Date       string `json:"date"`
}

// ReviewOnCommunityRequest представляет запрос на создание отзыва о сообществе
type ReviewOnCommunityRequest struct {
	Text     string `json:"text" binding:"required"`
	Date     string `json: "date"`
	AuthorTg string `json:"authorTg" binding:"required"`
}

// MemberModel представляет модель участника с дополнительной информацией
type MemberModel struct {
	Id       int64  `json:"id"`
	Tg       string `json:"tg"`
	Name     string `json:"name"`
	IsMentor bool   `json:"isMentor"`
}

func (ReviewOnCommunity) TableName() string {
	return "reviewOnCommunity"
}
