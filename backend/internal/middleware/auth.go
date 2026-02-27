package middleware

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	userRepo   *repository.AuthTokenRepository
	memberRepo *repository.MemberRepository
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:   repository.NewAuthTokenRepository(),
		memberRepo: repository.NewMemberRepository(),
	}
}

func (m *AuthMiddleware) RequireTGAuth(c *fiber.Ctx) error {
	telegramToken := c.Get("X-Telegram-User-Token")

	if telegramToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	authToken, err := m.userRepo.GetByToken(telegramToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Telegram User ID",
		})
	}

	if utils.CheckExpirationDate(authToken.ExpiredAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	member, err := m.memberRepo.GetByTelegramID(authToken.TelegramID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Member not found",
		})
	}

	// Добавляем информацию о пользователе в контекст
	c.Locals("member", member)

	return c.Next()
}

// RequireAuth проверяет авторизацию через Telegram и наличие разрешения на доступ к админ-панели
func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	tgToken := c.Get("X-Telegram-User-Token")
	if tgToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	authToken, err := m.userRepo.GetByToken(tgToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Telegram token",
		})
	}

	if utils.CheckExpirationDate(authToken.ExpiredAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token expired",
		})
	}

	// Get member and check if they can view admin panel
	member, err := m.memberRepo.GetByTelegramID(authToken.TelegramID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Member not found",
		})
	}

	// Check if user has permission to view admin panel
	if !m.memberRepo.HasPermission(member.Id, models.PermissionCanViewAdminPanel) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied. Insufficient permissions.",
		})
	}

	// Add member info to context
	c.Locals("member", member)
	return c.Next()
}

func (m *AuthMiddleware) RequirePermission(permission models.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by RequireAuth)
		member, ok := c.Locals("member").(*models.Member)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Check if user has the required permission
		if !m.memberRepo.HasPermission(member.Id, permission) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Insufficient permissions.",
			})
		}

		return c.Next()
	}
}
