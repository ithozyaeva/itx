package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync/atomic"
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// triggerCacheState — снапшот trigger_keys, попавших в сегодняшний набор.
// Используется в TrackDailyTrigger, чтобы не делать SQL для каждого
// view-эндпоинта, если соответствующего триггера сегодня нет в выборке.
type triggerCacheState struct {
	day  time.Time
	keys map[string]struct{}
}

var triggerCache atomic.Pointer[triggerCacheState]

// todayHasTrigger возвращает (есть_ли_в_сегодняшнем_сете, кэш_валиден).
// Если кэш пуст или устарел — возвращаем (true, false), чтобы fallback
// прошёл через SQL (надёжнее, чем потерять прогресс из-за прогрева).
func todayHasTrigger(key string) (bool, bool) {
	state := triggerCache.Load()
	if state == nil {
		return true, false
	}
	if !state.day.Equal(utils.MSKToday()) {
		return true, false
	}
	_, ok := state.keys[key]
	return ok, true
}

const dailyTasksPerDay = 5

type DailyTaskService struct {
	repo       *repository.DailyTaskRepository
	checkInSvc *CheckInService
	pointRepo  *repository.PointsRepository
}

func NewDailyTaskService() *DailyTaskService {
	return &DailyTaskService{
		repo:       repository.NewDailyTaskRepository(),
		checkInSvc: NewCheckInService(),
		pointRepo:  repository.NewPointsRepository(),
	}
}

// GenerateTodaySet формирует общий на всех набор из 5 задач для day.
// Идемпотентно (ON CONFLICT (day) DO NOTHING). Распределение:
// 2 engagement + 1 light + 1 meaningful + 1 случайная из всего пула.
// Detеrministic-rand на base day.Unix() — одинаковый набор на всех
// инстансах сервиса и при рестартах.
func (s *DailyTaskService) GenerateTodaySet(day time.Time) error {
	day = utils.MSKDay(day)

	if existing, err := s.repo.GetSetByDay(day); err == nil && existing != nil {
		return nil
	}

	tasks, err := s.repo.GetActive()
	if err != nil {
		return fmt.Errorf("get active tasks: %w", err)
	}
	if len(tasks) < dailyTasksPerDay {
		return fmt.Errorf("в пуле меньше %d активных задач", dailyTasksPerDay)
	}

	// Группировка по tier для квот.
	byTier := make(map[string][]models.DailyTask)
	for _, t := range tasks {
		byTier[t.Tier] = append(byTier[t.Tier], t)
	}

	rng := rand.New(rand.NewSource(day.Unix()))
	picked := make(map[int64]bool)
	pickIds := make([]int64, 0, dailyTasksPerDay)

	pickFromTier := func(tier string) {
		pool := filterUnpicked(byTier[tier], picked)
		if len(pool) == 0 {
			return
		}
		choice := pool[rng.Intn(len(pool))]
		picked[choice.Id] = true
		pickIds = append(pickIds, choice.Id)
	}

	pickFromTier(models.DailyTaskTierEngagement)
	pickFromTier(models.DailyTaskTierEngagement)
	pickFromTier(models.DailyTaskTierLight)
	pickFromTier(models.DailyTaskTierMeaningful)

	// Последний слот — случайная из всех оставшихся (любой tier).
	rest := filterUnpicked(tasks, picked)
	if len(rest) > 0 {
		choice := rest[rng.Intn(len(rest))]
		pickIds = append(pickIds, choice.Id)
		picked[choice.Id] = true
	}

	// Если каких-то tier'ов не хватило (пул мал), добиваем случайными.
	for len(pickIds) < dailyTasksPerDay {
		rest = filterUnpicked(tasks, picked)
		if len(rest) == 0 {
			break
		}
		choice := rest[rng.Intn(len(rest))]
		pickIds = append(pickIds, choice.Id)
		picked[choice.Id] = true
	}
	sort.Slice(pickIds, func(i, j int) bool { return pickIds[i] < pickIds[j] })

	set := &models.DailyTaskSet{Day: day, TaskIds: pickIds}
	return s.repo.UpsertSet(set)
}

func filterUnpicked(in []models.DailyTask, picked map[int64]bool) []models.DailyTask {
	out := in[:0:0]
	for _, t := range in {
		if !picked[t.Id] {
			out = append(out, t)
		}
	}
	return out
}

