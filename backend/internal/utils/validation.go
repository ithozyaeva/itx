package utils

import (
	"errors"
	"regexp"
	"time"
	"unicode/utf8"
)

const (
	MaxFirstNameLen = 64
	MaxLastNameLen  = 64
	MaxBioLen       = 500
	MaxGradeLen     = 100
	MaxCompanyLen   = 100
	MaxUsernameLen  = 32
	MinUsernameLen  = 1
)

// MinBirthday — самый ранний разрешённый день рождения. Пользователю не
// 100+ лет на момент регистрации.
var MinBirthday = time.Date(1920, 1, 1, 0, 0, 0, 0, time.UTC)

var usernameRegex = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// SanitizeTelegramUsername нормализует @username, пришедший из Telegram OAuth.
// Возвращает пустую строку, если значение явно поломано (длина или символы):
// в БД лучше пустота, чем мусор, который потом не пройдёт UNIQUE.
func SanitizeTelegramUsername(username string) string {
	if username == "" {
		return ""
	}
	n := utf8.RuneCountInString(username)
	if n < MinUsernameLen || n > MaxUsernameLen {
		return ""
	}
	if !usernameRegex.MatchString(username) {
		return ""
	}
	return username
}

// ValidateUsername проверяет admin-ввод username и возвращает понятное
// сообщение для 400/409 ответа.
func ValidateUsername(username string) error {
	if username == "" {
		return nil
	}
	n := utf8.RuneCountInString(username)
	if n < MinUsernameLen || n > MaxUsernameLen {
		return errors.New("Никнейм должен быть от 1 до 32 символов")
	}
	if !usernameRegex.MatchString(username) {
		return errors.New("Никнейм может содержать только латиницу, цифры и _")
	}
	return nil
}

// ValidateProfileLengths проверяет длины полей профиля. Передаются строки
// «как пришли с фронта». Считаем по rune-ам, чтобы кириллица не давала
// фальшивый перебор лимита.
func ValidateProfileLengths(firstName, lastName, bio, grade, company string) error {
	if utf8.RuneCountInString(firstName) > MaxFirstNameLen {
		return errors.New("Имя не должно превышать 64 символа")
	}
	if utf8.RuneCountInString(lastName) > MaxLastNameLen {
		return errors.New("Фамилия не должна превышать 64 символа")
	}
	if utf8.RuneCountInString(bio) > MaxBioLen {
		return errors.New("Био не должно превышать 500 символов")
	}
	if utf8.RuneCountInString(grade) > MaxGradeLen {
		return errors.New("Грейд не должен превышать 100 символов")
	}
	if utf8.RuneCountInString(company) > MaxCompanyLen {
		return errors.New("Компания не должна превышать 100 символов")
	}
	return nil
}

// ValidateBirthdayRange — день рождения должен попадать в реальный диапазон.
// nil-значение валидно (поле необязательное). «Сегодня» считаем по МСК:
// бизнес живёт в Москве, и юзер из России в полночь по своему времени не
// должен получать «дата в будущем» из-за того что в UTC ещё вчера.
func ValidateBirthdayRange(birthday *time.Time) error {
	if birthday == nil {
		return nil
	}
	if birthday.Before(MinBirthday) {
		return errors.New("Дата рождения слишком давняя")
	}
	todayMSK := time.Now().In(MSKLocation())
	endOfToday := time.Date(todayMSK.Year(), todayMSK.Month(), todayMSK.Day(), 23, 59, 59, 0, todayMSK.Location())
	if birthday.After(endOfToday) {
		return errors.New("Дата рождения не может быть в будущем")
	}
	return nil
}
