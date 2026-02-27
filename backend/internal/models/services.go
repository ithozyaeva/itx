package models

type Service struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	OwnerId int64  `json:"ownerId" gorm:"column:ownerId"`
}

// TableName указывает GORM использовать правильное имя таблицы
func (Service) TableName() string {
	return "services"
}

type ReviewOnService struct {
	Id        int     `json:"id" gorm:"primaryKey"`
	ServiceId int     `json:"serviceId" gorm:"column:serviceId"`
	Service   Service `json:"service" gorm:"foreignKey:ServiceId;references:Id"`
	Author    string  `json:"author"`
	Text      string  `json:"text"`
	Date      string  `json:"date"`
}

// TableName указывает GORM использовать правильное имя таблицы
func (ReviewOnService) TableName() string {
	return "reviewOnService"
}

// ReviewOnServiceWithMentor представляет отзыв на услугу с информацией о менторе
type ReviewOnServiceWithMentor struct {
	Id          int    `json:"id"`
	ServiceId   int    `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	MentorId    int64  `json:"mentorId"`
	MentorName  string `json:"mentorName"`
	Author      string `json:"author"`
	Text        string `json:"text"`
	Date        string `json:"date"`
}

// ReviewOnServiceRequest представляет запрос на создание отзыва на услугу
type ReviewOnServiceRequest struct {
	ServiceId int    `json:"serviceId" binding:"required"`
	Author    string `json:"author" binding:"required"`
	Text      string `json:"text" binding:"required"`
	Date      string `json:"date"`
}
