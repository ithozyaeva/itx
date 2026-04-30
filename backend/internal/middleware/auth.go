package middleware

import (
	"ithozyeva/config"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	userRepo         *repository.AuthTokenRepository
	memberRepo       *repository.MemberRepository
	subscriptionRepo *repository.SubscriptionRepository
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:         repository.NewAuthTokenRepository(),
		memberRepo:       repository.NewMemberRepository(),
		subscriptionRepo: repository.NewSubscriptionRepository(),
	}
}

func (m *AuthMiddleware) RequireTGAuth(c *fiber.Ctx) error {
	telegramToken := c.Get("X-Telegram-User-Token")
	if telegramToken == "" {
		telegramToken = c.Query("token")
	}

	if telegramToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	authToken, err := m.userRepo.GetByToken(telegramToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
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

// RequireSuperAdmin пропускает только супер-админа (Telegram-id из
// config.CFG.SuperAdminTelegramID). Используется для действий, которые не должны
// быть доступны ни одному другому админу платформы — например, ручное создание
// или удаление чатов в системе подписок.
func (m *AuthMiddleware) RequireSuperAdmin(c *fiber.Ctx) error {
	member, ok := c.Locals("member").(*models.Member)
	if !ok || member == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	superAdminID := config.CFG.SuperAdminTelegramID
	if superAdminID == 0 || member.TelegramID != superAdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only the super admin can perform this action.",
		})
	}

	return c.Next()
}

// RequireSubscription гейтит платформенные эндпоинты под активную подписку.
// Источник правды — subscription_users.EffectiveTierID() (manual override
// либо tier, разрешённый по членству в anchor-чате).
//
// Включается переменной SUBSCRIPTION_GATE_ENABLED — пока выключен, middleware
// no-op. Это позволяет смержить и задеплоить код, разослать анонс через бота,
// и только затем включить флаг в проде.
func (m *AuthMiddleware) RequireSubscription(c *fiber.Ctx) error {
	if !config.CFG.SubscriptionGateEnabled {
		return c.Next()
	}

	member, ok := c.Locals("member").(*models.Member)
	if !ok || member == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	user, err := m.subscriptionRepo.GetUser(member.TelegramID)
	if err != nil || user == nil || user.EffectiveTierID() == nil {
		// Кидаем на главную: на дашборде есть мягкий subscription-teaser,
		// а лобовое падение на /tariffs из любого премиум-эндпоинта (особенно
		// при переходе из бота по /members/:id) даёт ощущение «бот меня
		// продаёт» и сбивает контекст.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":    "subscription_required",
			"redirect": "/",
		})
	}

	c.Locals("subscription_tier_id", user.EffectiveTierID())
	return c.Next()
}

// RequireMinTier гейтит эндпоинты по минимальному уровню подписки. Используется
// для премиум-фич, которые открыты только начиная с конкретного тира (например,
// «AI-материалы» — только master+, level >= 3).
//
// Как и RequireSubscription, отключается флагом SUBSCRIPTION_GATE_ENABLED —
// пока флаг выключен, проверка не выполняется (единая точка переключения для
// всей системы подписок).
func (m *AuthMiddleware) RequireMinTier(minLevel int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !config.CFG.SubscriptionGateEnabled {
			return c.Next()
		}

		member, ok := c.Locals("member").(*models.Member)
		if !ok || member == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// ADMIN — универсальный модератор; пускаем без проверки тира,
		// синхронно с frontend hasMinTier и handler-level admin bypass'ами
		// (SetHidden, UpdateComment и т.п.). Без этого админ без master+
		// получал бы 403 на собственных moderation-эндпоинтах.
		for _, role := range member.Roles {
			if role == models.MemberRoleAdmin {
				return c.Next()
			}
		}

		level, ok := m.subscriptionRepo.GetUserEffectiveTierLevel(member.TelegramID)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":    "subscription_required",
				"redirect": "/",
			})
		}
		if level < minLevel {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":         "tier_too_low",
				"required_tier": minLevel,
				"current_tier":  level,
			})
		}

		c.Locals("subscription_tier_level", level)
		return c.Next()
	}
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
