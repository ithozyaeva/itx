package models

type MentorDbShortModel struct {
	Id         int64  `json:"id" gorm:"primaryKey"`
	MemberId   int64  `json:"memberId" gorm:"column:memberId"`
	Occupation string `json:"occupation"`
	Experience string `json:"experience"`
	Order      int    `json:"order"`
}

type MentorDbModel struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	MemberId   int64     `json:"memberId" gorm:"column:memberId"`
	Occupation string    `json:"occupation"`
	Experience string    `json:"experience"`
	Order      int       `json:"order"`
	Member     Member    `json:"member" gorm:"foreignKey:memberId;references:id"`
	ProfTags   []ProfTag `json:"profTags" gorm:"many2many:mentorsTags;foreignKey:id;joinForeignKey:mentorId;References:id;joinReferences:tagId"`
	Contacts   []Contact `json:"contacts" gorm:"foreignKey:ownerId;references:id"`
	Services   []Service `json:"services" gorm:"foreignKey:ownerId;references:id"`
}

type MentorModel struct {
	Id         int64     `json:"id"`
	Tg         string    `json:"tg"`
	Name       string    `json:"name"`
	Occupation string    `json:"occupation"`
	Experience string    `json:"experience"`
	Order      int       `json:"order"`
	MemberId   int       `json:"memberId"`
	ProfTags   []ProfTag `json:"profTags"`
	Contacts   []Contact `json:"contacts"`
	Services   []Service `json:"services"`
}

type MentorTagDbModel struct {
	MentorId int64 `gorm:"primaryKey;column:mentorId"`
	TagId    int64 `gorm:"primaryKey;column:tagId"`
}

// MentorCreateUpdateRequest представляет запрос на создание/обновление ментора
type MentorCreateUpdateRequest struct {
	Id         int64            `json:"id,omitempty"`
	MemberId   int64            `json:"memberId" binding:"required"`
	Occupation string           `json:"occupation"`
	Experience string           `json:"experience"` // Опечатка в названии поля
	Order      int              `json:"order"`
	ProfTags   []ProfTagRequest `json:"profTags,omitempty"`
	Contacts   []ContactRequest `json:"contacts,omitempty"`
	Services   []ServiceRequest `json:"services,omitempty"`
}

// ProfTagRequest представляет запрос на создание/обновление профессионального тега
type ProfTagRequest struct {
	Id    int64  `json:"id,omitempty"`
	Title string `json:"title" binding:"required"`
}

// ContactRequest представляет запрос на создание/обновление контакта
type ContactRequest struct {
	Id   int64  `json:"id,omitempty"`
	Type int16  `json:"type" binding:"required"`
	Link string `json:"link" binding:"required"`
}

// ServiceRequest представляет запрос на создание/обновление услуги
type ServiceRequest struct {
	Id    int64  `json:"id,omitempty"`
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price"`
}

func (MentorDbModel) TableName() string {
	return "mentors"
}

func (MentorDbShortModel) TableName() string {
	return "mentors"
}

func (MentorTagDbModel) TableName() string {
	return "mentorsTags"
}

func (m *MentorDbShortModel) SetID(id int64) {
	m.Id = id
}

func (m *MentorDbModel) SetID(id int64) {
	m.Id = id
}
