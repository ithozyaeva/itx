import type { CreateMarketplaceItemRequest, MarketplaceItem, MarketplaceSearchResponse } from '@/models/marketplace'
import { apiClient } from './api'

export const marketplaceService = {
  async getAll(params?: { status?: string, limit?: number, offset?: number }) {
    const searchParams = new URLSearchParams()
    if (params?.status)
      searchParams.set('status', params.status)
    if (params?.limit)
      searchParams.set('limit', params.limit.toString())
    if (params?.offset)
      searchParams.set('offset', params.offset.toString())

    return apiClient.get('marketplace', { searchParams }).json<MarketplaceSearchResponse>()
  },

  async getById(id: number) {
    return apiClient.get(`marketplace/${id}`).json<MarketplaceItem>()
  },

  async create(data: CreateMarketplaceItemRequest, image?: File) {
    const formData = new FormData()
    formData.append('title', data.title)
    formData.append('description', data.description)
    formData.append('price', data.price)
    formData.append('city', data.city)
    formData.append('canShip', data.canShip ? 'true' : 'false')
    formData.append('condition', data.condition)
    formData.append('defects', data.defects)
    formData.append('packageContents', data.packageContents)
    formData.append('contactTelegram', data.contactTelegram)
    formData.append('contactEmail', data.contactEmail)
    formData.append('contactPhone', data.contactPhone)
    if (image) {
      formData.append('image', image)
    }
    return apiClient.post('marketplace', { body: formData }).json<MarketplaceItem>()
  },

  async update(id: number, data: Partial<CreateMarketplaceItemRequest>) {
    return apiClient.patch(`marketplace/${id}`, { json: data }).json<MarketplaceItem>()
  },

  async requestPurchase(id: number) {
    return apiClient.post(`marketplace/${id}/request-purchase`).json<MarketplaceItem>()
  },

  async cancelPurchase(id: number) {
    return apiClient.post(`marketplace/${id}/cancel-purchase`).json<MarketplaceItem>()
  },

  async markSold(id: number) {
    return apiClient.post(`marketplace/${id}/sold`).json<MarketplaceItem>()
  },

  async remove(id: number) {
    return apiClient.delete(`marketplace/${id}`)
  },
}
