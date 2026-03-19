package handler

import (
	"ithozyeva/internal/service"
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

func (h *SubscriptionHandler) GetStats(c *fiber.Ctx) error {
	totalUsers, _ := h.svc.CountAllUsers()
	tiers, _ := h.svc.GetAllTiers()
	chats, _ := h.svc.GetAllChats()

	tierStats := make([]fiber.Map, 0, len(tiers))
	for _, t := range tiers {
		count, _ := h.svc.CountUsersByTier(t.ID)
		tierStats = append(tierStats, fiber.Map{
			"id":    t.ID,
			"slug":  t.Slug,
			"name":  t.Name,
			"level": t.Level,
			"users": count,
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

	items := make([]fiber.Map, 0, len(tiers))
	for _, t := range tiers {
		count, _ := h.svc.CountUsersByTier(t.ID)
		items = append(items, fiber.Map{
			"id":    t.ID,
			"slug":  t.Slug,
			"name":  t.Name,
			"level": t.Level,
			"users": count,
		})
	}

	return c.JSON(fiber.Map{"items": items, "total": len(items)})
}

func (h *SubscriptionHandler) GetChats(c *fiber.Ctx) error {
	chats, err := h.svc.GetAllChats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить чаты"})
	}

	items := make([]fiber.Map, 0, len(chats))
	for _, ch := range chats {
		item := fiber.Map{
			"id":       ch.ID,
			"title":    ch.Title,
			"chatType": ch.ChatType,
		}
		if ch.AnchorForTierID != nil {
			item["anchorForTierID"] = *ch.AnchorForTierID
			tier, err := h.svc.GetTier(*ch.AnchorForTierID)
			if err == nil {
				item["anchorTierName"] = tier.Name
			}
		}
		users, _ := h.svc.GetUsersWithAccessToChat(ch.ID)
		item["activeUsers"] = len(users)
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

	items := make([]fiber.Map, 0, len(users))
	for _, u := range users {
		item := fiber.Map{
			"id":       u.ID,
			"username": u.Username,
			"fullName": u.FullName,
			"isActive": u.IsActive,
		}

		effTierID := u.EffectiveTierID()
		if effTierID != nil {
			tier, err := h.svc.GetTier(*effTierID)
			if err == nil {
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

		access, _ := h.svc.GetActiveAccess(u.ID)
		item["activeChats"] = len(access)
		item["createdAt"] = u.CreatedAt

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

	result := fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"fullName": user.FullName,
		"isActive": user.IsActive,
	}

	if user.ResolvedTierID != nil {
		result["resolvedTierID"] = *user.ResolvedTierID
		tier, err := h.svc.GetTier(*user.ResolvedTierID)
		if err == nil {
			result["resolvedTierName"] = tier.Name
		}
	}
	if user.ManualTierID != nil {
		result["manualTierID"] = *user.ManualTierID
		tier, err := h.svc.GetTier(*user.ManualTierID)
		if err == nil {
			result["manualTierName"] = tier.Name
		}
	}

	effTierID := user.EffectiveTierID()
	if effTierID != nil {
		tier, err := h.svc.GetTier(*effTierID)
		if err == nil {
			result["effectiveTierName"] = tier.Name
		}
	}

	if user.LastCheckAt != nil {
		result["lastCheckAt"] = user.LastCheckAt
	}
	result["createdAt"] = user.CreatedAt

	access, _ := h.svc.GetActiveAccess(userID)
	chatAccess := make([]fiber.Map, 0, len(access))
	for _, a := range access {
		chatItem := fiber.Map{
			"chatID":    a.ChatID,
			"grantedAt": a.GrantedAt,
		}
		chat, err := h.svc.GetChat(a.ChatID)
		if err == nil {
			chatItem["chatTitle"] = chat.Title
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
