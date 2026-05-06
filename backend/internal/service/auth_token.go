package service

import (
	"errors"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
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

func (s *AuthTokenService) CreateOrUpdateToken(telegramID int64, token string) (*models.AuthToken, error) {
	authToken, err := s.authRepo.GetByTelegramID(telegramID)
	if err != nil {
		s.authRepo.Create(&models.AuthToken{TelegramID: telegramID, Token: token, ExpiredAt: time.Now().AddDate(0, 1, 0)})
	} else {
		s.authRepo.Update(&models.AuthToken{ID: authToken.ID, TelegramID: telegramID, Token: token, ExpiredAt: time.Now().AddDate(0, 1, 0)})
	}

	return authToken, nil
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
