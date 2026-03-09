package models

import (
	"time"
)

type CasinoBet struct {
	Id         int64   `json:"id" gorm:"primaryKey"`
	MemberId   int64   `json:"memberId" gorm:"column:member_id;not null"`
	Game       string  `json:"game" gorm:"column:game;size:20;not null"`
	BetAmount  int     `json:"betAmount" gorm:"column:bet_amount;not null"`
	BetChoice  string  `json:"betChoice" gorm:"column:bet_choice;size:50;not null"`
	Result     string  `json:"result" gorm:"column:result;size:50;not null"`
	Multiplier float64 `json:"multiplier" gorm:"column:multiplier;not null;default:0"`
	Payout     int     `json:"payout" gorm:"column:payout;not null;default:0"`
	Profit     int     `json:"profit" gorm:"column:profit;not null;default:0"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
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
	Id         int64   `json:"id"`
	Game       string  `json:"game"`
	BetAmount  int     `json:"betAmount"`
	BetChoice  string  `json:"betChoice"`
	Result     string  `json:"result"`
	Multiplier float64 `json:"multiplier"`
	Payout     int     `json:"payout"`
	Profit     int     `json:"profit"`
	Balance    int     `json:"balance"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CasinoHistoryItem struct {
	Id         int64   `json:"id"`
	Game       string  `json:"game"`
	BetAmount  int     `json:"betAmount"`
	BetChoice  string  `json:"betChoice"`
	Result     string  `json:"result"`
	Multiplier float64 `json:"multiplier"`
	Payout     int     `json:"payout"`
	Profit     int     `json:"profit"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CasinoStats struct {
	Balance      int   `json:"balance"`
	TotalBets    int64 `json:"totalBets"`
	TotalWagered int   `json:"totalWagered"`
	TotalPayout  int   `json:"totalPayout"`
	BiggestWin   int   `json:"biggestWin"`
}

type CasinoAdminStats struct {
	TotalBets     int64 `json:"totalBets"`
	TotalWagered  int   `json:"totalWagered"`
	TotalPayout   int   `json:"totalPayout"`
	HouseProfit   int   `json:"houseProfit"`
	UniquePlayers int64 `json:"uniquePlayers"`
}

type CasinoFeedItem struct {
	Id              int64     `json:"id"`
	MemberFirstName string    `json:"memberFirstName"`
	MemberUsername  string    `json:"memberUsername"`
	Game            string    `json:"game"`
	BetAmount       int       `json:"betAmount"`
	BetChoice       string    `json:"betChoice"`
	Result          string    `json:"result"`
	Multiplier      float64   `json:"multiplier"`
	Payout          int       `json:"payout"`
	Profit          int       `json:"profit"`
	CreatedAt       time.Time `json:"createdAt"`
}

type CasinoAdminBet struct {
	Id              int64   `json:"id"`
	MemberId        int64   `json:"memberId"`
	MemberFirstName string  `json:"memberFirstName"`
	MemberLastName  string  `json:"memberLastName"`
	MemberUsername  string  `json:"memberUsername"`
	Game            string  `json:"game"`
	BetAmount       int     `json:"betAmount"`
	BetChoice       string  `json:"betChoice"`
	Result          string  `json:"result"`
	Multiplier      float64 `json:"multiplier"`
	Payout          int     `json:"payout"`
	Profit          int     `json:"profit"`
	CreatedAt       time.Time `json:"createdAt"`
}
