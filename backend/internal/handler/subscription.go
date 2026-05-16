package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"ithozyeva/config"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
}

func NewSubscriptionHandler(redisClient *redis.Client) *SubscriptionHandler {
	return &SubscriptionHandler{
		svc: service.NewSubscriptionService(redisClient),
	}
}

// tierMap loads all tiers once and returns a map by ID.
func (h *SubscriptionHandler) tierMap() map[uint]models.SubscriptionTier {
	tiers, _ := h.svc.GetAllTiers()
	m := make(map[uint]models.SubscriptionTier, len(tiers))
	for _, t := range tiers {
		m[t.ID] = t
	}
	return m
}

func (h *SubscriptionHandler) GetStats(c *fiber.Ctx) error {
	totalUsers, _ := h.svc.CountAllUsers()
	tiers, _ := h.svc.GetAllTiers()
	chats, _ := h.svc.GetAllChats()
	tierCounts, _ := h.svc.CountAllUsersByTier()

	tierStats := make([]fiber.Map, 0, len(tiers))
	for _, t := range tiers {
		tierStats = append(tierStats, fiber.Map{
			"id":    t.ID,
			"slug":  t.Slug,
			"name":  t.Name,
			"level": t.Level,
			"users": tierCounts[t.ID],
		})
	}

	anchorCount := 0
	contentCount := 0
	for _, ch := range chats {
		if ch.AnchorForTierID != nil {
			anchorCount++
		} else {
			contentCount++
		}
	}

	return c.JSON(fiber.Map{
		"totalUsers":   totalUsers,
		"totalChats":   len(chats),
		"anchorChats":  anchorCount,
		"contentChats": contentCount,
		"tiers":        tierStats,
	})
}

// PublicTiers отдаёт публичные тарифы (is_public=true) для /tariffs страницы
// и прогрева в боте. Цены в рублях, features распарсены из JSONB.
func (h *SubscriptionHandler) PublicTiers(c *fiber.Ctx) error {
	tiers, err := h.svc.GetPublicTiers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить тарифы"})
	}
	return c.JSON(fiber.Map{"items": tiers})
}

func (h *SubscriptionHandler) GetTiers(c *fiber.Ctx) error {
	tiers, err := h.svc.GetAllTiers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить тиры"})
	}

	tierCounts, _ := h.svc.CountAllUsersByTier()

	items := make([]fiber.Map, 0, len(tiers))
	for _, t := range tiers {
		items = append(items, fiber.Map{
			"id":    t.ID,
			"slug":  t.Slug,
			"name":  t.Name,
			"level": t.Level,
			"users": tierCounts[t.ID],
		})
	}

	return c.JSON(fiber.Map{"items": items, "total": len(items)})
}

func (h *SubscriptionHandler) GetChats(c *fiber.Ctx) error {
	chats, err := h.svc.GetAllChats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить чаты"})
	}

	tm := h.tierMap()

	chatIDs := make([]int64, len(chats))
	for i, ch := range chats {
		chatIDs[i] = ch.ID
	}
	accessCounts, _ := h.svc.CountActiveAccessByChats(chatIDs)
	allTierChats, _ := h.svc.GetAllTierChats()

	items := make([]fiber.Map, 0, len(chats))
	for _, ch := range chats {
		item := fiber.Map{
			"id":          ch.ID,
			"title":       ch.Title,
			"chatType":    ch.ChatType,
			"activeUsers": accessCounts[ch.ID],
			"category":    ch.Category,
			"emoji":       ch.Emoji,
			"priority":    ch.Priority,
		}
		if ch.AnchorForTierID != nil {
			item["anchorForTierID"] = *ch.AnchorForTierID
			if tier, ok := tm[*ch.AnchorForTierID]; ok {
				item["anchorTierName"] = tier.Name
			}
		}
		if tierIDs, ok := allTierChats[ch.ID]; ok && len(tierIDs) > 0 {
			item["tierIDs"] = tierIDs
			names := make([]string, 0, len(tierIDs))
			for _, tid := range tierIDs {
				if tier, ok := tm[tid]; ok {
					names = append(names, tier.Name)
				}
			}
			item["tierNames"] = names
		}
		items = append(items, item)
	}

	return c.JSON(fiber.Map{"items": items, "total": len(items)})
}

