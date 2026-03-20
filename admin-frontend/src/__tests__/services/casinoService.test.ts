import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockHandleError = vi.fn()
vi.mock('@/services/errorService', () => ({
  handleError: (...args: any[]) => mockHandleError(...args),
}))

const mockJson = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson })),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { casinoService } = await import('@/services/casinoService')

describe('casinoService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    casinoService.items.value = { items: [], total: 0 }
    casinoService.stats.value = null
    casinoService.isLoading.value = false
    casinoService.pagination.value = { limit: 20, offset: 0 }
    casinoService.filters.value = {}
  })

  describe('getStats', () => {
    it('fetches and stores admin stats', async () => {
      const mockStats = {
        totalBets: 100,
        totalWagered: 5000,
        totalPayout: 4500,
        houseProfit: 500,
        uniquePlayers: 15,
        gameStats: [],
      }
      mockJson.mockResolvedValueOnce(mockStats)

      await casinoService.getStats()

      expect(mockApi.get).toHaveBeenCalledWith('minigames/stats')
      expect(casinoService.stats.value).toEqual(mockStats)
    })

    it('handles errors', async () => {
      const error = new Error('Stats fetch failed')
      mockJson.mockRejectedValueOnce(error)

      await casinoService.getStats()

      expect(mockHandleError).toHaveBeenCalledWith(error)
    })
  })

  describe('searchBets', () => {
    it('fetches bets with default pagination', async () => {
      const mockResponse = {
        items: [{ id: 1, game: 'coin_flip', betAmount: 50 }],
        total: 1,
      }
      mockJson.mockResolvedValueOnce(mockResponse)

      await casinoService.searchBets()

      expect(mockApi.get).toHaveBeenCalledWith('minigames/bets', {
        searchParams: { limit: 20, offset: 0 },
      })
      expect(casinoService.items.value).toEqual(mockResponse)
    })

    it('includes filters in search params', async () => {
      casinoService.filters.value = { username: 'testuser', game: 'wheel' }
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await casinoService.searchBets()

      expect(mockApi.get).toHaveBeenCalledWith('minigames/bets', {
        searchParams: { limit: 20, offset: 0, username: 'testuser', game: 'wheel' },
      })
    })

    it('sets isLoading during fetch', async () => {
      let loadingDuringFetch = false
      mockJson.mockImplementationOnce(() => {
        loadingDuringFetch = casinoService.isLoading.value
        return Promise.resolve({ items: [], total: 0 })
      })

      await casinoService.searchBets()

      expect(loadingDuringFetch).toBe(true)
      expect(casinoService.isLoading.value).toBe(false)
    })

    it('handles errors and resets isLoading', async () => {
      const error = new Error('Search failed')
      mockJson.mockRejectedValueOnce(error)

      await casinoService.searchBets()

      expect(mockHandleError).toHaveBeenCalledWith(error)
      expect(casinoService.isLoading.value).toBe(false)
    })
  })

  describe('changePagination', () => {
    it('updates offset based on page number', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      casinoService.changePagination(3)

      expect(casinoService.pagination.value.offset).toBe(40)
    })

    it('page 1 sets offset to 0', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      casinoService.changePagination(1)

      expect(casinoService.pagination.value.offset).toBe(0)
    })
  })

  describe('clearPagination', () => {
    it('resets offset to 0', () => {
      casinoService.pagination.value.offset = 100

      casinoService.clearPagination()

      expect(casinoService.pagination.value.offset).toBe(0)
    })
  })

  describe('applyFilters', () => {
    it('sets filters and resets offset', async () => {
      casinoService.pagination.value.offset = 40
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      casinoService.applyFilters({ username: 'player1', game: 'dice_roll' })

      expect(casinoService.filters.value).toEqual({ username: 'player1', game: 'dice_roll' })
      expect(casinoService.pagination.value.offset).toBe(0)
    })

    it('clears filters when empty object passed', async () => {
      casinoService.filters.value = { username: 'old' }
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      casinoService.applyFilters({})

      expect(casinoService.filters.value).toEqual({})
    })
  })
})
