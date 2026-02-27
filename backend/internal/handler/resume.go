package handler

import (
	"io"
	"mime"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
)

const maxResumeSize = 10 * 1024 * 1024 // 10 MB

type ResumeHandler struct {
	svc      *service.ResumeService
	auditSvc *service.AuditService
}

func NewResumeHandler() *ResumeHandler {
	return &ResumeHandler{
		svc:      service.NewResumeService(),
		auditSvc: service.NewAuditService(),
	}
}

func (h *ResumeHandler) Upload(c *fiber.Ctx) error {
	member, ok := c.Locals("member").(*models.Member)
	if !ok {
		return fiber.ErrUnauthorized
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Файл обязателен"})
	}
	if fileHeader.Size > maxResumeSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Файл превышает 10MB"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		if guessed := mime.TypeByExtension(filepath.Ext(fileHeader.Filename)); guessed != "" {
			contentType = guessed
		} else {
			contentType = "application/octet-stream"
		}
	}

	req := &models.CreateResumeRequest{
		WorkExperience:  c.FormValue("workExperience"),
		DesiredPosition: c.FormValue("desiredPosition"),
	}

	if wf := c.FormValue("workFormat"); wf != "" {
		value := models.WorkFormat(strings.ToUpper(wf))
		if value.IsValid() {
			req.WorkFormat = value
		}
	}

	resume, parsed, err := h.svc.UploadResume(member, fileHeader.Filename, contentType, data, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "resume", resume.Id, fileHeader.Filename)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"resume": resume,
		"parsed": parsed,
	})
}

func (h *ResumeHandler) ListMy(c *fiber.Ctx) error {
	member, ok := c.Locals("member").(*models.Member)
	if !ok {
		return fiber.ErrUnauthorized
	}

	resumes, err := h.svc.ListByTelegramID(member.TelegramID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(resumes)
}

func (h *ResumeHandler) UpdateMy(c *fiber.Ctx) error {
	member, ok := c.Locals("member").(*models.Member)
	if !ok {
		return fiber.ErrUnauthorized
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный идентификатор"})
	}

	payload := new(models.UpdateResumeRequest)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректное тело запроса"})
	}

	if payload.WorkFormat != nil {
		value := models.WorkFormat(strings.ToUpper(string(*payload.WorkFormat)))
		if value.IsValid() {
			payload.WorkFormat = &value
		} else {
			return fiber.NewError(fiber.StatusBadRequest, "Недопустимый формат работы")
		}
	}

	resume, err := h.svc.UpdateResume(id, member.TelegramID, payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(resume)
}

func (h *ResumeHandler) DeleteMy(c *fiber.Ctx) error {
	member, ok := c.Locals("member").(*models.Member)
	if !ok {
		return fiber.ErrUnauthorized
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный идентификатор"})
	}

	resume, err := h.svc.DeleteResume(id, member.TelegramID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionDelete, "resume", id, resume.FileName)

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ResumeHandler) AdminList(c *fiber.Ctx) error {
	limit := queryIntPointer(c.Query("limit"))
	offset := queryIntPointer(c.Query("offset"))

	filter := parseAdminResumeFilter(c)

	result, err := h.svc.SearchForAdmin(limit, offset, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(result)
}

func (h *ResumeHandler) AdminDownload(c *fiber.Ctx) error {
	filter := parseAdminResumeFilter(c)
	data, err := h.svc.GenerateArchive(filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=resumes.zip")
	return c.Send(data)
}

func (h *ResumeHandler) AdminGet(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный идентификатор"})
	}

	resume, err := h.svc.GetByIdWithMember(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(resume)
}

func parseAdminResumeFilter(c *fiber.Ctx) *models.ResumeFilter {
	filter := &models.ResumeFilter{}

	if wf := strings.TrimSpace(c.Query("workFormat")); wf != "" {
		value := models.WorkFormat(strings.ToUpper(wf))
		if value.IsValid() {
			filter.WorkFormat = &value
		}
	}

	if desired := strings.TrimSpace(c.Query("desiredPosition")); desired != "" {
		filter.DesiredPosition = &desired
	}

	if exp := strings.TrimSpace(c.Query("workExperience")); exp != "" {
		filter.WorkExperience = &exp
	}

	return filter
}

func queryIntPointer(value string) *int {
	if value == "" {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &parsed
}
