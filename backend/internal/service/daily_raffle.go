package service

import (
	"errors"
	"fmt"
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"gorm.io/gorm"
)

const (
	dailyRafflePrize       = "100 баллов"
	dailyRaffleTitle       = "Ежедневный розыгрыш"
	dailyRaffleDescription = "Каждый день в 23:59 МСК случайный участник получает 100 баллов. " +
		"Чтобы попасть в розыгрыш — просто сделай check-in."
)

type DailyRaffleService struct {
	repo *repository.RaffleRepository
}

func NewDailyRaffleService() *DailyRaffleService {
	return &DailyRaffleService{
		repo: repository.NewRaffleRepository(),
	}
}

// EnsureTodayRaffle создаёт сегодняшний daily-раффл, если его ещё нет.
// Идемпотентен через UNIQUE-индекс uniq_raffle_day_key WHERE kind='daily'.
// Срок действия — до 23:59:59 МСК того же дня; розыгрыш проводится
// существующим cron'ом DrawExpiredRaffles (5-минутный тикер).
func (s *DailyRaffleService) EnsureTodayRaffle() (*models.Raffle, error) {
	day := utils.MSKToday()

	if existing, err := s.GetByDayKey(day); err == nil && existing != nil {
		return existing, nil
	}

	endsAt := utils.MSKEndOfDay(day)
	raffle := &models.Raffle{
		Title:       dailyRaffleTitle,
		Description: dailyRaffleDescription,
		Prize:       dailyRafflePrize,
		TicketCost:  0,
		MaxTickets:  0,
		EndsAt:      endsAt,
		Status:      models.RaffleStatusActive,
		Kind:        models.RaffleKindDaily,
		EntryRule:   models.RaffleEntryRuleAutoCheckIn,
		DayKey:      &day,
	}

	// Используем raw INSERT с ON CONFLICT, чтобы UNIQUE-индекс
	// uniq_raffle_day_key корректно отрабатывал при гонке между cron'ом
	// и watchdog'ом. GORM Create на ON CONFLICT не предоставляет
	// одной строкой через partial-unique.
	res := database.DB.Exec(
		`INSERT INTO raffles (title, description, prize, ticket_cost, max_tickets, ends_at, status, kind, entry_rule, day_key)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (day_key) WHERE kind = 'daily' DO NOTHING`,
		raffle.Title, raffle.Description, raffle.Prize,
		raffle.TicketCost, raffle.MaxTickets, raffle.EndsAt,
		raffle.Status, raffle.Kind, raffle.EntryRule, raffle.DayKey,
	)
	if res.Error != nil {
		return nil, res.Error
	}

	return s.GetByDayKey(day)
}

// EnterRaffle вставляет билет от memberId в сегодняшний daily-раффл.
// Идемпотентно: повторный вызов в тот же день не плодит билеты
// (UNIQUE-проверка через выбор существующих билетов конкретного юзера).
func (s *DailyRaffleService) EnterRaffle(memberId int64) error {
	raffle, err := s.EnsureTodayRaffle()
	if err != nil {
		return err
	}
	if raffle == nil {
		return errors.New("daily-раффл недоступен")
	}

	// Проверяем существующий билет, чтобы не дублировать.
	count, err := s.repo.GetMemberTicketCount(raffle.Id, memberId)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	return s.repo.BuyTickets(raffle.Id, memberId, 1)
}

// GetByDayKey — выбор daily-раффла по конкретной МСК-дате.
func (s *DailyRaffleService) GetByDayKey(day time.Time) (*models.Raffle, error) {
	var raffle models.Raffle
	err := database.DB.
		Where("kind = ? AND day_key = ?", models.RaffleKindDaily, day).
		First(&raffle).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &raffle, nil
}

// GetTodayPublic — состояние сегодняшнего розыгрыша для фронта,
// включая счётчик участников и факт «я участвую».
func (s *DailyRaffleService) GetTodayPublic(memberId int64) (*models.RafflePublic, error) {
	raffle, err := s.EnsureTodayRaffle()
	if err != nil {
		return nil, err
	}
	if raffle == nil {
		return nil, nil
	}

	total, _ := s.repo.GetTicketCount(raffle.Id)
	myTickets, _ := s.repo.GetMemberTicketCount(raffle.Id, memberId)

	pub := &models.RafflePublic{
		Id:           raffle.Id,
		Title:        raffle.Title,
		Description:  raffle.Description,
		Prize:        raffle.Prize,
		TicketCost:   raffle.TicketCost,
		MaxTickets:   raffle.MaxTickets,
		EndsAt:       raffle.EndsAt,
		Status:       raffle.Status,
		Kind:         raffle.Kind,
		EntryRule:    raffle.EntryRule,
		DayKey:       raffle.DayKey,
		TotalTickets: int(total),
		MyTickets:    int(myTickets),
		WinnerId:     raffle.WinnerId,
	}
	return pub, nil
}

// AwardWinPoints начисляет 100 баллов победителю daily-раффла.
// Идемпотентно через source_type='raffle', source_id=raffle.id и reason
// 'daily_raffle_win' (отдельная пара ключей от 'raffle_spend').
func (s *DailyRaffleService) AwardWinPoints(memberId, raffleId int64) error {
	pointRepo := repository.NewPointsRepository()
	amount := models.PointValues[models.PointReasonDailyRaffleWin]
	tx := &models.PointTransaction{
		MemberId:    memberId,
		Amount:      amount,
		Reason:      models.PointReasonDailyRaffleWin,
		SourceType:  "daily_raffle",
		SourceId:    raffleId,
		Description: fmt.Sprintf("Победа в ежедневном розыгрыше #%d", raffleId),
	}
	if err := pointRepo.AwardPoints(tx); err != nil {
		return err
	}
	if amount > 0 {
		TrackChallengeMetric(memberId, "points_earned", amount)
	}
	return nil
}
