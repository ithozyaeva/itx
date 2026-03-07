import { beforeEach, describe, expect, it, vi } from 'vitest'

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

const { auditLogService } = await import('@/services/auditLogService')

describe('auditLogService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    auditLogService.pagination.value = { limit: 20, offset: 0 }
    auditLogService.filters.value = {}
    auditLogService.items.value = { items: [], total: 0 }
  })

  describe('search', () => {
    it('calls api.get with correct path and params', async () => {
      const mockResponse = { items: [{ id: 1, action: 'create' }], total: 1 }
      mockJson.mockResolvedValueOnce(mockResponse)

      await auditLogService.search()

      expect(mockApi.get).toHaveBeenCalledWith('audit-logs', {
        searchParams: { limit: 20, offset: 0 },
      })
      expect(auditLogService.items.value).toEqual(mockResponse)
    })

    it('includes filters in search params', async () => {
      auditLogService.filters.value = { action: 'create', entityType: 'member' }
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await auditLogService.search()

      expect(mockApi.get).toHaveBeenCalledWith('audit-logs', {
        searchParams: { limit: 20, offset: 0, action: 'create', entityType: 'member' },
      })
    })

    it('handles errors with handleError', async () => {
      const error = new Error('Network error')
      mockJson.mockRejectedValueOnce(error)

      await auditLogService.search()

      expect(mockHandleError).toHaveBeenCalledWith(error)
    })

    it('sets and resets isLoading', async () => {
      let loadingDuringRequest = false
      mockJson.mockImplementationOnce(() => {
        loadingDuringRequest = auditLogService.isLoading.value
        return Promise.resolve({ items: [], total: 0 })
      })

      await auditLogService.search()

      expect(loadingDuringRequest).toBe(true)
      expect(auditLogService.isLoading.value).toBe(false)
    })
  })

  describe('changePagination', () => {
    it('calculates correct offset for page number', () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      auditLogService.changePagination(3)

      // page 3 with limit 20 => offset 40
      expect(auditLogService.pagination.value.offset).toBe(40)
    })

    it('triggers a search', () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      auditLogService.changePagination(2)

      expect(mockApi.get).toHaveBeenCalled()
    })
  })

  describe('clearPagination', () => {
    it('resets offset to 0', () => {
      auditLogService.pagination.value.offset = 60

      auditLogService.clearPagination()

      expect(auditLogService.pagination.value.offset).toBe(0)
    })
  })

  describe('applyFilters', () => {
    it('sets filters and resets pagination offset', () => {
      auditLogService.pagination.value.offset = 40
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      auditLogService.applyFilters({ action: 'delete', actorType: 'admin' })

      expect(auditLogService.filters.value).toEqual({ action: 'delete', actorType: 'admin' })
      expect(auditLogService.pagination.value.offset).toBe(0)
    })

    it('triggers a search', () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      auditLogService.applyFilters({ entityType: 'event' })

      expect(mockApi.get).toHaveBeenCalled()
    })
  })
})
