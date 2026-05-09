package handler

import (
	"crypto/subtle"
	"log"

	"ithozyeva/config"
	"ithozyeva/internal/bot"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type TelegramAuthHandler struct {
	telegramService *service.TelegramService
	authService     *service.AuthTokenService
	memberService   *service.MemberService
	pendingReferral *service.PendingReferralService
}

func NewTelegramAuthHandler(redisClient *redis.Client) (*TelegramAuthHandler, error) {
	tgService, err := service.NewTelegramService()

	if err != nil {
		return nil, err
	}

	return &TelegramAuthHandler{
		telegramService: tgService,
		authService:     service.NewAuthTokenService(),
		memberService:   service.NewMemberService(),
		pendingReferral: service.NewPendingReferralService(redisClient),
	}, nil
}

type AuthRequest struct {
	Token string `json:"token"`
}

// mergeSubscriptionRole возвращает копию roles, в которой выставлен ровно один
// из флагов SUBSCRIBER/UNSUBSCRIBER в соответствии с isSubscriber, а все
// остальные роли (ADMIN, MENTOR, EVENT_MAKER, …) сохранены.
//
// changed=true означает «нужно сохранить в БД» — false когда состояние уже
// корректное (один правильный флаг, без лишних), чтобы лишний раз не
// дёргать MemberRepository.Update.
//
// Извлечено в чистую функцию специально под юнит-тест: см. auth_token_test.go.
// Регрессия #340 чинила побочный эффект periodic; этот хелпер закрывает
// исходный root cause — overwrite роль-слайса при флипе chat-membership
// затирал ADMIN/MENTOR/EVENT_MAKER при каждом /authenticate.
func mergeSubscriptionRole(roles []models.Role, isSubscriber bool) (newRoles []models.Role, changed bool) {
	desired := models.MemberRoleUnsubscriber
	if isSubscriber {
		desired = models.MemberRoleSubscriber
	}

	hasDesired := false
	newRoles = make([]models.Role, 0, len(roles)+1)
	for _, r := range roles {
		switch r {
		case models.MemberRoleSubscriber, models.MemberRoleUnsubscriber:
			if r == desired && !hasDesired {
				hasDesired = true
				newRoles = append(newRoles, r)
			} else {
				// либо противоположный флаг, либо дубликат desired — оба отбрасываем.
				changed = true
			}
		default:
			newRoles = append(newRoles, r)
		}
	}
	if !hasDesired {
		newRoles = append(newRoles, desired)
		changed = true
	}
	return newRoles, changed
}

type WebAppAuthRequest struct {
	InitData string `json:"init_data"`
}

// AuthenticateWebApp — авторизация из Telegram Mini App. На вход приходит
// сырая строка window.Telegram.WebApp.initData; HMAC-подпись валидируется
// бот-токеном. После успеха выпускается тот же tg_token, что и в обычном
// потоке через /telegram, поэтому фронт может пользоваться единым
// authService и не различать каналы. CheckUserInChat намеренно НЕ зовётся:
// основной API-сервер живёт в РФ, api.telegram.org заблокирован — роль
// доопределит периодический subscriptionChecker на боте (NL).
func (h *TelegramAuthHandler) AuthenticateWebApp(c *fiber.Ctx) error {
	var req WebAppAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	tgUser, err := service.ValidateInitData(req.InitData, config.CFG.TelegramToken, service.WebAppInitDataMaxAge)
	if err != nil {
		log.Printf("webapp auth: invalid init data: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid init data",
		})
	}
	tgUser.Username = utils.SanitizeTelegramUsername(tgUser.Username)

	token, err := h.telegramService.GenerateAuthToken(tgUser.ID)
	if err != nil {
		log.Printf("webapp auth: failed to generate token for tg_id=%d: %v", tgUser.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to issue token",
		})
	}

	user, err := h.authService.GetByTelegramID(tgUser.ID)
	if err != nil {
		if tgUser.Username != "" {
			if claimErr := h.memberService.ClaimUsername(tgUser.Username, 0); claimErr != nil {
				log.Printf("webapp auth: claim username on create failed (tg_id=%d, username=%s): %v", tgUser.ID, tgUser.Username, claimErr)
			}
		}
		newUser := &models.Member{
			TelegramID: tgUser.ID,
			Username:   tgUser.Username,
			FirstName:  tgUser.FirstName,
			LastName:   tgUser.LastName,
			Roles:      []models.Role{models.MemberRoleUnsubscriber},
		}
		created, createErr := h.authService.CreateNewMember(newUser, token)
		if createErr != nil {
			log.Printf("webapp auth: failed to create member tg_id=%d: %v", tgUser.ID, createErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
		user = created
	} else {
		user.FirstName = tgUser.FirstName
		user.LastName = tgUser.LastName
		if user.Username == "" && tgUser.Username != "" {
			if claimErr := h.memberService.ClaimUsername(tgUser.Username, user.Id); claimErr != nil {
				log.Printf("webapp auth: claim username on re-login failed (tg_id=%d, username=%s): %v", tgUser.ID, tgUser.Username, claimErr)
			}
			user.Username = tgUser.Username
		}
		h.memberService.Update(user)

		if _, err := h.authService.CreateOrUpdateToken(tgUser.ID, token); err != nil {
			log.Printf("webapp auth: failed to upsert token for tg_id=%d: %v", tgUser.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to issue token",
			})
		}
	}

	// Боб впервые открывает платформу — если у него в Redis-pending есть
	// referrer (бот /start ref_<code>), фиксируем атрибуцию + community-award.
	// Идемпотентно: повторный auth no-op'ится через ранний return на флаге.
	h.memberService.ApplyPendingReferral(c.Context(), user, h.pendingReferral)

	user.SubscriptionTier = h.memberService.GetEffectiveTier(user.TelegramID)

	c.Response().Header.Add("X-Telegram-User-Token", token)
	return c.JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
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
		newRoles, changed := mergeSubscriptionRole(user.Roles, isSubscriber)
		if !changed {
			return
		}
		user.Roles = newRoles
		h.memberService.Update(user)
	}(existingUser)

	// Добавляем заголовок
	c.Response().Header.Add("X-Telegram-User-Token", existingToken.Token)

	// Возвращаем информацию о пользователе и токен
	return c.JSON(fiber.Map{
		"user":  existingUser,
		"token": existingToken.Token,
	})
}

