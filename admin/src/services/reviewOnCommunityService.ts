import type { ReviewOnCommunity } from '@/models/reviewOnCommunity'
import { BaseService } from '@/services/api/baseService'

class ReviewOnCommunityService extends BaseService<ReviewOnCommunity> {
  constructor() {
    super('reviews')
  }
}

export const reviewOnCommunityService = new ReviewOnCommunityService()
