package service

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"
	"time"

	"gorm.io/gorm"
)

type PointsService struct {
	repo *repository.PointsRepository
}

func NewPointsService() *PointsService {
	return &PointsService{
		repo: repository.NewPointsRepository(),
	}
}

// AwardIdempotent начисляет баллы с защитой от дублирования (ON CONFLICT DO NOTHING).
// Используется для планировщика и одноразовых наград.
func (s *PointsService) AwardIdempotent(memberId int64, reason models.PointReason, sourceType string, sourceId int64, description string) {
	tx := &models.PointTransaction{
		MemberId:    memberId,
		Amount:      models.PointValues[reason],
		Reason:      reason,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
	}
	if err := s.repo.AwardPoints(tx); err != nil {
		log.Printf("Error awarding points (reason=%s, member=%d): %v", reason, memberId, err)
		return
	}
	if tx.Amount > 0 {
		TrackChallengeMetric(memberId, "points_earned", tx.Amount)
	}
}

// GiveForAction начисляет баллы за действие пользователя (обычный INSERT).
// Используется в хендлерах, где каждый вызов = одно уникальное действие.
func (s *PointsService) GiveForAction(memberId int64, reason models.PointReason, sourceType string, sourceId int64, description string) {
	tx := &models.PointTransaction{
		MemberId:    memberId,
		Amount:      models.PointValues[reason],
		Reason:      reason,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
	}
	if err := s.repo.GivePoints(tx); err != nil {
		log.Printf("Error giving points (reason=%s, member=%d): %v", reason, memberId, err)
		return
	}
	if tx.Amount > 0 {
		TrackChallengeMetric(memberId, "points_earned", tx.Amount)
	}
}

// GiveCustomPoints начисляет произвольное количество баллов (для квестов чатов)
func (s *PointsService) GiveCustomPoints(memberId int64, amount int, reason models.PointReason, sourceType string, sourceId int64, description string) {
	tx := &models.PointTransaction{
		MemberId:    memberId,
		Amount:      amount,
		Reason:      reason,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
	}
	if err := s.repo.GivePoints(tx); err != nil {
		log.Printf("Error giving custom points (reason=%s, member=%d, amount=%d): %v", reason, memberId, amount, err)
	}
}

// awardIdempotentTx — внутренний помощник для атомарного начисления внутри транзакции.
func (s *PointsService) awardIdempotentTx(db *gorm.DB, memberId int64, reason models.PointReason, sourceType string, sourceId int64, description string) error {
	tx := &models.PointTransaction{
		MemberId:    memberId,
		Amount:      models.PointValues[reason],
		Reason:      reason,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
	}
	return s.repo.AwardPointsTx(db, tx)
}

func (s *PointsService) AwardEventPoints(event *models.Event) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		for _, host := range event.Hosts {
			if err := s.awardIdempotentTx(tx, host.Id, models.PointReasonEventHost, "event", event.Id,
				fmt.Sprintf("Проведение события: %s", event.Title)); err != nil {
				return err
			}
		}

		for _, member := range event.Members {
			if err := s.awardIdempotentTx(tx, member.Id, models.PointReasonEventAttend, "event", event.Id,
				fmt.Sprintf("Участие в событии: %s", event.Title)); err != nil {
				return err
			}
			TrackDailyTrigger(member.Id, "attend_event", 1)
			TrackChallengeMetric(member.Id, "events_attended", 1)
		}

		return nil
	})
}

// CheckProfileComplete проверяет заполненность профиля и начисляет одноразовый бонус.
func (s *PointsService) CheckProfileComplete(member *models.Member) {
	if member.FirstName == "" || member.LastName == "" || member.Bio == "" || member.Birthday == nil {
		return
	}
	s.AwardIdempotent(member.Id, models.PointReasonProfileComplete, "profile", member.Id,
		"Полностью заполненный профиль")
}

func (s *PointsService) GetBalance(memberId int64) (int, error) {
	return s.repo.GetBalance(memberId)
}

func (s *PointsService) GetMyPoints(memberId int64) (*models.PointsSummary, error) {
	balance, err := s.repo.GetBalance(memberId)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.GetTransactions(memberId, 50)
	if err != nil {
		return nil, err
	}

	return &models.PointsSummary{
		Balance:      balance,
		Transactions: transactions,
	}, nil
}

func (s *PointsService) GetLeaderboard(limit int) ([]models.MemberPointsBalance, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.GetLeaderboard(limit)
}

func (s *PointsService) AwardPointsForPastEvents() {
	events, err := s.repo.GetPastEventsForAward(7)
	if err != nil {
		log.Printf("Error fetching past events for points: %v", err)
		return
	}

	for _, event := range events {
		if err := s.AwardEventPoints(&event); err != nil {
			log.Printf("Error awarding points for event %d: %v", event.Id, err)
		}
	}

	if len(events) > 0 {
		log.Printf("Processed points for %d past events", len(events))
	}
}

