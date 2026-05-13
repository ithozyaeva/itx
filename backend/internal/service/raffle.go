package service

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"

	"gorm.io/gorm"
)

// MaxBuyTicketsPerRequest — верхний лимит на одну покупку. Защита от
// случайного/злонамеренного запроса с count=1_000_000, который раньше вешал
// БД сотней тысяч раундтрипов и оставлял зомби-транзакции.
const MaxBuyTicketsPerRequest = 10_000

type RaffleService struct {
	repo      *repository.RaffleRepository
	pointRepo *repository.PointsRepository
}

func NewRaffleService() *RaffleService {
	return &RaffleService{
		repo:      repository.NewRaffleRepository(),
		pointRepo: repository.NewPointsRepository(),
	}
}

func (s *RaffleService) GetAll(memberId int64) ([]models.RafflePublic, error) {
	return s.repo.GetAll(memberId)
}

func (s *RaffleService) GetAllAdmin() ([]models.Raffle, error) {
	return s.repo.GetAllAdmin()
}

func (s *RaffleService) Create(raffle *models.Raffle) error {
	return s.repo.Create(raffle)
}

func (s *RaffleService) Update(raffle *models.Raffle) error {
	return s.repo.Update(raffle)
}

func (s *RaffleService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *RaffleService) BuyTickets(raffleId, memberId int64, count int) error {
	if count <= 0 {
		count = 1
	}
	if count > MaxBuyTicketsPerRequest {
		return fmt.Errorf("за один запрос нельзя купить больше %d билетов", MaxBuyTicketsPerRequest)
	}

	raffle, err := s.repo.GetById(raffleId)
	if err != nil {
		return fmt.Errorf("розыгрыш не найден")
	}

	if raffle.Status != models.RaffleStatusActive {
		return fmt.Errorf("розыгрыш завершён")
	}

	// Daily-раффлы (entry_rule=auto_check_in) выдают билеты автоматически
	// за активности; ручная покупка через этот эндпоинт запрещена.
	if raffle.EntryRule != models.RaffleEntryRulePurchase {
		return fmt.Errorf("этот розыгрыш не поддерживает покупку билетов")
	}

	if raffle.MaxTickets > 0 {
		total, _ := s.repo.GetTicketCount(raffleId)
		if int(total)+count > raffle.MaxTickets {
			return fmt.Errorf("превышен лимит билетов")
		}
	}

	totalCost := raffle.TicketCost * count

	// Балансовая проверка ВНУТРИ tx с pg_advisory_xact_lock(memberId) —
	// иначе параллельные BuyTickets читают одинаковый SUM (READ COMMITTED),
	// оба проходят check и оба INSERT'ят дебит, юзер уходит в минус
	// (INSERT'ы дебитов не конфликтуют сериализуемо). Mirroring casino.go
	// PlaceBet, где тот же класс гонки уже починен advisory-lock'ом.
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`SELECT pg_advisory_xact_lock(?)`, memberId).Error; err != nil {
			return err
		}
		var balance int
		if err := tx.Raw(
			`SELECT COALESCE(SUM(amount), 0) FROM point_transactions WHERE member_id = ?`,
			memberId,
		).Scan(&balance).Error; err != nil {
			return err
		}
		if balance < totalCost {
			return fmt.Errorf("недостаточно баллов (нужно %d, доступно %d)", totalCost, balance)
		}
		pt := &models.PointTransaction{
			MemberId:    memberId,
			Amount:      -totalCost,
			Reason:      models.PointReasonRaffleSpend,
			SourceType:  "raffle",
			SourceId:    raffleId,
			Description: fmt.Sprintf("Покупка %d билетов: %s", count, raffle.Title),
		}
		if err := tx.Create(pt).Error; err != nil {
			return err
		}
		return s.repo.BuyTicketsTx(tx, raffleId, memberId, count)
	})
	if err != nil {
		return err
	}

	// Дейлик «купить билет в обычный розыгрыш» — только за manual.
	if raffle.Kind != models.RaffleKindDaily {
		TrackDailyTrigger(memberId, "buy_raffle_ticket", 1)
	}

	return nil
}

func (s *RaffleService) DrawExpiredRaffles() {
	raffles, err := s.repo.GetExpired()
	if err != nil {
		log.Printf("Error getting expired raffles: %v", err)
		return
	}

	dailySvc := NewDailyRaffleService()

	for _, raffle := range raffles {
		ticketCount, _ := s.repo.GetTicketCount(raffle.Id)
		if ticketCount == 0 {
			s.repo.FinishRaffle(raffle.Id, 0)
			log.Printf("Raffle %d finished with no participants", raffle.Id)
			continue
		}

		winnerId, err := s.repo.PickRandomWinner(raffle.Id)
		if err != nil {
			log.Printf("Error picking winner for raffle %d: %v", raffle.Id, err)
			continue
		}

		if err := s.repo.FinishRaffle(raffle.Id, winnerId); err != nil {
			log.Printf("Error finishing raffle %d: %v", raffle.Id, err)
			continue
		}

		log.Printf("Raffle %d winner: member %d", raffle.Id, winnerId)

		// Daily-раффл — выдаём приз победителю автоматически.
		// Manual-розыгрыши обрабатываются админом отдельно (через UI/manual).
		if raffle.Kind == models.RaffleKindDaily {
			if err := dailySvc.AwardWinPoints(winnerId, raffle.Id); err != nil {
				log.Printf("award daily-raffle win points (raffle=%d, member=%d): %v",
					raffle.Id, winnerId, err)
			}
			GetSSEHub().Publish(winnerId, SSEEvent{Type: "points"})
			go PushDailyRaffleWin(winnerId, raffle.Prize)
		}
		GetSSEHub().Broadcast(SSEEvent{Type: "raffles"})
	}
}
