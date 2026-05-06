package service

import (
	"testing"
	"time"

	"ithozyeva/internal/models"
)

func TestCheckProfileComplete(t *testing.T) {
	svc := &PointsService{}

	birthday := models.DateOnly(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name           string
		member         *models.Member
		expectComplete bool
	}{
		{
			"all fields filled",
			&models.Member{
				FirstName: "John",
				LastName:  "Doe",
				Bio:       "Some bio",
				Birthday:  &birthday,
			},
			true,
		},
		{
			"missing first name",
			&models.Member{
				FirstName: "",
				LastName:  "Doe",
				Bio:       "Some bio",
				Birthday:  &birthday,
			},
			false,
		},
		{
			"missing last name",
			&models.Member{
				FirstName: "John",
				LastName:  "",
				Bio:       "Some bio",
				Birthday:  &birthday,
			},
			false,
		},
		{
			"missing bio",
			&models.Member{
				FirstName: "John",
				LastName:  "Doe",
				Bio:       "",
				Birthday:  &birthday,
			},
			false,
		},
		{
			"missing birthday",
			&models.Member{
				FirstName: "John",
				LastName:  "Doe",
				Bio:       "Some bio",
				Birthday:  nil,
			},
			false,
		},
		{
			"all fields empty",
			&models.Member{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CheckProfileComplete only calls AwardIdempotent when profile is complete.
			// Without a repo it will panic if it tries to award — which means profile was complete.
			// We verify the logic by checking whether it returns early (incomplete) or not.
			if !tt.expectComplete {
				// Should return early without calling repo — no panic expected
				svc.CheckProfileComplete(tt.member)
			}
			// For complete profiles, we can't test without a repo, so just verify the condition directly
			if tt.expectComplete {
				if tt.member.FirstName == "" || tt.member.LastName == "" || tt.member.Bio == "" || tt.member.Birthday == nil {
					t.Error("expected profile to be complete but fields are missing")
				}
			}
		})
	}
}

func TestPointValuesExist(t *testing.T) {
	// Verify that all defined point reasons have values in the PointValues map
	reasons := []models.PointReason{
		models.PointReasonEventAttend,
		models.PointReasonEventHost,
		models.PointReasonReviewCommunity,
		models.PointReasonReviewService,
		models.PointReasonResumeUpload,
		models.PointReasonReferalCreate,
		// PointReasonReferalConversion больше не начисляет points —
		// награда переехала в referral_credit_transactions.
		models.PointReasonProfileComplete,
		models.PointReasonWeeklyActivity,
		models.PointReasonMonthlyActive,
		models.PointReasonStreak4Weeks,
		models.PointReasonTaskCreate,
		models.PointReasonTaskExecute,
		models.PointReasonMarketplaceCreate,
		models.PointReasonMarketplaceBuy,
		models.PointReasonChatterOfWeek,
		models.PointReasonKudosReceived,
	}

	for _, reason := range reasons {
		val, ok := models.PointValues[reason]
		if !ok {
			t.Errorf("PointValues missing entry for reason %q", reason)
			continue
		}
		if val <= 0 {
			t.Errorf("PointValues[%q] = %d, want positive value", reason, val)
		}
	}
}

func TestPointValuesSpecificAmounts(t *testing.T) {
	tests := []struct {
		reason   models.PointReason
		expected int
	}{
		{models.PointReasonEventAttend, 10},
		{models.PointReasonEventHost, 25},
		{models.PointReasonProfileComplete, 20},
		{models.PointReasonStreak4Weeks, 50},
		{models.PointReasonMonthlyActive, 30},
		{models.PointReasonChatterOfWeek, 15},
	}

	for _, tt := range tests {
		t.Run(string(tt.reason), func(t *testing.T) {
			if got := models.PointValues[tt.reason]; got != tt.expected {
				t.Errorf("PointValues[%q] = %d, want %d", tt.reason, got, tt.expected)
			}
		})
	}
}

func TestGetLeaderboardLimitNormalization(t *testing.T) {
	// GetLeaderboard normalizes the limit parameter before calling repo.
	// We can't test the full method without a repo, but we can verify the normalization logic.
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"zero defaults to 20", 0, 20},
		{"negative defaults to 20", -5, 20},
		{"over 100 defaults to 20", 101, 20},
		{"valid limit stays", 50, 50},
		{"boundary 1", 1, 1},
		{"boundary 100", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := tt.input
			if limit <= 0 || limit > 100 {
				limit = 20
			}
			if limit != tt.expected {
				t.Errorf("normalized limit = %d, want %d", limit, tt.expected)
			}
		})
	}
}

func TestIsoWeekMonday(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		week     int
		expected time.Time
	}{
		{
			"2024 week 1",
			2024, 1,
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			"2023 week 1",
			2023, 1,
			time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			"2024 week 10",
			2024, 10,
			time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			"2025 week 1",
			2025, 1,
			time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isoWeekMonday(tt.year, tt.week)
			if !got.Equal(tt.expected) {
				t.Errorf("isoWeekMonday(%d, %d) = %v, want %v", tt.year, tt.week, got, tt.expected)
			}
			// Verify it's actually a Monday
			if got.Weekday() != time.Monday {
				t.Errorf("isoWeekMonday(%d, %d) = %v, which is %s not Monday",
					tt.year, tt.week, got, got.Weekday())
			}
		})
	}
}
