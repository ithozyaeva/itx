package handler

import (
	"log"
	"strconv"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type GuildHandler struct {
	svc *service.GuildService
}

func NewGuildHandler() *GuildHandler {
	return &GuildHandler{
		svc: service.NewGuildService(),
	}
}

func (h *GuildHandler) GetAll(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	items, err := h.svc.GetAll(member.Id)
	if err != nil {
		log.Printf("get all guilds error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки гильдий"})
	}
	return c.JSON(items)
}

func (h *GuildHandler) Create(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	req := new(models.CreateGuildRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	guild, err := h.svc.Create(member.Id, req)
	if err != nil {
		log.Printf("create guild error (member=%d): %v", member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось создать гильдию"})
	}

	BroadcastEvent("guilds")
	return c.Status(fiber.StatusCreated).JSON(guild)
}

func (h *GuildHandler) Join(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.svc.Join(id, member.Id); err != nil {
		log.Printf("join guild error (guild=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось вступить в гильдию"})
	}
	BroadcastEvent("guilds")
	return c.JSON(fiber.Map{"ok": true})
}

func (h *GuildHandler) Leave(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.svc.Leave(id, member.Id); err != nil {
		log.Printf("leave guild error (guild=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось покинуть гильдию"})
	}
	BroadcastEvent("guilds")
	return c.JSON(fiber.Map{"ok": true})
}

func (h *GuildHandler) Update(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	req := new(models.CreateGuildRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	if err := h.svc.Update(id, member.Id, req); err != nil {
		log.Printf("update guild error (guild=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось обновить гильдию"})
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (h *GuildHandler) Delete(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	if err := h.svc.Delete(id, member.Id); err != nil {
		log.Printf("delete guild error (guild=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось удалить гильдию"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *GuildHandler) GetMembers(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	members, err := h.svc.GetMembers(id)
	if err != nil {
		log.Printf("get guild members error (guild=%d): %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки участников гильдии"})
	}
	return c.JSON(members)
}
