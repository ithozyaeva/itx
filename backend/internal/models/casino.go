package models

import (
	"encoding/json"
	"time"
)

type CasinoBet struct {
	Id        int64           `json:"id" gorm:"primaryKey"`
	MemberId  int64           `json:"memberId" gorm:"column:member_id;not null"`
	Game      string          `json:"game" gorm:"column:game;size:20;not null"`
	BetAmount int             `json:"betAmount" gorm:"column:bet_amount;not null"`
	Multiplier float64        `json:"multiplier" gorm:"column:multiplier;not null;default:0"`
	Payout    int             `json:"payout" gorm:"column:payout;not null;default:0"`
	Result    json.RawMessage `json:"result" gorm:"column:result;type:jsonb;not null;default:'{}'"`
	CreatedAt time.Time       `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

type CoinFlipRequest struct {
	BetAmount int    `json:"betAmount"`
	Choice    string `json:"choice"`
}

type DiceRollRequest struct {
	BetAmount int    `json:"betAmount"`
	Target    int    `json:"target"`
	Direction string `json:"direction"`
}

type WheelRequest struct {
	BetAmount int `json:"betAmount"`
}

type CasinoBetResponse struct {
	Id         int64           `json:"id"`
	Game       string          `json:"game"`
	BetAmount  int             `json:"betAmount"`
	Multiplier float64         `json:"multiplier"`
	Payout     int             `json:"payout"`
	Result     json.RawMessage `json:"result"`
	Won        bool            `json:"won"`
	Balance    int             `json:"balance"`
	CreatedAt  time.Time       `json:"createdAt"`
}

type CasinoHistoryItem struct {
	Id         int64           `json:"id"`
	Game       string          `json:"game"`
	BetAmount  int             `json:"betAmount"`
	Multiplier float64         `json:"multiplier"`
	Payout     int             `json:"payout"`
	Result     json.RawMessage `json:"result"`
	CreatedAt  time.Time       `json:"createdAt"`
}

type CasinoStats struct {
	TotalBets   int64   `json:"totalBets"`
	TotalWagered int    `json:"totalWagered"`
	TotalPayout  int    `json:"totalPayout"`
	BiggestWin   int    `json:"biggestWin"`
}

type CasinoAdminStats struct {
	TotalBets    int64  `json:"totalBets"`
	TotalWagered int    `json:"totalWagered"`
	TotalPayout  int    `json:"totalPayout"`
	HouseProfit  int    `json:"houseProfit"`
	UniquePlayers int64 `json:"uniquePlayers"`
}

type CasinoAdminBet struct {
	Id              int64           `json:"id"`
	MemberId        int64           `json:"memberId"`
	MemberFirstName string          `json:"memberFirstName"`
	MemberLastName  string          `json:"memberLastName"`
	MemberUsername  string          `json:"memberUsername"`
	Game            string          `json:"game"`
	BetAmount       int             `json:"betAmount"`
	Multiplier      float64         `json:"multiplier"`
	Payout          int             `json:"payout"`
	Result          json.RawMessage `json:"result"`
	CreatedAt       time.Time       `json:"createdAt"`
}
