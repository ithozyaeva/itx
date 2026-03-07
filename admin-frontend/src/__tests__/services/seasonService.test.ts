import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock useToast
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
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { seasonService } = await import('@/services/seasonService')

describe('seasonService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    seasonService.items.value = []
    seasonService.isLoading.value = false
  })

  describe('getAll', () => {
    it('fetches all seasons', async () => {
      const mockSeasons = [
        { id: 1, title: 'Season 1', status: 'ACTIVE' },
        { id: 2, title: 'Season 2', status: 'FINISHED' },
      ]
      mockJson.mockResolvedValueOnce(mockSeasons)

      await seasonService.getAll()

      expect(mockApi.get).toHaveBeenCalledWith('seasons')
      expect(seasonService.items.value).toEqual(mockSeasons)
    })

    it('defaults to empty array when response is null', async () => {
      mockJson.mockResolvedValueOnce(null)

      await seasonService.getAll()

      expect(seasonService.items.value).toEqual([])
    })

    it('handles errors', async () => {
      const error = new Error('Fetch failed')
      mockJson.mockRejectedValueOnce(error)

      await seasonService.getAll()

      expect(mockHandleError).toHaveBeenCalledWith(error)
      expect(seasonService.isLoading.value).toBe(false)
    })
  })

  describe('create', () => {
    it('creates a season and shows toast', async () => {
      const createData = {
        title: 'Spring 2026',
        startDate: '2026-03-01',
        endDate: '2026-06-01',
      }
      mockJson
        .mockResolvedValueOnce({ id: 1, ...createData }) // create
        .mockResolvedValueOnce([{ id: 1, ...createData }]) // getAll refresh

      const result = await seasonService.create(createData)

      expect(result).toBe(true)
      expect(mockApi.post).toHaveBeenCalledWith('seasons', { json: createData })
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Сезон создан',
      })
    })

    it('returns false on error', async () => {
      const error = new Error('Create failed')
      mockJson.mockRejectedValueOnce(error)

      const result = await seasonService.create({
        title: 'Test',
        startDate: '',
        endDate: '',
      })

      expect(result).toBe(false)
      expect(mockHandleError).toHaveBeenCalledWith(error)
    })
  })

  describe('finish', () => {
    it('finishes a season and shows toast', async () => {
      // finish uses api.post without .json()
      mockApi.post.mockResolvedValueOnce(undefined)
      mockJson.mockResolvedValueOnce([]) // getAll refresh

      // Need to handle the case where finish calls api.post without .json()
      // Looking at the source: await api.post(`seasons/${id}/finish`) — no .json()
      const result = await seasonService.finish(3)

      expect(result).toBe(true)
      expect(mockApi.post).toHaveBeenCalledWith('seasons/3/finish')
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Сезон завершён',
      })
    })

    it('returns false on error', async () => {
      const error = new Error('Finish failed')
      mockApi.post.mockRejectedValueOnce(error)

      const result = await seasonService.finish(3)

      expect(result).toBe(false)
      expect(mockHandleError).toHaveBeenCalledWith(error)
    })
  })
})