// Logout инвалидирует текущий tg_token серверно: записывает в auth_tokens
// expired_at в прошлое, чтобы middleware.RequireTGAuth на следующий запрос
// этим токеном вернул 401. До появления этого хендлера клик «Выйти» только
// удалял токен из localStorage, но в БД он жил до natural expiry (~30 дней),
// и кто бы ни получил доступ к токену (украденный localStorage, sniff,
// расшаренный комп) мог использовать API дальше. Берём токен ТОЛЬКО из
// заголовка X-Telegram-User-Token — telegram-id юзера тут не нужен.
//
// Идемпотентно: если токен уже отсутствует в БД, отвечаем 204 — клиент
// в любом случае получит «logged out».
func (h *TelegramAuthHandler) Logout(c *fiber.Ctx) error {
	headerToken := c.Get("X-Telegram-User-Token")
	if headerToken == "" {
		return c.SendStatus(fiber.StatusNoContent)
	}
	if err := h.authService.InvalidateToken(headerToken); err != nil {
		log.Printf("logout: failed to invalidate token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// RefreshToken продлевает срок действия текущего токена. Принимает токен
// ТОЛЬКО из заголовка X-Telegram-User-Token: знание Telegram-ID не должно
// давать доступ к чужой сессии.
func (h *TelegramAuthHandler) RefreshToken(c *fiber.Ctx) error {
	headerToken := c.Get("X-Telegram-User-Token")
	if headerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	authToken, user, err := h.authService.GetByToken(headerToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Истёкший токен нельзя продлевать — иначе массовая инвалидация (через
	// сдвиг expired_at в прошлое) не вышибет уже скомпрометированные сессии.
	if utils.CheckExpirationDate(authToken.ExpiredAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token expired",
		})
	}

	if _, err := h.authService.CreateOrUpdateToken(authToken.TelegramID, authToken.Token); err != nil {
		log.Printf("refresh: failed to extend token for tg_id=%d: %v", authToken.TelegramID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	// Обогащаем user тем же набором, что отдаёт /me: SubscriptionTier через
	// GetEffectiveTier и mentor-поля при наличии. Иначе фронт после
	// proactive-refresh затирает subscriptionTier и схлопывает сайдбар до
	// следующего вызова /me.
	user.SubscriptionTier = h.memberService.GetEffectiveTier(user.TelegramID)

	c.Response().Header.Add("X-Telegram-User-Token", authToken.Token)

	if mentor, err := h.memberService.GetMentor(user.Id); err == nil {
		mentor.SubscriptionTier = user.SubscriptionTier
		return c.JSON(fiber.Map{
			"token": authToken.Token,
			"user":  mentor,
		})
	}

	return c.JSON(fiber.Map{
		"token": authToken.Token,
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
	// Telegram гарантирует валидный username, но защитимся от неожиданно длинных
	// или содержащих посторонние символы значений (прокси, фейковый бот, баг).
	req.Username = utils.SanitizeTelegramUsername(req.Username)

	// Проверяем, существует ли пользователь
	existingUser, err := h.authService.GetByTelegramID(req.UserID)
	if err != nil {
		// Определяем роль на основе проверки подписки, а не из тела запроса
		role := models.MemberRoleUnsubscriber
		if isSubscriber, checkErr := bot.CheckUserInChat(req.UserID); checkErr == nil && isSubscriber {
			role = models.MemberRoleSubscriber
		}

		// Если пользователь не существует, создаем нового.
		// Перед записью освобождаем username у чужих записей: тот же ник
		// мог «висеть» у заброшенного аккаунта, и UNIQUE-индекс заблокирует
		// создание. Источник истины — Telegram, поэтому приоритет за тем,
		// кто логинится сейчас.
		if req.Username != "" {
			if err := h.memberService.ClaimUsername(req.Username, 0); err != nil {
				log.Printf("claim username on create failed (tg_id=%d, username=%s): %v", req.UserID, req.Username, err)
			}
		}

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
		existingUser.FirstName = req.FirstName
		existingUser.LastName = req.LastName
		// Username не перезатираем при повторном логине, но если в БД пусто
		// (старая запись или username сбросили dedupe-миграцией) — забираем
		// его себе и попутно освобождаем у чужих записей.
		if existingUser.Username == "" && req.Username != "" {
			if err := h.memberService.ClaimUsername(req.Username, existingUser.Id); err != nil {
				log.Printf("claim username on re-login failed (tg_id=%d, username=%s): %v", req.UserID, req.Username, err)
			}
			existingUser.Username = req.Username
		}
		// AvatarURL — кастомный аватар (см. C2 в этом PR) не должен сноситься.
		if existingUser.AvatarURL == "" && req.AvatarURL != "" {
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

	// Реф-атрибуция: Боб впервые попал в систему через /start ref_<code>
	// в боте — переносим pending в БД и дёргаем community-award инвайтеру.
	h.memberService.ApplyPendingReferral(c.Context(), existingUser, h.pendingReferral)

	return c.JSON(existingUser)
}
