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
