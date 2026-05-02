package handler

import (
	"log"
	"strconv"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type KudosHandler struct {
	svc *service.KudosService
}

func NewKudosHandler() *KudosHandler {
	return &KudosHandler{
		svc: service.NewKudosService(),
	}
}

func (h *KudosHandler) Send(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	req := new(models.CreateKudosRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	kudos, err := h.svc.Send(member.Id, req.ToId, req.Message)
	if err != nil {
		log.Printf("send kudos error (from=%d, to=%d): %v", member.Id, req.ToId, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось отправить благодарность"})
	}

	BroadcastEvent("kudos")
	PublishToMember(req.ToId, "points")
	service.TrackDailyTrigger(member.Id, "send_kudos", 1)
	return c.Status(fiber.StatusCreated).JSON(kudos)
}

func (h *KudosHandler) GetRecent(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	items, total, err := h.svc.GetRecent(limit, offset)
	if err != nil {
		log.Printf("get recent kudos error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки благодарностей"})
	}

	if member, mErr := getMember(c); mErr == nil && member != nil {
		service.TrackDailyTrigger(member.Id, "view_kudos", 1)
	}

	return c.JSON(fiber.Map{
		"items": items,
		"total": total,
	})
}
