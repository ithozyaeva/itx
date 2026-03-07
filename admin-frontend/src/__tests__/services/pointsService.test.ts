import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock useToast before importing the service
const mockToast = vi.fn()
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({ toast: mockToast }),
}))

// Mock handleError
const mockHandleError = vi.fn()
vi.mock('@/services/errorService', () => ({
  handleError: (...args: any[]) => mockHandleError(...args),
}))

// Mock api
const mockJson = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson })),
  post: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { pointsService } = await import('@/services/pointsService')

describe('pointsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    pointsService.items.value = { items: [], total: 0 }
    pointsService.isLoading.value = false
    pointsService.pagination.value = { limit: 20, offset: 0 }
    pointsService.filters.value = {}
  })

  describe('search', () => {
    it('fetches points with default pagination', async () => {
      const mockResponse = {
        items: [{ id: 1, amount: 100, reason: 'bonus' }],
        total: 1,
      }
      mockJson.mockResolvedValueOnce(mockResponse)

      await pointsService.search()

      expect(mockApi.get).toHaveBeenCalledWith('points', {
        searchParams: { limit: 20, offset: 0 },
      })
      expect(pointsService.items.value).toEqual(mockResponse)
    })

    it('includes filters in search params', async () => {
      pointsService.filters.value = { username: 'testuser' }
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await pointsService.search()

      expect(mockApi.get).toHaveBeenCalledWith('points', {
        searchParams: { limit: 20, offset: 0, username: 'testuser' },
      })
    })

    it('handles errors', async () => {
      const error = new Error('Search failed')
      mockJson.mockRejectedValueOnce(error)

      await pointsService.search()

      expect(mockHandleError).toHaveBeenCalledWith(error)
      expect(pointsService.isLoading.value).toBe(false)
    })
  })

  describe('changePagination', () => {
    it('updates offset based on page number', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      pointsService.changePagination(3)

      // page 3 with limit 20 => offset 40
      expect(pointsService.pagination.value.offset).toBe(40)
    })
  })

  describe('clearPagination', () => {
    it('resets offset to 0', () => {
      pointsService.pagination.value.offset = 100

      pointsService.clearPagination()

      expect(pointsService.pagination.value.offset).toBe(0)
    })
  })

  describe('applyFilters', () => {
    it('sets filters and resets offset', async () => {
      pointsService.pagination.value.offset = 40
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      pointsService.applyFilters({ username: 'newuser' })

      expect(pointsService.filters.value).toEqual({ username: 'newuser' })
      expect(pointsService.pagination.value.offset).toBe(0)
    })
  })

  describe('award', () => {
    it('awards points and shows success toast', async () => {
      const awardData = { memberId: 1, amount: 50, description: 'Great work' }
      mockJson
        .mockResolvedValueOnce({}) // award response
        .mockResolvedValueOnce({ items: [], total: 0 }) // search refresh

      const result = await pointsService.award(awardData)

      expect(result).toBe(true)
      expect(mockApi.post).toHaveBeenCalledWith('points', { json: awardData })
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Баллы успешно начислены',
      })
    })

    it('returns false on failure', async () => {
      const error = new Error('Award failed')
      mockJson.mockRejectedValueOnce(error)

      const result = await pointsService.award({ memberId: 1, amount: 50, description: '' })

      expect(result).toBe(false)
      expect(mockHandleError).toHaveBeenCalledWith(error)
    })
  })

  describe('deleteTransaction', () => {
    it('deletes transaction and shows success toast', async () => {
      mockApi.delete.mockResolvedValueOnce(undefined)
      mockJson.mockResolvedValueOnce({ items: [], total: 0 }) // search refresh

      const result = await pointsService.deleteTransaction(42)

      expect(result).toBe(true)
      expect(mockApi.delete).toHaveBeenCalledWith('points/42')
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Транзакция удалена',
      })
    })

    it('returns false on failure', async () => {
      const error = new Error('Delete failed')
      mockApi.delete.mockRejectedValueOnce(error)

      const result = await pointsService.deleteTransaction(42)

      expect(result).toBe(false)
      expect(mockHandleError).toHaveBeenCalledWith(error)
    })
  })
})
