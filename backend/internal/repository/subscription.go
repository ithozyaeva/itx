package repository

import (
	"encoding/json"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{db: database.DB}
}

// --- Tiers ---

func (r *SubscriptionRepository) GetTier(id uint) (*models.SubscriptionTier, error) {
	var tier models.SubscriptionTier
	if err := r.db.First(&tier, id).Error; err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *SubscriptionRepository) GetTierBySlug(slug string) (*models.SubscriptionTier, error) {
	var tier models.SubscriptionTier
	if err := r.db.Where("slug = ?", slug).First(&tier).Error; err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *SubscriptionRepository) GetAllTiers() ([]models.SubscriptionTier, error) {
	var tiers []models.SubscriptionTier
	if err := r.db.Order("level ASC").Find(&tiers).Error; err != nil {
		return nil, err
	}
	return tiers, nil
}

func (r *SubscriptionRepository) GetPublicTiers() ([]models.SubscriptionTier, error) {
	var tiers []models.SubscriptionTier
	if err := r.db.Where("is_public = ?", true).Order("level ASC").Find(&tiers).Error; err != nil {
		return nil, err
	}
	return tiers, nil
}

func (r *SubscriptionRepository) GetAllTiersDesc() ([]models.SubscriptionTier, error) {
	var tiers []models.SubscriptionTier
	if err := r.db.Order("level DESC").Find(&tiers).Error; err != nil {
		return nil, err
	}
	return tiers, nil
}

// --- Chats ---

func (r *SubscriptionRepository) GetChat(chatID int64) (*models.SubscriptionChat, error) {
	var chat models.SubscriptionChat
	if err := r.db.First(&chat, chatID).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *SubscriptionRepository) GetAllChats() ([]models.SubscriptionChat, error) {
	var chats []models.SubscriptionChat
	if err := r.db.Order("id").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *SubscriptionRepository) GetAnchorChats() ([]models.SubscriptionChat, error) {
	var chats []models.SubscriptionChat
	if err := r.db.Where("anchor_for_tier_id IS NOT NULL").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *SubscriptionRepository) GetChatsForTierLevel(level int) ([]models.SubscriptionChat, error) {
	var chats []models.SubscriptionChat
	err := r.db.
		Joins("JOIN subscription_tier_chats ON subscription_tier_chats.chat_id = subscription_chats.id").
		Joins("JOIN subscription_tiers ON subscription_tiers.id = subscription_tier_chats.tier_id").
		Where("subscription_tiers.level <= ? AND subscription_chats.anchor_for_tier_id IS NULL", level).
		Distinct().
		Find(&chats).Error
	return chats, err
}

func (r *SubscriptionRepository) UpsertChat(chatID int64, title string, chatType string) error {
	var chat models.SubscriptionChat
	err := r.db.First(&chat, chatID).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(&models.SubscriptionChat{
			ID:       chatID,
			Title:    title,
			ChatType: chatType,
		}).Error
	}
	if err != nil {
		return err
	}
	chat.Title = title
	chat.ChatType = chatType
	return r.db.Save(&chat).Error
}

func (r *SubscriptionRepository) UpdateChatMeta(chatID int64, category, emoji *string) error {
	return r.db.Model(&models.SubscriptionChat{}).Where("id = ?", chatID).
		Updates(map[string]interface{}{"category": category, "emoji": emoji}).Error
}

func (r *SubscriptionRepository) UpdateChatPriority(chatID int64, priority int) error {
	return r.db.Model(&models.SubscriptionChat{}).Where("id = ?", chatID).
		Update("priority", priority).Error
}

func (r *SubscriptionRepository) SetAnchor(chatID int64, tierID *uint) error {
	return r.db.Model(&models.SubscriptionChat{}).Where("id = ?", chatID).
		Update("anchor_for_tier_id", tierID).Error
}

func (r *SubscriptionRepository) AddChatToTier(chatID int64, tierID uint) error {
	return r.db.Create(&models.SubscriptionTierChat{
		TierID: tierID,
		ChatID: chatID,
	}).Error
}

func (r *SubscriptionRepository) GetTierIDsForChat(chatID int64) ([]uint, error) {
	var tierChats []models.SubscriptionTierChat
	if err := r.db.Where("chat_id = ?", chatID).Find(&tierChats).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, len(tierChats))
	for i, tc := range tierChats {
		ids[i] = tc.TierID
	}
	return ids, nil
}

