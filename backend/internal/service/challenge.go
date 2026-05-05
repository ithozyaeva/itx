package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"gorm.io/gorm"
)

const (
	weeklyPickPerWeek = 3
	monthlyPickPerMonth = 1
)

type ChallengeService struct {
	repo           *repository.ChallengeRepository
	pointRepo      *repository.PointsRepository
	achievementSvc *AchievementService
}

func NewChallengeService() *ChallengeService {
	return &ChallengeService{
		repo:           repository.NewChallengeRepository(),
		pointRepo:      repository.NewPointsRepository(),
		achievementSvc: NewAchievementService(),
	}
}

// GenerateWeeklyChallenges создаёт 3 еженедельных инстанса на неделю,
// в которой лежит `at`. Идемпотентно: ON CONFLICT (template_id, period_key).
// Период — пн 00:00 МСК … вс 23:59:59 МСК. Выбор случайный, но
// детерминированный по ISO-неделе (rand.NewSource(year*100 + week)).
func (s *ChallengeService) GenerateWeeklyChallenges(at time.Time) error {
	at = at.In(utils.MSKLocation())
	monday := startOfISOWeekMSK(at)
	endsAt := monday.Add(7 * 24 * time.Hour).Add(-time.Second)
	year, week := monday.ISOWeek()
	periodKey := fmt.Sprintf("%d-W%02d", year, week)

	templates, err := s.repo.GetActiveTemplates(models.ChallengeKindWeekly)
	if err != nil {
		return err
	}
	if len(templates) == 0 {
		return nil
	}

	pickN := weeklyPickPerWeek
	if pickN > len(templates) {
		pickN = len(templates)
	}

	rng := rand.New(rand.NewSource(int64(year)*100 + int64(week)))
	picked := pickRandomTemplates(templates, pickN, rng)

	for _, t := range picked {
		inst := &models.ChallengeInstance{
			TemplateId: t.Id,
			Kind:       models.ChallengeKindWeekly,
			StartsAt:   monday,
			EndsAt:     endsAt,
			PeriodKey:  periodKey,
		}
		if err := s.repo.UpsertInstance(inst); err != nil {
			log.Printf("upsert weekly challenge instance (template=%d period=%s): %v", t.Id, periodKey, err)
		}
	}
	return nil
}

// GenerateMonthlyChallenge создаёт 1 ежемесячный инстанс на месяц `at`.
func (s *ChallengeService) GenerateMonthlyChallenge(at time.Time) error {
	at = at.In(utils.MSKLocation())
	startsAt := time.Date(at.Year(), at.Month(), 1, 0, 0, 0, 0, utils.MSKLocation())
	endsAt := startsAt.AddDate(0, 1, 0).Add(-time.Second)
	periodKey := fmt.Sprintf("%d-%02d", at.Year(), int(at.Month()))

	templates, err := s.repo.GetActiveTemplates(models.ChallengeKindMonthly)
	if err != nil {
		return err
	}
	if len(templates) == 0 {
		return nil
	}

	pickN := monthlyPickPerMonth
	if pickN > len(templates) {
		pickN = len(templates)
	}
	rng := rand.New(rand.NewSource(int64(at.Year())*100 + int64(at.Month())))
	picked := pickRandomTemplates(templates, pickN, rng)

	for _, t := range picked {
		inst := &models.ChallengeInstance{
			TemplateId: t.Id,
			Kind:       models.ChallengeKindMonthly,
			StartsAt:   startsAt,
			EndsAt:     endsAt,
			PeriodKey:  periodKey,
		}
		if err := s.repo.UpsertInstance(inst); err != nil {
			log.Printf("upsert monthly challenge instance (template=%d period=%s): %v", t.Id, periodKey, err)
		}
	}
	return nil
}

func pickRandomTemplates(in []models.ChallengeTemplate, n int, rng *rand.Rand) []models.ChallengeTemplate {
	if n >= len(in) {
		out := make([]models.ChallengeTemplate, len(in))
		copy(out, in)
		return out
	}
	indices := rng.Perm(len(in))[:n]
	sort.Ints(indices)
	out := make([]models.ChallengeTemplate, 0, n)
	for _, idx := range indices {
		out = append(out, in[idx])
	}
	return out
}

