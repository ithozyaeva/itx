package service

import (
	"testing"
	"time"
)

func TestParseHumanDuration(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
		bad  bool
	}{
		{"", 0, false},
		{"30s", 30 * time.Second, false},
		{"30m", 30 * time.Minute, false},
		{"1h", time.Hour, false},
		{"2h", 2 * time.Hour, false},
		{"1d", 24 * time.Hour, false},
		{"7d", 7 * 24 * time.Hour, false},
		{"1d12h", 36 * time.Hour, false},
		{"2h30m", 2*time.Hour + 30*time.Minute, false},
		{" 1H ", time.Hour, false},
		{"30", 0, true},
		{"abc", 0, true},
		{"10x", 0, true},
		{"d", 0, true},
	}
	for _, c := range cases {
		got, err := ParseHumanDuration(c.in)
		if c.bad {
			if err == nil {
				t.Errorf("ParseHumanDuration(%q): want error, got %v", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseHumanDuration(%q): unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseHumanDuration(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestFormatDurationHuman(t *testing.T) {
	cases := []struct {
		in   time.Duration
		want string
	}{
		{0, "навсегда"},
		{-time.Second, "навсегда"},
		{30 * time.Minute, "30м"},
		{time.Hour, "1ч"},
		{2 * time.Hour, "2ч"},
		{24 * time.Hour, "1д"},
		{7 * 24 * time.Hour, "7д"},
	}
	for _, c := range cases {
		if got := FormatDurationHuman(c.in); got != c.want {
			t.Errorf("FormatDurationHuman(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
