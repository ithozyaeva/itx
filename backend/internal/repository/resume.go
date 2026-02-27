package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

type ResumeRepository struct {
	BaseRepository[models.Resume]
	db *gorm.DB
}

func NewResumeRepository() *ResumeRepository {
	return &ResumeRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.Resume{}),
		db:             database.DB,
	}
}

func (r *ResumeRepository) ListByTelegramID(tgID int64) ([]models.Resume, error) {
	var resumes []models.Resume
	err := r.db.Where("tg_id = ?", tgID).Order("\"created_at\" DESC").Find(&resumes).Error
	return resumes, err
}

func (r *ResumeRepository) GetByIDAndTelegram(id, tgID int64) (*models.Resume, error) {
	resume := new(models.Resume)
	if err := r.db.Where("id = ? AND tg_id = ?", id, tgID).First(resume).Error; err != nil {
		return nil, err
	}
	return resume, nil
}

func (r *ResumeRepository) SearchForAdmin(limit *int, offset *int, filter *models.ResumeFilter) ([]models.Resume, int64, error) {
	query := r.db.Model(&models.Resume{}).Preload("Member")

	if filter != nil {
		if filter.WorkFormat != nil && *filter.WorkFormat != "" {
			query = query.Where("work_format = ?", *filter.WorkFormat)
		}
		if filter.DesiredPosition != nil && *filter.DesiredPosition != "" {
			query = query.Where("desired_position ILIKE ?", "%"+*filter.DesiredPosition+"%")
		}
		if filter.WorkExperience != nil && *filter.WorkExperience != "" {
			query = query.Where("work_experience ILIKE ?", "%"+*filter.WorkExperience+"%")
		}
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if limit != nil {
		query = query.Limit(*limit)
	}
	if offset != nil {
		query = query.Offset(*offset)
	}

	var items []models.Resume
	if err := query.Order("\"created_at\" DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

func (r *ResumeRepository) ListByIDs(ids []int64) ([]models.Resume, error) {
	var resumes []models.Resume
	if err := r.db.Where("id IN ?", ids).Find(&resumes).Error; err != nil {
		return nil, err
	}
	return resumes, nil
}

func (r *ResumeRepository) GetByIdWithMember(id int64) (*models.Resume, error) {
	resume := new(models.Resume)
	if err := r.db.Preload("Member").First(resume, id).Error; err != nil {
		return nil, err
	}
	return resume, nil
}
