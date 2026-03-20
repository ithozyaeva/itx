package handler

import (
	"log"
	"strconv"
	"sync"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CasinoHandler struct {
	svc      *service.CasinoService
	lastBet  sync.Map
}

func NewCasinoHandler() *CasinoHandler {
	return &CasinoHandler{
		svc: service.NewCasinoService(),
	}
}

func (h *CasinoHandler) checkRateLimit(memberId int64) error {
	now := time.Now()
	if last, ok := h.lastBet.Load(memberId); ok {
		if now.Sub(last.(time.Time)) < time.Second {
			return fiber.NewError(fiber.StatusTooManyRequests, "Подождите секунду между ставками")
		}
	}
	h.lastBet.Store(memberId, now)
	return nil
}

func (h *CasinoHandler) PlayCoinFlip(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	if err := h.checkRateLimit(member.Id); err != nil {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Подождите секунду между ставками"})
	}

	req := new(models.CoinFlipRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.PlayCoinFlip(member.Id, req)
	if err != nil {
		log.Printf("PlayCoinFlip error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось выполнить ставку"})
	}

	BroadcastEvent("minigames")
	return c.JSON(result)
}

func (h *CasinoHandler) PlayDiceRoll(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	if err := h.checkRateLimit(member.Id); err != nil {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Подождите секунду между ставками"})
	}

	req := new(models.DiceRollRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.PlayDiceRoll(member.Id, req)
	if err != nil {
		log.Printf("PlayDiceRoll error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось выполнить ставку"})
	}

	BroadcastEvent("minigames")
	return c.JSON(result)
}

func (h *CasinoHandler) PlayWheel(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	if err := h.checkRateLimit(member.Id); err != nil {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Подождите секунду между ставками"})
	}

	req := new(models.WheelRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.PlayWheel(member.Id, req)
	if err != nil {
		log.Printf("PlayWheel error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось выполнить ставку"})
	}

	BroadcastEvent("minigames")
	return c.JSON(result)
}

func (h *CasinoHandler) GetGlobalFeed(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	items, err := h.svc.GetGlobalFeed(limit)
	if err != nil {
		log.Printf("GetGlobalFeed error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки ленты"})
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *CasinoHandler) GetHistory(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	items, total, err := h.svc.GetHistory(member.Id, limit, offset)
	if err != nil {
		log.Printf("GetHistory error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки истории"})
	}

	return c.JSON(fiber.Map{"items": items, "total": total})
}

func (h *CasinoHandler) GetStats(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	stats, err := h.svc.GetStats(member.Id)
	if err != nil {
		log.Printf("GetStats error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки статистики"})
	}
	return c.JSON(stats)
}

func (h *CasinoHandler) GetAdminStats(c *fiber.Ctx) error {
	stats, err := h.svc.GetAdminStats()
	if err != nil {
		log.Printf("GetAdminStats error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки статистики"})
	}
	return c.JSON(stats)
}

func (h *CasinoHandler) GetAdminBets(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	var username *string
	if q := c.Query("username"); q != "" {
		username = &q
	}
	var game *string
	if q := c.Query("game"); q != "" {
		game = &q
	}

	items, total, err := h.svc.SearchBets(username, game, limit, offset)
	if err != nil {
		log.Printf("SearchBets error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка поиска ставок"})
	}

	return c.JSON(fiber.Map{"items": items, "total": total})
}
