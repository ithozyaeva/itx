import type { ReviewOnService } from '@/models/mentors'
import api from '@/lib/api'
import { BaseService } from '@/services/api/baseService'

class MentorsReviewService extends BaseService<ReviewOnService> {
  constructor() {
    super('reviews-on-service')
  }

  approve = async (id: number): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post(`reviews-on-service/${id}/approve`)
      await this.search()
      return true
    }
    catch {
      return false
    }
    finally {
      this.isLoading.value = false
    }
  }
}

export const mentorsReviewService = new MentorsReviewService()
