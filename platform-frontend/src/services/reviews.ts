import { apiClient } from './api'
import { handleError } from './errorService'

export const reviewService = {
  async createReview(text: string) {
    try {
      await apiClient.post('reviews/add', { json: { text } })
    }
    catch (err) {
      handleError(err)
    }
  },
}