func (h *SubscriptionHandler) GetUsers(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 {
		limit = 20
	}

	total, _ := h.svc.CountAllUsers()
	users, err := h.svc.GetPaginatedUsers(offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить пользователей"})
	}

	// Batch load tiers and access counts
	tm := h.tierMap()
	userIDs := make([]int64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}
	accessCounts, _ := h.svc.CountActiveAccessByUsers(userIDs)

	items := make([]fiber.Map, 0, len(users))
	for _, u := range users {
		item := fiber.Map{
			"id":          u.ID,
			"username":    u.Username,
			"fullName":    u.FullName,
			"isActive":    u.IsActive,
			"activeChats": accessCounts[u.ID],
			"createdAt":   u.CreatedAt,
		}

		effTierID := u.EffectiveTierID()
		if effTierID != nil {
			if tier, ok := tm[*effTierID]; ok {
				item["tierName"] = tier.Name
				item["tierSlug"] = tier.Slug
			}
		}
		if u.ManualTierID != nil {
			item["manualTierID"] = *u.ManualTierID
		}
		if u.ResolvedTierID != nil {
			item["resolvedTierID"] = *u.ResolvedTierID
		}
		if u.LastCheckAt != nil {
			item["lastCheckAt"] = u.LastCheckAt
		}

		items = append(items, item)
	}

	return c.JSON(fiber.Map{"items": items, "total": total})
}

func (h *SubscriptionHandler) GetUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	user, err := h.svc.GetUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Пользователь не найден"})
	}

	// Load all tiers once for name resolution
	tm := h.tierMap()

	result := fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"fullName": user.FullName,
		"isActive": user.IsActive,
	}

	if user.ResolvedTierID != nil {
		result["resolvedTierID"] = *user.ResolvedTierID
		if tier, ok := tm[*user.ResolvedTierID]; ok {
			result["resolvedTierName"] = tier.Name
		}
	}
	if user.ManualTierID != nil {
		result["manualTierID"] = *user.ManualTierID
		if tier, ok := tm[*user.ManualTierID]; ok {
			result["manualTierName"] = tier.Name
		}
	}
	if user.ManualTierExpiresAt != nil {
		result["manualTierExpiresAt"] = user.ManualTierExpiresAt
	}
	if effTierID := user.EffectiveTierID(); effTierID != nil {
		if tier, ok := tm[*effTierID]; ok {
			result["effectiveTierName"] = tier.Name
		}
	}

	if user.LastCheckAt != nil {
		result["lastCheckAt"] = user.LastCheckAt
	}
	result["createdAt"] = user.CreatedAt

	// Load access and batch-resolve chat titles
	access, _ := h.svc.GetActiveAccess(userID)
	allChats, _ := h.svc.GetAllChats()
	chatMap := make(map[int64]string, len(allChats))
	for _, ch := range allChats {
		chatMap[ch.ID] = ch.Title
	}

	chatAccess := make([]fiber.Map, 0, len(access))
	for _, a := range access {
		chatItem := fiber.Map{
			"chatID":    a.ChatID,
			"grantedAt": a.GrantedAt,
		}
		if title, ok := chatMap[a.ChatID]; ok {
			chatItem["chatTitle"] = title
		}
		chatAccess = append(chatAccess, chatItem)
	}
	result["access"] = chatAccess

	return c.JSON(result)
}

