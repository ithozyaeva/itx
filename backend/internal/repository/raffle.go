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

// BuyTicketsTx вставляет count purchase-билетов одним SQL-запросом.
// Каждому билету достаётся уникальный source_id из raffle_ticket_purchase_seq —
// иначе UNIQUE (raffle_id, member_id, source_type, source_id) схлопнулся бы
// уже на втором билете в одной покупке.
func (r *RaffleRepository) BuyTicketsTx(db *gorm.DB, raffleId, memberId int64, count int) error {
	if count <= 0 {
		return nil
	}
	return db.Exec(
		`INSERT INTO raffle_tickets (raffle_id, member_id, source_type, source_id, bought_at)
		 SELECT ?, ?, ?, nextval('raffle_ticket_purchase_seq'), NOW()
		 FROM generate_series(1, ?)`,
		raffleId, memberId, models.RaffleTicketSourcePurchase, count,
	).Error
}

// AwardTicketTx идемпотентно выдаёт один билет за конкретную активность.
// Повторный вызов с тем же (raffleId, memberId, sourceType, sourceId) ничего
// не делает благодаря UNIQUE-индексу uniq_raffle_ticket_source.
// Возвращает (true, nil), если билет реально был создан.
func (r *RaffleRepository) AwardTicketTx(db *gorm.DB, raffleId, memberId int64, sourceType string, sourceId int64) (bool, error) {
	res := db.Exec(
		`INSERT INTO raffle_tickets (raffle_id, member_id, source_type, source_id, bought_at)
		 VALUES (?, ?, ?, ?, NOW())
		 ON CONFLICT (raffle_id, member_id, source_type, source_id) DO NOTHING`,
		raffleId, memberId, sourceType, sourceId,
	)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
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

// GetMemberTicketSources возвращает уникальные source_type'ы билетов юзера
// в указанной раффле — фронт по этому списку показывает «✓ check-in»,
// «✓ дейлик» и какие способы ещё можно использовать сегодня.
func (r *RaffleRepository) GetMemberTicketSources(raffleId, memberId int64) ([]string, error) {
	sources := make([]string, 0)
	err := database.DB.Raw(
		`SELECT DISTINCT source_type FROM raffle_tickets
		 WHERE raffle_id = ? AND member_id = ?`,
		raffleId, memberId,
	).Scan(&sources).Error
	return sources, err
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
