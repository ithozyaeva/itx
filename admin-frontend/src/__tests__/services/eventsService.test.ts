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

const { eventsService } = await import('@/services/eventsService')

describe('eventsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "events"', () => {
    // Verify it uses the correct base path by calling search
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    eventsService.search()

    expect(mockApi.get).toHaveBeenCalledWith('events', expect.any(Object))
  })

  it('search returns events registry', async () => {
    const mockEvents = {
      items: [
        {
          id: 1,
          title: 'Event 1',
          description: 'Desc',
          date: '2026-04-01',
          timezone: 'UTC+3',
          placeType: 'ONLINE',
          place: 'Zoom',
          open: true,
          hosts: [],
          eventTags: [],
        },
      ],
      total: 1,
    }
    mockJson.mockResolvedValueOnce(mockEvents)

    const result = await eventsService.search()

    expect(result).toEqual(mockEvents)
  })

  it('getById fetches a single event', async () => {
    const event = { id: 5, title: 'Event 5' }
    mockJson.mockResolvedValueOnce(event)

    const result = await eventsService.getById(5)

    expect(mockApi.get).toHaveBeenCalledWith('events/5')
    expect(result).toEqual(event)
  })

  it('create creates a new event', async () => {
    const newEvent = { title: 'New Event', description: 'Desc' }
    mockJson
      .mockResolvedValueOnce({ id: 1, ...newEvent })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await eventsService.create(newEvent as any)

    expect(mockApi.post).toHaveBeenCalledWith('events', { json: newEvent })
    expect(result).toEqual({ id: 1, ...newEvent })
  })

  it('delete removes an event', async () => {
    mockApi.delete.mockResolvedValueOnce(undefined)
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    const result = await eventsService.delete(3)

    expect(result).toBe(true)
    expect(mockApi.delete).toHaveBeenCalledWith('events/3')
  })

  it('update modifies an event', async () => {
    mockJson
      .mockResolvedValueOnce({ id: 2, title: 'Updated' })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await eventsService.update(2, { title: 'Updated' } as any)

    expect(mockApi.put).toHaveBeenCalledWith('events/2', {
      json: { title: 'Updated', id: 2 },
    })
    expect(result).toEqual({ id: 2, title: 'Updated' })
  })
})
