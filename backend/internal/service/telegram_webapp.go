package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// WebAppInitDataMaxAge — после такого возраста initData считается просроченной.
// Telegram-клиент пересоздаёт initData при каждом запуске miniapp; 24 часа —
// safety net на случай скриншот-атак с украденной строкой.
const WebAppInitDataMaxAge = 24 * time.Hour

// WebAppUser — поля Telegram WebApp.initDataUnsafe.user.
// Описание: https://core.telegram.org/bots/webapps#webappuser
type WebAppUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	PhotoURL     string `json:"photo_url"`
	LanguageCode string `json:"language_code"`
}

// ValidateInitData проверяет HMAC-подпись и свежесть initData по протоколу
// Telegram WebApp (https://core.telegram.org/bots/webapps#validating-data-received-via-the-mini-app).
//
// data_check_string = "\n".join(sort_by_key(k=v for k≠"hash"))
// secret_key        = HMAC_SHA256("WebAppData", bot_token)
// expected_hash     = HMAC_SHA256(secret_key, data_check_string)
func ValidateInitData(initData, botToken string, maxAge time.Duration) (*WebAppUser, error) {
	if initData == "" {
		return nil, fmt.Errorf("init data is empty")
	}
	if botToken == "" {
		return nil, fmt.Errorf("bot token is empty")
	}

	parsed, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("parse init data: %w", err)
	}

	receivedHash := parsed.Get("hash")
	if receivedHash == "" {
		return nil, fmt.Errorf("hash missing")
	}
	parsed.Del("hash")

	keys := make([]string, 0, len(parsed))
	for k := range parsed {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(parsed.Get(k))
	}

	secretMac := hmac.New(sha256.New, []byte("WebAppData"))
	secretMac.Write([]byte(botToken))
	secretKey := secretMac.Sum(nil)

	dataMac := hmac.New(sha256.New, secretKey)
	dataMac.Write([]byte(sb.String()))
	expectedHash := hex.EncodeToString(dataMac.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(receivedHash)) != 1 {
		return nil, fmt.Errorf("invalid hash")
	}

	authDateStr := parsed.Get("auth_date")
	if authDateStr == "" {
		return nil, fmt.Errorf("auth_date missing")
	}
	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid auth_date: %w", err)
	}
	if maxAge > 0 && time.Since(time.Unix(authDate, 0)) > maxAge {
		return nil, fmt.Errorf("init data expired")
	}

	userJSON := parsed.Get("user")
	if userJSON == "" {
		return nil, fmt.Errorf("user data missing")
	}
	var user WebAppUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, fmt.Errorf("parse user: %w", err)
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user id missing")
	}

	return &user, nil
}
