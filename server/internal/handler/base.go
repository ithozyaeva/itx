package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Identifiable интерфейс для сущностей с ID
type Identifiable interface {
	SetID(id int64)
}

// BaseHandler базовый обработчик для CRUD операций
type BaseHandler[T any] struct {
	service service.BaseService[T]
}

// NewBaseHandler создает новый экземпляр базового обработчика
func NewBaseHandler[T any](service service.BaseService[T]) *BaseHandler[T] {
	return &BaseHandler[T]{
		service: service,
	}
}

// Search выполняет поиск сущностей с пагинацией
func (h *BaseHandler[T]) Search(c *fiber.Ctx) error {
	req := new(models.SearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	// Передаем указатели в сервис
	result, err := h.service.Search(req.Limit, req.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// Create создает новую сущность
func (h *BaseHandler[T]) Create(c *fiber.Ctx) error {
	entity := new(T)
	if err := c.BodyParser(entity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.service.Create(entity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Update обновляет существующую сущность
func (h *BaseHandler[T]) Update(c *fiber.Ctx) error {
	entity := new(T)
	if err := c.BodyParser(entity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.service.Update(entity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// Delete удаляет сущность по ID
func (h *BaseHandler[T]) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity := new(T)

	// Пробуем использовать интерфейс Identifiable
	if identifiable, ok := any(entity).(Identifiable); ok {
		identifiable.SetID(int64(id))
	} else {
		// Используем рефлексию для установки ID
		val := reflect.ValueOf(entity).Elem()
		idField := val.FieldByName("Id")
		if idField.IsValid() && idField.CanSet() {
			idField.SetInt(int64(id))
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Невозможно установить ID для сущности"})
		}
	}

	if err := h.service.Delete(entity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetById получает сущность по ID
func (h *BaseHandler[T]) GetById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity, err := h.service.GetById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Сущность не найдена"})
	}

	return c.JSON(entity)
}
