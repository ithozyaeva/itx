import type { TelegramUser } from '@/models/profile'

export type MarketplaceItemStatus = 'ACTIVE' | 'RESERVED' | 'SOLD' | 'ARCHIVED'
export type MarketplaceItemCondition = 'NEW' | 'USED'

export interface MarketplaceItem {
  id: number
  title: string
  description: string
  price: string
  city: string
  canShip: boolean
  condition: MarketplaceItemCondition
  defects: string
  packageContents: string
  contactTelegram: string
  contactEmail: string
  contactPhone: string
  imagePath: string
  sellerId: number
  seller: TelegramUser
  buyerId: number | null
  buyer: TelegramUser | null
  status: MarketplaceItemStatus
  createdAt: string
  updatedAt: string
}

export interface MarketplaceSearchResponse {
  items: MarketplaceItem[]
  total: number
}

export interface CreateMarketplaceItemRequest {
  title: string
  description: string
  price: string
  city: string
  canShip: boolean
  condition: MarketplaceItemCondition
  defects: string
  packageContents: string
  contactTelegram: string
  contactEmail: string
  contactPhone: string
}
