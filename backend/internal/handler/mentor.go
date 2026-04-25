package handler

import (
	"fmt"
	"log"
	"strconv"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

// MentorHandler обработчик для работы с менторами
type MentorHandler struct {
	BaseHandler[models.MentorDbShortModel]
	svc       *service.MentorService
	auditSvc  *service.AuditService
	pointsSvc *service.PointsService
}

// NewMentorHandler создает новый экземпляр обработчика менторов
func NewMentorHandler() *MentorHandler {
	svc := service.NewMentorService()
	return &MentorHandler{
		BaseHandler: *NewBaseHandler[models.MentorDbShortModel](svc),
		svc:         svc,
		auditSvc:    service.NewAuditService(),
		pointsSvc:   service.NewPointsService(),
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
		log.Printf("add review to service error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка добавления отзыва"})
	}

	return c.JSON(result)
}

// AddReviewFromPlatform добавляет отзыв от авторизованного пользователя платформы
func (h *MentorHandler) AddReviewFromPlatform(c *fiber.Ctx) error {
	type Request struct {
		ServiceId int    `json:"serviceId"`
		Text      string `json:"text"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member, err := getMember(c)
	if err != nil {
		return err
	}
	author := member.FirstName + " " + member.LastName

	review := &models.ReviewOnService{
		ServiceId:      req.ServiceId,
		Author:         author,
		Text:           req.Text,
		Date:           c.Context().Time().Format("2006-01-02"),
		Status:         models.ReviewOnServiceStatusDraft,
		AuthorMemberId: &member.Id,
	}

	result, err := h.svc.AddReviewToService(review)
	if err != nil {
		log.Printf("add review from platform error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка добавления отзыва"})
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
		log.Printf("create mentor error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка создания ментора"})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "mentor", result.Id, fmt.Sprintf("%s %s", result.FirstName, result.LastName))

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
		log.Printf("update mentor error (id=%d): %v", request.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления ментора"})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionUpdate, "mentor", result.Id, fmt.Sprintf("%s %s", result.FirstName, result.LastName))

	return c.JSON(result)
}

// Delete переопределяет базовый метод Delete для добавления аудит-лога
func (h *MentorHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	fullEntity, err := h.svc.GetByIdFull(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ментор не найден"})
	}

	entity, err := h.service.GetById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ментор не найден"})
	}

	entityName := fmt.Sprintf("%s %s", fullEntity.FirstName, fullEntity.LastName)

	if err := h.service.Delete(entity); err != nil {
		log.Printf("delete mentor error (id=%d): %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления ментора"})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionDelete, "mentor", int64(id), entityName)

	return c.SendStatus(fiber.StatusNoContent)
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
		log.Printf("get all mentors error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки менторов"})
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

	member, err := getMember(c)
	if err != nil {
		return err
	}
	existedMentor, err := h.svc.GetByMemberID(member.Id)

	if err != nil {
		log.Printf("get mentor by member ID error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка получения данных ментора"})
	}

	existedMentor.Occupation = req.Occupation
	existedMentor.Experience = req.Experience

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		log.Printf("update mentor info error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления информации ментора"})
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

	member, err := getMember(c)
	if err != nil {
		return err
	}
	existedMentor, err := h.svc.GetByMemberID(member.Id)

	if err != nil {
		log.Printf("get mentor by member ID error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка получения данных ментора"})
	}

	existedMentor.ProfTags = req.ProfTags

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		log.Printf("update mentor prof tags error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления тегов ментора"})
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

	member, err := getMember(c)
	if err != nil {
		return err
	}
	existedMentor, err := h.svc.GetByMemberID(member.Id)

	if err != nil {
		log.Printf("get mentor by member ID error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка получения данных ментора"})
	}

	existedMentor.Contacts = req.Contacts

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		log.Printf("update mentor contacts error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления контактов ментора"})
	}

	return c.JSON(result)
}

type UpdateOrderRequest struct {
	Order int `json:"order"`
}

func (h *MentorHandler) UpdateOrder(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	req := new(UpdateOrderRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := database.DB.Model(&models.MentorDbShortModel{}).Where("id = ?", id).Update("order", req.Order).Error; err != nil {
		log.Printf("update mentor order error (id=%d): %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления порядка ментора"})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionUpdate, "mentor", id, fmt.Sprintf("order=%d", req.Order))

	return c.JSON(fiber.Map{"ok": true})
}

type UpdateServicesRequest struct {
	Services []models.Service `json:"services"`
}

func (h *MentorHandler) UpdateServices(c *fiber.Ctx) error {
	req := new(UpdateServicesRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member, err := getMember(c)
	if err != nil {
		return err
	}
	existedMentor, err := h.svc.GetByMemberID(member.Id)

	if err != nil {
		log.Printf("get mentor by member ID error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка получения данных ментора"})
	}

	existedMentor.Services = req.Services

	result, err := h.svc.UpdateWithRelations(existedMentor)
	if err != nil {
		log.Printf("update mentor services error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления услуг ментора"})
	}

	return c.JSON(result)
}
