package repository

import (
	"encoding/json"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"

	"gorm.io/gorm"
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

func (r *SubscriptionRepository) DeleteChat(chatID int64) error {
	r.db.Where("chat_id = ?", chatID).Delete(&models.SubscriptionTierChat{})
	r.db.Where("chat_id = ?", chatID).Delete(&models.SubscriptionUserChatAccess{})
	return r.db.Delete(&models.SubscriptionChat{}, chatID).Error
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

func (r *SubscriptionRepository) UpdateResolvedTier(userID int64, tierID *uint) error {
	now := time.Now()
	return r.db.Model(&models.SubscriptionUser{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"resolved_tier_id": tierID,
			"last_check_at":    now,
			"updated_at":       now,
		}).Error
}

func (r *SubscriptionRepository) SetManualTier(userID int64, tierID *uint) error {
	return r.db.Model(&models.SubscriptionUser{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"manual_tier_id": tierID,
			"updated_at":     time.Now(),
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

func (r *SubscriptionRepository) GrantAccess(userID int64, chatID int64) error {
	var existing models.SubscriptionUserChatAccess
	err := r.db.Where("user_id = ? AND chat_id = ?", userID, chatID).First(&existing).Error
	if err == nil {
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"granted_at": time.Now(),
			"revoked_at": nil,
		}).Error
	}
	return r.db.Create(&models.SubscriptionUserChatAccess{
		UserID:    userID,
		ChatID:    chatID,
		GrantedAt: time.Now(),
	}).Error
}

func (r *SubscriptionRepository) RevokeAccess(userID int64, chatID int64) error {
	now := time.Now()
	return r.db.Model(&models.SubscriptionUserChatAccess{}).
		Where("user_id = ? AND chat_id = ? AND revoked_at IS NULL", userID, chatID).
		Update("revoked_at", now).Error
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

func (r *SubscriptionRepository) GetUsersByTier(tierID uint) ([]models.SubscriptionUser, error) {
	var users []models.SubscriptionUser
	err := r.db.Where(
		"is_active = ? AND ((manual_tier_id = ?) OR (manual_tier_id IS NULL AND resolved_tier_id = ?))",
		true, tierID, tierID,
	).Find(&users).Error
	return users, err
}

func (r *SubscriptionRepository) CountUsersByTier(tierID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.SubscriptionUser{}).Where(
		"is_active = ? AND ((manual_tier_id = ?) OR (manual_tier_id IS NULL AND resolved_tier_id = ?))",
		true, tierID, tierID,
	).Count(&count).Error
	return count, err
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
