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

const { eventTagService } = await import('@/services/eventTagService')

describe('eventTagService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "eventTags"', () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    eventTagService.search()

    expect(mockApi.get).toHaveBeenCalledWith('eventTags', expect.any(Object))
  })

  it('search returns eventTags registry', async () => {
    const mockTags = {
      items: [{ id: 1, name: 'Event Tag 1' }],
      total: 1,
    }
    mockJson.mockResolvedValueOnce(mockTags)

    const result = await eventTagService.search()

    expect(result).toEqual(mockTags)
  })

  it('getById fetches a single eventTag', async () => {
    const tag = { id: 5, name: 'Event Tag 5' }
    mockJson.mockResolvedValueOnce(tag)

    const result = await eventTagService.getById(5)

    expect(mockApi.get).toHaveBeenCalledWith('eventTags/5')
    expect(result).toEqual(tag)
  })

  it('create creates a new eventTag', async () => {
    const newTag = { name: 'New Event Tag' }
    mockJson
      .mockResolvedValueOnce({ id: 1, ...newTag })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await eventTagService.create(newTag as any)

    expect(mockApi.post).toHaveBeenCalledWith('eventTags', { json: newTag })
    expect(result).toEqual({ id: 1, ...newTag })
  })

  it('update modifies an eventTag', async () => {
    mockJson
      .mockResolvedValueOnce({ id: 2, name: 'Updated' })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await eventTagService.update(2, { name: 'Updated' } as any)

    expect(mockApi.put).toHaveBeenCalledWith('eventTags/2', {
      json: { name: 'Updated', id: 2 },
    })
    expect(result).toEqual({ id: 2, name: 'Updated' })
  })

  it('delete removes an eventTag', async () => {
    mockApi.delete.mockResolvedValueOnce(undefined)
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    const result = await eventTagService.delete(3)

    expect(result).toBe(true)
    expect(mockApi.delete).toHaveBeenCalledWith('eventTags/3')
  })
})
