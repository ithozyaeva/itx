package service

import (
	"testing"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/utils"
)

// mskDay — short helper для тестов: фиксированная МСК-дата без time-of-day.
func mskDay(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, utils.MSKLocation())
}

// fixedThresholds — урезанный набор для предсказуемости тестов.
var fixedThresholds = []StreakThreshold{
	{Days: 3, Reason: models.PointReasonDailyStreak3, Reward: 15},
	{Days: 7, Reason: models.PointReasonDailyStreak7, Reward: 50},
	{Days: 14, Reason: models.PointReasonDailyStreak14, Reward: 150},
	{Days: 30, Reason: models.PointReasonDailyStreak30, Reward: 500},
}

func ptrTime(t time.Time) *time.Time { return &t }
func ptrInt(i int) *int              { return &i }

func TestRecalcStreak(t *testing.T) {
	day := mskDay(2026, 5, 4) // понедельник, ISO-неделя 19
	yr, wk := day.ISOWeek()

	type want struct {
		current  int
		longest  int
		freezes  int
		crossed  []int
		lastDate *time.Time
	}

	tests := []struct {
		name   string
		streak *models.MemberStreak
		day    time.Time
		want   want
	}{
		{
			name:   "first check-in (no prior data)",
			streak: &models.MemberStreak{},
			day:    day,
			want: want{
				current:  1,
				longest:  1,
				freezes:  1,
				lastDate: ptrTime(day),
			},
		},
		{
			name: "same-day check-in is a no-op",
			streak: &models.MemberStreak{
				CurrentStreak:    5,
				LongestStreak:    7,
				LastCheckInDate:  ptrTime(day),
				FreezesAvailable: 1,
				FreezeWeekYear:   ptrInt(yr),
				FreezeWeekNum:    ptrInt(wk),
			},
			day: day,
			want: want{
				current:  5,
				longest:  7,
				freezes:  1,
				lastDate: ptrTime(day),
			},
		},
		{
			name: "consecutive day increments streak",
			streak: &models.MemberStreak{
				CurrentStreak:    2,
				LongestStreak:    2,
				LastCheckInDate:  ptrTime(mskDay(2026, 5, 3)),
				FreezesAvailable: 0,
				FreezeWeekYear:   ptrInt(yr),
				FreezeWeekNum:    ptrInt(wk - 1),
			},
			day: day,
			want: want{
				current: 3,
				longest: 3,
				// new ISO-week → freeze restored
				freezes:  1,
				crossed:  []int{3},
				lastDate: ptrTime(day),
			},
		},
		{
			name: "gap=2 with freeze available consumes freeze",
			streak: &models.MemberStreak{
				CurrentStreak:    6,
				LongestStreak:    10,
				LastCheckInDate:  ptrTime(mskDay(2026, 5, 2)),
				FreezesAvailable: 1,
				FreezeWeekYear:   ptrInt(yr),
				FreezeWeekNum:    ptrInt(wk),
			},
			day: day,
			want: want{
				current:  7,
				longest:  10,
				freezes:  0,
				crossed:  []int{7},
				lastDate: ptrTime(day),
			},
		},
		{
			name: "gap=2 without freeze resets to 1",
			streak: &models.MemberStreak{
				CurrentStreak:    9,
				LongestStreak:    9,
				LastCheckInDate:  ptrTime(mskDay(2026, 5, 2)),
				FreezesAvailable: 0,
				FreezeWeekYear:   ptrInt(yr),
				FreezeWeekNum:    ptrInt(wk),
			},
			day: day,
			want: want{
				current:  1,
				longest:  9,
				freezes:  0,
				lastDate: ptrTime(day),
			},
		},
		{
			name: "gap=3 always resets even with freeze",
			streak: &models.MemberStreak{
				CurrentStreak:    4,
				LongestStreak:    4,
				LastCheckInDate:  ptrTime(mskDay(2026, 5, 1)),
				FreezesAvailable: 1,
				FreezeWeekYear:   ptrInt(yr),
				FreezeWeekNum:    ptrInt(wk),
			},
			day: day,
			want: want{
				current: 1,
				longest: 4,
				// freeze не тратится при reset
				freezes:  1,
				lastDate: ptrTime(day),
			},
		},
		{
			name: "crossing 7-day threshold from 6→7",
			streak: &models.MemberStreak{
				CurrentStreak:    6,
				LongestStreak:    6,
				LastCheckInDate:  ptrTime(mskDay(2026, 5, 3)),
				FreezesAvailable: 0,
				FreezeWeekYear:   ptrInt(yr),
				FreezeWeekNum:    ptrInt(wk - 1),
			},
			day: day,
			want: want{
				current:  7,
				longest:  7,
				freezes:  1,
				crossed:  []int{7},
				lastDate: ptrTime(day),
			},
		},
		{
			name: "no double-award on revisit (current already past threshold)",
			streak: &models.MemberStreak{
				CurrentStreak:    7,
				LongestStreak:    7,
				LastCheckInDate:  ptrTime(mskDay(2026, 5, 3)),
				FreezesAvailable: 0,
				FreezeWeekYear:   ptrInt(yr),
				FreezeWeekNum:    ptrInt(wk - 1),
			},
			day: day,
			want: want{
				current:  8,
				longest:  8,
				freezes:  1,
				crossed:  []int{}, // 7 уже пересечён, новых нет
				lastDate: ptrTime(day),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, crossed := recalcStreak(tt.streak, tt.day, fixedThresholds)
			if got.CurrentStreak != tt.want.current {
				t.Errorf("CurrentStreak = %d, want %d", got.CurrentStreak, tt.want.current)
			}
			if got.LongestStreak != tt.want.longest {
				t.Errorf("LongestStreak = %d, want %d", got.LongestStreak, tt.want.longest)
			}
			if got.FreezesAvailable != tt.want.freezes {
				t.Errorf("FreezesAvailable = %d, want %d", got.FreezesAvailable, tt.want.freezes)
			}
			if tt.want.lastDate != nil && (got.LastCheckInDate == nil || !got.LastCheckInDate.Equal(*tt.want.lastDate)) {
				t.Errorf("LastCheckInDate = %v, want %v", got.LastCheckInDate, tt.want.lastDate)
			}

			gotCrossed := make([]int, 0, len(crossed))
			for _, c := range crossed {
				gotCrossed = append(gotCrossed, c.Days)
			}
			if !equalIntSlice(gotCrossed, tt.want.crossed) {
				t.Errorf("crossed thresholds = %v, want %v", gotCrossed, tt.want.crossed)
			}
		})
	}
}

func equalIntSlice(a, b []int) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
