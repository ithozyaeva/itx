package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// getMemberFromContext извлекает member из контекста
func getMemberFromContext(c *fiber.Ctx) *models.Member {
	if member, ok := c.Locals("member").(*models.Member); ok {
		return member
	}
	return nil
}

type ChatQuestHandler struct {
	service *service.ChatQuestService
}

func NewChatQuestHandler() *ChatQuestHandler {
	return &ChatQuestHandler{
		service: service.NewChatQuestService(),
	}
}

// GetAll возвращает все квесты для админки с пагинацией
func (h *ChatQuestHandler) GetAll(c *fiber.Ctx) error {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	quests, total, err := h.service.GetAllQuests(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки заданий"})
	}

	return c.JSON(fiber.Map{
		"items": quests,
		"total": total,
	})
}

// Create создаёт новый квест
func (h *ChatQuestHandler) Create(c *fiber.Ctx) error {
	var quest models.ChatQuest
	if err := c.BodyParser(&quest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	if quest.Title == "" || quest.TargetCount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Необходимо указать название и целевое количество"})
	}

	if err := h.service.CreateQuest(&quest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка создания задания"})
	}

	return c.Status(fiber.StatusCreated).JSON(quest)
}

// Update обновляет квест
func (h *ChatQuestHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	existing, err := h.service.GetQuestByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Задание не найдено"})
	}

	if err := c.BodyParser(existing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
	}
	existing.Id = id

	if err := h.service.UpdateQuest(existing); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления задания"})
	}

	return c.JSON(existing)
}

// Delete удаляет квест
func (h *ChatQuestHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.service.DeleteQuest(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления задания"})
	}

	return c.JSON(fiber.Map{"success": true})
}

// GetAllForMember возвращает все квесты с прогрессом текущего пользователя
func (h *ChatQuestHandler) GetAllForMember(c *fiber.Ctx) error {
	member := getMemberFromContext(c)
	if member == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Не авторизован"})
	}

	filter := c.Query("filter", "all")

	quests, err := h.service.GetAllQuestsForMember(member.Id, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки заданий"})
	}

	return c.JSON(quests)
}

// GetActive возвращает активные квесты с прогрессом текущего пользователя
func (h *ChatQuestHandler) GetActive(c *fiber.Ctx) error {
	member := getMemberFromContext(c)
	if member == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Не авторизован"})
	}

	quests, err := h.service.GetActiveQuestsForMember(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки заданий"})
	}

	return c.JSON(quests)
}
