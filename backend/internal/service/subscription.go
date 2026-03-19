package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"

	"github.com/redis/go-redis/v9"
)

const membershipCacheTTL = 5 * time.Minute

type SubscriptionService struct {
	repo  *repository.SubscriptionRepository
	redis *redis.Client
}

func NewSubscriptionService(redisClient *redis.Client) *SubscriptionService {
	return &SubscriptionService{
		repo:  repository.NewSubscriptionRepository(),
		redis: redisClient,
	}
}

// IsMember checks if a user is a member of a chat, with Redis caching.
// botCheckFunc should call the Telegram Bot API getChatMember.
func (s *SubscriptionService) IsMember(chatID int64, userID int64, botCheckFunc func(chatID, userID int64) bool) bool {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("sub:member:%d:%d", chatID, userID)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		return cached == "1"
	}

	result := botCheckFunc(chatID, userID)

	val := "0"
	if result {
		val = "1"
	}
	s.redis.Set(ctx, cacheKey, val, membershipCacheTTL)

	return result
}

// InvalidateMemberCache removes the membership cache for a specific user/chat combo.
func (s *SubscriptionService) InvalidateMemberCache(chatID int64, userID int64) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("sub:member:%d:%d", chatID, userID)
	s.redis.Del(ctx, cacheKey)
}

// ResolveTierID checks anchor chats from highest tier downward, returns first match.
func (s *SubscriptionService) ResolveTierID(userID int64, botCheckFunc func(chatID, userID int64) bool) *uint {
	anchorChats, err := s.repo.GetAnchorChats()
	if err != nil {
		log.Printf("Error getting anchor chats: %v", err)
		return nil
	}

	tiersDesc, err := s.repo.GetAllTiersDesc()
	if err != nil {
		log.Printf("Error getting tiers: %v", err)
		return nil
	}

	// Build anchor map: tierID -> chatID
	anchorMap := make(map[uint]int64)
	for _, c := range anchorChats {
		if c.AnchorForTierID != nil {
			anchorMap[*c.AnchorForTierID] = c.ID
		}
	}

	for _, tier := range tiersDesc {
		anchorChatID, ok := anchorMap[tier.ID]
		if !ok {
			continue
		}
		if s.IsMember(anchorChatID, userID, botCheckFunc) {
			id := tier.ID
			return &id
		}
	}
	return nil
}

type SyncResult struct {
	UserID          int64
	OldTierID       *uint
	NewTierID       *uint
	EffectiveTierID *uint
	Granted         []GrantedChat
	Revoked         []int64
}

type GrantedChat struct {
	ChatID int64
	Link   string
}

// CheckAndSyncUser performs a full subscription check and sync for a user.
func (s *SubscriptionService) CheckAndSyncUser(
	userID int64,
	botCheckFunc func(chatID, userID int64) bool,
	createInviteLink func(chatID int64) (string, error),
	kickUser func(chatID, userID int64),
) (*SyncResult, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	newTierID := s.ResolveTierID(userID, botCheckFunc)
	oldTierID := user.ResolvedTierID

	if !tierIDsEqual(newTierID, oldTierID) {
		s.repo.UpdateResolvedTier(userID, newTierID)
		s.repo.AddAudit(userID, "tier_change", map[string]interface{}{
			"old_tier_id": oldTierID,
			"new_tier_id": newTierID,
		})
	}

	// Reload user
	user, _ = s.repo.GetUser(userID)
	effectiveTierID := user.EffectiveTierID()

	// Determine entitled chats
	var entitledChats []models.SubscriptionChat
	if effectiveTierID != nil {
		tier, err := s.repo.GetTier(*effectiveTierID)
		if err == nil {
			entitledChats, _ = s.repo.GetChatsForTierLevel(tier.Level)
		}
	}

	entitledIDs := make(map[int64]bool)
	for _, c := range entitledChats {
		entitledIDs[c.ID] = true
	}

	// Current active access
	currentAccess, _ := s.repo.GetActiveAccess(userID)
	currentIDs := make(map[int64]bool)
	for _, a := range currentAccess {
		currentIDs[a.ChatID] = true
	}

	result := &SyncResult{
		UserID:          userID,
		OldTierID:       oldTierID,
		NewTierID:       newTierID,
		EffectiveTierID: effectiveTierID,
	}

	// Grant missing
	for chatID := range entitledIDs {
		if !currentIDs[chatID] {
			link, err := createInviteLink(chatID)
			if err != nil {
				log.Printf("Failed to create invite link for chat %d: %v", chatID, err)
				continue
			}
			s.repo.GrantAccess(userID, chatID)
			s.repo.AddAudit(userID, "grant", map[string]interface{}{
				"chat_id": chatID,
			})
			result.Granted = append(result.Granted, GrantedChat{ChatID: chatID, Link: link})
		}
	}

	// Revoke extra
	for chatID := range currentIDs {
		if !entitledIDs[chatID] {
			kickUser(chatID, userID)
			s.repo.RevokeAccess(userID, chatID)
			s.repo.AddAudit(userID, "revoke", map[string]interface{}{
				"chat_id": chatID,
			})
			result.Revoked = append(result.Revoked, chatID)
		}
	}

	return result, nil
}

