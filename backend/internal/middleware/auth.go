package middleware

import (
	"ithozyeva/config"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func (m *AuthMiddleware) RequireJWTAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.CFG.JwtSecret), nil
	})

	if err == nil && token.Valid {
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}

func (m *AuthMiddleware) RequireTGAuth(c *fiber.Ctx) error {
	// Если JWT токен не валиден или отсутствует, проверяем Telegram токен
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

func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	// Try JWT first
	jwtToken := c.Get("Authorization")
	if jwtToken != "" {
		tokenStr := strings.TrimPrefix(jwtToken, "Bearer ")

		if tokenStr != "" {
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(config.CFG.JwtSecret), nil
			})

			if err == nil && token.Valid {
				// JWT is valid, continue to next handler
				return c.Next()
			}
		}
	}

	// Try Telegram token if JWT fails
	tgToken := c.Get("X-Telegram-User-Token")
	if tgToken != "" {
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

	// Return unauthorized if both fail
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}

func (m *AuthMiddleware) RequirePermission(permission models.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by RequireAuth)
		member, ok := c.Locals("member").(*models.Member)
		if !ok {
			// For JWT auth, we don't have member in context, so we need to check if it's JWT auth
			jwtToken := c.Get("Authorization")
			if jwtToken != "" {
				// For now, allow JWT auth to pass through permission checks
				// In a real implementation, you would want to fetch user permissions for JWT users too
				return c.Next()
			}

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
