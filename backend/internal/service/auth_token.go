package service

import (
	"errors"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"log"
	"time"

	"gorm.io/gorm"
)

type AuthTokenService struct {
	authRepo *repository.AuthTokenRepository
	userRepo *repository.MemberRepository
}

func NewAuthTokenService() *AuthTokenService {
	return &AuthTokenService{
		authRepo: repository.NewAuthTokenRepository(),
		userRepo: repository.NewMemberRepository(),
	}
}

func (s *AuthTokenService) GetByToken(token string) (*models.AuthToken, *models.Member, error) {
	authToken, err := s.authRepo.GetByToken(token)

	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.GetByTelegramID(authToken.TelegramID)

	if err != nil {
		return nil, nil, err
	}

	return authToken, user, nil
}

func (s *AuthTokenService) GetByTelegramID(telegramID int64) (*models.Member, error) {
	return s.userRepo.GetByTelegramID(telegramID)
}

func (s *AuthTokenService) CreateNewMember(user *models.Member, token string) (*models.Member, error) {
	createdUser, err := s.userRepo.Create(user)

	if err != nil {
		return nil, err
	}

	// Сразу присваиваем referral_code новому юзеру — иначе не сможет
	// поделиться ссылкой пока не откроет /referral хотя бы раз.
	// AssignReferralCode идемпотентен (WHERE IS NULL), так что повторные
	// вызовы безопасны.
	if code, _, err := s.userRepo.AssignReferralCode(createdUser.Id); err == nil {
		createdUser.ReferralCode = &code
	} else {
		log.Printf("CreateNewMember: AssignReferralCode failed for %d: %v", createdUser.Id, err)
	}

	_, err = s.authRepo.Create(&models.AuthToken{
		TelegramID: createdUser.TelegramID,
		ExpiredAt:  time.Now().AddDate(0, 1, 0),
		Token:      token,
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// CreateOrUpdateToken создаёт новую запись auth_tokens либо продляет
// существующую для данного telegramID. Возвращает результирующий токен
// и любую ошибку Create/Update — раньше эти ошибки молча терялись
// (возврат был always (authToken, nil) с authToken=pre-write значение),
// и хендлеры (AuthenticateWebApp/HandleBotMessage/RefreshToken) на DB-сбое
// отвечали юзеру 200 со свежим token-стрингом, который в БД не был
// сохранён → следующий запрос middleware-401 → пользователь застревал
// между «логин ок» и «токен невалиден».
func (s *AuthTokenService) CreateOrUpdateToken(telegramID int64, token string) (*models.AuthToken, error) {
	authToken, err := s.authRepo.GetByTelegramID(telegramID)
	if err != nil {
		created, cerr := s.authRepo.Create(&models.AuthToken{TelegramID: telegramID, Token: token, ExpiredAt: time.Now().AddDate(0, 1, 0)})
		if cerr != nil {
			return nil, cerr
		}
		return created, nil
	}
	updated, uerr := s.authRepo.Update(&models.AuthToken{ID: authToken.ID, TelegramID: telegramID, Token: token, ExpiredAt: time.Now().AddDate(0, 1, 0)})
	if uerr != nil {
		return nil, uerr
	}
	return updated, nil
}

func (s *AuthTokenService) GetTokenByTelegramID(telegramID int64) (*models.AuthToken, error) {
	return s.authRepo.GetByTelegramID(telegramID)
}

// InvalidateToken помечает токен как истёкший. Идемпотентно по «токена нет
// в БД» — на ErrRecordNotFound возвращаем nil (клиент в любом случае получает
// «logged out», 404 на logout не нужен). Все остальные ошибки (соединение,
// SQL) пробрасываем — иначе DB-сбой превращается в тихий «полу-logout»:
// хендлер вернёт 204, юзер думает что вышел, а токен в БД жив.
//
// Используется хендлером /api/auth/telegram/logout, чтобы клик «Выйти»
// действительно инвалидировал сессию серверно. До этого токен жил в БД
// до своего natural expiry (~30 дней), и тот, у кого он есть (украденный
// localStorage, sniff, расшаренный комп), мог пользоваться API дальше.
func (s *AuthTokenService) InvalidateToken(token string) error {
	authToken, err := s.authRepo.GetByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	authToken.ExpiredAt = time.Now().Add(-time.Hour)
	_, err = s.authRepo.Update(authToken)
	return err
}
