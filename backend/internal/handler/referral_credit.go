package handler

import (
	"errors"
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReferralCreditHandler struct {
	svc      *service.ReferralCreditService
	auditSvc *service.AuditService
}

func NewReferralCreditHandler() *ReferralCreditHandler {
	return &ReferralCreditHandler{
		svc:      service.NewReferralCreditService(),
		auditSvc: service.NewAuditService(),
	}
}

// GetMine — баланс кредитов и последние 50 транзакций.
// Доступен любому авторизованному (даже UNSUBSCRIBER'у), чтобы юзер
// мог увидеть, хватает ли ему кредитов до покупки.
func (h *ReferralCreditHandler) GetMine(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	summary, err := h.svc.GetSummary(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить кредиты"})
	}
	return c.JSON(summary)
}

func (h *ReferralCreditHandler) AdminSearch(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")
	username := c.Query("username")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 {
		limit = 20
	}

	var usernamePtr *string
	if username != "" {
		usernamePtr = &username
	}
	items, total, err := h.svc.SearchTransactions(usernamePtr, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить транзакции"})
	}
	return c.JSON(fiber.Map{"items": items, "total": total})
}

func (h *ReferralCreditHandler) AdminAward(c *fiber.Ctx) error {
	var req models.AdminAwardCreditsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	if req.MemberId <= 0 || req.Amount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "memberId обязателен; amount должен быть != 0"})
	}
	if err := h.svc.AdminAward(req.MemberId, req.Amount, req.Description); err != nil {
		if errors.Is(err, repository.ErrInsufficientCredits) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Списание невозможно: на балансе меньше указанной суммы"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось начислить кредиты"})
	}

	// Audit + notification: credits — реальная стоимость (через них покупается
	// подписка), поэтому ручное движение баланса должно быть отслеживаемо
	// и явно сообщено владельцу.
	auditDesc := fmt.Sprintf("amount=%d %s", req.Amount, req.Description)
	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "referral_credits", req.MemberId, auditDesc)
	notifTitle := "Начисление кредитов"
	notifBody := fmt.Sprintf("Администратор начислил вам %d кредитов", req.Amount)
	if req.Amount < 0 {
		notifTitle = "Списание кредитов"
		notifBody = fmt.Sprintf("Администратор списал у вас %d кредитов", -req.Amount)
	}
	if req.Description != "" {
		notifBody += ": " + req.Description
	}
	go CreateNotification(req.MemberId, "credits_admin", notifTitle, notifBody)

	return c.JSON(fiber.Map{"success": true})
}
