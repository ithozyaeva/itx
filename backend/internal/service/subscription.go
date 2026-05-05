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

// MemberCheckFunc — Telegram getChatMember с распространением ошибки.
// При (false, error) ResolveTierID не понижает resolved_tier пользователя:
// иначе при rate-limit или таймауте Telegram юзер «терял» бы тир и при
// SUBSCRIPTION_AUTO_KICK_ENABLED=true получал бы ложный кик из content-чатов.
type MemberCheckFunc func(chatID, userID int64) (bool, error)

// SubscriptionContext — снэпшот anchor-чатов и тиров, разделяемый между
// per-user итерациями PeriodicCheck/DryRunPeriodicCheck. Anchor-чаты и
// тиры меняются раз в недели; читать их из БД на каждого юзера было
// чистой воды лишний трафик NL→РФ (250+ юзеров × несколько SELECT).
type SubscriptionContext struct {
	AnchorChatsByTier map[uint][]int64 // tierID -> anchor chat IDs
	AnchorChatIDs     map[int64]bool   // set всех anchor chat IDs
	TiersDesc         []models.SubscriptionTier
}

// BuildContext — строит SubscriptionContext одним проходом по БД.
// Используется PeriodicCheck/DryRunPeriodicCheck до loop'а, и единичными
// точками входа (CheckAndSyncUser/ResolveTierID) — для них накладные
// расходы те же, что были до фикса.
func (s *SubscriptionService) BuildContext() (*SubscriptionContext, error) {
	anchors, err := s.repo.GetAnchorChats()
	if err != nil {
		return nil, fmt.Errorf("get anchor chats: %w", err)
	}
	tiers, err := s.repo.GetAllTiersDesc()
	if err != nil {
		return nil, fmt.Errorf("get tiers desc: %w", err)
	}

	byTier := make(map[uint][]int64)
	ids := make(map[int64]bool, len(anchors))
	for _, c := range anchors {
		if c.AnchorForTierID != nil {
			byTier[*c.AnchorForTierID] = append(byTier[*c.AnchorForTierID], c.ID)
		}
		ids[c.ID] = true
	}
	return &SubscriptionContext{
		AnchorChatsByTier: byTier,
		AnchorChatIDs:     ids,
		TiersDesc:         tiers,
	}, nil
}

// IsMember checks if a user is a member of a chat, with Redis caching.
// botCheckFunc should call the Telegram Bot API getChatMember.
//
// При ошибке botCheckFunc (rate-limit, таймаут, network) возвращаем
// (false, err) и НЕ кэшируем — иначе на 5 минут зависал бы false-позитив
// и юзер на следующих проходах получил бы ложный кик из content-чатов.
func (s *SubscriptionService) IsMember(chatID int64, userID int64, botCheckFunc MemberCheckFunc) (bool, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("sub:member:%d:%d", chatID, userID)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		return cached == "1", nil
	}

	result, callErr := botCheckFunc(chatID, userID)
	if callErr != nil {
		return false, callErr
	}

	val := "0"
	if result {
		val = "1"
	}
	s.redis.Set(ctx, cacheKey, val, membershipCacheTTL)

	return result, nil
}

// InvalidateMemberCache removes the membership cache for a specific user/chat combo.
func (s *SubscriptionService) InvalidateMemberCache(chatID int64, userID int64) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("sub:member:%d:%d", chatID, userID)
	s.redis.Del(ctx, cacheKey)
}

// ResolveTierID checks anchor chats from highest tier downward, returns first match.
//
// Внешняя точка входа: каждый раз строит свежий SubscriptionContext.
// Используется единичными сценариями (/sub onboarding, anchor-join,
// /suboverride). Для loop'а — resolveTierIDFromContext с переиспользованием
// снэпшота.
//
// При ошибке IsMember (Telegram API упал/rate-limit) пропускает остаток
// тиров и возвращает (nil, err). Каскадный fail-stop важен: если до
// master-anchor мы не дозвонились, нельзя «понизить» юзера до beginner —
// иначе получим ложный кик из master-only чатов.
func (s *SubscriptionService) ResolveTierID(userID int64, botCheckFunc MemberCheckFunc) (*uint, error) {
	ctx, err := s.BuildContext()
	if err != nil {
		return nil, err
	}
	return s.resolveTierIDFromContext(userID, botCheckFunc, ctx)
}

