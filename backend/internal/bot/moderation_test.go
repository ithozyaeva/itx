package bot

import "testing"

func TestComputeVotebanThreshold(t *testing.T) {
	cases := []struct {
		active int64
		want   int
	}{
		{0, votebanRequiredVotesMin},
		{1, votebanRequiredVotesMin},
		{10, votebanRequiredVotesMin},  // 1.5 → 2 → clamp 3
		{20, votebanRequiredVotesMin},  // 3
		{30, 5},                        // 4.5 → 5
		{50, 8},                        // 7.5 → 8
		{67, votebanRequiredVotesMax},  // 10.05 → 10
		{100, votebanRequiredVotesMax}, // 15 → clamp 10
		{500, votebanRequiredVotesMax}, // clamp 10
	}
	for _, c := range cases {
		if got := computeVotebanThreshold(c.active); got != c.want {
			t.Errorf("computeVotebanThreshold(%d) = %d, want %d", c.active, got, c.want)
		}
	}
}
