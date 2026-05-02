package service

import (
	"fmt"
	"log"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
)

// CheckInResult — внутренний результат успешного check-in для последующих
// побочных эффектов (SSE-публикация, Telegram-пуши).
type CheckInResult struct {
	Day              time.Time
	Inserted         bool
	PrevStreak       int
	CurrentStreak    int
	CrossedThreshold []StreakThreshold
	BasePoints       int
	RaffleEntered    bool
}

type CheckInService struct {
	repo            *repository.CheckInRepository
	pointRepo       *repository.PointsRepository
	streakSvc       *StreakService
	dailyRaffleSvc  *DailyRaffleService
}

func NewCheckInService() *CheckInService {
	return &CheckInService{
		repo:           repository.NewCheckInRepository(),
		pointRepo:      repository.NewPointsRepository(),
		streakSvc:      NewStreakService(),
		dailyRaffleSvc: NewDailyRaffleService(),
	}
}

// CheckIn — основной flow ежедневного check-in. Идемпотентен:
// повторный вызов в тот же МСК-день возвращает текущее состояние без
// дублирования начислений или билетов.
func (s *CheckInService) CheckIn(memberId int64) (*CheckInResult, error) {
	day := utils.MSKToday()

	inserted, err := s.repo.Insert(memberId, day)
	if err != nil {
		return nil, fmt.Errorf("записать check-in: %w", err)
	}

	if !inserted {
		// Уже был check-in сегодня — отдадим текущее состояние стрика.
		streak, err := s.streakSvc.repo.Get(memberId)
		if err != nil {
			return nil, err
		}
		return &CheckInResult{
			Day:           day,
			Inserted:      false,
			PrevStreak:    streak.CurrentStreak,
			CurrentStreak: streak.CurrentStreak,
		}, nil
	}

	prev, current, crossed, err := s.streakSvc.ApplyCheckIn(memberId, day)
	if err != nil {
		log.Printf("apply check-in streak (member=%d): %v", memberId, err)
	}

	if err := s.awardBase(memberId, day); err != nil {
		log.Printf("award base check-in points (member=%d): %v", memberId, err)
	}

	// Пуш на пересечение порога стрика — асинхронно, чтобы не блокировать
	// ответ хендлера; функция сама уважает MuteAll/DailyStreak.
	for _, th := range crossed {
		go PushStreakThreshold(memberId, th.Days, th.Reward)
	}

	raffleEntered := false
	if err := s.dailyRaffleSvc.EnterRaffle(memberId); err != nil {
		log.Printf("enter daily raffle (member=%d): %v", memberId, err)
	} else {
		raffleEntered = true
	}

	// Челлендж-метрики, привязанные к check-in.
	TrackChallengeMetric(memberId, "check_ins", 1)

	return &CheckInResult{
		Day:              day,
		Inserted:         true,
		PrevStreak:       prev,
		CurrentStreak:    current,
		CrossedThreshold: crossed,
		BasePoints:       models.PointValues[models.PointReasonDailyCheckIn],
		RaffleEntered:    raffleEntered,
	}, nil
}

// awardBase идемпотентно начисляет базовые +5 за check-in (ключ — unix-день).
func (s *CheckInService) awardBase(memberId int64, day time.Time) error {
	amount := models.PointValues[models.PointReasonDailyCheckIn]
	tx := &models.PointTransaction{
		MemberId:    memberId,
		Amount:      amount,
		Reason:      models.PointReasonDailyCheckIn,
		SourceType:  "check_in",
		SourceId:    day.Unix(),
		Description: "Ежедневный check-in",
	}
	if err := s.pointRepo.AwardPoints(tx); err != nil {
		return err
	}
	if amount > 0 {
		TrackChallengeMetric(memberId, "points_earned", amount)
	}
	return nil
}

// HasCheckedInToday — проксирует в repo, удобно для GET /dailies/today.
func (s *CheckInService) HasCheckedInToday(memberId int64) (bool, *time.Time, error) {
	day := utils.MSKToday()
	c, err := s.repo.Get(memberId, day)
	if err != nil {
		return false, nil, nil
	}
	return true, &c.CreatedAt, nil
}
