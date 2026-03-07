import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockJson = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson })),
  post: vi.fn(() => ({ json: mockJson })),
  put: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { statsService } = await import('@/services/statsService')

describe('statsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getStats', () => {
    it('calls api.get with "stats" and returns dashboard stats', async () => {
      const mockStats = {
        totalMembers: 100,
        totalMentors: 10,
        upcomingEvents: 5,
        pastEvents: 20,
        pendingReviews: 3,
        approvedReviews: 15,
        referralLinks: 8,
        resumes: 12,
        openTasks: 4,
        inProgressTasks: 2,
        doneTasks: 6,
        approvedTasks: 5,
      }
      mockJson.mockResolvedValueOnce(mockStats)

      const result = await statsService.getStats()

      expect(mockApi.get).toHaveBeenCalledWith('stats')
      expect(result).toEqual(mockStats)
    })
  })

  describe('getChartStats', () => {
    it('calls api.get with "stats/charts" and returns chart stats', async () => {
      const mockChartStats = {
        memberGrowth: [{ month: '2026-01', count: 10 }],
        eventAttendance: [{ month: '2026-01', count: 25 }],
      }
      mockJson.mockResolvedValueOnce(mockChartStats)

      const result = await statsService.getChartStats()

      expect(mockApi.get).toHaveBeenCalledWith('stats/charts')
      expect(result).toEqual(mockChartStats)
    })
  })
})