// resolveTierIDFromContext — без БД-обращений к anchor-чатам/тирам.
// Использует переданный snapshot. Для PeriodicCheck/DryRunPeriodicCheck —
// один SELECT на проход вместо одного на пользователя.
func (s *SubscriptionService) resolveTierIDFromContext(
	userID int64,
	botCheckFunc MemberCheckFunc,
	ctx *SubscriptionContext,
) (*uint, error) {
	for _, tier := range ctx.TiersDesc {
		chatIDs, ok := ctx.AnchorChatsByTier[tier.ID]
		if !ok {
			continue
		}
		for _, chatID := range chatIDs {
			isMember, err := s.IsMember(chatID, userID, botCheckFunc)
			if err != nil {
				return nil, fmt.Errorf("check anchor chat %d: %w", chatID, err)
			}
			if isMember {
				id := tier.ID
				return &id, nil
			}
		}
	}
	return nil, nil
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
//
// kickUser возвращает bool — реально ли произошёл kick (false при
// SUBSCRIPTION_AUTO_KICK_ENABLED=false). Запись revoke в БД, audit и
// добавление чата в result.Revoked происходят ТОЛЬКО когда kickUser
// вернул true. Это превращает выключенный auto-kick в полноценный
// dry-run всей цепочки: ни БД-state, ни нотификации не меняются.
//
// Anchor-чаты явно пропускаются при revoke: они определяют тир, а не
// являются объектом доступа. В entitled они не попадают (по дизайну
// GetChatsForTierLevel), и без явного skip каждый periodic-check
// пытался бы их revoke'нуть.
//
// Внешняя точка входа: строит свой SubscriptionContext и грузит user
// из БД. Используется онбордингом (/sub, anchor-join), /suboverride.
// Для loop'а PeriodicCheck — checkAndSyncUserCtx с переиспользованием
// снэпшота.
func (s *SubscriptionService) CheckAndSyncUser(
	userID int64,
	botCheckFunc MemberCheckFunc,
	createInviteLink func(chatID int64) (string, error),
	kickUser func(chatID, userID int64) bool,
) (*SyncResult, error) {
	subCtx, err := s.BuildContext()
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return s.checkAndSyncUserCtx(user, botCheckFunc, createInviteLink, kickUser, subCtx)
}

// checkAndSyncUserCtx — внутренняя версия. user уже загружен, ctx —
// shared snapshot anchor/tiers. Принимает указатель на user, чтобы
// in-place обновить ResolvedTierID после UpdateResolvedTier (как было
// в публичной версии).
func (s *SubscriptionService) checkAndSyncUserCtx(
	user *models.SubscriptionUser,
	botCheckFunc MemberCheckFunc,
	createInviteLink func(chatID int64) (string, error),
	kickUser func(chatID, userID int64) bool,
	subCtx *SubscriptionContext,
) (*SyncResult, error) {
	userID := user.ID
	newTierID, err := s.resolveTierIDFromContext(userID, botCheckFunc, subCtx)
	if err != nil {
		return nil, fmt.Errorf("resolve tier: %w", err)
	}
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

	// Revoke extra. Skip anchor-чаты (см. doc-комментарий выше). Запись
	// revoke в БД делаем только если kickUser реально кикнул — иначе
	// при auto_kick=false мы бы тихо снимали access и слали юзерам
	// «уровень изменился», хотя в Telegram они остаются на месте.
	for chatID := range currentIDs {
		if entitledIDs[chatID] {
			continue
		}
		if subCtx.AnchorChatIDs[chatID] {
			continue
		}
		if !kickUser(chatID, userID) {
			continue
		}
		s.repo.RevokeAccess(userID, chatID)
		s.repo.AddAudit(userID, "revoke", map[string]interface{}{
			"chat_id": chatID,
		})
		result.Revoked = append(result.Revoked, chatID)
	}

	return result, nil
}

// OnboardUser creates/updates user and syncs access.
func (s *SubscriptionService) OnboardUser(
	userID int64,
	username *string,
	fullName string,
	botCheckFunc MemberCheckFunc,
	createInviteLink func(chatID int64) (string, error),
	kickUser func(chatID, userID int64) bool,
) (*SyncResult, error) {
	_, err := s.repo.GetOrCreateUser(userID, username, fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create user: %w", err)
	}
	return s.CheckAndSyncUser(userID, botCheckFunc, createInviteLink, kickUser)
}

// EnsureUser — лёгкий upsert в subscription_users без полного sync.
// Используется для онбординга юзеров, которых мы видим в content-чатах
// (через chat_member updates или backfill-sweep), но которые сами /start
// в боте ещё не нажимали — иначе их нет в таблице и PeriodicCheck про них
// не знает.
func (s *SubscriptionService) EnsureUser(userID int64, username *string, fullName string) error {
	return s.repo.EnsureUser(userID, username, fullName)
}

// SyncContentJoin — пользователь зашёл в content-чат. Заводим запись
// в subscription_users (если ещё нет) и проставляем access. Тир не
// пересчитываем: его обновит ближайший PeriodicCheck — здесь нам нужно
// только зафиксировать факт членства, чтобы DryRun/PeriodicCheck потом
// корректно посчитал «лишних».
func (s *SubscriptionService) SyncContentJoin(userID int64, chatID int64, username *string, fullName string) error {
	if err := s.repo.EnsureUser(userID, username, fullName); err != nil {
		return err
	}
	return s.repo.GrantAccess(userID, chatID)
}

// SyncContentLeave — пользователь вышел/кикнут из content-чата. Снимаем
// access; saving в audit здесь не пишем, чтобы не засорять — это «реакция
// на естественный уход», не действие системы.
func (s *SubscriptionService) SyncContentLeave(userID int64, chatID int64) error {
	return s.repo.RevokeAccess(userID, chatID)
}

// SweepStats — счётчики однопроходного backfill-обхода реального членства.
type SweepStats struct {
	UsersScanned    int
	UsersCreated    int
	AccessGranted   int
	AccessRevoked   int
	ChecksPerformed int
}

// SweepRealMembership — однопроходный обход «кто реально сидит в наших
// чатах». Для каждого user_id из subscription_users ∪ chat_messages
// делает getChatMember по каждому subscription_chat и приводит
// subscription_user_chat_access в соответствие.
//
// botCheckFunc: true для member/administrator/creator/restricted, false
// для left/kicked/PARTICIPANT_ID_INVALID/любой ошибки. Кэш membership
// обходить не нужно — sweep редкий, кэш живёт 5 минут, всё равно
// обновится за время прохода.
//
// rateDelay: пауза между getChatMember-вызовами, чтобы не упереться в
// rate-limit Telegram (в Bot API глобально ~30 rps на разные чаты).
//
// Anchor-чаты пропускаются: их роль — определять тир, а не выдавать
// доступ. Если запоминать членство в anchor как access, PeriodicCheck
// потом увидит anchor-access как «лишний» (он не входит в entitled by
// design) и попытается revoke — будет шлать юзеру «уровень изменился»
// при каждом проходе. Membership в anchor читается через ResolveTierID
// напрямую из Telegram, БД для этого не нужна.
func (s *SubscriptionService) SweepRealMembership(
	botCheckFunc MemberCheckFunc,
	rateDelay time.Duration,
) (*SweepStats, error) {
	chats, err := s.repo.GetAllChats()
	if err != nil {
		return nil, fmt.Errorf("get chats: %w", err)
	}

	userIDs, err := s.repo.GetSweepUserIDs()
	if err != nil {
		return nil, fmt.Errorf("get sweep user ids: %w", err)
	}

	// Отфильтровываем anchor-чаты сразу — иначе их membership попал бы в
	// access-таблицу и portal-effect: каждый periodic-check после sweep
	// сносил бы их обратно с уведомлением юзеру.
	contentChats := make([]models.SubscriptionChat, 0, len(chats))
	for _, c := range chats {
		if c.AnchorForTierID == nil {
			contentChats = append(contentChats, c)
		}
	}

	stats := &SweepStats{UsersScanned: len(userIDs)}

	for _, uid := range userIDs {
		// Если юзера нет в subscription_users — заводим легковесно.
		// FullName пустой: chat_messages хранит только first_name,
		// но это не критично — UI/PeriodicCheck оперируют id и username.
		if _, err := s.repo.GetUser(uid); err != nil {
			if err := s.repo.EnsureUser(uid, nil, ""); err == nil {
				stats.UsersCreated++
			}
		}

		// Текущие active-access — чтобы корректно считать diff.
		current, _ := s.repo.GetActiveAccess(uid)
		currentSet := make(map[int64]bool, len(current))
		for _, a := range current {
			currentSet[a.ChatID] = true
		}

		for _, chat := range contentChats {
			isMember, err := botCheckFunc(chat.ID, uid)
			stats.ChecksPerformed++
			if err != nil {
				// Telegram API упал/rate-limit — лучше пропустить пару (chat,user)
				// на этом проходе, чем ошибочно revoke'нуть active access. Sweep
				// — best-effort, разойдётся при следующем суточном тикере.
				log.Printf("sweep: skip chat=%d user=%d: %v", chat.ID, uid, err)
				time.Sleep(rateDelay)
				continue
			}

			if isMember && !currentSet[chat.ID] {
				if err := s.repo.GrantAccess(uid, chat.ID); err == nil {
					stats.AccessGranted++
				}
			} else if !isMember && currentSet[chat.ID] {
				if err := s.repo.RevokeAccess(uid, chat.ID); err == nil {
					stats.AccessRevoked++
				}
			}

			time.Sleep(rateDelay)
		}
	}

	return stats, nil
}

// DryRunUserResult — что сделал бы PeriodicCheck для одного пользователя,
// если бы был запущен сейчас. Действий в БД и Telegram не выполняется.
type DryRunUserResult struct {
	UserID        int64
	Username      *string
	OldTierID     *uint
	NewTierID     *uint
	EffectiveTier *uint
	WouldGrant    []int64
	WouldRevoke   []int64
}

// DryRunCheckUser — повторяет логику CheckAndSyncUser, но без записи в БД
// и без вызовов createInviteLink/kickUser. Используется командой
// /subkickdry, чтобы сначала посмотреть, кого бот удалил бы из чатов.
//
// Внешняя точка входа: строит свой SubscriptionContext и грузит user.
// Для loop'а DryRunPeriodicCheck — dryRunCheckUserCtx.
func (s *SubscriptionService) DryRunCheckUser(
	userID int64,
	botCheckFunc MemberCheckFunc,
) (*DryRunUserResult, error) {
	subCtx, err := s.BuildContext()
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return s.dryRunCheckUserCtx(user, botCheckFunc, subCtx)
}

func (s *SubscriptionService) dryRunCheckUserCtx(
	user *models.SubscriptionUser,
	botCheckFunc MemberCheckFunc,
	subCtx *SubscriptionContext,
) (*DryRunUserResult, error) {
	userID := user.ID
	newTierID, err := s.resolveTierIDFromContext(userID, botCheckFunc, subCtx)
	if err != nil {
		return nil, fmt.Errorf("resolve tier: %w", err)
	}

	// EffectiveTierID учитывает manual override; если ручной — он же и итог.
	var effective *uint
	if user.ManualTierID != nil {
		effective = user.ManualTierID
	} else {
		effective = newTierID
	}

	entitled := make(map[int64]bool)
	if effective != nil {
		tier, err := s.repo.GetTier(*effective)
		if err == nil {
			chats, _ := s.repo.GetChatsForTierLevel(tier.Level)
			for _, c := range chats {
				entitled[c.ID] = true
			}
		}
	}

	currentAccess, _ := s.repo.GetActiveAccess(userID)
	current := make(map[int64]bool, len(currentAccess))
	for _, a := range currentAccess {
		current[a.ChatID] = true
	}

	res := &DryRunUserResult{
		UserID:        userID,
		Username:      user.Username,
		OldTierID:     user.ResolvedTierID,
		NewTierID:     newTierID,
		EffectiveTier: effective,
	}
	for cid := range entitled {
		if !current[cid] {
			res.WouldGrant = append(res.WouldGrant, cid)
		}
	}
	for cid := range current {
		if entitled[cid] {
			continue
		}
		if subCtx.AnchorChatIDs[cid] {
			continue
		}
		res.WouldRevoke = append(res.WouldRevoke, cid)
	}
	return res, nil
}

// DryRunPeriodicCheck — обходит всех active subscription_users и собирает
// список действий, которые сделал бы PeriodicCheck. Действий не выполняет.
// Используется /subkickdry для предварительного отчёта перед включением
// SUBSCRIPTION_AUTO_KICK_ENABLED.
func (s *SubscriptionService) DryRunPeriodicCheck(
	botCheckFunc MemberCheckFunc,
	rateDelay time.Duration,
) ([]DryRunUserResult, error) {
	subCtx, err := s.BuildContext()
	if err != nil {
		return nil, err
	}
	users, err := s.repo.GetAllActiveUsers()
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	results := make([]DryRunUserResult, 0, len(users))
	for i := range users {
		u := &users[i]
		r, err := s.dryRunCheckUserCtx(u, botCheckFunc, subCtx)
		if err != nil {
			log.Printf("dry-run: user %d: %v", u.ID, err)
			continue
		}
		results = append(results, *r)
		time.Sleep(rateDelay)
	}
	return results, nil
}

// PeriodicCheck checks all active users and syncs their access.
//
// В отличие от интерактивных сценариев (/sub onboarding, anchor-join,
// /suboverride) periodic-проход НЕ шлёт юзерам нотификаций — это массовый
// фоновой синк, и любой шум там воспринимается как спам. Возвращаем
// сводку: вызывающий код (bot) шлёт её админу для ручной проверки.
//
// Юзеры с пустыми Granted и Revoked в результат не включаются, чтобы
// сводка содержала только то, по чему действительно были действия.
//
// Anchor-чаты и тиры читаются из БД один раз до loop'а: на 250+ юзерах
// это сокращает NL→РФ-трафик с ~3 SELECT/юзера до constant-2.
//
// Если Telegram API упал на anchor-проверке конкретного юзера, юзер
// этим проходом пропускается — лучше отложить sync на следующий тикер,
// чем ложно понизить тир и кикнуть из master-only чатов на rate-limit'е.
func (s *SubscriptionService) PeriodicCheck(
	botCheckFunc MemberCheckFunc,
	createInviteLink func(chatID int64) (string, error),
	kickUser func(chatID, userID int64) bool,
	rateDelay time.Duration,
) []SyncResult {
	log.Println("Starting periodic subscription check")

	subCtx, err := s.BuildContext()
	if err != nil {
		log.Printf("periodic: build context failed: %v", err)
		return nil
	}

	users, err := s.repo.GetAllActiveUsers()
	if err != nil {
		log.Printf("Error getting active users: %v", err)
		return nil
	}

	var changed []SyncResult
	for i := range users {
		user := &users[i]
		result, err := s.checkAndSyncUserCtx(user, botCheckFunc, createInviteLink, kickUser, subCtx)
		if err != nil {
			log.Printf("periodic: skip user %d: %v", user.ID, err)
		} else if len(result.Granted) > 0 || len(result.Revoked) > 0 {
			log.Printf("User %d: granted=%d revoked=%d", user.ID, len(result.Granted), len(result.Revoked))
			changed = append(changed, *result)
		}
		time.Sleep(rateDelay)
	}

	log.Println("Periodic subscription check complete")
	return changed
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
