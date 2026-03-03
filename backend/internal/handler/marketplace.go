package handler

import (
	"fmt"
	"io"
	"log"
	"mime"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
)

type MarketplaceHandler struct {
	svc      *service.MarketplaceService
	pointSvc *service.PointsService
}

func NewMarketplaceHandler() *MarketplaceHandler {
	return &MarketplaceHandler{
		svc:      service.NewMarketplaceService(),
		pointSvc: service.NewPointsService(),
	}
}

func (h *MarketplaceHandler) Search(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")
	status := c.Query("status")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	items, total, err := h.svc.Search(statusPtr, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить объявления"})
	}

	return c.JSON(fiber.Map{"items": items, "total": total})
}

func (h *MarketplaceHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	item, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Объявление не найдено"})
	}

	return c.JSON(item)
}

func (h *MarketplaceHandler) Create(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	req := &models.CreateMarketplaceItemRequest{
		Title:           c.FormValue("title"),
		Description:     c.FormValue("description"),
		Price:           c.FormValue("price"),
		City:            c.FormValue("city"),
		Condition:       models.MarketplaceItemCondition(c.FormValue("condition")),
		Defects:         c.FormValue("defects"),
		PackageContents: c.FormValue("packageContents"),
		ContactTelegram: c.FormValue("contactTelegram"),
		ContactEmail:    c.FormValue("contactEmail"),
		ContactPhone:    c.FormValue("contactPhone"),
	}

	if canShip := c.FormValue("canShip"); canShip == "true" || canShip == "1" {
		req.CanShip = true
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Название обязательно"})
	}

	var imageContent []byte
	var imageFileName string
	var imageContentType string

	fileHeader, err := c.FormFile("image")
	if err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось открыть файл"})
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось прочитать файл"})
		}

		ct := fileHeader.Header.Get("Content-Type")
		if ct == "" {
			if guessed := mime.TypeByExtension(filepath.Ext(fileHeader.Filename)); guessed != "" {
				ct = guessed
			} else {
				ct = "application/octet-stream"
			}
		}

		imageContent = data
		imageFileName = fileHeader.Filename
		imageContentType = ct
	}

	item, err := h.svc.Create(req, member.Id, imageFileName, imageContent, imageContentType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать объявление"})
	}

	go h.pointSvc.GiveForAction(member.Id, models.PointReasonMarketplaceCreate, "marketplace", item.Id,
		fmt.Sprintf("Публикация объявления: %s", item.Title))

	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *MarketplaceHandler) Update(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	isAdmin := false
	for _, role := range member.Roles {
		if role == models.MemberRoleAdmin {
			isAdmin = true
			break
		}
	}

	var req models.CreateMarketplaceItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	item, err := h.svc.Update(id, &req, member.Id, isAdmin)
	if err != nil {
		log.Printf("Marketplace update error (item=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(item)
}

func (h *MarketplaceHandler) RequestPurchase(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	item, err := h.svc.RequestPurchase(id, member.Id)
	if err != nil {
		log.Printf("Marketplace requestPurchase error (item=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	go func() {
		h.pointSvc.GiveForAction(member.Id, models.PointReasonMarketplaceBuy, "marketplace", item.Id,
			fmt.Sprintf("Заявка на покупку: %s", item.Title))

		if err := CreateNotification(item.SellerId, "marketplace", "Новая заявка на покупку",
			fmt.Sprintf("На ваше объявление «%s» поступила заявка на покупку", item.Title)); err != nil {
			log.Printf("Error creating notification: %v", err)
		}
	}()

	return c.JSON(item)
}

func (h *MarketplaceHandler) CancelPurchase(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	isAdmin := false
	for _, role := range member.Roles {
		if role == models.MemberRoleAdmin {
			isAdmin = true
			break
		}
	}

	item, err := h.svc.CancelPurchase(id, member.Id, isAdmin)
	if err != nil {
		log.Printf("Marketplace cancelPurchase error (item=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(item)
}

func (h *MarketplaceHandler) MarkSold(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	isAdmin := false
	for _, role := range member.Roles {
		if role == models.MemberRoleAdmin {
			isAdmin = true
			break
		}
	}

	item, err := h.svc.MarkSold(id, member.Id, isAdmin)
	if err != nil {
		log.Printf("Marketplace markSold error (item=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(item)
}

func (h *MarketplaceHandler) Delete(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	isAdmin := false
	for _, role := range member.Roles {
		if role == models.MemberRoleAdmin {
			isAdmin = true
			break
		}
	}

	if err := h.svc.Delete(id, member.Id, isAdmin); err != nil {
		log.Printf("Marketplace delete error (item=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