func (h *SubscriptionHandler) SetOverride(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	var req struct {
		TierSlug string `json:"tierSlug"`
		// Months > 0 — выдать на N месяцев (как платная подписка в обход Boosty).
		// 0 / отсутствует — бессрочный grant (триггерит ErrBessrochnyGrantExists
		// и блокирует апгрейд за кредиты, см. service/subscription.go:1147).
		Months int `json:"months"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	if req.Months < 0 || req.Months > 60 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Срок должен быть от 0 до 60 месяцев"})
	}

	tier, err := h.svc.GetTierBySlug(req.TierSlug)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Тир не найден"})
	}

	var expiresAt *time.Time
	if req.Months > 0 {
		t := time.Now().AddDate(0, req.Months, 0)
		expiresAt = &t
	}

	if err := h.svc.SetManualTierWithExpiry(userID, &tier.ID, expiresAt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось установить тир"})
	}

	auditPayload := map[string]interface{}{
		"tier_id": tier.ID, "tier_slug": tier.Slug, "source": "admin_panel",
	}
	if expiresAt != nil {
		auditPayload["months"] = req.Months
		auditPayload["expires_at"] = *expiresAt
	}
	h.svc.AddAudit(userID, "manual_override", auditPayload)

	return c.JSON(fiber.Map{"success": true})
}

func (h *SubscriptionHandler) ClearOverride(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.svc.SetManualTier(userID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось снять тир"})
	}

	h.svc.AddAudit(userID, "manual_override", map[string]interface{}{
		"tier": nil, "source": "admin_panel",
	})

	return c.JSON(fiber.Map{"success": true})
}

func (h *SubscriptionHandler) CreateChat(c *fiber.Ctx) error {
	var req struct {
		ID              int64   `json:"id"`
		Title           string  `json:"title"`
		ChatType        string  `json:"chatType"`
		AnchorForTierID *uint   `json:"anchorForTierID"`
		TierIDs         []uint  `json:"tierIDs"`
		Category        *string `json:"category"`
		Emoji           *string `json:"emoji"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	if req.ID == 0 || req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID и название обязательны"})
	}
	if req.ChatType == "" {
		req.ChatType = "supergroup"
	}

	if err := h.svc.UpsertChat(req.ID, req.Title, req.ChatType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать чат"})
	}

	if req.Category != nil || req.Emoji != nil {
		if err := h.svc.UpdateChatMeta(req.ID, req.Category, req.Emoji); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось сохранить категорию"})
		}
	}

	if req.AnchorForTierID != nil {
		if err := h.svc.SetAnchor(req.ID, req.AnchorForTierID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось установить anchor"})
		}
	}

	if len(req.TierIDs) > 0 {
		if err := h.svc.SetChatTiers(req.ID, req.TierIDs); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось привязать тиры"})
		}
		// Новый чат → все его тиры «новые». Шлём одно событие с min level.
		if level, ok := h.minTierLevel(req.TierIDs); ok {
			if err := h.svc.PublishNewChatAccess(c.Context(), req.ID, level); err != nil {
				log.Printf("CreateChat: publish new-chat-access chat=%d level=%d failed: %v", req.ID, level, err)
			}
		}
	}

	return c.JSON(fiber.Map{"success": true})
}

// minTierLevel возвращает минимальный level среди указанных tierID, чтобы
// publisher мог уведомить всех пользователей, у кого effective tier >= этого
// значения. Если ни один tier не найден — возвращает (0, false).
func (h *SubscriptionHandler) minTierLevel(tierIDs []uint) (int, bool) {
	min, found := 0, false
	for _, tid := range tierIDs {
		tier, err := h.svc.GetTier(tid)
		if err != nil {
			continue
		}
		if !found || tier.Level < min {
			min = tier.Level
			found = true
		}
	}
	return min, found
}

