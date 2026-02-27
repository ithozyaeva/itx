package models

import "time"

type WorkFormat string

const (
	WorkFormatRemote WorkFormat = "REMOTE"
	WorkFormatHybrid WorkFormat = "HYBRID"
	WorkFormatOffice WorkFormat = "OFFICE"
)

type Resume struct {
	Id               int64      `json:"id" gorm:"primaryKey"`
	TgID             int64      `json:"tgId" gorm:"column:tg_id"`
	FilePath         string     `json:"filePath" gorm:"column:file_path"`
	FileName         string     `json:"fileName" gorm:"column:file_name"`
	WorkExperience   string     `json:"workExperience" gorm:"column:work_experience"`
	DesiredPosition  string     `json:"desiredPosition" gorm:"column:desired_position"`
	WorkFormat       WorkFormat `json:"workFormat" gorm:"column:work_format"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	Member           *Member    `json:"member,omitempty" gorm:"foreignKey:TgID;references:TelegramID"`
	ParsedAt         *time.Time `json:"parsedAt,omitempty" gorm:"-:all"`
	ParsedConfidence float64    `json:"parsedConfidence,omitempty" gorm:"-:all"`
}

func (Resume) TableName() string {
	return "resumes"
}

type ResumeFilter struct {
	WorkFormat      *WorkFormat `query:"workFormat"`
	DesiredPosition *string     `query:"desiredPosition"`
	WorkExperience  *string     `query:"workExperience"`
}

type CreateResumeRequest struct {
	WorkExperience  string     `form:"workExperience"`
	DesiredPosition string     `form:"desiredPosition"`
	WorkFormat      WorkFormat `form:"workFormat"`
}

type UpdateResumeRequest struct {
	WorkExperience  *string     `json:"workExperience"`
	DesiredPosition *string     `json:"desiredPosition"`
	WorkFormat      *WorkFormat `json:"workFormat"`
}

func (wf WorkFormat) IsValid() bool {
	switch wf {
	case WorkFormatRemote, WorkFormatHybrid, WorkFormatOffice, "":
		return true
	default:
		return false
	}
}
