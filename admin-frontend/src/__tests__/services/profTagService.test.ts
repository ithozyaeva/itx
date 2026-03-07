import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockToast = vi.fn()
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({ toast: mockToast }),
}))

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

const { profTagService } = await import('@/services/profTagService')

describe('profTagService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "profTags"', () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    profTagService.search()

    expect(mockApi.get).toHaveBeenCalledWith('profTags', expect.any(Object))
  })

  it('search returns profTags registry', async () => {
    const mockTags = {
      items: [{ id: 1, name: 'Tag 1' }],
      total: 1,
    }
    mockJson.mockResolvedValueOnce(mockTags)

    const result = await profTagService.search()

    expect(result).toEqual(mockTags)
  })

  it('getById fetches a single profTag', async () => {
    const tag = { id: 5, name: 'Tag 5' }
    mockJson.mockResolvedValueOnce(tag)

    const result = await profTagService.getById(5)

    expect(mockApi.get).toHaveBeenCalledWith('profTags/5')
    expect(result).toEqual(tag)
  })

  it('create creates a new profTag', async () => {
    const newTag = { name: 'New Tag' }
    mockJson
      .mockResolvedValueOnce({ id: 1, ...newTag })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await profTagService.create(newTag as any)

    expect(mockApi.post).toHaveBeenCalledWith('profTags', { json: newTag })
    expect(result).toEqual({ id: 1, ...newTag })
  })

  it('update modifies a profTag', async () => {
    mockJson
      .mockResolvedValueOnce({ id: 2, name: 'Updated' })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await profTagService.update(2, { name: 'Updated' } as any)

    expect(mockApi.put).toHaveBeenCalledWith('profTags/2', {
      json: { name: 'Updated', id: 2 },
    })
    expect(result).toEqual({ id: 2, name: 'Updated' })
  })

  it('delete removes a profTag', async () => {
    mockApi.delete.mockResolvedValueOnce(undefined)
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    const result = await profTagService.delete(3)

    expect(result).toBe(true)
    expect(mockApi.delete).toHaveBeenCalledWith('profTags/3')
  })
})
