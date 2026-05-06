package models

import "time"

type ReferralCreditReason string

const (
	CreditReasonReferalConversion         ReferralCreditReason = "referal_conversion"
	CreditReasonReferralPurchaseFirst     ReferralCreditReason = "referral_purchase_first"
	CreditReasonReferralPurchaseRecurring ReferralCreditReason = "referral_purchase_recurring"
	CreditReasonAdminManual               ReferralCreditReason = "admin_manual"
	CreditReasonSubscriptionPurchase      ReferralCreditReason = "subscription_purchase"
)

type ReferralCreditTransaction struct {
	Id          int64                `json:"id" gorm:"primaryKey"`
	MemberId    int64                `json:"memberId" gorm:"column:member_id;not null"`
	Amount      int                  `json:"amount" gorm:"not null"`
	Reason      ReferralCreditReason `json:"reason" gorm:"column:reason;size:50;not null"`
	SourceType  string               `json:"sourceType" gorm:"column:source_type;size:50;not null"`
	SourceId    int64                `json:"sourceId" gorm:"column:source_id;not null;default:0"`
	Description string               `json:"description" gorm:"column:description;default:''"`
	CreatedAt   time.Time            `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (ReferralCreditTransaction) TableName() string { return "referral_credit_transactions" }

type ReferralCreditSummary struct {
	Balance      int                         `json:"balance"`
	Transactions []ReferralCreditTransaction `json:"transactions"`
}

type AdminCreditTransaction struct {
	Id              int64                `json:"id"`
	MemberId        int64                `json:"memberId"`
	MemberFirstName string               `json:"memberFirstName"`
	MemberLastName  string               `json:"memberLastName"`
	MemberUsername  string               `json:"memberUsername"`
	Amount          int                  `json:"amount"`
	Reason          ReferralCreditReason `json:"reason"`
	SourceType      string               `json:"sourceType"`
	Description     string               `json:"description"`
	CreatedAt       time.Time            `json:"createdAt"`
}

type AdminAwardCreditsRequest struct {
	MemberId    int64  `json:"memberId"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}
