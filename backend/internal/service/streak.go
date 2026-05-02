package service

import (
	"fmt"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
)

// StreakThreshold описывает порог стрика и связанный PointReason.
type StreakThreshold struct {
	Days   int
	Reason models.PointReason
	Reward int
}

// StreakThresholds — пороговые значения стрика, отсортированы по возрастанию.
var StreakThresholds = []StreakThreshold{
	{Days: 3, Reason: models.PointReasonDailyStreak3, Reward: models.PointValues[models.PointReasonDailyStreak3]},
	{Days: 7, Reason: models.PointReasonDailyStreak7, Reward: models.PointValues[models.PointReasonDailyStreak7]},
	{Days: 14, Reason: models.PointReasonDailyStreak14, Reward: models.PointValues[models.PointReasonDailyStreak14]},
	{Days: 30, Reason: models.PointReasonDailyStreak30, Reward: models.PointValues[models.PointReasonDailyStreak30]},
}

type StreakService struct {
	repo      *repository.StreakRepository
	pointRepo *repository.PointsRepository
}

func NewStreakService() *StreakService {
	return &StreakService{
		repo:      repository.NewStreakRepository(),
		pointRepo: repository.NewPointsRepository(),
	}
}

// ApplyCheckIn пересчитывает стрик после факта успешного check-in за day.
// Возвращает (старый_streak, новый_streak, пересечённые_пороги).
func (s *StreakService) ApplyCheckIn(memberId int64, day time.Time) (int, int, []StreakThreshold, error) {
	streak, err := s.repo.Get(memberId)
	if err != nil {
		return 0, 0, nil, err
	}

	day = utils.MSKDay(day)
	prev := streak.CurrentStreak

	// Перевыдача freeze в начале новой ISO-недели.
	yr, wk := day.ISOWeek()
	if streak.FreezeWeekYear == nil || streak.FreezeWeekNum == nil ||
		*streak.FreezeWeekYear != yr || *streak.FreezeWeekNum != wk {
		streak.FreezesAvailable = 1
		streak.FreezeWeekYear = &yr
		streak.FreezeWeekNum = &wk
	}

	switch {
	case streak.LastCheckInDate == nil:
		streak.CurrentStreak = 1
	case utils.MSKDay(*streak.LastCheckInDate).Equal(day):
		// уже был check-in сегодня — ничего не пересчитываем
		return prev, streak.CurrentStreak, nil, nil
	default:
		gap := utils.DaysBetweenMSK(*streak.LastCheckInDate, day)
		switch {
		case gap == 1:
			streak.CurrentStreak++
		case gap == 2 && streak.FreezesAvailable > 0:
			// одна "заморозка" — пропуск 1 дня не сбрасывает стрик
			streak.CurrentStreak++
			streak.FreezesAvailable--
		default:
			streak.CurrentStreak = 1
		}
	}

	if streak.CurrentStreak > streak.LongestStreak {
		streak.LongestStreak = streak.CurrentStreak
	}
	streak.LastCheckInDate = &day
	streak.MemberId = memberId

	if err := s.repo.Save(streak); err != nil {
		return prev, streak.CurrentStreak, nil, err
	}

	crossed := make([]StreakThreshold, 0)
	for _, th := range StreakThresholds {
		if prev < th.Days && streak.CurrentStreak >= th.Days {
			crossed = append(crossed, th)
		}
	}
	if err := s.awardThresholdBonuses(memberId, day, crossed); err != nil {
		return prev, streak.CurrentStreak, crossed, err
	}

	return prev, streak.CurrentStreak, crossed, nil
}

// awardThresholdBonuses идемпотентно выдаёт баллы за пересечённые пороги.
// source_id = unix-timestamp дня пересечения, чтобы повторный набор того же
// порога после сброса стрика снова давал бонус.
func (s *StreakService) awardThresholdBonuses(memberId int64, day time.Time, crossed []StreakThreshold) error {
	for _, th := range crossed {
		tx := &models.PointTransaction{
			MemberId:    memberId,
			Amount:      th.Reward,
			Reason:      th.Reason,
			SourceType:  fmt.Sprintf("streak_%d", th.Days),
			SourceId:    day.Unix(),
			Description: fmt.Sprintf("Стрик %d дней", th.Days),
		}
		if err := s.pointRepo.AwardPoints(tx); err != nil {
			return err
		}
	}
	return nil
}

// BuildResponse возвращает StreakResponse, удобно сериализуемый в JSON.
func (s *StreakService) BuildResponse(memberId int64) (models.StreakResponse, error) {
	streak, err := s.repo.Get(memberId)
	if err != nil {
		return models.StreakResponse{}, err
	}

	resp := models.StreakResponse{
		Current:          streak.CurrentStreak,
		Longest:          streak.LongestStreak,
		FreezesAvailable: streak.FreezesAvailable,
		LastCheckIn:      streak.LastCheckInDate,
		Milestones:       make([]models.StreakMilestone, 0, len(StreakThresholds)),
	}

	for _, th := range StreakThresholds {
		resp.Milestones = append(resp.Milestones, models.StreakMilestone{
			Days:    th.Days,
			Reward:  th.Reward,
			Reached: streak.CurrentStreak >= th.Days,
		})
		if resp.NextThreshold == nil && streak.CurrentStreak < th.Days {
			next := th.Days
			to := th.Days - streak.CurrentStreak
			resp.NextThreshold = &next
			resp.DaysToNext = &to
		}
	}

	return resp, nil
}