// GetMyToday собирает DailyTodayResponse для конкретного юзера.
// При отсутствии сегодняшнего набора пытается его сгенерировать (защита
// от случая «пользователь зашёл раньше watchdog'а»).
func (s *DailyTaskService) GetMyToday(memberId int64) (models.DailyTodayResponse, error) {
	day := utils.MSKToday()

	set, err := s.repo.GetSetByDay(day)
	if errors.Is(err, gorm.ErrRecordNotFound) || set == nil {
		if genErr := s.GenerateTodaySet(day); genErr != nil {
			return models.DailyTodayResponse{}, genErr
		}
		set, err = s.repo.GetSetByDay(day)
	}
	if err != nil {
		return models.DailyTodayResponse{}, err
	}

	tasks, err := s.repo.GetByIds(set.TaskIds)
	if err != nil {
		return models.DailyTodayResponse{}, err
	}
	progress, err := s.repo.GetMyProgress(memberId, day)
	if err != nil {
		return models.DailyTodayResponse{}, err
	}
	progressByTask := make(map[int64]models.DailyTaskProgress, len(progress))
	for _, p := range progress {
		progressByTask[p.TaskId] = p
	}

	checkInDone, checkInAt, _ := s.checkInSvc.HasCheckedInToday(memberId)

	resp := models.DailyTodayResponse{
		CheckIn: models.DailyCheckInState{Done: checkInDone, At: checkInAt},
		Tasks:   make([]models.DailyTaskWithProgress, 0, len(tasks)),
	}

	tasksById := make(map[int64]models.DailyTask, len(tasks))
	for _, t := range tasks {
		tasksById[t.Id] = t
	}
	// Порядок — как в task_ids набора (стабильный для всех юзеров).
	awarded := 0
	bonusAwarded := false
	for _, id := range set.TaskIds {
		t, ok := tasksById[id]
		if !ok {
			continue
		}
		p := progressByTask[id]
		if p.Awarded {
			awarded++
		}
		if p.BonusAwarded {
			bonusAwarded = true
		}
		resp.Tasks = append(resp.Tasks, models.DailyTaskWithProgress{
			DailyTask:   t,
			Progress:    p.Progress,
			CompletedAt: p.CompletedAt,
			Awarded:     p.Awarded,
		})
	}

	resp.AllBonus = models.DailyAllBonusState{
		Points:  models.PointValues[models.PointReasonDailyAllTasksBonus],
		Awarded: bonusAwarded,
	}

	_ = awarded
	return resp, nil
}

// IncrementProgress инкрементирует прогресс юзера по всем сегодняшним задачам
// с указанным triggerKey. Атомарно через одну SQL-транзакцию: upsert,
// инкремент, фиксация awarded и выдача баллов. После — проверка all-bonus.
//
// Запускайте в goroutine из существующих хендлеров — функция сама
// публикует SSE и не должна блокировать пользовательский ответ.
func (s *DailyTaskService) IncrementProgress(memberId int64, triggerKey string, n int) {
	if n <= 0 {
		n = 1
	}
	day := utils.MSKToday()

	set, err := s.repo.GetSetByDay(day)
	if err != nil || set == nil {
		return
	}

	matching, err := s.repo.GetActiveByTrigger(triggerKey)
	if err != nil || len(matching) == 0 {
		return
	}

	// Пересечение matching ∩ today's set.
	idsInSet := make(map[int64]bool, len(set.TaskIds))
	for _, id := range set.TaskIds {
		idsInSet[id] = true
	}

	awardedAny := false
	for _, t := range matching {
		if !idsInSet[t.Id] {
			continue
		}
		didAward, err := s.bumpAndAward(memberId, day, t, n)
		if err != nil {
			log.Printf("daily increment (member=%d, code=%s): %v", memberId, t.Code, err)
			continue
		}
		if didAward {
			awardedAny = true
		}
	}

	if awardedAny {
		s.maybeAwardAllBonus(memberId, day, set)
		GetSSEHub().Publish(memberId, SSEEvent{Type: "dailies"})
		GetSSEHub().Publish(memberId, SSEEvent{Type: "points"})
	}
}