func (r *SubscriptionRepository) GetAllTierChats() (map[int64][]uint, error) {
	var tierChats []models.SubscriptionTierChat
	if err := r.db.Find(&tierChats).Error; err != nil {
		return nil, err
	}
	m := make(map[int64][]uint)
	for _, tc := range tierChats {
		m[tc.ChatID] = append(m[tc.ChatID], tc.TierID)
	}
	return m, nil
}

func (r *SubscriptionRepository) SetChatTiers(chatID int64, tierIDs []uint) error {
	r.db.Where("chat_id = ?", chatID).Delete(&models.SubscriptionTierChat{})
	for _, tierID := range tierIDs {
		if err := r.db.Create(&models.SubscriptionTierChat{TierID: tierID, ChatID: chatID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *SubscriptionRepository) DeleteChat(chatID int64) error {
	r.db.Where("chat_id = ?", chatID).Delete(&models.SubscriptionTierChat{})
	r.db.Where("chat_id = ?", chatID).Delete(&models.SubscriptionUserChatAccess{})
	return r.db.Delete(&models.SubscriptionChat{}, chatID).Error
}

// GetUserEffectiveTierLevel возвращает уровень эффективного тира пользователя
// (manual override либо resolved). Если tier не назначен — 0, ok=false.
// Используется middleware'ом RequireMinTier для гейта по уровню подписки.
func (r *SubscriptionRepository) GetUserEffectiveTierLevel(userID int64) (int, bool) {
	type row struct{ Level int }
	var res row
	// effective_tier = manual_tier_id если manual ещё не истёк, иначе
	// resolved_tier_id. Без проверки expires юзер с просроченной покупкой
	// держит API-доступ к платным ручкам до ближайшего PeriodicCheck (~30 мин).
	err := r.db.Raw(`
		SELECT st.level FROM subscription_users su
		JOIN subscription_tiers st ON st.id = COALESCE(
			CASE WHEN su.manual_tier_expires_at IS NULL OR su.manual_tier_expires_at > NOW()
			     THEN su.manual_tier_id END,
			su.resolved_tier_id
		)
		WHERE su.id = ? AND su.is_active = TRUE
	`, userID).Scan(&res).Error
	if err != nil || res.Level == 0 {
		return 0, false
	}
	return res.Level, true
}

// --- Users ---

func (r *SubscriptionRepository) GetUser(userID int64) (*models.SubscriptionUser, error) {
	var user models.SubscriptionUser
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SubscriptionRepository) GetOrCreateUser(userID int64, username *string, fullName string) (*models.SubscriptionUser, error) {
	var user models.SubscriptionUser
	err := r.db.First(&user, userID).Error
	if err == gorm.ErrRecordNotFound {
		user = models.SubscriptionUser{
			ID:       userID,
			Username: username,
			FullName: fullName,
			IsActive: true,
		}
		if err := r.db.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	if err != nil {
		return nil, err
	}
	user.Username = username
	user.FullName = fullName
	r.db.Save(&user)
	return &user, nil
}

func (r *SubscriptionRepository) GetAllActiveUsers() ([]models.SubscriptionUser, error) {
	var users []models.SubscriptionUser
	if err := r.db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// EnsureUser — лёгкий upsert: создаёт запись с минимальным набором полей,
// если её ещё нет. В отличие от GetOrCreateUser не апдейтит username/full_name
// у уже существующего пользователя — это позволяет sweeper-у дёшево
// заводить «не онбордингенных» юзеров (видим только telegram_user_id и
// иногда username из chat_messages), не затирая данные, которые мог обновить
// /start.
func (r *SubscriptionRepository) EnsureUser(userID int64, username *string, fullName string) error {
	var u models.SubscriptionUser
	err := r.db.First(&u, userID).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(&models.SubscriptionUser{
		ID:       userID,
		Username: username,
		FullName: fullName,
		IsActive: true,
	}).Error
}

// GetSweepUserIDs — все telegram_user_id, которых имеет смысл прогонять
// через sweep: записи в subscription_users (is_active=true) ∪ авторы
// сообщений из chat_messages ∪ зарегистрированные на платформе members.
//
// chat_messages нужен, чтобы покрыть людей, которые сидят в чатах сообщества,
// но никогда не открывали бота. members.telegram_id — чтобы покрыть silent
// readers: участники подписочных чатов, которые зарегистрировались на
// платформе, но ни разу не написали в подписочных чатах и которых бот
// не видел в join-event (например, попали в чат до того, как туда зашёл
// бот). Без этого их PeriodicCheck не видит и из контент-чатов не вычищает.
func (r *SubscriptionRepository) GetSweepUserIDs() ([]int64, error) {
	var ids []int64
	err := r.db.Raw(`
		SELECT id FROM subscription_users WHERE is_active = TRUE
		UNION
		SELECT DISTINCT telegram_user_id FROM chat_messages
		WHERE telegram_user_id IS NOT NULL
		UNION
		SELECT telegram_id FROM members WHERE telegram_id > 0
	`).Scan(&ids).Error
	return ids, err
}

func (r *SubscriptionRepository) UpdateResolvedTier(userID int64, tierID *uint) error {
	now := time.Now()
	return r.db.Model(&models.SubscriptionUser{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"resolved_tier_id": tierID,
			"last_check_at":    now,
			"updated_at":       now,
		}).Error
}

// SetManualTier — админский override (бессрочный grant из /suboverride).
// Атомарно устанавливает manual_tier_id и зануляет manual_tier_expires_at,
// чтобы новый grant не унаследовал stale expires от прошлой credits-покупки.
// Без этого зануления админский «бессрочный» grant получал бы случайный
// 30-дневный таймер от предыдущей записи юзера.
func (r *SubscriptionRepository) SetManualTier(userID int64, tierID *uint) error {
	return r.SetManualTierWithExpiry(userID, tierID, nil)
}

// SetManualTierWithExpiry атомарно записывает manual_tier_id и
// manual_tier_expires_at одной UPDATE-командой. Используется покупкой
// тарифа за credits и сбросом истёкшего manual.
//
// Через raw Exec, потому что GORM v2 Updates(map) с nil-value не всегда
// генерирует SET col = NULL — а сброс manual_tier_expires_at нужен и
// при «бессрочной» админской выдаче, и при истечении (там нужно занулить
// и tierID, и expiresAt одновременно).
func (r *SubscriptionRepository) SetManualTierWithExpiry(userID int64, tierID *uint, expiresAt *time.Time) error {
	return r.SetManualTierWithExpiryTx(r.db, userID, tierID, expiresAt)
}

func (r *SubscriptionRepository) SetManualTierWithExpiryTx(db *gorm.DB, userID int64, tierID *uint, expiresAt *time.Time) error {
	return db.Exec(
		`UPDATE subscription_users
		 SET manual_tier_id = ?, manual_tier_expires_at = ?, updated_at = NOW()
		 WHERE id = ?`,
		tierID, expiresAt, userID,
	).Error
}

// GetUserTx — версия GetUser в рамках переданной транзакции.
func (r *SubscriptionRepository) GetUserTx(db *gorm.DB, userID int64) (*models.SubscriptionUser, error) {
	var user models.SubscriptionUser
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// EnsureUserTx — версия EnsureUser в рамках переданной транзакции.
// Возвращает (created bool, err) — true, если запись была создана.
//
// Использует ON CONFLICT DO NOTHING вместо First→Create — без этого
// две одновременные PurchaseTier для одного нового юзера (например,
// клик в двух табах) обе видели бы NotFound и обе делали Create →
// одна получала unique_violation, и её транзакция падала.
func (r *SubscriptionRepository) EnsureUserTx(db *gorm.DB, userID int64, username *string, fullName string) (bool, error) {
	res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.SubscriptionUser{
		ID:       userID,
		Username: username,
		FullName: fullName,
		IsActive: true,
	})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// AddAuditTx — версия AddAudit в рамках переданной транзакции.
func (r *SubscriptionRepository) AddAuditTx(db *gorm.DB, userID int64, action string, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)
	return db.Create(&models.SubscriptionAuditLog{
		UserID:  userID,
		Action:  action,
		Details: string(detailsJSON),
	}).Error
}

// --- Access ---

func (r *SubscriptionRepository) GetActiveAccess(userID int64) ([]models.SubscriptionUserChatAccess, error) {
	var access []models.SubscriptionUserChatAccess
	if err := r.db.Where("user_id = ? AND revoked_at IS NULL", userID).Find(&access).Error; err != nil {
		return nil, err
	}
	return access, nil
}

// GrantAccess upsert'ит активный access. isManual=true ставится при ручном
// добавлении админом (chat_member event from!=user); чистый auto-grant
// (бот выдал invite-link) идёт с isManual=false.
//
// Семантика is_manual: повышается, но не понижается в рамках одной
// «жизни» записи (между Create и RevokeAccess). Если запись уже помечена
// как manual, повторный auto-grant (например, sweep увидел юзера в чате)
// сохраняет защиту. После RevokeAccess флаг сбрасывается — следующий
// grant начинает «новую жизнь» с чистого листа.
func (r *SubscriptionRepository) GrantAccess(userID, chatID int64, isManual bool) error {
	var existing models.SubscriptionUserChatAccess
	err := r.db.Where("user_id = ? AND chat_id = ?", userID, chatID).First(&existing).Error
	if err == nil {
		updates := map[string]interface{}{
			"granted_at": time.Now(),
		}
		if isManual && !existing.IsManual {
			updates["is_manual"] = true
		}
		if err := r.db.Model(&models.SubscriptionUserChatAccess{}).
			Where("user_id = ? AND chat_id = ?", userID, chatID).
			Updates(updates).Error; err != nil {
			return err
		}
		// revoked_at сбрасываем отдельным Update через gorm.Expr —
		// в GORM v2 Updates(map) с nil-value не всегда генерирует
		// UPDATE ... SET col = NULL.
		return r.db.Model(&models.SubscriptionUserChatAccess{}).
			Where("user_id = ? AND chat_id = ?", userID, chatID).
			Update("revoked_at", gorm.Expr("NULL")).Error
	}
	return r.db.Create(&models.SubscriptionUserChatAccess{
		UserID:    userID,
		ChatID:    chatID,
		GrantedAt: time.Now(),
		IsManual:  isManual,
	}).Error
}

// RevokeAccess помечает запись как отозванную и сбрасывает is_manual=false.
// Сброс важен для семантики: manual-защита покрывает текущее членство
// в чате; если юзер ушёл (chat_member leave) или его кикнул periodic, при
// следующем вступлении по invite-link запись «всплывёт» как обычный
// auto-grant. Если потом админ снова добавит вручную — chat_member event
// от админа поднимет is_manual=true заново через GrantAccess.
func (r *SubscriptionRepository) RevokeAccess(userID int64, chatID int64) error {
	now := time.Now()
	return r.db.Model(&models.SubscriptionUserChatAccess{}).
		Where("user_id = ? AND chat_id = ? AND revoked_at IS NULL", userID, chatID).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"is_manual":  false,
		}).Error
}

func (r *SubscriptionRepository) GetUsersWithAccessToChat(chatID int64) ([]models.SubscriptionUser, error) {
	var users []models.SubscriptionUser
	err := r.db.
		Joins("JOIN subscription_user_chat_access ON subscription_user_chat_access.user_id = subscription_users.id").
		Where("subscription_user_chat_access.chat_id = ? AND subscription_user_chat_access.revoked_at IS NULL", chatID).
		Find(&users).Error
	return users, err
}

func (r *SubscriptionRepository) CountAllUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.SubscriptionUser{}).Count(&count).Error
	return count, err
}

// GetEligibleUsersWithoutAccessForChat возвращает пользователей, у которых
// эффективный тир (manual, иначе resolved) имеет уровень >= tierLevel и
// при этом нет активной записи в subscription_user_chat_access для chatID.
// Используется для рассылки новых чат-доступов, когда чат только что
// привязали к тиру.
//
// Anchor-чаты исключаются: их роль — определять тир, а не выдавать access.
// Симметрично с GetChatsForTierLevel; без skip'а после привязки anchor-чата
// к тиру шёл бы массовый «новый чат» спам в ЛС всем eligible.
func (r *SubscriptionRepository) GetEligibleUsersWithoutAccessForChat(
	chatID int64, tierLevel int,
) ([]models.SubscriptionUser, error) {
	var anchor bool
	err := r.db.
		Table("subscription_chats").
		Select("anchor_for_tier_id IS NOT NULL").
		Where("id = ?", chatID).
		Scan(&anchor).Error
	if err != nil {
		return nil, err
	}
	if anchor {
		return nil, nil
	}

	var users []models.SubscriptionUser
	err = r.db.
		Table("subscription_users AS su").
		Select("su.*").
		Joins(`JOIN subscription_tiers st ON st.id = COALESCE(
			CASE WHEN su.manual_tier_expires_at IS NULL OR su.manual_tier_expires_at > NOW()
			     THEN su.manual_tier_id END,
			su.resolved_tier_id
		)`).
		Joins(`LEFT JOIN subscription_user_chat_access sa ON sa.user_id = su.id AND sa.chat_id = ? AND sa.revoked_at IS NULL`, chatID).
		Where("su.is_active = ? AND st.level >= ? AND sa.user_id IS NULL", true, tierLevel).
		Find(&users).Error
	return users, err
}

func (r *SubscriptionRepository) GetUsersByTier(tierID uint) ([]models.SubscriptionUser, error) {
	var users []models.SubscriptionUser
	err := r.db.Where(
		`is_active = ? AND (
			(manual_tier_id = ? AND (manual_tier_expires_at IS NULL OR manual_tier_expires_at > NOW()))
			OR
			((manual_tier_id IS NULL OR manual_tier_expires_at <= NOW()) AND resolved_tier_id = ?)
		)`,
		true, tierID, tierID,
	).Find(&users).Error
	return users, err
}

func (r *SubscriptionRepository) CountUsersByTier(tierID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.SubscriptionUser{}).Where(
		`is_active = ? AND (
			(manual_tier_id = ? AND (manual_tier_expires_at IS NULL OR manual_tier_expires_at > NOW()))
			OR
			((manual_tier_id IS NULL OR manual_tier_expires_at <= NOW()) AND resolved_tier_id = ?)
		)`,
		true, tierID, tierID,
	).Count(&count).Error
	return count, err
}

