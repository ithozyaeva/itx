import { apiClient } from './api'
import { handleError } from './errorService'

export const feedbackService = {
  async submit(score: number, comment?: string) {
    try {
      await apiClient.post('feedback', {
        json: {
          score,
          comment: comment?.trim() || undefined,
        },
      })
    }
    catch (error) {
      handleError(error)
      throw error
    }
  },
}
