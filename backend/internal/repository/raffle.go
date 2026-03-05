package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

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
	var items []models.RafflePublic
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
	return items, err
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
	tx := database.DB.Begin()
	for i := 0; i < count; i++ {
		ticket := &models.RaffleTicket{
			RaffleId: raffleId,
			MemberId: memberId,
		}
		if err := tx.Create(ticket).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
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
	return database.DB.Model(&models.Raffle{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    models.RaffleStatusFinished,
			"winner_id": winnerId,
		}).Error
}

func (r *RaffleRepository) Update(raffle *models.Raffle) error {
	return database.DB.Save(raffle).Error
}

func (r *RaffleRepository) Delete(id int64) error {
	return database.DB.Delete(&models.Raffle{}, id).Error
}