func (s *PointsService) SearchTransactions(username *string, limit, offset int) ([]models.AdminPointTransaction, int64, error) {
	return s.repo.SearchTransactions(username, limit, offset)
}

func (s *PointsService) AdminAwardPoints(memberId int64, amount int, description string) error {
	tx := &models.PointTransaction{
		MemberId:    memberId,
		Amount:      amount,
		Reason:      models.PointReasonAdminManual,
		SourceType:  "admin",
		SourceId:    0,
		Description: description,
	}
	return s.repo.CreateManualTransaction(tx)
}

func (s *PointsService) DeleteTransaction(id int64) error {
	return s.repo.DeleteTransaction(id)
}

// isoWeekMonday возвращает дату понедельника ISO-недели (UTC 00:00:00).
func isoWeekMonday(year, week int) time.Time {
	// 4 января всегда в первой ISO-неделе года
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	// weekday: Monday=0 ... Sunday=6
	weekday := int(jan4.Weekday()+6) % 7
	// Monday of week 1
	week1Monday := jan4.AddDate(0, 0, -weekday)
	return week1Monday.AddDate(0, 0, (week-1)*7)
}

// AwardWeeklyChatter находит самого активного участника чата за прошлую неделю и начисляет ему 15 баллов.
func (s *PointsService) AwardWeeklyChatter() {
	now := time.Now().UTC()
	prevWeek := now.AddDate(0, 0, -7)
	prevYear, prevWeekNum := prevWeek.ISOWeek()
	monday := isoWeekMonday(prevYear, prevWeekNum)

	memberId, err := s.repo.GetTopChatterForWeek(monday)
	if err != nil {
		log.Printf("AwardWeeklyChatter: error querying top chatter: %v", err)
		return
	}
	if memberId == 0 {
		return
	}

	sourceId := int64(prevYear*100 + prevWeekNum)
	s.AwardIdempotent(memberId, models.PointReasonChatterOfWeek, "weekly_chatter", sourceId,
		fmt.Sprintf("Чаттер недели %d/%d", prevWeekNum, prevYear))
	log.Printf("AwardWeeklyChatter: awarded member %d for week %d/%d", memberId, prevWeekNum, prevYear)
}

// AwardActivityBonuses начисляет бонусы за активность: еженедельная активность, 3+ события в месяц, серия 4 недели.
// Каждый бонус-блок выполняется в собственной транзакции, чтобы сбой одного типа
// не откатывал успешные начисления другого.
func (s *PointsService) AwardActivityBonuses() {
	now := time.Now()
	year, week := now.ISOWeek()

	// Еженедельная активность (проверяем предыдущую неделю)
	prevWeekTime := now.AddDate(0, 0, -7)
	prevYear, prevWeek := prevWeekTime.ISOWeek()

	if weeklyMembers, err := s.repo.GetMembersWithEventsInWeek(prevYear, prevWeek); err != nil {
		log.Printf("Error getting weekly active members: %v", err)
	} else {
		sourceId := int64(prevYear*100 + prevWeek)
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			for _, memberId := range weeklyMembers {
				if err := s.awardIdempotentTx(tx, memberId, models.PointReasonWeeklyActivity, "weekly", sourceId,
					fmt.Sprintf("Активность на неделе %d/%d", prevWeek, prevYear)); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("AwardActivityBonuses weekly tx error: %v", err)
		}
	}

	// 3+ события за прошлый месяц
	prevMonthTime := now.AddDate(0, -1, 0)
	monthYear := prevMonthTime.Year()
	prevMonth := int(prevMonthTime.Month())

	if monthlyMembers, err := s.repo.GetMembersWithMonthlyEvents(monthYear, prevMonth, 3); err != nil {
		log.Printf("Error getting monthly active members: %v", err)
	} else {
		sourceId := int64(monthYear*100 + prevMonth)
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			for _, memberId := range monthlyMembers {
				if err := s.awardIdempotentTx(tx, memberId, models.PointReasonMonthlyActive, "monthly", sourceId,
					fmt.Sprintf("3+ событий за %d/%d", prevMonth, monthYear)); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("AwardActivityBonuses monthly tx error: %v", err)
		}
	}

	// Серия 4 недели подряд
	if streakMembers, err := s.repo.GetMembersWithStreak(4); err != nil {
		log.Printf("Error getting streak members: %v", err)
	} else {
		sourceId := int64(year*100 + week)
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			for _, memberId := range streakMembers {
				if err := s.awardIdempotentTx(tx, memberId, models.PointReasonStreak4Weeks, "streak", sourceId,
					"Серия: 4 недели подряд с событиями"); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("AwardActivityBonuses streak tx error: %v", err)
		}
	}
}
