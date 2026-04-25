package handler

import (
	"ithozyeva/internal/models"

	"github.com/gofiber/fiber/v2"
)

// getMember безопасно достаёт *models.Member из c.Locals.
// Возвращает 401, если member отсутствует или некорректного типа —
// вместо паники из-за прямого type assertion.
func getMember(c *fiber.Ctx) (*models.Member, error) {
	member, ok := c.Locals("member").(*models.Member)
	if !ok || member == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "member context not found")
	}
	return member, nil
}
