package handler

import (
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "неверный запрос"})
	}

	token, err := h.svc.Login(req.Login, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(TokenResponse{Token: token})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(RefreshTokenRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "неверный запрос"})
	}

	newToken, err := h.svc.RefreshToken(req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(TokenResponse{Token: newToken})
}
