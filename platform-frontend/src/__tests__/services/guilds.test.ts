import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
      put: vi.fn(() => ({ json: mockJson })),
      delete: vi.fn(() => ({ json: mockJson })),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { guildService } from '@/services/guilds'

describe('guildService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('should call GET guilds', async () => {
      const guilds = [{ id: 1, name: 'Alpha' }]
      mockJson.mockResolvedValue(guilds)

      const result = await guildService.getAll()

      expect(mockApiClient.get).toHaveBeenCalledWith('guilds')
      expect(result).toEqual(guilds)
    })
  })

  describe('create', () => {
    it('should call POST guilds with data', async () => {
      const data = { name: 'Beta', description: 'A guild', icon: 'star', color: '#ff0000' }
      const created = { id: 2, ...data }
      mockJson.mockResolvedValue(created)

      const result = await guildService.create(data)

      expect(mockApiClient.post).toHaveBeenCalledWith('guilds', { json: data })
      expect(result).toEqual(created)
    })
  })

  describe('update', () => {
    it('should call PUT guilds/:id with data', async () => {
      const data = { name: 'Beta Updated', description: 'Updated', icon: 'star', color: '#00ff00' }
      mockJson.mockResolvedValue({})

      await guildService.update(2, data)

      expect(mockApiClient.put).toHaveBeenCalledWith('guilds/2', { json: data })
    })
  })

  describe('join', () => {
    it('should call POST guilds/:id/join', async () => {
      mockJson.mockResolvedValue({})

      await guildService.join(3)

      expect(mockApiClient.post).toHaveBeenCalledWith('guilds/3/join')
    })
  })

  describe('leave', () => {
    it('should call POST guilds/:id/leave', async () => {
      mockJson.mockResolvedValue({})

      await guildService.leave(3)

      expect(mockApiClient.post).toHaveBeenCalledWith('guilds/3/leave')
    })
  })

  describe('remove', () => {
    it('should call DELETE guilds/:id', async () => {
      await guildService.remove(4)

      expect(mockApiClient.delete).toHaveBeenCalledWith('guilds/4')
    })
  })

  describe('getMembers', () => {
    it('should call GET guilds/:id/members', async () => {
      const members = [{ id: 1, name: 'User1' }]
      mockJson.mockResolvedValue(members)

      const result = await guildService.getMembers(5)

      expect(mockApiClient.get).toHaveBeenCalledWith('guilds/5/members')
      expect(result).toEqual(members)
    })
  })
})
