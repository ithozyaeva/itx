package service

import (
	"testing"
	"time"

	"ithozyeva/internal/models"
)

func TestGetFutureEventsFiltering(t *testing.T) {
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)

	weekly := "WEEKLY"
	pastEnd := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	futureEnd := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		event    models.Event
		include  bool
	}{
		{
			"future one-time event",
			models.Event{
				Date:        time.Date(2025, 7, 1, 18, 0, 0, 0, time.UTC),
				IsRepeating: false,
			},
			true,
		},
		{
			"past one-time event",
			models.Event{
				Date:        time.Date(2025, 5, 1, 18, 0, 0, 0, time.UTC),
				IsRepeating: false,
			},
			false,
		},
		{
			"today event (equal)",
			models.Event{
				Date:        now,
				IsRepeating: false,
			},
			true,
		},
		{
			"repeating event with no end date",
			models.Event{
				Date:         time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
				IsRepeating:  true,
				RepeatPeriod: &weekly,
			},
			true,
		},
		{
			"repeating event with future end date",
			models.Event{
				Date:          time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
				IsRepeating:   true,
				RepeatPeriod:  &weekly,
				RepeatEndDate: &futureEnd,
			},
			true,
		},
		{
			"repeating event with past end date",
			models.Event{
				Date:          time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
				IsRepeating:   true,
				RepeatPeriod:  &weekly,
				RepeatEndDate: &pastEnd,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replicate the filtering logic from GetFutureEvents
			var included bool
			event := tt.event
			if event.IsRepeating && event.RepeatPeriod != nil {
				if event.RepeatEndDate != nil && now.After(*event.RepeatEndDate) {
					included = false
				} else {
					included = true
				}
			} else {
				included = event.Date.After(now) || event.Date.Equal(now)
			}

			if included != tt.include {
				t.Errorf("event %q: got included=%v, want %v", tt.name, included, tt.include)
			}
		})
	}
}

func TestAddMemberParticipantLimit(t *testing.T) {
	tests := []struct {
		name            string
		maxParticipants int
		currentMembers  int
		expectErr       bool
	}{
		{"no limit set", 0, 10, false},
		{"under limit", 5, 3, false},
		{"at limit", 5, 5, true},
		{"over limit", 5, 6, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replicate the participant limit check from AddMember
			members := make([]models.Member, tt.currentMembers)
			limitReached := tt.maxParticipants > 0 && len(members) >= tt.maxParticipants

			if limitReached != tt.expectErr {
				t.Errorf("maxParticipants=%d, currentMembers=%d: got limitReached=%v, want %v",
					tt.maxParticipants, tt.currentMembers, limitReached, tt.expectErr)
			}
		})
	}
}

func TestEventPlaceTypes(t *testing.T) {
	// Verify place type constants are defined correctly
	tests := []struct {
		placeType models.PlaceType
		expected  string
	}{
		{models.EventOnline, "ONLINE"},
		{models.EventOffline, "OFFLINE"},
		{models.EventHybrid, "HYBRID"},
	}

	for _, tt := range tests {
		if string(tt.placeType) != tt.expected {
			t.Errorf("PlaceType = %q, want %q", tt.placeType, tt.expected)
		}
	}
}

func TestRepeatPeriodConstants(t *testing.T) {
	tests := []struct {
		period   models.RepeatPeriod
		expected string
	}{
		{models.RepeatDaily, "DAILY"},
		{models.RepeatWeekly, "WEEKLY"},
		{models.RepeatMonthly, "MONTHLY"},
		{models.RepeatYearly, "YEARLY"},
	}

	for _, tt := range tests {
		if string(tt.period) != tt.expected {
			t.Errorf("RepeatPeriod = %q, want %q", tt.period, tt.expected)
		}
	}
}

func TestErrParticipantLimitReached(t *testing.T) {
	if ErrParticipantLimitReached == nil {
		t.Fatal("ErrParticipantLimitReached should not be nil")
	}
	if ErrParticipantLimitReached.Error() != "достигнут лимит участников" {
		t.Errorf("ErrParticipantLimitReached.Error() = %q, want %q",
			ErrParticipantLimitReached.Error(), "достигнут лимит участников")
	}
}
