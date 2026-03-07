import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient, mockKyGet } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
    },
    mockKyGet: vi.fn(),
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

vi.mock('ky', () => ({
  default: {
    get: mockKyGet,
  },
}))

import { mentorsService } from '@/services/mentors'

describe('mentorsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('should call GET /api/mentors with default params using raw ky', async () => {
      const data = { items: [], total: 0 }
      const mockKyJson = vi.fn().mockResolvedValue(data)
      mockKyGet.mockReturnValue({ json: mockKyJson })

      const result = await mentorsService.getAll()

      expect(mockKyGet).toHaveBeenCalledWith('/api/mentors', { searchParams: { limit: 100, offset: 0 } })
      expect(result).toEqual(data)
    })

    it('should call GET /api/mentors with custom params', async () => {
      const data = { items: [{ id: 1 }], total: 1 }
      const mockKyJson = vi.fn().mockResolvedValue(data)
      mockKyGet.mockReturnValue({ json: mockKyJson })

      const result = await mentorsService.getAll(50, 10)

      expect(mockKyGet).toHaveBeenCalledWith('/api/mentors', { searchParams: { limit: 50, offset: 10 } })
      expect(result).toEqual(data)
    })
  })

  describe('getById', () => {
    it('should call GET mentors/:id via apiClient', async () => {
      const mentor = { id: 5, name: 'John' }
      mockJson.mockResolvedValue(mentor)

      const result = await mentorsService.getById(5)

      expect(mockApiClient.get).toHaveBeenCalledWith('mentors/5')
      expect(result).toEqual(mentor)
    })
  })

  describe('getServices', () => {
    it('should call GET /api/mentors/:id/services using raw ky', async () => {
      const services = [{ id: 1, name: 'Consulting' }]
      const mockKyJson = vi.fn().mockResolvedValue(services)
      mockKyGet.mockReturnValue({ json: mockKyJson })

      const result = await mentorsService.getServices(7)

      expect(mockKyGet).toHaveBeenCalledWith('/api/mentors/7/services')
      expect(result).toEqual(services)
    })
  })

  describe('addReview', () => {
    it('should call POST mentors/:id/reviews via apiClient', async () => {
      const review = { id: 1, serviceId: 2, text: 'Excellent' }
      mockJson.mockResolvedValue(review)

      const result = await mentorsService.addReview(3, 2, 'Excellent')

      expect(mockApiClient.post).toHaveBeenCalledWith('mentors/3/reviews', {
        json: { serviceId: 2, text: 'Excellent' },
      })
      expect(result).toEqual(review)
    })
  })
})
