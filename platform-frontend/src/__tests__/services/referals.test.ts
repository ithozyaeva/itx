import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
      put: vi.fn(() => ({ json: mockJson })),
      delete: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { referalLinkService } from '@/services/referals'

describe('referalLinkService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('search', () => {
    it('should call GET referals with limit and offset', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      const result = await referalLinkService.search(10, 0)

      expect(mockApiClient.get).toHaveBeenCalledWith('referals', {
        searchParams: { limit: 10, offset: 0 },
      })
      expect(result).toEqual(data)
    })

    it('should include grade and company filters', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await referalLinkService.search(10, 0, { grade: 'Senior', company: 'Google' })

      expect(mockApiClient.get).toHaveBeenCalledWith('referals', {
        searchParams: { limit: 10, offset: 0, grade: 'Senior', company: 'Google' },
      })
    })

    it('should include status filter', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await referalLinkService.search(10, 0, { status: 'active' })

      expect(mockApiClient.get).toHaveBeenCalledWith('referals', {
        searchParams: { limit: 10, offset: 0, status: 'active' },
      })
    })

    it('should join profTagIds with comma', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await referalLinkService.search(10, 0, { profTagIds: [1, 2, 3] })

      expect(mockApiClient.get).toHaveBeenCalledWith('referals', {
        searchParams: { limit: 10, offset: 0, profTagIds: '1,2,3' },
      })
    })

    it('should exclude empty profTagIds array', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await referalLinkService.search(10, 0, { profTagIds: [] })

      expect(mockApiClient.get).toHaveBeenCalledWith('referals', {
        searchParams: { limit: 10, offset: 0 },
      })
    })

    it('should exclude undefined filter values', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await referalLinkService.search(10, 0, { grade: undefined, company: undefined })

      expect(mockApiClient.get).toHaveBeenCalledWith('referals', {
        searchParams: { limit: 10, offset: 0 },
      })
    })
  })

  describe('addLink', () => {
    it('should call POST referals/add-link with link data', async () => {
      const link = { url: 'https://example.com', company: 'ACME' }
      const responseData = { id: 1, ...link }
      mockJson.mockResolvedValue(responseData)

      const result = await referalLinkService.addLink(link as any)

      expect(mockApiClient.post).toHaveBeenCalledWith('referals/add-link', { json: link })
      expect(result).toEqual(responseData)
    })
  })

  describe('updateLink', () => {
    it('should call PUT referals/update-link with link data', async () => {
      const link = { id: 1, url: 'https://updated.com' }
      const responseData = { ...link, company: 'ACME' }
      mockJson.mockResolvedValue(responseData)

      const result = await referalLinkService.updateLink(link as any)

      expect(mockApiClient.put).toHaveBeenCalledWith('referals/update-link', { json: link })
      expect(result).toEqual(responseData)
    })
  })

  describe('deleteLink', () => {
    it('should call DELETE referals/delete-link with id in body', async () => {
      const responseData = { id: 5 }
      mockJson.mockResolvedValue(responseData)

      const result = await referalLinkService.deleteLink(5)

      expect(mockApiClient.delete).toHaveBeenCalledWith('referals/delete-link', { json: { id: 5 } })
      expect(result).toEqual(responseData)
    })
  })

  describe('trackConversion', () => {
    it('should call POST referals/track-conversion with referralLinkId', async () => {
      mockApiClient.post.mockResolvedValue(undefined)

      await referalLinkService.trackConversion(42)

      expect(mockApiClient.post).toHaveBeenCalledWith('referals/track-conversion', { json: { referralLinkId: 42 } })
    })
  })
})
