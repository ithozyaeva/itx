package handler

import (
	"crypto/subtle"
	"encoding/base64"
	"log"
	"strconv"

	"ithozyeva/config"
	"ithozyeva/internal/bot"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type TelegramAuthHandler struct {
	telegramService *service.TelegramService
	authService     *service.AuthTokenService
	memberService   *service.MemberService
}

func NewTelegramAuthHandler() (*TelegramAuthHandler, error) {
	tgService, err := service.NewTelegramService()

	if err != nil {
		return nil, err
	}

	return &TelegramAuthHandler{
		telegramService: tgService,
		authService:     service.NewAuthTokenService(),
		memberService:   service.NewMemberService(),
	}, nil
}

type AuthRequest struct {
	Token string `json:"token"`
}

func (h *TelegramAuthHandler) Authenticate(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем, существует ли токен
	existingToken, existingUser, err := h.authService.GetByToken(req.Token)
	if err != nil {
		// Если токена не существует, создаем нового
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	go func(user *models.Member) {
		isSubscriber, err := bot.CheckUserInChat(user.TelegramID)
		if err != nil {
			log.Printf("failed to check user %d in chat (skipping role update): %v", user.TelegramID, err)
			return
		}
		if isSubscriber && utils.HasRole(user.Roles, models.MemberRoleUnsubscriber) {
			user.Roles = []models.Role{models.MemberRoleSubscriber}
			h.memberService.Update(user)
		}
		if !isSubscriber && utils.HasRole(user.Roles, models.MemberRoleSubscriber) {
			user.Roles = []models.Role{models.MemberRoleUnsubscriber}
			h.memberService.Update(user)
		}
	}(existingUser)

	// Добавляем заголовок
	c.Response().Header.Add("X-Telegram-User-Token", existingToken.Token)

	// Возвращаем информацию о пользователе и токен
	return c.JSON(fiber.Map{
		"user":  existingUser,
		"token": existingToken.Token,
	})
}

func (h *TelegramAuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	decodedToken, err := base64.StdEncoding.DecodeString(req.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	decodedString := string(decodedToken)

	tgId, err := strconv.ParseInt(decodedString, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	existedToken, err := h.authService.GetTokenByTelegramID(tgId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	existedToken, err = h.authService.CreateOrUpdateToken(tgId, existedToken.Token)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	user, err := h.authService.GetByTelegramID(tgId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	c.Response().Header.Add("X-Telegram-User-Token", existedToken.Token)

	return c.JSON(fiber.Map{
		"token": existedToken.Token,
		"user":  user,
	})
}

type HandleBotMessageReq struct {
	Token     string      `json:"token"`
	UserID    int64       `json:"user_id"`
	Username  string      `json:"username"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Role      models.Role `json:"role"`
	AvatarURL string      `json:"avatar_url"`
}

func (h *TelegramAuthHandler) HandleBotMessage(c *fiber.Ctx) error {
	// Проверяем shared secret для защиты от неавторизованных вызовов
	secret := c.Get("X-Bot-Secret")
	if config.CFG.BotSharedSecret == "" || subtle.ConstantTimeCompare([]byte(secret), []byte(config.CFG.BotSharedSecret)) != 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req HandleBotMessageReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем, существует ли пользователь
	existingUser, err := h.authService.GetByTelegramID(req.UserID)
	if err != nil {
		// Определяем роль на основе проверки подписки, а не из тела запроса
		role := models.MemberRoleUnsubscriber
		if isSubscriber, checkErr := bot.CheckUserInChat(req.UserID); checkErr == nil && isSubscriber {
			role = models.MemberRoleSubscriber
		}

		// Если пользователь не существует, создаем нового
		newUser := &models.Member{
			TelegramID: req.UserID,
			Username:   req.Username,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			AvatarURL:  req.AvatarURL,
			Roles:      []models.Role{role},
		}

		createdUser, err := h.authService.CreateNewMember(newUser, req.Token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
		existingUser = createdUser
	} else {
		existingUser.Username = req.Username
		existingUser.FirstName = req.FirstName
		existingUser.LastName = req.LastName
		if req.AvatarURL != "" {
			existingUser.AvatarURL = req.AvatarURL
		}
		h.memberService.Update(existingUser)

		_, err := h.authService.CreateOrUpdateToken(req.UserID, req.Token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get auth token",
			})
		}
	}

	return c.JSON(existingUser)
}
