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
	request := new(models.MentorDbModel)
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
	request := new(models.MentorDbModel)
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

type UpdateInfoRequest struct {
	Occupation string `json:"occupation"`
	Experience string `json:"experience"`
}

func (h *MentorHandler) UpdateInfo(c *fiber.Ctx) error {
	req := new(UpdateInfoRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	existedMentor, err := h.svc.GetByMemberID(c.Locals("member").(*models.Member).Id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	existedMentor.Occupation = req.Occupation
	existedMentor.Experience = req.Experience

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

type UpdateProfTagsRequest struct {
	ProfTags []models.ProfTag `json:"profTags"`
}

func (h *MentorHandler) UpdateProfTags(c *fiber.Ctx) error {
	req := new(UpdateProfTagsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	existedMentor, err := h.svc.GetByMemberID(c.Locals("member").(*models.Member).Id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	existedMentor.ProfTags = req.ProfTags

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

type UpdateContactsRequest struct {
	Contacts []models.Contact `json:"contacts"`
}

func (h *MentorHandler) UpdateContacts(c *fiber.Ctx) error {
	req := new(UpdateContactsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	existedMentor, err := h.svc.GetByMemberID(c.Locals("member").(*models.Member).Id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	existedMentor.Contacts = req.Contacts

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

type UpdateServicesRequest struct {
	Services []models.Service `json:"services"`
}

func (h *MentorHandler) UpdateServices(c *fiber.Ctx) error {
	req := new(UpdateServicesRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	existedMentor, err := h.svc.GetByMemberID(c.Locals("member").(*models.Member).Id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	existedMentor.Services = req.Services

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
