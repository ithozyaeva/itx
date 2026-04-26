package handler

import (
	"strconv"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type ModerationHandler struct {
	svc *service.ModerationService
}

func NewModerationHandler(redisClient *redis.Client) *ModerationHandler {
	return &ModerationHandler{
		svc: service.NewModerationServiceWithRedis(redisClient),
	}
}

// GetActiveSanctions GET /api/admin/moderation/sanctions
func (h *ModerationHandler) GetActiveSanctions(c *fiber.Ctx) error {
	rows, err := h.svc.ListActiveSanctionsView()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"items": rows, "total": len(rows)})
}

// GetRecentActions GET /api/admin/moderation/actions
func (h *ModerationHandler) GetRecentActions(c *fiber.Ctx) error {
	rows, err := h.svc.ListRecentActionsView()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"items": rows, "total": len(rows)})
}

// GetGlobalBans GET /api/admin/moderation/global-bans
func (h *ModerationHandler) GetGlobalBans(c *fiber.Ctx) error {
	bans, err := h.svc.ListActiveGlobalBans()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"items": bans, "total": len(bans)})
}

// GetOpenVotebans GET /api/admin/moderation/votebans
func (h *ModerationHandler) GetOpenVotebans(c *fiber.Ctx) error {
	rows, err := h.svc.ListOpenVotebansEnriched()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"items": rows, "total": len(rows)})
}

// RevokeSanction POST /api/admin/moderation/sanctions/:id/revoke
// Удаляет действующую санкцию (ban/mute/voteban_kick) — публикует событие
// в Redis для бота, который выполнит UnbanChatMember/RestrictChatMember.
func (h *ModerationHandler) RevokeSanction(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
	}
	action, err := h.svc.GetActionByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if action == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}

	actorMember := actorMemberID(c)
	if err := h.svc.PublishRevoke(c.Context(), service.ModerationRevokeEvent{
		Kind:         service.RevokeKindSanction,
		ActionID:     action.Id,
		ChatID:       action.ChatID,
		TargetUserID: action.TargetUserID,
		ActorMember:  actorMember,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "publish failed: " + err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true})
}

// RevokeGlobalBan DELETE /api/admin/moderation/global-bans/:user_id
func (h *ModerationHandler) RevokeGlobalBan(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad user_id"})
	}
	if err := h.svc.PublishRevoke(c.Context(), service.ModerationRevokeEvent{
		Kind:         service.RevokeKindGlobalBan,
		TargetUserID: userID,
		ActorMember:  actorMemberID(c),
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "publish failed: " + err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true})
}

// CancelVoteban POST /api/admin/moderation/votebans/:id/cancel
func (h *ModerationHandler) CancelVoteban(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
	}
	vb, changed, err := h.svc.CancelOpenVoteban(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if vb == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	if changed {
		// Просим бота отредактировать poll-сообщение.
		_ = h.svc.PublishRevoke(c.Context(), service.ModerationRevokeEvent{
			Kind:         service.RevokeKindVoteban,
			VotebanID:    vb.Id,
			ChatID:       vb.ChatID,
			TargetUserID: vb.TargetUserID,
			ActorMember:  actorMemberID(c),
		})
	}
	return c.JSON(fiber.Map{"ok": true, "changed": changed})
}

// actorMemberID извлекает Member.Id из context, проставленного RequireAuth.
// 0 — если что-то пошло не так (не должен происходить в защищённой группе).
func actorMemberID(c *fiber.Ctx) int64 {
	if v := c.Locals("member"); v != nil {
		if m, ok := v.(*models.Member); ok && m != nil {
			return m.Id
		}
	}
	return 0
}
