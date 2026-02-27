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
	ProfTags   []ProfTag `json:"profTags" gorm:"many2many:mentors_tags;foreignKey:id;joinForeignKey:mentor_id;References:id;joinReferences:tagId"`
	Contacts   []Contact `json:"contacts" gorm:"foreignKey:ownerId;references:id"`
	Services   []Service `json:"services" gorm:"foreignKey:ownerId;references:id"`
}

type MentorModel struct {
	Id         int64     `json:"id"`
	Username   string    `json:"tg"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Occupation string    `json:"occupation"`
	Experience string    `json:"experience"`
	Birthday   *DateOnly `json:"birthday"`
	Roles      []Role    `json:"roles"`
	Order      int       `json:"order"`
	MemberId   int       `json:"memberId"`
	ProfTags   []ProfTag `json:"profTags"`
	Contacts   []Contact `json:"contacts"`
	Services   []Service `json:"services"`
}

type MentorsTag struct {
	MentorId int64 `gorm:"primaryKey;column:mentor_id"`
	TagId    int64 `gorm:"primaryKey;column:tag_id"`
}

func (MentorDbModel) TableName() string {
	return "mentors"
}

func (MentorDbShortModel) TableName() string {
	return "mentors"
}

func (MentorsTag) TableName() string {
	return "mentors_tags"
}

func (m *MentorDbShortModel) SetID(id int64) {
	m.Id = id
}

func (m *MentorDbModel) SetID(id int64) {
	m.Id = id
}

func (m *MentorDbModel) ToModel() MentorModel {
	return MentorModel{
		Id:         m.Id,
		Username:   m.Member.Username,
		FirstName:  m.Member.FirstName,
		LastName:   m.Member.LastName,
		Birthday:   m.Member.Birthday,
		Roles:      m.Member.Roles,
		Occupation: m.Occupation,
		Experience: m.Experience,
		Order:      m.Order,
		MemberId:   int(m.MemberId),
		ProfTags:   m.ProfTags,
		Contacts:   m.Contacts,
		Services:   m.Services,
	}
}
