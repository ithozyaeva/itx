import api from '@/lib/api'

export const bulkService = {
  deleteEvents: (ids: number[]) => api.post('bulk/events/delete', { json: { ids } }).json(),
  deleteMentors: (ids: number[]) => api.post('bulk/mentors/delete', { json: { ids } }).json(),
  deleteMembers: (ids: number[]) => api.post('bulk/members/delete', { json: { ids } }).json(),
  deleteReviews: (ids: number[]) => api.post('bulk/reviews/delete', { json: { ids } }).json(),
  approveReviews: (ids: number[]) => api.post('bulk/reviews/approve', { json: { ids } }).json(),
  deleteMentorsReviews: (ids: number[]) => api.post('bulk/mentors-reviews/delete', { json: { ids } }).json(),
}
