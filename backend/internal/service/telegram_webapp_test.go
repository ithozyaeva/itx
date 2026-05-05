package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"
)

const testBotToken = "1234567890:AAH-fakeBotTokenForUnitTestsOnly_xxxxxxxx"

// signInitData собирает корректную initData-строку и подписывает её тем же
// алгоритмом, что использует Telegram-клиент. Используется в тестах вместо
// статичной фикстуры, чтобы любую корректную пару (data, hash) проверить
// против ValidateInitData.
func signInitData(t *testing.T, params map[string]string, botToken string) string {
	t.Helper()

	keys := make([]string, 0, len(params))
	for k := range params {
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
		sb.WriteString(params[k])
	}

	secretMac := hmac.New(sha256.New, []byte("WebAppData"))
	secretMac.Write([]byte(botToken))
	secretKey := secretMac.Sum(nil)

	dataMac := hmac.New(sha256.New, secretKey)
	dataMac.Write([]byte(sb.String()))
	hash := hex.EncodeToString(dataMac.Sum(nil))

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	values.Set("hash", hash)
	return values.Encode()
}

func TestValidateInitData_Valid(t *testing.T) {
	params := map[string]string{
		"auth_date": fmt.Sprintf("%d", time.Now().Unix()),
		"query_id":  "AAHdF6IQAAAAAN0XohDhrOrc",
		"user":      `{"id":42,"first_name":"Иван","last_name":"Петров","username":"ivan","language_code":"ru"}`,
	}
	initData := signInitData(t, params, testBotToken)

	user, err := ValidateInitData(initData, testBotToken, WebAppInitDataMaxAge)
	if err != nil {
		t.Fatalf("expected valid initData, got error: %v", err)
	}
	if user.ID != 42 {
		t.Errorf("user.ID: want 42, got %d", user.ID)
	}
	if user.Username != "ivan" {
		t.Errorf("user.Username: want ivan, got %q", user.Username)
	}
	if user.FirstName != "Иван" {
		t.Errorf("user.FirstName: want Иван, got %q", user.FirstName)
	}
}

func TestValidateInitData_TamperedHash(t *testing.T) {
	params := map[string]string{
		"auth_date": fmt.Sprintf("%d", time.Now().Unix()),
		"user":      `{"id":42,"first_name":"X"}`,
	}
	initData := signInitData(t, params, testBotToken)

	// Подменяем hash на чужой.
	values, _ := url.ParseQuery(initData)
	values.Set("hash", strings.Repeat("0", 64))
	tampered := values.Encode()

	if _, err := ValidateInitData(tampered, testBotToken, WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on tampered hash, got nil")
	}
}

func TestValidateInitData_TamperedPayload(t *testing.T) {
	params := map[string]string{
		"auth_date": fmt.Sprintf("%d", time.Now().Unix()),
		"user":      `{"id":42,"first_name":"X"}`,
	}
	initData := signInitData(t, params, testBotToken)

	// Подменяем user, не пересчитывая hash → подпись становится невалидной.
	values, _ := url.ParseQuery(initData)
	values.Set("user", `{"id":999,"first_name":"Hacker"}`)
	tampered := values.Encode()

	if _, err := ValidateInitData(tampered, testBotToken, WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on tampered payload, got nil")
	}
}

func TestValidateInitData_ExpiredAuthDate(t *testing.T) {
	expired := time.Now().Add(-48 * time.Hour).Unix()
	params := map[string]string{
		"auth_date": fmt.Sprintf("%d", expired),
		"user":      `{"id":42,"first_name":"X"}`,
	}
	initData := signInitData(t, params, testBotToken)

	if _, err := ValidateInitData(initData, testBotToken, WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on expired auth_date, got nil")
	}
}

func TestValidateInitData_WrongBotToken(t *testing.T) {
	params := map[string]string{
		"auth_date": fmt.Sprintf("%d", time.Now().Unix()),
		"user":      `{"id":42,"first_name":"X"}`,
	}
	initData := signInitData(t, params, testBotToken)

	if _, err := ValidateInitData(initData, "different-bot-token", WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on wrong bot token, got nil")
	}
}

func TestValidateInitData_MissingHash(t *testing.T) {
	if _, err := ValidateInitData("auth_date=1&user=%7B%7D", testBotToken, WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on missing hash, got nil")
	}
}

func TestValidateInitData_MissingUser(t *testing.T) {
	params := map[string]string{
		"auth_date": fmt.Sprintf("%d", time.Now().Unix()),
	}
	initData := signInitData(t, params, testBotToken)

	if _, err := ValidateInitData(initData, testBotToken, WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on missing user, got nil")
	}
}

func TestValidateInitData_EmptyInputs(t *testing.T) {
	if _, err := ValidateInitData("", testBotToken, WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on empty initData, got nil")
	}
	if _, err := ValidateInitData("hash=x", "", WebAppInitDataMaxAge); err == nil {
		t.Fatal("expected error on empty bot token, got nil")
	}
}
