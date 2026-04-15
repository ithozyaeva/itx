package handler

import (
	"encoding/json"
	"fmt"
	"ithozyeva/config"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"net/http"
	"strconv"

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
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	tier, err := h.svc.GetTierBySlug(req.TierSlug)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Тир не найден"})
	}

	if err := h.svc.SetManualTier(userID, &tier.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось установить тир"})
	}

	h.svc.AddAudit(userID, "manual_override", map[string]interface{}{
		"tier_id": tier.ID, "tier_slug": tier.Slug, "source": "admin_panel",
	})

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
	}

	return c.JSON(fiber.Map{"success": true})
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

	if req.TierIDs != nil {
		if err := h.svc.SetChatTiers(chatID, *req.TierIDs); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось обновить тиры"})
		}
	}

	return c.JSON(fiber.Map{"success": true})
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
