package handler

import (
	"log"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type DailiesHandler struct {
	checkInSvc *service.CheckInService
	streakSvc  *service.StreakService
	taskSvc    *service.DailyTaskService
}

func NewDailiesHandler() *DailiesHandler {
	return &DailiesHandler{
		checkInSvc: service.NewCheckInService(),
		streakSvc:  service.NewStreakService(),
		taskSvc:    service.NewDailyTaskService(),
	}
}

// Today — GET /api/platform/dailies/today
func (h *DailiesHandler) Today(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	resp, err := h.taskSvc.GetMyToday(member.Id)
	if err != nil {
		log.Printf("dailies today (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки заданий"})
	}
	return c.JSON(resp)
}

// CheckIn — POST /api/platform/dailies/check-in
func (h *DailiesHandler) CheckIn(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}

	result, err := h.checkInSvc.CheckIn(member.Id)
	if err != nil {
		log.Printf("check-in failed (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось зафиксировать check-in"})
	}

	streak, err := h.streakSvc.BuildResponse(member.Id)
	if err != nil {
		log.Printf("streak build (member=%d): %v", member.Id, err)
	}

	if result.Inserted {
		PublishToMember(member.Id, "streak")
		PublishToMember(member.Id, "points")
		PublishToMember(member.Id, "dailies")
		PublishToMember(member.Id, "raffles")
	}

	return c.JSON(models.CheckInResponse{
		CheckInDone:   true,
		AlreadyToday:  !result.Inserted,
		Streak:        streak,
		RaffleEntered: result.RaffleEntered,
	})
}

// MyStreak — GET /api/platform/streak/me
func (h *DailiesHandler) MyStreak(c *fiber.Ctx) error {
	member, err := getMember(c)
	if err != nil {
		return err
	}
	resp, err := h.streakSvc.BuildResponse(member.Id)
	if err != nil {
		log.Printf("streak fetch (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки стрика"})
	}
	return c.JSON(resp)
}
