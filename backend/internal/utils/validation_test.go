package utils

import (
	"strings"
	"testing"
	"time"
)

func TestSanitizeTelegramUsername(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"short", "a", "a"},
		{"max length", strings.Repeat("a", 32), strings.Repeat("a", 32)},
		{"too long", strings.Repeat("a", 33), ""},
		{"with underscore", "joint_imer", "joint_imer"},
		{"digits", "user123", "user123"},
		{"@-prefix forbidden", "@username", ""},
		{"cyrillic", "Александр", ""},
		{"dash forbidden", "joint-imer", ""},
		{"space forbidden", "two words", ""},
		{"emoji forbidden", "user🔥", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeTelegramUsername(tt.in); got != tt.want {
				t.Errorf("SanitizeTelegramUsername(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{"empty allowed (admin can clear)", "", false},
		{"valid", "user_42", false},
		{"too long", strings.Repeat("a", 33), true},
		{"bad chars", "user@host", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q) err=%v, wantErr=%v", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestValidateProfileLengths(t *testing.T) {
	long := strings.Repeat("я", 65)
	bigBio := strings.Repeat("a", 501)
	if err := ValidateProfileLengths(long, "Doe", "bio", "g", "c"); err == nil {
		t.Fatal("expected error on long firstName")
	}
	if err := ValidateProfileLengths("Jane", long, "bio", "g", "c"); err == nil {
		t.Fatal("expected error on long lastName")
	}
	if err := ValidateProfileLengths("Jane", "Doe", bigBio, "g", "c"); err == nil {
		t.Fatal("expected error on long bio")
	}
	// Кириллица 60 символов = 60 рун, должна пройти.
	if err := ValidateProfileLengths(strings.Repeat("я", 60), "Doe", "bio", "g", "c"); err != nil {
		t.Fatalf("60 cyrillic runes should pass, got %v", err)
	}
}

func TestValidateBirthdayRange(t *testing.T) {
	// nil — корректно (поле необязательное).
	if err := ValidateBirthdayRange(nil); err != nil {
		t.Fatalf("nil birthday must pass, got %v", err)
	}

	// Слишком давно.
	old := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := ValidateBirthdayRange(&old); err == nil {
		t.Fatal("1900 must fail")
	}

	// «Сегодня» по МСК — даже если в UTC ещё вчера, не должно падать.
	// Берём сегодняшнюю МСК-дату полночь — это валидный день рождения.
	todayMSK := time.Now().In(MSKLocation())
	birthdayToday := time.Date(todayMSK.Year(), todayMSK.Month(), todayMSK.Day(), 0, 0, 0, 0, MSKLocation())
	if err := ValidateBirthdayRange(&birthdayToday); err != nil {
		t.Fatalf("today (MSK midnight) must pass, got %v", err)
	}

	// Завтра — должно падать.
	tomorrow := birthdayToday.Add(48 * time.Hour)
	if err := ValidateBirthdayRange(&tomorrow); err == nil {
		t.Fatal("tomorrow must fail")
	}

	// Граничный кейс: 1920-01-01 — допустим.
	min := MinBirthday
	if err := ValidateBirthdayRange(&min); err != nil {
		t.Fatalf("MinBirthday must pass, got %v", err)
	}
}
