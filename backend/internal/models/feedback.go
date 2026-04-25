package models

import "time"

type Feedback struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	UserId    *int64    `json:"userId" gorm:"column:user_id"`
	Score     int       `json:"score" gorm:"column:score"`
	Comment   *string   `json:"comment" gorm:"column:comment"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}

type CreateFeedbackRequest struct {
	Score   int     `json:"score"`
	Comment *string `json:"comment"`
}

type FeedbackPublic struct {
	Id            int64     `json:"id"`
	UserId        *int64    `json:"userId"`
	UserFirstName *string   `json:"userFirstName"`
	UserLastName  *string   `json:"userLastName"`
	UserUsername  *string   `json:"userUsername"`
	Score         int       `json:"score"`
	Comment       *string   `json:"comment"`
	CreatedAt     time.Time `json:"createdAt"`
}
