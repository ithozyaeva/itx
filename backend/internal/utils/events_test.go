package utils

import (
	"ithozyeva/internal/models"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }

func TestNextOccurrence_NonRepeating(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	past := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	event := &models.Event{Date: past, IsRepeating: false}
	if got := NextOccurrence(event, now); !got.Equal(past) {
		t.Errorf("non-repeating: got %v, want %v", got, past)
	}
}

func TestNextOccurrence_RepeatingButFutureStart(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	future := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	weekly := string(models.RepeatWeekly)
	event := &models.Event{Date: future, IsRepeating: true, RepeatPeriod: &weekly}
	if got := NextOccurrence(event, now); !got.Equal(future) {
		t.Errorf("future start: got %v, want %v", got, future)
	}
}

func TestNextOccurrence_Weekly(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	start := time.Date(2026, 3, 15, 11, 0, 0, 0, time.UTC) // воскресенье
	weekly := string(models.RepeatWeekly)
	event := &models.Event{Date: start, IsRepeating: true, RepeatPeriod: &weekly, RepeatInterval: ptr(1)}
	want := time.Date(2026, 4, 26, 11, 0, 0, 0, time.UTC) // ближайшее воскресенье после 24 апр
	if got := NextOccurrence(event, now); !got.Equal(want) {
		t.Errorf("weekly: got %v, want %v", got, want)
	}
}

func TestNextOccurrence_WeeklyWithInterval(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	start := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	weekly := string(models.RepeatWeekly)
	event := &models.Event{Date: start, IsRepeating: true, RepeatPeriod: &weekly, RepeatInterval: ptr(2)}
	// каждые 2 недели от 1 марта: 15 мар, 29 мар, 12 апр, 26 апр...
	want := time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC)
	if got := NextOccurrence(event, now); !got.Equal(want) {
		t.Errorf("weekly interval=2: got %v, want %v", got, want)
	}
}

func TestNextOccurrence_Daily(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	start := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	daily := string(models.RepeatDaily)
	event := &models.Event{Date: start, IsRepeating: true, RepeatPeriod: &daily, RepeatInterval: ptr(1)}
	want := time.Date(2026, 4, 25, 10, 0, 0, 0, time.UTC)
	if got := NextOccurrence(event, now); !got.Equal(want) {
		t.Errorf("daily: got %v, want %v", got, want)
	}
}

func TestNextOccurrence_Monthly(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	start := time.Date(2026, 3, 7, 11, 0, 0, 0, time.UTC)
	monthly := string(models.RepeatMonthly)
	event := &models.Event{Date: start, IsRepeating: true, RepeatPeriod: &monthly, RepeatInterval: ptr(1)}
	// 7 мар → 7 апр (прошёл) → 7 мая
	want := time.Date(2026, 5, 7, 11, 0, 0, 0, time.UTC)
	if got := NextOccurrence(event, now); !got.Equal(want) {
		t.Errorf("monthly: got %v, want %v", got, want)
	}
}

func TestNextOccurrence_Yearly(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	start := time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC)
	yearly := string(models.RepeatYearly)
	event := &models.Event{Date: start, IsRepeating: true, RepeatPeriod: &yearly, RepeatInterval: ptr(1)}
	want := time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC)
	if got := NextOccurrence(event, now); !got.Equal(want) {
		t.Errorf("yearly: got %v, want %v", got, want)
	}
}

func TestNextOccurrence_NilPeriod(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	start := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	event := &models.Event{Date: start, IsRepeating: true, RepeatPeriod: nil}
	if got := NextOccurrence(event, now); !got.Equal(start) {
		t.Errorf("nil period: got %v, want %v", got, start)
	}
}
