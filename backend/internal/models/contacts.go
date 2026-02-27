package models

type Contact struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	Type    int16  `json:"type"`
	Link    string `json:"link"`
	OwnerId int64  `json:"ownerId" gorm:"column:ownerId"` // Изменено на int64, чтобы соответствовать типу ID ментора
}

// TableName указывает GORM использовать правильное имя таблицы
func (Contact) TableName() string {
	return "contacts"
}