// CountAllUsersByTier returns user counts for all tiers in a single query.
func (r *SubscriptionRepository) CountAllUsersByTier() (map[uint]int64, error) {
	type result struct {
		TierID uint
		Count  int64
	}
	var results []result
	err := r.db.Raw(`
		SELECT tier_id, COUNT(*) as count FROM (
			SELECT COALESCE(
				CASE WHEN manual_tier_expires_at IS NULL OR manual_tier_expires_at > NOW()
				     THEN manual_tier_id END,
				resolved_tier_id
			) as tier_id
			FROM subscription_users
			WHERE is_active = true
			  AND COALESCE(
				CASE WHEN manual_tier_expires_at IS NULL OR manual_tier_expires_at > NOW()
				     THEN manual_tier_id END,
				resolved_tier_id
			  ) IS NOT NULL
		) t GROUP BY tier_id
	`).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[uint]int64, len(results))
	for _, v := range results {
		m[v.TierID] = v.Count
	}
	return m, nil
}

// CountUsersWithAccessToChat returns the number of users with active access to a chat.
func (r *SubscriptionRepository) CountUsersWithAccessToChat(chatID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.SubscriptionUserChatAccess{}).
		Where("chat_id = ? AND revoked_at IS NULL", chatID).
		Count(&count).Error
	return count, err
}