// OnboardUser creates/updates user and syncs access.
func (s *SubscriptionService) OnboardUser(
	userID int64,
	username *string,
	fullName string,
	botCheckFunc func(chatID, userID int64) bool,
	createInviteLink func(chatID int64) (string, error),
	kickUser func(chatID, userID int64),
) (*SyncResult, error) {
	_, err := s.repo.GetOrCreateUser(userID, username, fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create user: %w", err)
	}
	return s.CheckAndSyncUser(userID, botCheckFunc, createInviteLink, kickUser)
}

// PeriodicCheck checks all active users and syncs their access.
func (s *SubscriptionService) PeriodicCheck(
	botCheckFunc func(chatID, userID int64) bool,
	createInviteLink func(chatID int64) (string, error),
	kickUser func(chatID, userID int64),
	notifyUser func(userID int64, result *SyncResult),
	rateDelay time.Duration,
) {
	log.Println("Starting periodic subscription check")

	users, err := s.repo.GetAllActiveUsers()
	if err != nil {
		log.Printf("Error getting active users: %v", err)
		return
	}

	for _, user := range users {
		result, err := s.CheckAndSyncUser(user.ID, botCheckFunc, createInviteLink, kickUser)
		if err != nil {
			log.Printf("Error checking user %d: %v", user.ID, err)
		} else if len(result.Granted) > 0 || len(result.Revoked) > 0 {
			log.Printf("User %d: granted=%d revoked=%d", user.ID, len(result.Granted), len(result.Revoked))
			notifyUser(user.ID, result)
		}
		time.Sleep(rateDelay)
	}

	log.Println("Periodic subscription check complete")
}

// --- Repo delegation methods ---

func (s *SubscriptionService) GetAllTiers() ([]models.SubscriptionTier, error) {
	return s.repo.GetAllTiers()
}

func (s *SubscriptionService) GetTierBySlug(slug string) (*models.SubscriptionTier, error) {
	return s.repo.GetTierBySlug(slug)
}

func (s *SubscriptionService) GetTier(id uint) (*models.SubscriptionTier, error) {
	return s.repo.GetTier(id)
}

func (s *SubscriptionService) GetAllChats() ([]models.SubscriptionChat, error) {
	return s.repo.GetAllChats()
}

func (s *SubscriptionService) GetChat(chatID int64) (*models.SubscriptionChat, error) {
	return s.repo.GetChat(chatID)
}

func (s *SubscriptionService) UpsertChat(chatID int64, title, chatType string) error {
	return s.repo.UpsertChat(chatID, title, chatType)
}

func (s *SubscriptionService) SetAnchor(chatID int64, tierID *uint) error {
	return s.repo.SetAnchor(chatID, tierID)
}

func (s *SubscriptionService) AddChatToTier(chatID int64, tierID uint) error {
	return s.repo.AddChatToTier(chatID, tierID)
}

func (s *SubscriptionService) DeleteChat(chatID int64) error {
	return s.repo.DeleteChat(chatID)
}

func (s *SubscriptionService) GetUser(userID int64) (*models.SubscriptionUser, error) {
	return s.repo.GetUser(userID)
}

func (s *SubscriptionService) SetManualTier(userID int64, tierID *uint) error {
	return s.repo.SetManualTier(userID, tierID)
}

func (s *SubscriptionService) AddAudit(userID int64, action string, details map[string]interface{}) error {
	return s.repo.AddAudit(userID, action, details)
}

func (s *SubscriptionService) GetActiveAccess(userID int64) ([]models.SubscriptionUserChatAccess, error) {
	return s.repo.GetActiveAccess(userID)
}

func (s *SubscriptionService) GetUsersWithAccessToChat(chatID int64) ([]models.SubscriptionUser, error) {
	return s.repo.GetUsersWithAccessToChat(chatID)
}

func (s *SubscriptionService) RevokeAccess(userID int64, chatID int64) error {
	return s.repo.RevokeAccess(userID, chatID)
}

func (s *SubscriptionService) CountAllUsers() (int64, error) {
	return s.repo.CountAllUsers()
}

func (s *SubscriptionService) GetUsersByTier(tierID uint) ([]models.SubscriptionUser, error) {
	return s.repo.GetUsersByTier(tierID)
}

func (s *SubscriptionService) GetPaginatedUsers(offset, limit int) ([]models.SubscriptionUser, error) {
	return s.repo.GetPaginatedUsers(offset, limit)
}

func tierIDsEqual(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
