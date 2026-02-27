import { apiClient } from './api'

export const reviewService = {
  async createReview(text: string) {
    try {
      await apiClient.post('reviews/add', { json: { text } })
    }
    catch (err) {
      console.error('Произошла ошибка', err)
    }
  },
}
