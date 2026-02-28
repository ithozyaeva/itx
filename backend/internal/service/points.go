package service

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"
	"time"
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
	}
}

func (s *PointsService) AwardEventPoints(event *models.Event) error {
	for _, host := range event.Hosts {
		s.AwardIdempotent(host.Id, models.PointReasonEventHost, "event", event.Id,
			fmt.Sprintf("Проведение события: %s", event.Title))
	}

	for _, member := range event.Members {
		s.AwardIdempotent(member.Id, models.PointReasonEventAttend, "event", event.Id,
			fmt.Sprintf("Участие в событии: %s", event.Title))
	}

	return nil
}

// CheckProfileComplete проверяет заполненность профиля и начисляет одноразовый бонус.
func (s *PointsService) CheckProfileComplete(member *models.Member) {
	if member.FirstName == "" || member.LastName == "" || member.Bio == "" || member.Birthday == nil || member.AvatarURL == "" {
		return
	}
	s.AwardIdempotent(member.Id, models.PointReasonProfileComplete, "profile", member.Id,
		"Полностью заполненный профиль")
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

// AwardActivityBonuses начисляет бонусы за активность: еженедельная активность, 3+ события в месяц, серия 4 недели.
func (s *PointsService) AwardActivityBonuses() {
	now := time.Now()
	year, week := now.ISOWeek()

	// Еженедельная активность (проверяем предыдущую неделю)
	prevWeekTime := now.AddDate(0, 0, -7)
	prevYear, prevWeek := prevWeekTime.ISOWeek()

	weeklyMembers, err := s.repo.GetMembersWithEventsInWeek(prevYear, prevWeek)
	if err != nil {
		log.Printf("Error getting weekly active members: %v", err)
	} else {
		sourceId := int64(prevYear*100 + prevWeek)
		for _, memberId := range weeklyMembers {
			s.AwardIdempotent(memberId, models.PointReasonWeeklyActivity, "weekly", sourceId,
				fmt.Sprintf("Активность на неделе %d/%d", prevWeek, prevYear))
		}
	}

	// 3+ события за прошлый месяц
	prevMonthTime := now.AddDate(0, -1, 0)
	monthYear := prevMonthTime.Year()
	prevMonth := int(prevMonthTime.Month())

	monthlyMembers, err := s.repo.GetMembersWithMonthlyEvents(monthYear, prevMonth, 3)
	if err != nil {
		log.Printf("Error getting monthly active members: %v", err)
	} else {
		sourceId := int64(monthYear*100 + prevMonth)
		for _, memberId := range monthlyMembers {
			s.AwardIdempotent(memberId, models.PointReasonMonthlyActive, "monthly", sourceId,
				fmt.Sprintf("3+ событий за %d/%d", prevMonth, monthYear))
		}
	}

	// Серия 4 недели подряд
	streakMembers, err := s.repo.GetMembersWithStreak(4)
	if err != nil {
		log.Printf("Error getting streak members: %v", err)
	} else {
		sourceId := int64(year*100 + week)
		for _, memberId := range streakMembers {
			s.AwardIdempotent(memberId, models.PointReasonStreak4Weeks, "streak", sourceId,
				"Серия: 4 недели подряд с событиями")
		}
	}
}
