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
  put: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { BaseService } = await import('@/services/api/baseService')

interface TestEntity {
  id: number
  name: string
}

describe('BaseService', () => {
  let service: InstanceType<typeof BaseService<TestEntity>>

  beforeEach(() => {
    vi.clearAllMocks()
    service = new BaseService<TestEntity>('test-entities')
  })

  describe('search', () => {
    it('fetches with default pagination', async () => {
      const mockResponse = { items: [{ id: 1, name: 'Item 1' }], total: 1 }
      mockJson.mockResolvedValueOnce(mockResponse)

      const result = await service.search()

      expect(mockApi.get).toHaveBeenCalledWith('test-entities', {
        searchParams: { limit: 10, offset: 0 },
      })
      expect(result).toEqual(mockResponse)
      expect(service.items.value).toEqual(mockResponse)
    })

    it('merges additional params', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await service.search({ status: 'ACTIVE', search: 'test' })

      expect(mockApi.get).toHaveBeenCalledWith('test-entities', {
        searchParams: { limit: 10, offset: 0, status: 'ACTIVE', search: 'test' },
      })
    })

    it('returns empty registry on error', async () => {
      const error = new Error('Network error')
      mockJson.mockRejectedValueOnce(error)

      const result = await service.search()

      expect(result).toEqual({ items: [], total: 0 })
      expect(mockHandleError).toHaveBeenCalledWith(error)
    })

    it('sets and resets isLoading', async () => {
      let loadingDuringRequest = false
      mockJson.mockImplementationOnce(() => {
        loadingDuringRequest = service.isLoading.value
        return Promise.resolve({ items: [], total: 0 })
      })

      await service.search()

      expect(loadingDuringRequest).toBe(true)
      expect(service.isLoading.value).toBe(false)
    })
  })

  describe('getAll', () => {
    it('fetches without pagination params', async () => {
      const mockResponse = { items: [{ id: 1, name: 'All' }], total: 1 }
      mockJson.mockResolvedValueOnce(mockResponse)

      const result = await service.getAll()

      expect(mockApi.get).toHaveBeenCalledWith('test-entities')
      expect(result).toEqual(mockResponse)
    })

    it('returns empty registry on error', async () => {
      mockJson.mockRejectedValueOnce(new Error('fail'))

      const result = await service.getAll()

      expect(result).toEqual({ items: [], total: 0 })
    })
  })

  describe('getById', () => {
    it('fetches a single entity', async () => {
      const entity = { id: 5, name: 'Entity 5' }
      mockJson.mockResolvedValueOnce(entity)

      const result = await service.getById(5)

      expect(mockApi.get).toHaveBeenCalledWith('test-entities/5')
      expect(result).toEqual(entity)
    })

    it('returns null on error', async () => {
      mockJson.mockRejectedValueOnce(new Error('Not found'))

      const result = await service.getById(999)

      expect(result).toBeNull()
      expect(mockHandleError).toHaveBeenCalled()
    })
  })

  describe('create', () => {
    it('creates an entity and shows success toast', async () => {
      const newEntity = { id: 10, name: 'New Entity' }
      mockJson
        .mockResolvedValueOnce(newEntity) // create response
        .mockResolvedValueOnce({ items: [newEntity], total: 1 }) // search refresh

      const result = await service.create({ name: 'New Entity' })

      expect(mockApi.post).toHaveBeenCalledWith('test-entities', {
        json: { name: 'New Entity' },
      })
      expect(result).toEqual(newEntity)
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Запись успешно создана',
      })
    })

    it('returns null on error', async () => {
      mockJson.mockRejectedValueOnce(new Error('Create failed'))

      const result = await service.create({ name: 'Fail' })

      expect(result).toBeNull()
    })

    it('refreshes list after create', async () => {
      mockJson
        .mockResolvedValueOnce({ id: 1, name: 'X' }) // create
        .mockResolvedValueOnce({ items: [], total: 0 }) // search

      await service.create({ name: 'X' })

      // search is called after create
      expect(mockApi.get).toHaveBeenCalledWith('test-entities', {
        searchParams: { limit: 10, offset: 0 },
      })
    })
  })

  describe('update', () => {
    it('updates an entity and shows success toast', async () => {
      const updated = { id: 3, name: 'Updated' }
      mockJson
        .mockResolvedValueOnce(updated) // update response
        .mockResolvedValueOnce({ items: [updated], total: 1 }) // search refresh

      const result = await service.update(3, { name: 'Updated' })

      expect(mockApi.put).toHaveBeenCalledWith('test-entities/3', {
        json: { name: 'Updated', id: 3 },
      })
      expect(result).toEqual(updated)
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Запись успешно обновлена',
      })
    })

    it('returns null on error', async () => {
      mockJson.mockRejectedValueOnce(new Error('Update failed'))

      const result = await service.update(3, { name: 'Fail' })

      expect(result).toBeNull()
    })
  })

  describe('delete', () => {
    it('deletes an entity and shows success toast', async () => {
      mockApi.delete.mockResolvedValueOnce(undefined)
      mockJson.mockResolvedValueOnce({ items: [], total: 0 }) // search refresh

      const result = await service.delete(7)

      expect(result).toBe(true)
      expect(mockApi.delete).toHaveBeenCalledWith('test-entities/7')
      expect(mockToast).toHaveBeenCalledWith({
        title: 'Успешно',
        description: 'Запись успешно удалена',
      })
    })

    it('returns false on error', async () => {
      mockApi.delete.mockRejectedValueOnce(new Error('Delete failed'))

      const result = await service.delete(7)

      expect(result).toBe(false)
    })
  })

  describe('changePagination', () => {
    it('calculates correct offset for page number', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      service.changePagination(5)

      // page 5 with limit 10 => offset 40
      expect(service.pagination.value.offset).toBe(40)
    })

    it('triggers a search', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      service.changePagination(2)

      expect(mockApi.get).toHaveBeenCalled()
    })
  })

  describe('clearPagination', () => {
    it('resets offset to 0', () => {
      service.pagination.value.offset = 50

      service.clearPagination()

      expect(service.pagination.value.offset).toBe(0)
    })
  })
})
