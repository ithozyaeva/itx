package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

// patchMeApp поднимает минимальный fiber-app с PATCH /members/me и
// фейковой авторизацией: указанный member кладётся в c.Locals.
// Нужен только для проверки enrichment ответа — реальную auth-цепочку
// здесь дублировать не имеет смысла.
func patchMeApp(member *models.Member) *fiber.App {
	app := fiber.New()
	h := NewMembersHandler()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("member", member)
		return c.Next()
	})
	app.Patch("/members/me", h.UpdateProfile)
	return app
}

func decodePatchResponse(t *testing.T, body io.Reader) map[string]any {
	t.Helper()
	raw, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal %s: %v", string(raw), err)
	}
	return out
}

// TestUpdateProfile_ResponseIncludesSubscriptionTier — регрессия:
// после редактирования профиля фронт должен видеть subscriptionTier,
// иначе уровень подписки слетает в стейте платформы. Поле помечено
// gorm:"-:all", поэтому каждый хендлер обязан подмешивать его руками.
func TestUpdateProfile_ResponseIncludesSubscriptionTier(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db,
		"member_roles", "members",
		"subscription_users",
	)

	const tgID int64 = 991001
	m := &models.Member{
		TelegramID: tgID,
		Username:   "patch_user",
		FirstName:  "Old",
		LastName:   "Name",
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}

	// foreman (id=2, level=2) seeded миграцией 20260319000000.
	const foremanTierID uint = 2
	tierID := foremanTierID
	if err := db.Create(&models.SubscriptionUser{
		ID:             tgID,
		FullName:       "Old Name",
		ResolvedTierID: &tierID,
		IsActive:       true,
	}).Error; err != nil {
		t.Fatalf("create subscription_user: %v", err)
	}

	body, _ := json.Marshal(map[string]any{
		"firstName": "New",
		"lastName":  "Name",
		"bio":       "",
		"grade":     "Middle",
		"company":   "",
		"tg":        "patch_user",
	})
	req, _ := http.NewRequest(http.MethodPatch, "/members/me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := patchMeApp(m).Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	out := decodePatchResponse(t, resp.Body)

	if got := out["firstName"]; got != "New" {
		t.Errorf("firstName = %v, want New (PATCH должен применить изменения)", got)
	}

	tier, ok := out["subscriptionTier"].(map[string]any)
	if !ok {
		t.Fatalf("subscriptionTier отсутствует в ответе PATCH /members/me; got keys = %v", keysOf(out))
	}
	if tier["slug"] != "foreman" {
		t.Errorf("subscriptionTier.slug = %v, want foreman", tier["slug"])
	}
	if level, _ := tier["level"].(float64); int(level) != 2 {
		t.Errorf("subscriptionTier.level = %v, want 2", tier["level"])
	}
}

// TestUpdateProfile_NoSubscription — у member'а без записи в
// subscription_users поле subscriptionTier должно отсутствовать (omitempty),
// а сам PATCH должен отрабатывать без ошибок.
func TestUpdateProfile_NoSubscription(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db,
		"member_roles", "members",
		"subscription_users",
	)

	m := &models.Member{
		TelegramID: 991002,
		Username:   "patch_no_sub",
		FirstName:  "Free",
		LastName:   "User",
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}

	body, _ := json.Marshal(map[string]any{
		"firstName": "Free",
		"lastName":  "User",
		"bio":       "",
		"grade":     "",
		"company":   "",
		"tg":        "patch_no_sub",
	})
	req, _ := http.NewRequest(http.MethodPatch, "/members/me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := patchMeApp(m).Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	out := decodePatchResponse(t, resp.Body)
	if _, present := out["subscriptionTier"]; present {
		t.Errorf("subscriptionTier должен отсутствовать (omitempty) для юзера без подписки; got = %v", out["subscriptionTier"])
	}
}

func keysOf(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
