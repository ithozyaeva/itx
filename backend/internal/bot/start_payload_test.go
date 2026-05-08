package bot

import "testing"

func TestParseReferralPayload(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int64
	}{
		{"valid simple", "ref_42", 42},
		{"valid large", "ref_1000000", 1000000},
		{"empty", "", 0},
		{"plain sub", "sub", 0},
		{"prefix without id", "ref_", 0},
		{"non-numeric", "ref_abc", 0},
		{"negative", "ref_-1", 0},
		{"zero", "ref_0", 0},
		{"different prefix", "referer_42", 0},
		{"trailing junk", "ref_42abc", 0},
		{"hyphen instead of underscore", "ref-42", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseReferralPayload(tt.in); got != tt.want {
				t.Errorf("parseReferralPayload(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
