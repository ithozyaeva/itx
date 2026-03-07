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

const { mentorsReviewService } = await import('@/services/mentorsReviewService')

describe('mentorsReviewService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "reviews-on-service"', () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    mentorsReviewService.search()

    expect(mockApi.get).toHaveBeenCalledWith('reviews-on-service', expect.any(Object))
  })

  describe('approve', () => {
    it('returns true on success', async () => {
      mockApi.post.mockResolvedValueOnce(undefined) // approve call (no .json())
      mockJson.mockResolvedValueOnce({ items: [], total: 0 }) // search refresh

      const result = await mentorsReviewService.approve(1)

      expect(mockApi.post).toHaveBeenCalledWith('reviews-on-service/1/approve')
      expect(result).toBe(true)
    })

    it('returns false on failure', async () => {
      mockApi.post.mockRejectedValueOnce(new Error('Approve failed'))

      const result = await mentorsReviewService.approve(1)

      expect(result).toBe(false)
    })

    it('sets isLoading to false after completion', async () => {
      mockApi.post.mockResolvedValueOnce(undefined)
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await mentorsReviewService.approve(1)

      expect(mentorsReviewService.isLoading.value).toBe(false)
    })

    it('calls search after successful approve', async () => {
      mockApi.post.mockResolvedValueOnce(undefined)
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await mentorsReviewService.approve(5)

      expect(mockApi.get).toHaveBeenCalledWith('reviews-on-service', expect.any(Object))
    })
  })
})
