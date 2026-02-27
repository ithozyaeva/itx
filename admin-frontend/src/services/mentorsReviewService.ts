import type { ReviewOnService } from '@/models/mentors'
import { BaseService } from '@/services/api/baseService'

class MentorsReviewService extends BaseService<ReviewOnService> {
  constructor() {
    super('reviews-on-service')
  }
}

export const mentorsReviewService = new MentorsReviewService()
