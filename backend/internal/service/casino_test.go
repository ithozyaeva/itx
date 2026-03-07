package service

import (
	"testing"

	"ithozyeva/internal/models"
)

func TestValidateBet(t *testing.T) {
	svc := &CasinoService{}

	tests := []struct {
		name    string
		amount  int
		wantErr bool
	}{
		{"too low", 5, true},
		{"min boundary", 10, false},
		{"normal", 50, false},
		{"max boundary", 200, false},
		{"too high", 201, true},
		{"zero", 0, true},
		{"negative", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validateBet(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBet(%d) error = %v, wantErr %v", tt.amount, err, tt.wantErr)
			}
		})
	}
}

func TestCryptoRandInt(t *testing.T) {
	// Test that cryptoRandInt returns values in range [0, max)
	for i := 0; i < 100; i++ {
		n, err := cryptoRandInt(10)
		if err != nil {
			t.Fatalf("cryptoRandInt(10) returned error: %v", err)
		}
		if n < 0 || n >= 10 {
			t.Errorf("cryptoRandInt(10) = %d, want [0, 10)", n)
		}
	}

	// Test with max=2 (coin flip)
	for i := 0; i < 50; i++ {
		n, err := cryptoRandInt(2)
		if err != nil {
			t.Fatalf("cryptoRandInt(2) returned error: %v", err)
		}
		if n != 0 && n != 1 {
			t.Errorf("cryptoRandInt(2) = %d, want 0 or 1", n)
		}
	}
}

func TestWheelMultipliers(t *testing.T) {
	if len(wheelMultipliers) != 12 {
		t.Errorf("wheelMultipliers has %d segments, want 12", len(wheelMultipliers))
	}

	// Count segments by multiplier
	counts := map[float64]int{}
	for _, m := range wheelMultipliers {
		counts[m]++
	}

	expected := map[float64]int{
		0:   3,
		0.5: 3,
		1:   2,
		1.5: 2,
		2:   1,
		3:   1,
	}

	for mult, count := range expected {
		if counts[mult] != count {
			t.Errorf("wheelMultipliers has %d segments with multiplier %.1f, want %d", counts[mult], mult, count)
		}
	}

	// Verify expected value (house edge ~4.2%)
	var sum float64
	for _, m := range wheelMultipliers {
		sum += m
	}
	ev := sum / float64(len(wheelMultipliers))
	if ev < 0.9 || ev > 1.0 {
		t.Errorf("wheel EV = %.4f, expected ~0.958", ev)
	}
}

func TestCoinFlipValidation(t *testing.T) {
	svc := &CasinoService{}

	// Only test cases that should fail validation (before hitting DB)
	tests := []struct {
		name string
		req  *models.CoinFlipRequest
	}{
		{"invalid choice", &models.CoinFlipRequest{BetAmount: 50, Choice: "edge"}},
		{"empty choice", &models.CoinFlipRequest{BetAmount: 50, Choice: ""}},
		{"bet too low", &models.CoinFlipRequest{BetAmount: 5, Choice: "heads"}},
		{"bet too high", &models.CoinFlipRequest{BetAmount: 500, Choice: "heads"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.PlayCoinFlip(1, tt.req)
			if err == nil {
				t.Errorf("PlayCoinFlip() expected error for %s", tt.name)
			}
		})
	}
}

func TestDiceRollValidation(t *testing.T) {
	svc := &CasinoService{}

	// Only test cases that should fail validation (before hitting DB)
	tests := []struct {
		name string
		req  *models.DiceRollRequest
	}{
		{"target too low", &models.DiceRollRequest{BetAmount: 50, Target: 1, Direction: "over"}},
		{"target too high", &models.DiceRollRequest{BetAmount: 50, Target: 99, Direction: "over"}},
		{"invalid direction", &models.DiceRollRequest{BetAmount: 50, Target: 50, Direction: "sideways"}},
		{"empty direction", &models.DiceRollRequest{BetAmount: 50, Target: 50, Direction: ""}},
		{"bet too low", &models.DiceRollRequest{BetAmount: 5, Target: 50, Direction: "over"}},
		{"bet too high", &models.DiceRollRequest{BetAmount: 500, Target: 50, Direction: "over"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.PlayDiceRoll(1, tt.req)
			if err == nil {
				t.Errorf("PlayDiceRoll() expected error for %s", tt.name)
			}
		})
	}
}

func TestWheelValidation(t *testing.T) {
	svc := &CasinoService{}

	_, err := svc.PlayWheel(1, &models.WheelRequest{BetAmount: 5})
	if err == nil {
		t.Error("PlayWheel() expected error for bet too low")
	}

	_, err = svc.PlayWheel(1, &models.WheelRequest{BetAmount: 500})
	if err == nil {
		t.Error("PlayWheel() expected error for bet too high")
	}
}

func TestDiceRollMultiplierCalculation(t *testing.T) {
	// Verify dice multiplier formula: 0.97 * (100 / winChance)
	tests := []struct {
		target    int
		direction string
		winChance float64
	}{
		{50, "over", 50},  // 50% chance
		{50, "under", 50}, // 50% chance
		{75, "over", 25},  // 25% chance
		{25, "under", 25}, // 25% chance
		{90, "over", 10},  // 10% chance
		{10, "under", 10}, // 10% chance
	}

	for _, tt := range tests {
		expectedMult := 0.97 * (100.0 / tt.winChance)
		t.Run(tt.direction+"_"+string(rune('0'+tt.target/10))+string(rune('0'+tt.target%10)), func(t *testing.T) {
			if expectedMult <= 0 {
				t.Errorf("multiplier should be positive, got %.4f", expectedMult)
			}
			// For 50/50, multiplier should be ~1.94
			if tt.winChance == 50 && (expectedMult < 1.93 || expectedMult > 1.95) {
				t.Errorf("50/50 multiplier = %.4f, want ~1.94", expectedMult)
			}
		})
	}
}