// CountActiveAccessByChats returns active user counts per chat in a single query.
func (r *SubscriptionRepository) CountActiveAccessByChats(chatIDs []int64) (map[int64]int64, error) {
	if len(chatIDs) == 0 {
		return map[int64]int64{}, nil
	}
	type row struct {
		ChatID int64
		Count  int64
	}
	var rows []row
	err := r.db.Model(&models.SubscriptionUserChatAccess{}).
		Select("chat_id, COUNT(*) as count").
		Where("chat_id IN ? AND revoked_at IS NULL", chatIDs).
		Group("chat_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[int64]int64, len(rows))
	for _, v := range rows {
		m[v.ChatID] = v.Count
	}
	return m, nil
}

// CountActiveAccessByUsers returns active access counts for multiple users in a single query.
func (r *SubscriptionRepository) CountActiveAccessByUsers(userIDs []int64) (map[int64]int64, error) {
	if len(userIDs) == 0 {
		return map[int64]int64{}, nil
	}
	type result struct {
		UserID int64
		Count  int64
	}
	var results []result
	err := r.db.Model(&models.SubscriptionUserChatAccess{}).
		Select("user_id, COUNT(*) as count").
		Where("user_id IN ? AND revoked_at IS NULL", userIDs).
		Group("user_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[int64]int64, len(results))
	for _, v := range results {
		m[v.UserID] = v.Count
	}
	return m, nil
}

func (r *SubscriptionRepository) GetPaginatedUsers(offset, limit int) ([]models.SubscriptionUser, error) {
	var users []models.SubscriptionUser
	err := r.db.Order("id").Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// --- Audit ---

func (r *SubscriptionRepository) AddAudit(userID int64, action string, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)
	return r.db.Create(&models.SubscriptionAuditLog{
		UserID:  userID,
		Action:  action,
		Details: string(detailsJSON),
	}).Error
}