// bumpAndAward в одной транзакции: upsert progress, инкремент, и если
// порог достигнут впервые — фиксация awarded + начисление баллов.
// Возвращает didAward=true только при переходе awarded false→true.
func (s *DailyTaskService) bumpAndAward(memberId int64, day time.Time, t models.DailyTask, n int) (bool, error) {
	var didAward bool
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Upsert: создаём строку или подтягиваем существующую.
		var p models.DailyTaskProgress
		err := tx.Where("member_id = ? AND day = ? AND task_id = ?", memberId, day, t.Id).First(&p).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			p = models.DailyTaskProgress{
				MemberId: memberId,
				Day:      day,
				TaskId:   t.Id,
				Progress: 0,
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
		}
		if err := tx.Model(&models.DailyTaskProgress{}).
			Where("id = ?", p.Id).Updates(updates).Error; err != nil {
			return err
		}

		if newProgress >= t.Target {
			didAward = true
			pt := &models.PointTransaction{
				MemberId:    memberId,
				Amount:      t.Points,
				Reason:      models.PointReasonDailyTaskComplete,
				SourceType:  "daily_task",
				SourceId:    p.Id,
				Description: fmt.Sprintf("Дейлик: %s", t.Title),
			}
			if err := s.pointRepo.AwardPointsTx(tx, pt); err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil && didAward && t.Points > 0 {
		// Покрываем агрегатную метрику m_owner и т.п.
		TrackChallengeMetric(memberId, "points_earned", t.Points)
	}
	return didAward, err
}

func (s *DailyTaskService) maybeAwardAllBonus(memberId int64, day time.Time, set *models.DailyTaskSet) {
	if len(set.TaskIds) == 0 {
		return
	}
	var awardedCount int64
	if err := database.DB.Model(&models.DailyTaskProgress{}).
		Where("member_id = ? AND day = ? AND task_id IN ? AND awarded = TRUE", memberId, day, []int64(set.TaskIds)).
		Count(&awardedCount).Error; err != nil {
		log.Printf("count awarded dailies (member=%d): %v", memberId, err)
		return
	}
	if int(awardedCount) < len(set.TaskIds) {
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// idempotent: source_type='daily_all_bonus', source_id=day.Unix(); per-member.
		pt := &models.PointTransaction{
			MemberId:    memberId,
			Amount:      models.PointValues[models.PointReasonDailyAllTasksBonus],
			Reason:      models.PointReasonDailyAllTasksBonus,
			SourceType:  "daily_all_bonus",
			SourceId:    day.Unix(),
			Description: "Бонус за все 5 дейликов",
		}
		if err := s.pointRepo.AwardPointsTx(tx, pt); err != nil {
			return err
		}
		// Помечаем bonus_awarded=true на всех записях прогресса этого
		// дня — фронт показывает через первую запись.
		return tx.Model(&models.DailyTaskProgress{}).
			Where("member_id = ? AND day = ?", memberId, day).
			Update("bonus_awarded", true).Error
	})
	if err != nil {
		log.Printf("award all-dailies bonus (member=%d): %v", memberId, err)
		return
	}
	GetSSEHub().Publish(memberId, SSEEvent{Type: "points"})

	// Челлендж-метрики: «день, когда выполнены все 5 дейликов» и
	// агрегат points_earned для m_owner.
	TrackChallengeMetric(memberId, "all_dailies_days", 1)
	if bonus := models.PointValues[models.PointReasonDailyAllTasksBonus]; bonus > 0 {
		TrackChallengeMetric(memberId, "points_earned", bonus)
	}
}

// EnsureTodaySet — публичный хук для cron/watchdog: гарантирует, что
// набор на сегодня существует (idempotent). Заодно обновляет
// triggerCache, чтобы hot-path вызовы TrackDailyTrigger не ходили
// в БД, если триггер не выпал в сегодняшний сет.
func (s *DailyTaskService) EnsureTodaySet() error {
	if err := s.GenerateTodaySet(utils.MSKToday()); err != nil {
		return err
	}
	s.RefreshTriggerCache()
	return nil
}

// RefreshTriggerCache подгружает trigger_keys сегодняшнего набора в
// in-memory atomic-кэш. Дёшево (одна выборка на 5 строк) и сильно
// сокращает число SQL на view-хендлерах вне сегодняшних задач.
func (s *DailyTaskService) RefreshTriggerCache() {
	day := utils.MSKToday()
	set, err := s.repo.GetSetByDay(day)
	if err != nil || set == nil {
		// Не сбрасываем существующий кэш — повторим в следующий watchdog.
		return
	}
	tasks, err := s.repo.GetByIds([]int64(set.TaskIds))
	if err != nil {
		return
	}
	keys := make(map[string]struct{}, len(tasks))
	for _, t := range tasks {
		if t.Active {
			keys[t.TriggerKey] = struct{}{}
		}
	}
	triggerCache.Store(&triggerCacheState{day: day, keys: keys})
}

// debug-only — выгрузить ID набора для проверки.
func (s *DailyTaskService) TaskIdsForDay(day time.Time) ([]int64, error) {
	set, err := s.repo.GetSetByDay(utils.MSKDay(day))
	if err != nil {
		return nil, err
	}
	return []int64(pq.Int64Array(set.TaskIds)), nil
}
