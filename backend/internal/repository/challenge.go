package repository

import (
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type ChallengeRepository struct{}

func NewChallengeRepository() *ChallengeRepository {
	return &ChallengeRepository{}
}

func (r *ChallengeRepository) GetActiveTemplates(kind string) ([]models.ChallengeTemplate, error) {
	var tpl []models.ChallengeTemplate
	q := database.DB.Where("active = TRUE")
	if kind != "" {
		q = q.Where("kind = ?", kind)
	}
	err := q.Order("title").Find(&tpl).Error
	return tpl, err
}

func (r *ChallengeRepository) GetTemplatesByMetric(metricKey string) ([]models.ChallengeTemplate, error) {
	var tpl []models.ChallengeTemplate
	err := database.DB.Where("active = TRUE AND metric_key = ?", metricKey).Find(&tpl).Error
	return tpl, err
}

func (r *ChallengeRepository) UpsertInstance(inst *models.ChallengeInstance) error {
	return database.DB.Exec(
		`INSERT INTO challenge_instances (template_id, kind, starts_at, ends_at, period_key)
		 VALUES (?, ?, ?, ?, ?) ON CONFLICT (template_id, period_key) DO NOTHING`,
		inst.TemplateId, inst.Kind, inst.StartsAt, inst.EndsAt, inst.PeriodKey,
	).Error
}

func (r *ChallengeRepository) GetActiveInstances(at time.Time) ([]models.ChallengeInstance, error) {
	var inst []models.ChallengeInstance
	err := database.DB.Where("starts_at <= ? AND ends_at >= ?", at, at).Find(&inst).Error
	return inst, err
}

func (r *ChallengeRepository) GetMyProgress(memberId int64, instanceIds []int64) ([]models.ChallengeProgress, error) {
	var progress []models.ChallengeProgress
	if len(instanceIds) == 0 {
		return progress, nil
	}
	err := database.DB.Where("member_id = ? AND instance_id IN ?", memberId, instanceIds).Find(&progress).Error
	return progress, err
}

// Admin CRUD ---------------------------------------------------------------

func (r *ChallengeRepository) GetAllTemplatesAdmin() ([]models.ChallengeTemplate, error) {
	var tpls []models.ChallengeTemplate
	err := database.DB.Order("active DESC, kind, code").Find(&tpls).Error
	return tpls, err
}

func (r *ChallengeRepository) GetTemplateById(id int64) (*models.ChallengeTemplate, error) {
	var t models.ChallengeTemplate
	if err := database.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ChallengeRepository) CreateTemplate(t *models.ChallengeTemplate) error {
	return database.DB.Create(t).Error
}

func (r *ChallengeRepository) UpdateTemplate(t *models.ChallengeTemplate) error {
	return database.DB.Save(t).Error
}

func (r *ChallengeRepository) DeleteTemplate(id int64) error {
	return database.DB.Delete(&models.ChallengeTemplate{}, id).Error
}

// GetRecentInstances — последние N запущенных инстансов независимо от
// активности. Используется в админке для аудита истории.
func (r *ChallengeRepository) GetRecentInstances(limit int) ([]models.ChallengeInstance, error) {
	var insts []models.ChallengeInstance
	err := database.DB.Order("starts_at DESC").Limit(limit).Find(&insts).Error
	return insts, err
}
