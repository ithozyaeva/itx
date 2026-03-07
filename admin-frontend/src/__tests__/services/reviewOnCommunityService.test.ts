import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockToast = vi.fn()
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({ toast: mockToast }),
}))

const mockHandleError = vi.fn()
vi.mock('@/services/errorService', () => ({
  handleError: (...args: any[]) => mockHandleError(...args),
}))

const mockJson = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson })),
  post: vi.fn(() => ({ json: mockJson })),
  put: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { reviewOnCommunityService } = await import('@/services/reviewOnCommunityService')

describe('reviewOnCommunityService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "reviews"', () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    reviewOnCommunityService.search()

    expect(mockApi.get).toHaveBeenCalledWith('reviews', expect.any(Object))
  })

  describe('approve', () => {
    it('calls api.post with correct path on success', async () => {
      const approveResponse = { id: 1, approved: true }
      mockJson
        .mockResolvedValueOnce(approveResponse) // approve response
        .mockResolvedValueOnce({ items: [], total: 0 }) // search refresh in finally

      const result = await reviewOnCommunityService.approve(1)

      expect(mockApi.post).toHaveBeenCalledWith('reviews/1/approve')
      expect(result).toEqual(approveResponse)
    })

    it('calls handleError on failure and returns null', async () => {
      const error = new Error('Approve failed')
      mockApi.post.mockReturnValueOnce({ json: vi.fn().mockRejectedValueOnce(error) })
      mockJson.mockResolvedValueOnce({ items: [], total: 0 }) // search refresh in finally

      const result = await reviewOnCommunityService.approve(1)

      expect(mockHandleError).toHaveBeenCalledWith(error)
      expect(result).toBeNull()
    })

    it('sets isLoading during request and resets after', async () => {
      let loadingDuringRequest = false
      mockApi.post.mockReturnValueOnce({
        json: vi.fn().mockImplementation(() => {
          loadingDuringRequest = reviewOnCommunityService.isLoading.value
          return Promise.resolve({ approved: true })
        }),
      })
      mockJson.mockResolvedValueOnce({ items: [], total: 0 }) // search refresh

      await reviewOnCommunityService.approve(1)

      expect(loadingDuringRequest).toBe(true)
      expect(reviewOnCommunityService.isLoading.value).toBe(false)
    })

    it('calls search in finally block', async () => {
      mockJson
        .mockResolvedValueOnce({ approved: true }) // approve
        .mockResolvedValueOnce({ items: [], total: 0 }) // search

      await reviewOnCommunityService.approve(1)

      expect(mockApi.get).toHaveBeenCalledWith('reviews', expect.any(Object))
    })
  })
})
