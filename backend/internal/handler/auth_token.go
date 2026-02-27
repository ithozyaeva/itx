package handler

import (
	"encoding/base64"
	"ithozyeva/internal/bot"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TelegramAuthHandler struct {
	telegramService *service.TelegramService
	authService     *service.AuthTokenService
	memberService   *service.MemberService
}

func NewTelegramAuthHandler() *TelegramAuthHandler {
	tgService, err := service.NewTelegramService()

	if err != nil {
		panic(err)
	}

	return &TelegramAuthHandler{
		telegramService: tgService,
		authService:     service.NewAuthTokenService(),
		memberService:   service.NewMemberService(),
	}
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

	isSubcriber, err := bot.CheckUserInChat(existingUser.TelegramID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if isSubcriber && utils.HasRole(existingUser.Roles, models.MemberRoleUnsubscriber) {
		existingUser.Roles = []models.Role{models.MemberRoleSubscriber}
		existingUser, _ = h.memberService.Update(existingUser)
	}

	if !isSubcriber && utils.HasRole(existingUser.Roles, models.MemberRoleSubscriber) {
		existingUser.Roles = []models.Role{models.MemberRoleUnsubscriber}
		existingUser, _ = h.memberService.Update(existingUser)
	}

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
}

func (h *TelegramAuthHandler) HandleBotMessage(c *fiber.Ctx) error {
	var req HandleBotMessageReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем, существует ли пользователь
	existingUser, err := h.authService.GetByTelegramID(req.UserID)
	if err != nil {
		// Если пользователь не существует, создаем нового
		newUser := &models.Member{
			TelegramID: req.UserID,
			Username:   req.Username,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Roles:      []models.Role{req.Role},
		}

		createdUser, err := h.authService.CreateNewMember(newUser, req.Token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
		existingUser = createdUser
	} else {
		_, err := h.authService.CreateOrUpdateToken(req.UserID, req.Token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get auth token",
			})
		}
	}

	return c.JSON(existingUser)
}
