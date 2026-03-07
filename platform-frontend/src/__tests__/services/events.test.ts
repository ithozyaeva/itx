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

import { eventsService } from '@/services/events'

describe('eventsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-01-15T12:00:00.000Z'))
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('searchOld', () => {
    it('should call GET events with dateTo param', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      const result = await eventsService.searchOld(10, 0)

      expect(mockApiClient.get).toHaveBeenCalledWith('events', {
        searchParams: {
          limit: 10,
          offset: 0,
          dateTo: '2026-01-15T12:00:00.000Z',
        },
      })
      expect(result).toEqual(data)
    })

    it('should include filters in search params', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await eventsService.searchOld(10, 0, { title: 'Meetup', placeType: 'online' })

      expect(mockApiClient.get).toHaveBeenCalledWith('events', {
        searchParams: {
          limit: 10,
          offset: 0,
          dateTo: '2026-01-15T12:00:00.000Z',
          title: 'Meetup',
          placeType: 'online',
        },
      })
    })

    it('should exclude empty filter values', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await eventsService.searchOld(10, 0, { title: '', placeType: undefined })

      expect(mockApiClient.get).toHaveBeenCalledWith('events', {
        searchParams: {
          limit: 10,
          offset: 0,
          dateTo: '2026-01-15T12:00:00.000Z',
        },
      })
    })
  })

  describe('searchNext', () => {
    it('should call GET events with dateFrom param', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      const result = await eventsService.searchNext(5, 10)

      expect(mockApiClient.get).toHaveBeenCalledWith('events', {
        searchParams: {
          limit: 5,
          offset: 10,
          dateFrom: '2026-01-15T12:00:00.000Z',
        },
      })
      expect(result).toEqual(data)
    })

    it('should include filters in search params', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      await eventsService.searchNext(10, 0, { placeType: 'offline' })

      expect(mockApiClient.get).toHaveBeenCalledWith('events', {
        searchParams: {
          limit: 10,
          offset: 0,
          dateFrom: '2026-01-15T12:00:00.000Z',
          placeType: 'offline',
        },
      })
    })
  })

  describe('getById', () => {
    it('should call GET events/:id', async () => {
      const event = { id: 1, title: 'Conference' }
      mockJson.mockResolvedValue(event)

      const result = await eventsService.getById(1)

      expect(mockApiClient.get).toHaveBeenCalledWith('events/1')
      expect(result).toEqual(event)
    })
  })

  describe('applyEvent', () => {
    it('should call POST events/apply with eventId', async () => {
      const event = { id: 5, title: 'Workshop' }
      mockJson.mockResolvedValue(event)

      const result = await eventsService.applyEvent(5)

      expect(mockApiClient.post).toHaveBeenCalledWith('events/apply', { json: { eventId: 5 } })
      expect(result).toEqual(event)
    })
  })

  describe('declineEvent', () => {
    it('should call POST events/decline with eventId', async () => {
      const event = { id: 3, title: 'Talk' }
      mockJson.mockResolvedValue(event)

      const result = await eventsService.declineEvent(3)

      expect(mockApiClient.post).toHaveBeenCalledWith('events/decline', { json: { eventId: 3 } })
      expect(result).toEqual(event)
    })
  })
})
