import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { kudosService } from '@/services/kudos'

describe('kudosService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getRecent', () => {
    it('should call GET kudos with default params', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      const result = await kudosService.getRecent()

      expect(mockApiClient.get).toHaveBeenCalledWith('kudos', { searchParams: { limit: 20, offset: 0 } })
      expect(result).toEqual(data)
    })

    it('should call GET kudos with custom params', async () => {
      const data = { items: [{ id: 1 }], total: 1 }
      mockJson.mockResolvedValue(data)

      const result = await kudosService.getRecent(10, 5)

      expect(mockApiClient.get).toHaveBeenCalledWith('kudos', { searchParams: { limit: 10, offset: 5 } })
      expect(result).toEqual(data)
    })
  })

  describe('send', () => {
    it('should call POST kudos with toId and message', async () => {
      const kudos = { id: 1, toId: 42, message: 'Great work!' }
      mockJson.mockResolvedValue(kudos)

      const result = await kudosService.send(42, 'Great work!')

      expect(mockApiClient.post).toHaveBeenCalledWith('kudos', { json: { toId: 42, message: 'Great work!' } })
      expect(result).toEqual(kudos)
    })
  })
})
