package bot

import (
	"testing"
	"time"
)

func TestFormatRemainingHuman(t *testing.T) {
	cases := []struct {
		in   time.Duration
		want string
	}{
		{0, "0с"},
		{-time.Second, "0с"},
		{30 * time.Second, "30с"},
		{15 * time.Minute, "15м"},
		{13*time.Minute + 18*time.Second + 12*time.Millisecond, "13м 18с"},
		{time.Hour, "1ч"},
		{time.Hour + 5*time.Minute, "1ч 5м"},
	}
	for _, c := range cases {
		if got := formatRemainingHuman(c.in); got != c.want {
			t.Errorf("formatRemainingHuman(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestComputeVotebanThreshold(t *testing.T) {
	cases := []struct {
		active int64
		want   int
	}{
		{0, votebanRequiredVotesMin},   // clamp 3
		{20, votebanRequiredVotesMin},  // 1.6 → 2 → clamp 3
		{30, votebanRequiredVotesMin},  // 2.4 → 2 → clamp 3
		{38, votebanRequiredVotesMin},  // 3.04 → 3
		{50, 4},                        // 4.0
		{63, 5},                        // 5.04
		{80, 6},                        // 6.4 → 6
		{94, 8},                        // 7.52 → 8
		{100, votebanRequiredVotesMax}, // 8.0
		{500, votebanRequiredVotesMax}, // clamp 8
	}
	for _, c := range cases {
		if got := computeVotebanThreshold(c.active); got != c.want {
			t.Errorf("computeVotebanThreshold(%d) = %d, want %d", c.active, got, c.want)
		}
	}
}
