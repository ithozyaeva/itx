package bot

import "testing"

func TestParseReferralPayload(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"valid 8 chars", "ref_A8F3K2P9", "A8F3K2P9"},
		{"valid 4 min", "ref_A8F3", "A8F3"},
		{"valid digits only", "ref_23456789", "23456789"},
		{"empty", "", ""},
		{"plain sub", "sub", ""},
		{"prefix without code", "ref_", ""},
		{"too short", "ref_AB1", ""},
		{"too long", "ref_AAAAAAAAAAAAAAAAA", ""},
		{"contains forbidden 0", "ref_A0F3K2P9", ""},
		{"contains forbidden O", "ref_A8F3KOP9", ""},
		{"contains forbidden 1", "ref_A8F3K1P9", ""},
		{"contains forbidden I", "ref_A8F3KIP9", ""},
		{"contains forbidden L", "ref_A8F3KLP9", ""},
		{"contains forbidden U", "ref_A8F3KUP9", ""},
		{"lowercase rejected", "ref_a8f3k2p9", ""},
		{"hyphen prefix", "ref-A8F3K2P9", ""},
		{"different prefix", "referer_A8F3K2", ""},
		{"trailing junk", "ref_A8F3K2P9!", ""},
		{"leading space", "ref_ A8F3K2P9", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseReferralPayload(tt.in); got != tt.want {
				t.Errorf("parseReferralPayload(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
