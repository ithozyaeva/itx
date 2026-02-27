import type { ReviewOnCommunity } from '@/models/reviewOnCommunity'
import api from '@/lib/api'
import { BaseService } from '@/services/api/baseService'
import { handleError } from './errorService'

class ReviewOnCommunityService extends BaseService<ReviewOnCommunity> {
  constructor() {
    super('reviews')
  }

  approve = async (id: number) => {
    try {
      this.isLoading.value = true
      return await api.post(`${this.basePath}/${id}/approve`).json()
    }
    catch (error) {
      handleError(error)
      return null
    }
    finally {
      this.isLoading.value = false
      this.search()
    }
  }
}

export const reviewOnCommunityService = new ReviewOnCommunityService()
