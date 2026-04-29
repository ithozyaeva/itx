package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"

	"github.com/redis/go-redis/v9"
)

const membershipCacheTTL = 5 * time.Minute

// NewChatAccessChannel — Redis pub/sub канал для уведомлений о том, что чат
// стал доступен новому тиру подписки. Publisher — backend-handler (UI),
// subscriber — бот на NL (он единственный, кто может дойти до Telegram API).
const NewChatAccessChannel = "subscription:new_chat_access"

// NewChatAccessEvent — payload события «чат стал доступен новой аудитории».
// MinTierLevel — минимальный уровень тира среди только что добавленных
// привязок; подписчик уведомляет всех пользователей с level >= этого
// значения. Одно событие на чат, чтобы избежать кратных рассылок, когда
// чат одновременно привязан к нескольким тирам.
type NewChatAccessEvent struct {
	ChatID       int64 `json:"chat_id"`
	MinTierLevel int   `json:"min_tier_level"`
}

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

	// tierID -> все чаты, которые для него anchor. Раньше был map[uint]int64
	// и при дубликате anchor'ов на один тир последний перезаписывал первые,
	// из-за чего часть юзеров из anchor-чата мимо «выживающего» теряли тир.
	anchorMap := make(map[uint][]int64)
	for _, c := range anchorChats {
		if c.AnchorForTierID != nil {
			anchorMap[*c.AnchorForTierID] = append(anchorMap[*c.AnchorForTierID], c.ID)
		}
	}

	for _, tier := range tiersDesc {
		chatIDs, ok := anchorMap[tier.ID]
		if !ok {
			continue
		}
		for _, chatID := range chatIDs {
			if s.IsMember(chatID, userID, botCheckFunc) {
				id := tier.ID
				return &id
			}
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
		user.ResolvedTierID = newTierID
	}

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

// TierPublic — публичная карточка тарифа для UI лендинга/платформы и сообщений
// бота. Цена отдаётся в рублях (price_cents переведён). Features — массив строк.
type TierPublic struct {
	ID          uint     `json:"id"`
	Slug        string   `json:"slug"`
	Name        string   `json:"name"`
	Level       int      `json:"level"`
	Price       int      `json:"price"`
	BoostyURL   string   `json:"boosty_url"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
}

// GetPublicTiers возвращает только тарифы с is_public=true, отсортированные
// по level. Используется как единый источник правды для /tariffs, прогрева
// в боте и SEO-блока на лендинге.
func (s *SubscriptionService) GetPublicTiers() ([]TierPublic, error) {
	tiers, err := s.repo.GetPublicTiers()
	if err != nil {
		return nil, err
	}
	result := make([]TierPublic, 0, len(tiers))
	for _, t := range tiers {
		features := []string{}
		if t.Features != "" {
			_ = json.Unmarshal([]byte(t.Features), &features)
		}
		public := TierPublic{
			ID:       t.ID,
			Slug:     t.Slug,
			Name:     t.Name,
			Level:    t.Level,
			Features: features,
		}
		if t.PriceCents != nil {
			public.Price = *t.PriceCents / 100
		}
		if t.BoostyURL != nil {
			public.BoostyURL = *t.BoostyURL
		}
		if t.PublicDescription != nil {
			public.Description = *t.PublicDescription
		}
		result = append(result, public)
	}
	return result, nil
}

// GetSubscriptionUser возвращает запись subscription_users по telegram-id.
// Используется RequireSubscription middleware и онбордингом в боте.
func (s *SubscriptionService) GetSubscriptionUser(userID int64) (*models.SubscriptionUser, error) {
	return s.repo.GetUser(userID)
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

func (s *SubscriptionService) UpdateChatMeta(chatID int64, category, emoji *string) error {
	return s.repo.UpdateChatMeta(chatID, category, emoji)
}

func (s *SubscriptionService) SetChatPriority(chatID int64, priority int) error {
	return s.repo.UpdateChatPriority(chatID, priority)
}

func (s *SubscriptionService) SetAnchor(chatID int64, tierID *uint) error {
	return s.repo.SetAnchor(chatID, tierID)
}

func (s *SubscriptionService) AddChatToTier(chatID int64, tierID uint) error {
	return s.repo.AddChatToTier(chatID, tierID)
}

func (s *SubscriptionService) GetAllTierChats() (map[int64][]uint, error) {
	return s.repo.GetAllTierChats()
}

func (s *SubscriptionService) GetTierIDsForChat(chatID int64) ([]uint, error) {
	return s.repo.GetTierIDsForChat(chatID)
}

func (s *SubscriptionService) SetChatTiers(chatID int64, tierIDs []uint) error {
	return s.repo.SetChatTiers(chatID, tierIDs)
}

// GetEligibleUsersWithoutAccessForChat — пользователи с эффективным тиром
// уровня >= tierLevel, которым доступ к этому чату ещё не выдан.
func (s *SubscriptionService) GetEligibleUsersWithoutAccessForChat(
	chatID int64, tierLevel int,
) ([]models.SubscriptionUser, error) {
	return s.repo.GetEligibleUsersWithoutAccessForChat(chatID, tierLevel)
}

// GetChatsForTierLevel — все content-чаты, привязанные к тирам с level <= tierLevel.
// Anchor-чаты не включены (членство в них определяет сам тир).
func (s *SubscriptionService) GetChatsForTierLevel(tierLevel int) ([]models.SubscriptionChat, error) {
	return s.repo.GetChatsForTierLevel(tierLevel)
}

// PublishNewChatAccess сигналит боту, что чат chatID стал доступен новой
// аудитории — пользователей с эффективным тиром >= minTierLevel надо
// пригласить. Бэкенд в РФ не может сам пойти в Telegram (i/o timeout),
// поэтому рассылку делает бот на NL, подписанный на этот канал.
// Шлём одно событие на чат, а не на каждый tier — иначе при привязке
// чата сразу к нескольким тирам рассылка повторялась бы N раз.
func (s *SubscriptionService) PublishNewChatAccess(ctx context.Context, chatID int64, minTierLevel int) error {
	payload, err := json.Marshal(NewChatAccessEvent{ChatID: chatID, MinTierLevel: minTierLevel})
	if err != nil {
		return err
	}
	return s.redis.Publish(ctx, NewChatAccessChannel, payload).Err()
}

// SubscribeNewChatAccess запускает горутину, которая читает события pub/sub
// и для каждого вызывает handler. Вызывается при старте бота.
func (s *SubscriptionService) SubscribeNewChatAccess(ctx context.Context, handler func(ev NewChatAccessEvent)) {
	pubsub := s.redis.Subscribe(ctx, NewChatAccessChannel)
	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var ev NewChatAccessEvent
				if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
					log.Printf("new-chat-access: bad payload: %v", err)
					continue
				}
				handler(ev)
			}
		}
	}()
	log.Printf("Subscribed to %s for new-chat-access events", NewChatAccessChannel)
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

func (s *SubscriptionService) GrantAccess(userID int64, chatID int64) error {
	return s.repo.GrantAccess(userID, chatID)
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

func (s *SubscriptionService) CountUsersByTier(tierID uint) (int64, error) {
	return s.repo.CountUsersByTier(tierID)
}

func (s *SubscriptionService) CountAllUsersByTier() (map[uint]int64, error) {
	return s.repo.CountAllUsersByTier()
}

func (s *SubscriptionService) CountUsersWithAccessToChat(chatID int64) (int64, error) {
	return s.repo.CountUsersWithAccessToChat(chatID)
}

func (s *SubscriptionService) CountActiveAccessByUsers(userIDs []int64) (map[int64]int64, error) {
	return s.repo.CountActiveAccessByUsers(userIDs)
}

func (s *SubscriptionService) CountActiveAccessByChats(chatIDs []int64) (map[int64]int64, error) {
	return s.repo.CountActiveAccessByChats(chatIDs)
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
