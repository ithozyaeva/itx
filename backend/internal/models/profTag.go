package models

type ProfTag struct {
	Id    int64  `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}

// GORM рофл
func (ProfTag) TableName() string {
	return "profTags"
}