func (h *SubscriptionHandler) UpdateChat(c *fiber.Ctx) error {
	chatID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	existing, err := h.svc.GetChat(chatID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Чат не найден"})
	}

	var req struct {
		Title           *string `json:"title"`
		AnchorForTierID *uint   `json:"anchorForTierID"`
		ClearAnchor     bool    `json:"clearAnchor"`
		TierIDs         *[]uint `json:"tierIDs"`
		Category        *string `json:"category"`
		Emoji           *string `json:"emoji"`
		ClearCategory   bool    `json:"clearCategory"`
		Priority        *int    `json:"priority"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	if req.Title != nil && *req.Title != "" {
		if err := h.svc.UpsertChat(chatID, *req.Title, existing.ChatType); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось обновить чат"})
		}
	}

	if req.Category != nil || req.Emoji != nil || req.ClearCategory {
		cat := req.Category
		emoji := req.Emoji
		if req.ClearCategory {
			cat = nil
			emoji = nil
		}
		if err := h.svc.UpdateChatMeta(chatID, cat, emoji); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось сохранить категорию"})
		}
	}

	if req.ClearAnchor {
		if err := h.svc.SetAnchor(chatID, nil); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось снять anchor"})
		}
	} else if req.AnchorForTierID != nil {
		if err := h.svc.SetAnchor(chatID, req.AnchorForTierID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось установить anchor"})
		}
	}

	var addedTiers []uint
	if req.TierIDs != nil {
		// Фиксируем старые привязки до SetChatTiers, чтобы вычислить diff и
		// опубликовать рассылку только по реально добавленным тирам.
		oldTierIDs, _ := h.svc.GetTierIDsForChat(chatID)
		if err := h.svc.SetChatTiers(chatID, *req.TierIDs); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось обновить тиры"})
		}
		addedTiers = diffTiers(*req.TierIDs, oldTierIDs)
	}

	if req.Priority != nil {
		if err := h.svc.SetChatPriority(chatID, *req.Priority); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось сохранить приоритет"})
		}
	}

	// Публикуем события уже ПОСЛЕ всех БД-изменений, иначе подписчик
	// может прибежать с invite раньше, чем чат стал доступным тиру.
	// Одно событие на чат — бот рассылает всем с level >= min(added).
	if level, ok := h.minTierLevel(addedTiers); ok {
		if err := h.svc.PublishNewChatAccess(c.Context(), chatID, level); err != nil {
			log.Printf("UpdateChat: publish new-chat-access chat=%d level=%d failed: %v", chatID, level, err)
		}
	}

	return c.JSON(fiber.Map{"success": true})
}

// diffTiers возвращает tierID, которые есть в next, но отсутствуют в prev.
func diffTiers(next, prev []uint) []uint {
	prevSet := make(map[uint]struct{}, len(prev))
	for _, t := range prev {
		prevSet[t] = struct{}{}
	}
	added := make([]uint, 0, len(next))
	for _, t := range next {
		if _, ok := prevSet[t]; !ok {
			added = append(added, t)
		}
	}
	return added
}

func (h *SubscriptionHandler) DeleteChat(c *fiber.Ctx) error {
	chatID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.svc.DeleteChat(chatID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось удалить чат"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *SubscriptionHandler) GetChatDetail(c *fiber.Ctx) error {
	chatID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	chat, err := h.svc.GetChat(chatID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Чат не найден"})
	}

	tierIDs, _ := h.svc.GetTierIDsForChat(chatID)

	result := fiber.Map{
		"id":       chat.ID,
		"title":    chat.Title,
		"chatType": chat.ChatType,
		"tierIDs":  tierIDs,
		"category": chat.Category,
		"emoji":    chat.Emoji,
		"priority": chat.Priority,
	}
	if chat.AnchorForTierID != nil {
		result["anchorForTierID"] = *chat.AnchorForTierID
	}

	return c.JSON(result)
}

func (h *SubscriptionHandler) ResolveChat(c *fiber.Ctx) error {
	chatID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChat?chat_id=%d", config.CFG.TelegramToken, chatID)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Не удалось связаться с Telegram"})
	}
	defer resp.Body.Close()

	var tgResp struct {
		OK     bool `json:"ok"`
		Result struct {
			Title string `json:"title"`
			Type  string `json:"type"`
		} `json:"result"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tgResp); err != nil || !tgResp.OK {
		desc := tgResp.Description
		if desc == "" {
			desc = "Чат не найден или бот не является участником"
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": desc})
	}

	return c.JSON(fiber.Map{
		"id":       chatID,
		"title":    tgResp.Result.Title,
		"chatType": tgResp.Result.Type,
	})
}

// GetInternalUserSubscription отдаёт статус подписки по telegram-id для
// server-to-server вызовов от соседних сервисов (защищён shared secret через
// RequireInternalSecret). Если пользователь не найден или нет активного тира —
// возвращает 200 с is_subscriber=false и tier=null, чтобы консьюмеру не нужно
// было различать 404 и пустой ответ.
func (h *SubscriptionHandler) GetInternalUserSubscription(c *fiber.Ctx) error {
	tgID, err := strconv.ParseInt(c.Params("tg_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный tg_id"})
	}

	resp := fiber.Map{
		"tg_id":         tgID,
		"is_subscriber": false,
		"tier":          nil,
	}

	user, err := h.svc.GetUser(tgID)
	if err != nil || user == nil {
		return c.JSON(resp)
	}

	effTierID := user.EffectiveTierID()
	if effTierID == nil {
		return c.JSON(resp)
	}

	tier, err := h.svc.GetTier(*effTierID)
	if err != nil {
		return c.JSON(resp)
	}

	resp["is_subscriber"] = true
	resp["tier"] = fiber.Map{
		"slug":  tier.Slug,
		"name":  tier.Name,
		"level": tier.Level,
	}
	return c.JSON(resp)
}

