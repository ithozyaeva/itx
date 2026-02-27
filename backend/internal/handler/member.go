package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// MembersHandler обработчик для работы с участниками
type MembersHandler struct {
	svc      *service.MemberService
	auditSvc *service.AuditService
}

// NewMembersHandler создает новый экземпляр обработчика участников
func NewMembersHandler() *MembersHandler {
	svc := service.NewMemberService()
	return &MembersHandler{
		svc:      svc,
		auditSvc: service.NewAuditService(),
	}
}

type SearchMembersRequest struct {
	Limit    *int     `query:"limit"`
	Offset   *int     `query:"offset"`
	Username *string  `query:"username"`
	Roles    []string `query:"roles"`
}

// Search выполняет поиск участников с пагинацией
func (h *MembersHandler) Search(c *fiber.Ctx) error {
	req := new(SearchMembersRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	filter := make(repository.SearchFilter)

	if req.Username != nil {
		filter["username ILIKE ?"] = "%" + *req.Username + "%"
	}

	if len(req.Roles) > 0 {
		filter["EXISTS (SELECT 1 FROM member_roles WHERE member_id = members.id AND role IN ?)"] = req.Roles
	}

	var finalFilter *repository.SearchFilter
	if len(filter) > 0 {
		finalFilter = &filter
	} else {
		finalFilter = nil
	}

	result, err := h.svc.Search(req.Limit, req.Offset, finalFilter, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// GetById получает участника по ID
func (h *MembersHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	result, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// Реализуем все необходимые методы напрямую
func (h *MembersHandler) Create(c *fiber.Ctx) error {
	request := new(models.Member)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.Create(request)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "member", result.Id, result.FirstName+" "+result.LastName)

	return c.JSON(result)
}

type UpdateRequest struct {
	Id        int64         `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Birthday  *string       `json:"birthday"`
	Roles     []models.Role `json:"roles"`
	Username  string        `json:"tg"`
}

func (h *MembersHandler) Update(c *fiber.Ctx) error {
	request := new(UpdateRequest)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member, err := h.svc.GetById(request.Id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Участник не найден"})
	}

	member.FirstName = request.FirstName
	member.LastName = request.LastName
	member.Roles = request.Roles
	member.Username = request.Username

	parsedDate, err := utils.ParseDate(request.Birthday)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	member.Birthday = parsedDate

	result, err := h.svc.Update(member)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionUpdate, "member", result.Id, result.FirstName+" "+result.LastName)

	return c.JSON(result)
}

func (h *MembersHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity, err := h.svc.GetById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Участник не найден"})
	}

	entityName := entity.FirstName + " " + entity.LastName

	if err := h.svc.Delete(entity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionDelete, "member", int64(id), entityName)

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MembersHandler) Me(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	mentor, err := h.svc.GetMentor(member.Id)

	if err != nil {
		return c.JSON(member)
	}

	return c.JSON(mentor)
}

func (h *MembersHandler) UpdateProfile(c *fiber.Ctx) error {
	request := new(UpdateRequest)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)
	member.FirstName = request.FirstName
	member.LastName = request.LastName

	parsedDate, err := utils.ParseDate(request.Birthday)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	member.Birthday = parsedDate

	result, err := h.svc.Update(member)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	mentor, err := h.svc.GetMentor(member.Id)

	if err != nil {
		return c.JSON(result)
	}

	return c.JSON(mentor)
}

func (h *MembersHandler) GetPermissions(c *fiber.Ctx) error {
	member, ok := c.Locals("member").(*models.Member)
	if !ok || member == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	permissions, err := h.svc.GetPermissions(member.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(permissions)
}
