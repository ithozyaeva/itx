import type { Feedback } from '@/models/feedback'
import { BaseService } from '@/services/api/baseService'

class FeedbackService extends BaseService<Feedback> {
  constructor() {
    super('feedback')
  }
}

export const feedbackService = new FeedbackService()