// PurchaseWithCredits — покупка тарифа за реферальные кредиты.
// Доступен любому авторизованному (UNSUBSCRIBER тоже): это и есть точка
// входа в подписку для тех, кто накопил кредиты, но не платил через Boosty.
//
// Ошибки:
//   400 — пустой tier_slug или tier.price_credits IS NULL (не покупаем за credits)
//   402 — недостаточно кредитов на балансе
//   500 — прочее
func (h *SubscriptionHandler) PurchaseWithCredits(c *fiber.Ctx) error {
	var req struct {
		TierSlug string `json:"tier_slug"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	if req.TierSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tier_slug обязателен"})
	}

	member, err := getMember(c)
	if err != nil {
		return err
	}

	fullName := strings.TrimSpace(member.FirstName + " " + member.LastName)
	var username *string
	if member.Username != "" {
		u := member.Username
		username = &u
	}

	result, err := h.svc.PurchaseTierWithCredits(member.Id, member.TelegramID, username, fullName, req.TierSlug)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientCredits) {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{"error": "Недостаточно кредитов"})
		}
		if errors.Is(err, service.ErrTierNotPurchasable) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Этот тариф нельзя купить за кредиты"})
		}
		if errors.Is(err, service.ErrBessrochnyGrantExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "У вас уже бессрочная подписка от администратора"})
		}
		if errors.Is(err, service.ErrTierDowngrade) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Нельзя купить тариф ниже текущего"})
		}
		log.Printf("PurchaseWithCredits error (member=%d, slug=%s): %v", member.Id, req.TierSlug, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось купить подписку"})
	}
	return c.JSON(result)
}

func (h *SubscriptionHandler) RevokeAccess(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный user ID"})
	}

	chatID, err := strconv.ParseInt(c.Params("chatId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный chat ID"})
	}

	if err := h.svc.RevokeAccess(userID, chatID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось отозвать доступ"})
	}

	h.svc.AddAudit(userID, "revoke", map[string]interface{}{
		"chat_id": chatID, "source": "admin_panel",
	})

	return c.JSON(fiber.Map{"success": true})
}
