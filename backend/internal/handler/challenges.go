package handler

import (
	"log"

	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ChallengesHandler struct {
	svc *service.ChallengeService
}

func NewChallengesHandler() *ChallengesHandler {
	return &ChallengesHandler{
		svc: service.NewChallengeService(),
	}
}

// GetMyChallenges — GET /api/platform/challenges
func (h *ChallengesHandler) GetMyChallenges(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.GetMyChallenges(member.Id)
	if err != nil {
		log.Printf("get my challenges (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки челленджей"})
	}
	return c.JSON(resp)
}