// startOfISOWeekMSK — понедельник 00:00 МСК для недели, в которой лежит t.
func startOfISOWeekMSK(t time.Time) time.Time {
	t = t.In(utils.MSKLocation())
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, utils.MSKLocation())
	return monday.AddDate(0, 0, -(weekday - 1))
}

// EnsureCurrentInstances — публичный хук для cron/watchdog: гарантирует,
// что текущие еженедельные и ежемесячный инстансы существуют. Идемпотентно.
func (s *ChallengeService) EnsureCurrentInstances() {
	now := time.Now().In(utils.MSKLocation())
	if err := s.GenerateWeeklyChallenges(now); err != nil {
		log.Printf("generate weekly challenges: %v", err)
	}
	if err := s.GenerateMonthlyChallenge(now); err != nil {
		log.Printf("generate monthly challenge: %v", err)
	}
}

// IncrementProgress инкрементирует прогресс юзера по всем активным
// instances с указанным metric_key. Идемпотентно по (instance, member)
// через UPSERT. Если порог достигнут — выдача reward_points (и опционально
// разблокировка ачивки).
//
// Запускайте в goroutine: TrackChallengeMetric делает это автоматически.
func (s *ChallengeService) IncrementProgress(memberId int64, metricKey string, n int) {
	if n <= 0 {
		n = 1
	}
	now := time.Now()

	// Активные instances вне зависимости от kind.
	insts, err := s.repo.GetActiveInstances(now)
	if err != nil {
		log.Printf("challenge active instances: %v", err)
		return
	}
	if len(insts) == 0 {
		return
	}

	// Подгружаем templates для проверки metric_key и target.
	templateIds := make([]int64, 0, len(insts))
	for _, inst := range insts {
		templateIds = append(templateIds, inst.TemplateId)
	}
	templatesMap := make(map[int64]models.ChallengeTemplate, len(templateIds))
	{
		var tpls []models.ChallengeTemplate
		if err := database.DB.Where("id IN ?", templateIds).Find(&tpls).Error; err != nil {
			log.Printf("challenge templates lookup: %v", err)
			return
		}
		for _, t := range tpls {
			templatesMap[t.Id] = t
		}
	}

	for _, inst := range insts {
		t, ok := templatesMap[inst.TemplateId]
		if !ok || t.MetricKey != metricKey {
			continue
		}
		if err := s.bumpAndAward(memberId, inst, t, n); err != nil {
			log.Printf("challenge bump (member=%d, code=%s): %v", memberId, t.Code, err)
		}
	}
}

func (s *ChallengeService) bumpAndAward(memberId int64, inst models.ChallengeInstance, t models.ChallengeTemplate, n int) error {
	awarded := false
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var p models.ChallengeProgress
		err := tx.Where("instance_id = ? AND member_id = ?", inst.Id, memberId).First(&p).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			p = models.ChallengeProgress{
				InstanceId: inst.Id,
				MemberId:   memberId,
				Progress:   0,
			}
			if cErr := tx.Create(&p).Error; cErr != nil {
				return cErr
			}
		case err != nil:
			return err
		}

		if p.Awarded {
			return nil
		}

		newProgress := p.Progress + n
		if newProgress > t.Target {
			newProgress = t.Target
		}
		updates := map[string]interface{}{"progress": newProgress}
		if newProgress >= t.Target {
			now := time.Now()
			updates["awarded"] = true
			updates["completed_at"] = now
			awarded = true
		}
		if err := tx.Model(&models.ChallengeProgress{}).
			Where("id = ?", p.Id).Updates(updates).Error; err != nil {
			return err
		}
		if awarded {
			pt := &models.PointTransaction{
				MemberId:    memberId,
				Amount:      t.RewardPoints,
				Reason:      models.PointReasonChallengeComplete,
				SourceType:  "challenge",
				SourceId:    p.Id,
				Description: fmt.Sprintf("Челлендж: %s", t.Title),
			}
			if err := s.pointRepo.AwardPointsTx(tx, pt); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if awarded {
		GetSSEHub().Publish(memberId, SSEEvent{Type: "challenges"})
		GetSSEHub().Publish(memberId, SSEEvent{Type: "points"})
		if t.RewardPoints > 0 {
			TrackChallengeMetric(memberId, "points_earned", t.RewardPoints)
		}
		// За завершение челленджа — билет в сегодняшний daily-раффл.
		// source_id=inst.Id уникально per member per challenge instance.
		AwardRaffleTicket(memberId, models.RaffleTicketSourceChallenge, inst.Id)
		// Если у шаблона прописан achievement_code — выдаём явную ачивку.
		// Идемпотентно (PRIMARY KEY на (member_id, code)).
		if t.AchievementCode != nil && *t.AchievementCode != "" {
			if grantErr := s.achievementSvc.GrantExplicit(memberId, *t.AchievementCode); grantErr != nil {
				log.Printf("grant achievement %q for member=%d: %v", *t.AchievementCode, memberId, grantErr)
			}
		}
	}
	return nil
}

