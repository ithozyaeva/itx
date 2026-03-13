import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
      patch: vi.fn(() => ({ json: mockJson })),
      delete: vi.fn(() => Promise.resolve()),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { reviewService } from '@/services/reviews'

describe('reviewService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('createReview', () => {
    it('should call POST reviews/add with text', async () => {
      mockApiClient.post.mockReturnValue({ json: mockJson })

      await reviewService.createReview('Great community!')

      expect(mockApiClient.post).toHaveBeenCalledWith('reviews/add', { json: { text: 'Great community!' } })
    })

    it('should throw on failure', async () => {
      const error = new Error('Network error')
      mockApiClient.post.mockImplementation(() => { throw error })

      await expect(reviewService.createReview('Great community!')).rejects.toThrow('Network error')
    })
  })

  describe('getMyReviews', () => {
    it('should call GET reviews/my', async () => {
      const reviews = [{ id: 1, text: 'Nice' }]
      mockJson.mockResolvedValue(reviews)

      const result = await reviewService.getMyReviews()

      expect(mockApiClient.get).toHaveBeenCalledWith('reviews/my')
      expect(result).toEqual(reviews)
    })
  })

  describe('updateReview', () => {
    it('should call PATCH reviews/:id with text', async () => {
      const updated = { id: 5, text: 'Updated review' }
      mockJson.mockResolvedValue(updated)
      mockApiClient.patch.mockReturnValue({ json: mockJson })

      const result = await reviewService.updateReview(5, 'Updated review')

      expect(mockApiClient.patch).toHaveBeenCalledWith('reviews/5', { json: { text: 'Updated review' } })
      expect(result).toEqual(updated)
    })
  })

  describe('deleteReview', () => {
    it('should call DELETE reviews/:id', async () => {
      mockApiClient.delete.mockResolvedValue(undefined)

      await reviewService.deleteReview(3)

      expect(mockApiClient.delete).toHaveBeenCalledWith('reviews/3')
    })
  })
})
