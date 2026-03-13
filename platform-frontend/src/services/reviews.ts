import { apiClient } from './api'

export interface ReviewOnCommunity {
  id: number
  authorId: number
  text: string
  date: string
  status: 'DRAFT' | 'APPROVED'
}

export const reviewService = {
  async createReview(text: string) {
    await apiClient.post('reviews/add', { json: { text } })
  },

  async getMyReviews() {
    return apiClient.get('reviews/my').json<ReviewOnCommunity[]>()
  },

  async updateReview(id: number, text: string) {
    return apiClient.patch(`reviews/${id}`, { json: { text } }).json<ReviewOnCommunity>()
  },

  async deleteReview(id: number) {
    await apiClient.delete(`reviews/${id}`)
  },
}
