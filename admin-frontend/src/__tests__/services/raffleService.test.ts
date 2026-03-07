import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock useToast before importing the service (used at module scope)
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
  put: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

// Import after mocks
const { raffleService } = await import('@/services/raffleService')

describe('raffleService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    raffleService.items.value = []
    raffleService.isLoading.value = false
  })

  describe('getAll', () => {
    it('fetches all raffles and updates items', async () => {
      const mockRaffles = [
        { id: 1, title: 'Raffle 1', status: 'ACTIVE' },
        { id: 2, title: 'Raffle 2', status: 'FINISHED' },
      ]
      mockJson.mockResolvedValueOnce(mockRaffles)

      await raffleService.getAll()

      expect(mockApi.get).toHaveBeenCalledWith('raffles')
      expect(raffleService.items.value).toEqual(mockRaffles)
      expect(raffleService.isLoading.value).toBe(false)
    })

    it('sets items to empty array when response is null', async () => {
      mockJson.mockResolvedValueOnce(null)

      await raffleService.getAll()

      expect(raffleService.items.value).toEqual([])
    })

    it('handles errors via handleError', async () => {
      const error = new Error('Network error')
      mockJson.mockRejectedValueOnce(error)

      await raffleService.getAll()

      expect(mockHandleError).toHaveBeenCalledWith(error)
      expect(raffleService.isLoading.value).toBe(false)
    })

    it('sets isLoading to true during request', async () => {
      let loadingDuringRequest = false
      mockJson.mockImplementationOnce(() => {
        loadingDuringRequest = raffleService.isLoading.value
        return Promise.resolve([])
      })

      await raffleService.getAll()

      expect(loadingDuringRequest).toBe(true)
      expect(raffleService.isLoading.value).toBe(false)
    })
  })

  describe('create', () => {
    it('creates a raffle and shows success toast', async () => {
      const createData = {
        title: 'New Raffle',
        description: 'Desc',
        prize: 'A prize',
        ticketCost: 10,
        maxTickets: 100,
        endsAt: '2026-05-01T00:00:00Z',
      }
      mockJson.mockResolvedValueOnce({ id: 1, ...createData }) // create response
      mockJson.mockResolvedValueOnce([]) // getAll after create

      const result = await raffleService.create(createData)

      expect(result).toBe(true)
      expect(mockApi.post).toHaveBeenCalledWith('raffles', { json: createData })
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Розыгрыш создан',
      })
    })

    it('returns false and calls handleError on failure', async () => {
      const error = new Error('Create failed')
      mockJson.mockRejectedValueOnce(error)

      const result = await raffleService.create({
        title: 'Test',
        description: '',
        prize: '',
        ticketCost: 0,
        maxTickets: 0,
        endsAt: '',
      })

      expect(result).toBe(false)
      expect(mockHandleError).toHaveBeenCalledWith(error)
    })

    it('refreshes list after successful create', async () => {
      mockJson.mockResolvedValueOnce({}) // create
      mockJson.mockResolvedValueOnce([{ id: 1 }]) // getAll

      await raffleService.create({
        title: 'Test',
        description: '',
        prize: '',
        ticketCost: 0,
        maxTickets: 0,
        endsAt: '',
      })

      // getAll is called after create, so get is called once
      expect(mockApi.get).toHaveBeenCalledWith('raffles')
    })
  })

  describe('delete', () => {
    it('deletes a raffle by id and shows success toast', async () => {
      mockApi.delete.mockResolvedValueOnce(undefined)
      mockJson.mockResolvedValueOnce([]) // getAll after delete

      const result = await raffleService.delete(5)

      expect(result).toBe(true)
      expect(mockApi.delete).toHaveBeenCalledWith('raffles/5')
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Розыгрыш удалён',
      })
    })

    it('returns false and calls handleError on failure', async () => {
      const error = new Error('Delete failed')
      mockApi.delete.mockRejectedValueOnce(error)

      const result = await raffleService.delete(5)

      expect(result).toBe(false)
      expect(mockHandleError).toHaveBeenCalledWith(error)
    })

    it('resets isLoading after delete completes', async () => {
      mockApi.delete.mockResolvedValueOnce(undefined)
      mockJson.mockResolvedValueOnce([])

      await raffleService.delete(1)

      expect(raffleService.isLoading.value).toBe(false)
    })
  })
})
