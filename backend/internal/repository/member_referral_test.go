package repository

import (
	"strings"
	"testing"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

// TestGenerateReferralCode_Format — проверка формата сгенерированного кода:
// длина = referralCodeLength, символы только из allowed-алфавита (без 0/O/1/I/L/U).
func TestGenerateReferralCode_Format(t *testing.T) {
	for i := 0; i < 1000; i++ {
		code, err := generateReferralCode()
		if err != nil {
			t.Fatalf("generateReferralCode: %v", err)
		}
		if len(code) != referralCodeLength {
			t.Fatalf("code length = %d, want %d", len(code), referralCodeLength)
		}
		for _, c := range code {
			if !strings.ContainsRune(referralCodeAlphabet, c) {
				t.Fatalf("forbidden char %q in code %q", c, code)
			}
		}
	}
}

// TestAssignReferralCode_Idempotent — повторный вызов на том же juзере
// возвращает уже установленный код, не генерирует новый.
func TestAssignReferralCode_Idempotent(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "members")

	m := &models.Member{TelegramID: 11001, Username: "ref_idem"}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}

	repo := NewMemberRepository()
	first, created, err := repo.AssignReferralCode(m.Id)
	if err != nil {
		t.Fatalf("first AssignReferralCode: %v", err)
	}
	if !created {
		t.Errorf("first call: created = false, want true")
	}
	if first == "" {
		t.Errorf("first code empty")
	}

	second, created, err := repo.AssignReferralCode(m.Id)
	if err != nil {
		t.Fatalf("second AssignReferralCode: %v", err)
	}
	if created {
		t.Errorf("second call: created = true, want false (already had code)")
	}
	if second != first {
		t.Errorf("second code %q != first %q", second, first)
	}
}

// TestSetReferredByMemberID_FirstWriteWins — повторный вызов с другим
// referrerID не перетирает первого инвайтера.
func TestSetReferredByMemberID_FirstWriteWins(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "members")

	alice := &models.Member{TelegramID: 12001, Username: "alice"}
	charlie := &models.Member{TelegramID: 12002, Username: "charlie"}
	bob := &models.Member{TelegramID: 12003, Username: "bob"}
	for _, m := range []*models.Member{alice, charlie, bob} {
		if err := db.Create(m).Error; err != nil {
			t.Fatalf("create member: %v", err)
		}
	}

	repo := NewMemberRepository()

	// Alice как первый инвайтер Боба.
	written, err := repo.SetReferredByMemberID(bob.Id, alice.Id)
	if err != nil {
		t.Fatalf("first SetReferredByMemberID: %v", err)
	}
	if !written {
		t.Errorf("first call: written = false, want true")
	}

	// Charlie пытается перезаписать — должен fail (RowsAffected = 0).
	written, err = repo.SetReferredByMemberID(bob.Id, charlie.Id)
	if err != nil {
		t.Fatalf("second SetReferredByMemberID: %v", err)
	}
	if written {
		t.Errorf("second call: written = true, want false (Alice already there)")
	}

	// Сверим что инвайтер всё ещё Alice.
	got, err := repo.GetReferredByMemberID(bob.Id)
	if err != nil {
		t.Fatalf("GetReferredByMemberID: %v", err)
	}
	if got == nil || *got != alice.Id {
		t.Errorf("referrer = %v, want %d (Alice)", got, alice.Id)
	}
}

// TestSetReferredByMemberID_RejectsBadInput — defense-in-depth: невалидный
// referrerID (0, отрицательный, self) отвергается на уровне репозитория.
func TestSetReferredByMemberID_RejectsBadInput(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "members")

	bob := &models.Member{TelegramID: 12100, Username: "bob_bad"}
	if err := db.Create(bob).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}
	repo := NewMemberRepository()

	cases := []struct {
		name       string
		referrerID int64
	}{
		{"zero", 0},
		{"negative", -1},
		{"self", bob.Id},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.SetReferredByMemberID(bob.Id, tc.referrerID)
			if err == nil {
				t.Errorf("expected error for referrerID=%d, got nil", tc.referrerID)
			}
		})
	}
}

// TestSetReferralWelcomeSeenAt_FirstWriteWins — повторный вызов не перезатирает
// timestamp, чтобы аналитика «когда впервые закрыл баннер» сохранялась.
func TestSetReferralWelcomeSeenAt_FirstWriteWins(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	testutil.TruncateAll(t, db, "members")

	bob := &models.Member{TelegramID: 12200, Username: "bob_seen"}
	if err := db.Create(bob).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}
	repo := NewMemberRepository()

	if err := repo.SetReferralWelcomeSeenAt(bob.Id); err != nil {
		t.Fatalf("first call: %v", err)
	}
	var firstSeen models.Member
	if err := db.First(&firstSeen, bob.Id).Error; err != nil {
		t.Fatalf("reload bob: %v", err)
	}
	if firstSeen.ReferralWelcomeSeenAt == nil {
		t.Fatalf("first call did not set timestamp")
	}
	t1 := *firstSeen.ReferralWelcomeSeenAt

	// Повторный вызов не должен изменить timestamp.
	if err := repo.SetReferralWelcomeSeenAt(bob.Id); err != nil {
		t.Fatalf("second call: %v", err)
	}
	var secondSeen models.Member
	if err := db.First(&secondSeen, bob.Id).Error; err != nil {
		t.Fatalf("reload bob: %v", err)
	}
	if !secondSeen.ReferralWelcomeSeenAt.Equal(t1) {
		t.Errorf("timestamp changed: %v → %v", t1, *secondSeen.ReferralWelcomeSeenAt)
	}
}