// SetProgress используется метриками-агрегатами (например, points_earned —
// сумма за период), которые не инкрементируются по событию, а пересчитываются
// сводным запросом. Не использован в первой версии — оставлен под будущее
// расширение.
func (s *ChallengeService) SetProgress(memberId int64, metricKey string, value int) {
	// no-op заглушка для расширения метриками-агрегатами.
	_ = memberId
	_ = metricKey
	_ = value
}

// GetMyChallenges собирает ChallengesResponse для текущего юзера:
// все активные instances + мой прогресс по каждому.
func (s *ChallengeService) GetMyChallenges(memberId int64) (models.ChallengesResponse, error) {
	now := time.Now()
	insts, err := s.repo.GetActiveInstances(now)
	if err != nil {
		return models.ChallengesResponse{}, err
	}
	if len(insts) == 0 {
		return models.ChallengesResponse{Weekly: []models.ChallengeWithProgress{}, Monthly: []models.ChallengeWithProgress{}}, nil
	}

	templateIds := make([]int64, 0, len(insts))
	instanceIds := make([]int64, 0, len(insts))
	for _, inst := range insts {
		templateIds = append(templateIds, inst.TemplateId)
		instanceIds = append(instanceIds, inst.Id)
	}
	var tpls []models.ChallengeTemplate
	if err := database.DB.Where("id IN ?", templateIds).Find(&tpls).Error; err != nil {
		return models.ChallengesResponse{}, err
	}
	tplsMap := make(map[int64]models.ChallengeTemplate, len(tpls))
	for _, t := range tpls {
		tplsMap[t.Id] = t
	}

	progressList, err := s.repo.GetMyProgress(memberId, instanceIds)
	if err != nil {
		return models.ChallengesResponse{}, err
	}
	progressMap := make(map[int64]models.ChallengeProgress, len(progressList))
	for _, p := range progressList {
		progressMap[p.InstanceId] = p
	}

	resp := models.ChallengesResponse{
		Weekly:  make([]models.ChallengeWithProgress, 0),
		Monthly: make([]models.ChallengeWithProgress, 0),
	}
	for _, inst := range insts {
		t, ok := tplsMap[inst.TemplateId]
		if !ok {
			continue
		}
		p := progressMap[inst.Id]
		entry := models.ChallengeWithProgress{
			InstanceId:      inst.Id,
			TemplateId:      t.Id,
			Code:            t.Code,
			Title:           t.Title,
			Description:     t.Description,
			Icon:            t.Icon,
			Kind:            inst.Kind,
			MetricKey:       t.MetricKey,
			Target:          t.Target,
			RewardPoints:    t.RewardPoints,
			AchievementCode: t.AchievementCode,
			StartsAt:        inst.StartsAt,
			EndsAt:          inst.EndsAt,
			Progress:        p.Progress,
			CompletedAt:     p.CompletedAt,
			Awarded:         p.Awarded,
		}
		if inst.Kind == models.ChallengeKindMonthly {
			resp.Monthly = append(resp.Monthly, entry)
		} else {
			resp.Weekly = append(resp.Weekly, entry)
		}
	}
	return resp, nil
}
