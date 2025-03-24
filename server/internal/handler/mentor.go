package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// MentorHandler обработчик для работы с менторами
type MentorHandler struct {
	BaseHandler[models.MentorDbShortModel]
	svc *service.MentorService
}

// NewMentorHandler создает новый экземпляр обработчика менторов
func NewMentorHandler() *MentorHandler {
	svc := service.NewMentorService()
	return &MentorHandler{
		BaseHandler: *NewBaseHandler[models.MentorDbShortModel](svc),
		svc:         svc,
	}
}

// GetById получает ментора по ID с полной информацией
func (h *MentorHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity, err := h.svc.GetByIdFull(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ментор не найден"})
	}

	return c.JSON(entity)
}

// FindByTagRequest структура запроса для поиска менторов по тегу
type FindByTagRequest struct {
	TagId  int64 `json:"tagId"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// FindByTag находит менторов по тегу
func (h *MentorHandler) FindByTag(c *fiber.Ctx) error {
	req := new(FindByTagRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.FindByTag(req.TagId, req.Limit, req.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// AddReviewToService добавляет отзыв к услуге ментора
func (h *MentorHandler) AddReviewToService(c *fiber.Ctx) error {
	review := new(models.ReviewOnService)
	if err := c.BodyParser(review); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.AddReviewToService(review)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// Create создает нового ментора со всеми связанными сущностями
func (h *MentorHandler) Create(c *fiber.Ctx) error {
	request := new(models.MentorCreateUpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.CreateWithRelations(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Update обновляет ментора со всеми связанными сущностями
func (h *MentorHandler) Update(c *fiber.Ctx) error {
	request := new(models.MentorCreateUpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	// Проверяем, что ID указан
	if request.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID ментора не указан"})
	}

	result, err := h.svc.UpdateWithRelations(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *MentorHandler) GetServices(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.GetServices(int64(id))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	return c.JSON(result)
}

// GetAllWithRelations получает всех менторов с полной информацией
func (h *MentorHandler) GetAllWithRelations(c *fiber.Ctx) error {
	req := new(models.SearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.GetAllWithRelations(req.Limit, req.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
