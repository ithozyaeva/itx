import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockJson = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson })),
  post: vi.fn(() => ({ json: mockJson })),
  put: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { bulkService } = await import('@/services/bulkService')

describe('bulkService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deleteEvents posts to correct endpoint', async () => {
    mockJson.mockResolvedValueOnce({ success: true })

    await bulkService.deleteEvents([1, 2, 3])

    expect(mockApi.post).toHaveBeenCalledWith('bulk/events/delete', { json: { ids: [1, 2, 3] } })
  })

  it('deleteMentors posts to correct endpoint', async () => {
    mockJson.mockResolvedValueOnce({ success: true })

    await bulkService.deleteMentors([4, 5])

    expect(mockApi.post).toHaveBeenCalledWith('bulk/mentors/delete', { json: { ids: [4, 5] } })
  })

  it('deleteMembers posts to correct endpoint', async () => {
    mockJson.mockResolvedValueOnce({ success: true })

    await bulkService.deleteMembers([6, 7])

    expect(mockApi.post).toHaveBeenCalledWith('bulk/members/delete', { json: { ids: [6, 7] } })
  })

  it('deleteReviews posts to correct endpoint', async () => {
    mockJson.mockResolvedValueOnce({ success: true })

    await bulkService.deleteReviews([8, 9])

    expect(mockApi.post).toHaveBeenCalledWith('bulk/reviews/delete', { json: { ids: [8, 9] } })
  })

  it('approveReviews posts to correct endpoint', async () => {
    mockJson.mockResolvedValueOnce({ success: true })

    await bulkService.approveReviews([10, 11])

    expect(mockApi.post).toHaveBeenCalledWith('bulk/reviews/approve', { json: { ids: [10, 11] } })
  })

  it('deleteMentorsReviews posts to correct endpoint', async () => {
    mockJson.mockResolvedValueOnce({ success: true })

    await bulkService.deleteMentorsReviews([12, 13])

    expect(mockApi.post).toHaveBeenCalledWith('bulk/mentors-reviews/delete', { json: { ids: [12, 13] } })
  })

  it('approveMentorsReviews posts to correct endpoint', async () => {
    mockJson.mockResolvedValueOnce({ success: true })

    await bulkService.approveMentorsReviews([14, 15])

    expect(mockApi.post).toHaveBeenCalledWith('bulk/mentors-reviews/approve', { json: { ids: [14, 15] } })
  })
})
