package handler

import (
	"errors"
	"log"
	"strconv"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AdminDailyTaskHandler — CRUD пула шаблонов дейликов.
// Сами наборы (daily_task_sets) формируются cron-горутиной автоматически.
type AdminDailyTaskHandler struct {
	repo *repository.DailyTaskRepository
}

func NewAdminDailyTaskHandler() *AdminDailyTaskHandler {
	return &AdminDailyTaskHandler{repo: repository.NewDailyTaskRepository()}
}

func (h *AdminDailyTaskHandler) List(c *fiber.Ctx) error {
	tasks, err := h.repo.GetAllAdmin()
	if err != nil {
		log.Printf("admin daily tasks list: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки"})
	}
	return c.JSON(fiber.Map{"items": tasks})
}

func (h *AdminDailyTaskHandler) Create(c *fiber.Ctx) error {
	t := new(models.DailyTask)
	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if err := validateDailyTask(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	t.Id = 0
	if err := h.repo.Create(t); err != nil {
		log.Printf("admin daily task create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка создания"})
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}

func (h *AdminDailyTaskHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	existing, err := h.repo.GetById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Не найдено"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка"})
	}
	patch := new(models.DailyTask)
	if err := c.BodyParser(patch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	// Обновляем только поля, разрешённые админу. Code оставляем неизменным
	// для предотвращения коллизий с уже выпавшими сетами (UNIQUE по code).
	existing.Title = patch.Title
	existing.Description = patch.Description
	existing.Icon = patch.Icon
	existing.Tier = patch.Tier
	existing.Points = patch.Points
	existing.Target = patch.Target
	existing.TriggerKey = patch.TriggerKey
	existing.Active = patch.Active
	if err := validateDailyTask(existing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.repo.Update(existing); err != nil {
		log.Printf("admin daily task update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления"})
	}
	return c.JSON(existing)
}

func (h *AdminDailyTaskHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	if err := h.repo.Delete(id); err != nil {
		log.Printf("admin daily task delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AdminDailyTaskHandler) RecentSets(c *fiber.Ctx) error {
	limit := 14
	if v := c.QueryInt("limit", 0); v > 0 && v <= 90 {
		limit = v
	}
	sets, err := h.repo.GetRecentSets(limit)
	if err != nil {
		log.Printf("admin daily task sets: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки"})
	}
	return c.JSON(fiber.Map{"items": sets})
}

func validateDailyTask(t *models.DailyTask) error {
	if t.Code == "" {
		return errors.New("code обязателен")
	}
	if t.Title == "" {
		return errors.New("title обязателен")
	}
	if t.Tier != models.DailyTaskTierEngagement &&
		t.Tier != models.DailyTaskTierLight &&
		t.Tier != models.DailyTaskTierMeaningful &&
		t.Tier != models.DailyTaskTierBig {
		return errors.New("неизвестный tier")
	}
	if t.Points < 0 {
		return errors.New("points >= 0")
	}
	if t.Target <= 0 {
		t.Target = 1
	}
	if t.TriggerKey == "" {
		return errors.New("trigger_key обязателен")
	}
	return nil
}

// AdminChallengeHandler — CRUD шаблонов челленджей.
type AdminChallengeHandler struct {
	repo *repository.ChallengeRepository
}

func NewAdminChallengeHandler() *AdminChallengeHandler {
	return &AdminChallengeHandler{repo: repository.NewChallengeRepository()}
}

func (h *AdminChallengeHandler) ListTemplates(c *fiber.Ctx) error {
	tpls, err := h.repo.GetAllTemplatesAdmin()
	if err != nil {
		log.Printf("admin challenge templates list: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки"})
	}
	return c.JSON(fiber.Map{"items": tpls})
}

func (h *AdminChallengeHandler) CreateTemplate(c *fiber.Ctx) error {
	t := new(models.ChallengeTemplate)
	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if err := validateChallengeTemplate(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	t.Id = 0
	if err := h.repo.CreateTemplate(t); err != nil {
		log.Printf("admin challenge template create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка создания"})
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}

func (h *AdminChallengeHandler) UpdateTemplate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	existing, err := h.repo.GetTemplateById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Не найдено"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка"})
	}
	patch := new(models.ChallengeTemplate)
	if err := c.BodyParser(patch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	existing.Title = patch.Title
	existing.Description = patch.Description
	existing.Icon = patch.Icon
	existing.Kind = patch.Kind
	existing.MetricKey = patch.MetricKey
	existing.Target = patch.Target
	existing.RewardPoints = patch.RewardPoints
	existing.AchievementCode = patch.AchievementCode
	existing.Active = patch.Active
	if err := validateChallengeTemplate(existing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.repo.UpdateTemplate(existing); err != nil {
		log.Printf("admin challenge template update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления"})
	}
	return c.JSON(existing)
}

func (h *AdminChallengeHandler) DeleteTemplate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}
	if err := h.repo.DeleteTemplate(id); err != nil {
		log.Printf("admin challenge template delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AdminChallengeHandler) ListInstances(c *fiber.Ctx) error {
	limit := 30
	if v := c.QueryInt("limit", 0); v > 0 && v <= 200 {
		limit = v
	}
	insts, err := h.repo.GetRecentInstances(limit)
	if err != nil {
		log.Printf("admin challenge instances: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки"})
	}
	return c.JSON(fiber.Map{"items": insts})
}

func validateChallengeTemplate(t *models.ChallengeTemplate) error {
	if t.Code == "" {
		return errors.New("code обязателен")
	}
	if t.Title == "" {
		return errors.New("title обязателен")
	}
	if t.Kind != models.ChallengeKindWeekly && t.Kind != models.ChallengeKindMonthly {
		return errors.New("kind должен быть weekly или monthly")
	}
	if t.MetricKey == "" {
		return errors.New("metric_key обязателен")
	}
	if t.Target <= 0 {
		return errors.New("target > 0")
	}
	if t.RewardPoints < 0 {
		return errors.New("reward_points >= 0")
	}
	return nil
}
