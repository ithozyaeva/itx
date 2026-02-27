package service

import (
	"fmt"
	"ithozyeva/config"
	"ithozyeva/internal/utils"
)

type TelegramService struct {
	botToken string
}

func NewTelegramService() (*TelegramService, error) {
	botToken := config.CFG.TelegramToken
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	return &TelegramService{
		botToken: botToken,
	}, nil
}

// GenerateAuthToken генерирует токен для авторизации
func (s *TelegramService) GenerateAuthToken(userID int64) string {
	token := utils.HashToken(fmt.Sprintf("%d", userID))

	return token
}
