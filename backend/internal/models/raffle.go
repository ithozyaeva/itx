package models

import (
	"ithozyeva/internal/s3resolve"
	"time"

	"gorm.io/gorm"
)

type RaffleStatus string

const (
	RaffleStatusActive   RaffleStatus = "ACTIVE"
	RaffleStatusFinished RaffleStatus = "FINISHED"
)

type Raffle struct {
	Id          int64        `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title" gorm:"column:title;not null"`
	Description string       `json:"description" gorm:"column:description;default:''"`
	Prize       string       `json:"prize" gorm:"column:prize;not null"`
	TicketCost  int          `json:"ticketCost" gorm:"column:ticket_cost;not null;default:10"`
	MaxTickets  int          `json:"maxTickets" gorm:"column:max_tickets;default:0"`
	EndsAt      time.Time    `json:"endsAt" gorm:"column:ends_at;not null"`
	Status      RaffleStatus `json:"status" gorm:"column:status;default:'ACTIVE'"`
	WinnerId    *int64       `json:"winnerId" gorm:"column:winner_id"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

type RaffleTicket struct {
	Id       int64     `json:"id" gorm:"primaryKey"`
	RaffleId int64     `json:"raffleId" gorm:"column:raffle_id;not null"`
	MemberId int64     `json:"memberId" gorm:"column:member_id;not null"`
	BoughtAt time.Time `json:"boughtAt" gorm:"column:bought_at;autoCreateTime"`
}

type RafflePublic struct {
	Id              int64        `json:"id"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Prize           string       `json:"prize"`
	TicketCost      int          `json:"ticketCost"`
	MaxTickets      int          `json:"maxTickets"`
	EndsAt          time.Time    `json:"endsAt"`
	Status          RaffleStatus `json:"status"`
	TotalTickets    int          `json:"totalTickets"`
	MyTickets       int          `json:"myTickets"`
	WinnerId        *int64       `json:"winnerId"`
	WinnerFirstName string       `json:"winnerFirstName,omitempty"`
	WinnerLastName  string       `json:"winnerLastName,omitempty"`
	WinnerUsername  string       `json:"winnerUsername,omitempty"`
	WinnerAvatarURL string       `json:"winnerAvatarUrl,omitempty"`
}

func (r *RafflePublic) AfterFind(tx *gorm.DB) (err error) {
	r.WinnerAvatarURL = s3resolve.ResolveS3URL(r.WinnerAvatarURL)
	return nil
}

type BuyTicketRequest struct {
	Count int `json:"count"`
}
