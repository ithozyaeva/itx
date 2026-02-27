package service

import (
	"errors"
	"ithozyeva/config"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	Login(login string, password string) (string, error)
	RefreshToken(token string) (string, error)
}

type userService struct {
	rp repository.UserRepository
}

func NewUserService() UserService {
	return &userService{
		rp: repository.NewUserRepository(),
	}
}

func (u *userService) Login(login string, password string) (string, error) {
	// Поиск пользователя в БД
	user, err := u.rp.GetUserByLogin(login)

	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	// Проверка пароля
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("неверный пароль")
	}

	// Генерация JWT
	token := utils.GenerateJWT(login)

	tokenString, err := token.SignedString(config.CFG.JwtSecret)

	if err != nil {
		return "", errors.New("не удалось сгенерировать токен")
	}

	return tokenString, nil
}

func (u *userService) RefreshToken(tokenString string) (string, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return config.CFG.JwtSecret, nil
	})

	if err != nil {
		return "", errors.New("недействительный токен")
	}

	// Проверяем валидность токена
	if !token.Valid {
		return "", errors.New("недействительный токен")
	}

	// Получаем claims из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("недействительные данные токена")
	}

	// Проверяем, что токен не истек
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("недействительная дата истечения токена")
	}

	// Если токен истек более чем 30 дней назад, не обновляем его
	if time.Now().Unix() > int64(exp)+30*24*60*60 {
		return "", errors.New("токен слишком старый для обновления")
	}

	// Получаем логин пользователя из токена
	login, ok := claims["login"].(string)
	if !ok {
		return "", errors.New("недействительный логин в токене")
	}

	// Проверяем существование пользователя
	_, err = u.rp.GetUserByLogin(login)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	// Генерируем новый токен
	newToken := utils.GenerateJWT(login)
	newTokenString, err := newToken.SignedString(config.CFG.JwtSecret)
	if err != nil {
		return "", errors.New("не удалось сгенерировать новый токен")
	}

	return newTokenString, nil
}
