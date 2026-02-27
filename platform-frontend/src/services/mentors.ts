import type { Mentor, Service } from '@/models/profile'
import ky from 'ky'
import { apiClient } from './api'

export interface ReviewOnService {
  id: number
  serviceId: number
  service: Service
  author: string
  text: string
  date: string
}

export interface MentorWithReviews extends Mentor {
  reviews?: ReviewOnService[]
}

export const mentorsService = {
  getAll: async (limit: number = 100, offset: number = 0) => {
    return ky.get('/api/mentors', {
      searchParams: { limit, offset },
    }).json<{ items: Mentor[], total: number }>()
  },

  getById: async (id: number) => {
    return apiClient.get(`mentors/${id}`).json<MentorWithReviews>()
  },

  getServices: async (mentorId: number) => {
    return ky.get(`/api/mentors/${mentorId}/services`).json<Service[]>()
  },

  addReview: async (mentorId: number, serviceId: number, text: string) => {
    return apiClient.post(`mentors/${mentorId}/reviews`, {
      json: { serviceId, text },
    }).json<ReviewOnService>()
  },
}
