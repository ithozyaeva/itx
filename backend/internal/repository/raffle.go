package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

const buyTicketsBatchSize = 1000

type RaffleRepository struct{}

func NewRaffleRepository() *RaffleRepository {
	return &RaffleRepository{}
}

func (r *RaffleRepository) Create(raffle *models.Raffle) error {
	return database.DB.Create(raffle).Error
}

func (r *RaffleRepository) GetById(id int64) (*models.Raffle, error) {
	var raffle models.Raffle
	err := database.DB.First(&raffle, id).Error
	if err != nil {
		return nil, err
	}
	return &raffle, nil
}

func (r *RaffleRepository) GetAll(memberId int64) ([]models.RafflePublic, error) {
	items := make([]models.RafflePublic, 0)
	err := database.DB.Raw(`
		SELECT r.id, r.title, r.description, r.prize, r.ticket_cost, r.max_tickets,
			r.ends_at, r.status, r.winner_id,
			w.first_name as winner_first_name, w.last_name as winner_last_name,
			w.username as winner_username, w.avatar_url as winner_avatar_url,
			(SELECT COUNT(*) FROM raffle_tickets WHERE raffle_id = r.id) as total_tickets,
			(SELECT COUNT(*) FROM raffle_tickets WHERE raffle_id = r.id AND member_id = ?) as my_tickets
		FROM raffles r
		LEFT JOIN members w ON w.id = r.winner_id
		ORDER BY r.status ASC, r.ends_at ASC
	`, memberId).Scan(&items).Error

	for i := range items {
		items[i].AfterFind(nil)
	}

	return items, err
}

func (r *RaffleRepository) GetAllAdmin() ([]models.Raffle, error) {
	var raffles []models.Raffle
	err := database.DB.Order("status ASC, ends_at ASC").Find(&raffles).Error
	return raffles, err
}

func (r *RaffleRepository) GetActive() ([]models.Raffle, error) {
	var raffles []models.Raffle
	err := database.DB.Where("status = ? AND ends_at > NOW()", models.RaffleStatusActive).Find(&raffles).Error
	return raffles, err
}

func (r *RaffleRepository) GetExpired() ([]models.Raffle, error) {
	var raffles []models.Raffle
	err := database.DB.Where("status = ? AND ends_at <= NOW()", models.RaffleStatusActive).Find(&raffles).Error
	return raffles, err
}

func (r *RaffleRepository) BuyTickets(raffleId, memberId int64, count int) error {
	return r.BuyTicketsTx(database.DB, raffleId, memberId, count)
}

// BuyTicketsTx вставляет count билетов одним батчем в указанной транзакции.
// CreateInBatches шлёт INSERT ... VALUES (...),(...),... вместо count отдельных
// раундтрипов — критично для больших count (раньше count=1M вешало транзакцию).
func (r *RaffleRepository) BuyTicketsTx(db *gorm.DB, raffleId, memberId int64, count int) error {
	if count <= 0 {
		return nil
	}
	tickets := make([]models.RaffleTicket, count)
	for i := range tickets {
		tickets[i] = models.RaffleTicket{RaffleId: raffleId, MemberId: memberId}
	}
	return db.CreateInBatches(tickets, buyTicketsBatchSize).Error
}

func (r *RaffleRepository) GetTicketCount(raffleId int64) (int64, error) {
	var count int64
	err := database.DB.Model(&models.RaffleTicket{}).Where("raffle_id = ?", raffleId).Count(&count).Error
	return count, err
}

func (r *RaffleRepository) GetMemberTicketCount(raffleId, memberId int64) (int64, error) {
	var count int64
	err := database.DB.Model(&models.RaffleTicket{}).
		Where("raffle_id = ? AND member_id = ?", raffleId, memberId).
		Count(&count).Error
	return count, err
}

func (r *RaffleRepository) PickRandomWinner(raffleId int64) (int64, error) {
	var memberId int64
	err := database.DB.Raw(`
		SELECT member_id FROM raffle_tickets
		WHERE raffle_id = ?
		ORDER BY RANDOM() LIMIT 1
	`, raffleId).Scan(&memberId).Error
	return memberId, err
}

func (r *RaffleRepository) FinishRaffle(id int64, winnerId int64) error {
	updates := map[string]interface{}{
		"status": models.RaffleStatusFinished,
	}
	if winnerId != 0 {
		updates["winner_id"] = winnerId
	}
	return database.DB.Model(&models.Raffle{}).Where("id = ?", id).
		Updates(updates).Error
}

func (r *RaffleRepository) Update(raffle *models.Raffle) error {
	return database.DB.Save(raffle).Error
}

func (r *RaffleRepository) Delete(id int64) error {
	return database.DB.Delete(&models.Raffle{}, id).Error
}
