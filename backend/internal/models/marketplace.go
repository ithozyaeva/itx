package models

import (
	"time"

	"ithozyeva/internal/s3resolve"

	"gorm.io/gorm"
)

type MarketplaceItemStatus string
type MarketplaceItemCondition string

const (
	MarketplaceStatusActive   MarketplaceItemStatus = "ACTIVE"
	MarketplaceStatusReserved MarketplaceItemStatus = "RESERVED"
	MarketplaceStatusSold     MarketplaceItemStatus = "SOLD"
	MarketplaceStatusArchived MarketplaceItemStatus = "ARCHIVED"

	MarketplaceConditionNew  MarketplaceItemCondition = "NEW"
	MarketplaceConditionUsed MarketplaceItemCondition = "USED"
)

type MarketplaceItem struct {
	Id              int64                    `json:"id" gorm:"primaryKey"`
	Title           string                   `json:"title" gorm:"column:title;size:255;not null"`
	Description     string                   `json:"description" gorm:"column:description;default:''"`
	Price           string                   `json:"price" gorm:"column:price;size:100;default:''"`
	City            string                   `json:"city" gorm:"column:city;size:255;default:''"`
	CanShip         bool                     `json:"canShip" gorm:"column:can_ship;default:false"`
	Condition       MarketplaceItemCondition `json:"condition" gorm:"column:condition;size:20;default:'USED'"`
	Defects         string                   `json:"defects" gorm:"column:defects;default:''"`
	PackageContents string                   `json:"packageContents" gorm:"column:package_contents;default:''"`
	ContactTelegram string                   `json:"contactTelegram" gorm:"column:contact_telegram;size:255;default:''"`
	ContactEmail    string                   `json:"contactEmail" gorm:"column:contact_email;size:255;default:''"`
	ContactPhone    string                   `json:"contactPhone" gorm:"column:contact_phone;size:255;default:''"`
	ImagePath       string                   `json:"imagePath" gorm:"column:image_path;default:''"`
	SellerId        int64                    `json:"sellerId" gorm:"column:seller_id;not null"`
	Seller          Member                   `json:"seller" gorm:"foreignKey:SellerId"`
	BuyerId         *int64                   `json:"buyerId" gorm:"column:buyer_id"`
	Buyer           *Member                  `json:"buyer" gorm:"foreignKey:BuyerId"`
	Status          MarketplaceItemStatus    `json:"status" gorm:"column:status;size:20;default:'ACTIVE'"`
	CreatedAt       time.Time                `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time                `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (MarketplaceItem) TableName() string {
	return "marketplace_items"
}

func (m *MarketplaceItem) AfterFind(tx *gorm.DB) (err error) {
	m.ImagePath = s3resolve.ResolveS3URL(m.ImagePath)
	return nil
}

type CreateMarketplaceItemRequest struct {
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	Price           string                   `json:"price"`
	City            string                   `json:"city"`
	CanShip         bool                     `json:"canShip"`
	Condition       MarketplaceItemCondition `json:"condition"`
	Defects         string                   `json:"defects"`
	PackageContents string                   `json:"packageContents"`
	ContactTelegram string                   `json:"contactTelegram"`
	ContactEmail    string                   `json:"contactEmail"`
	ContactPhone    string                   `json:"contactPhone"`
}
