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

import { taskExchangeService } from '@/services/taskExchange'

describe('taskExchangeService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('should call GET tasks with no params', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      const result = await taskExchangeService.getAll()

      expect(mockApiClient.get).toHaveBeenCalledWith('tasks', {
        searchParams: expect.any(URLSearchParams),
      })
      const callArgs = mockApiClient.get.mock.calls[0]
      const searchParams = callArgs[1].searchParams as URLSearchParams
      expect(searchParams.toString()).toBe('')
      expect(result).toEqual(data)
    })

    it('should call GET tasks with all params', async () => {
      const data = { items: [{ id: 1 }], total: 1 }
      mockJson.mockResolvedValue(data)

      const result = await taskExchangeService.getAll({ status: 'open', limit: 5, offset: 10 })

      const callArgs = mockApiClient.get.mock.calls[0]
      const searchParams = callArgs[1].searchParams as URLSearchParams
      expect(searchParams.get('status')).toBe('open')
      expect(searchParams.get('limit')).toBe('5')
      expect(searchParams.get('offset')).toBe('10')
      expect(result).toEqual(data)
    })
  })

  describe('getById', () => {
    it('should call GET tasks/:id', async () => {
      const task = { id: 3, title: 'Fix bug' }
      mockJson.mockResolvedValue(task)

      const result = await taskExchangeService.getById(3)

      expect(mockApiClient.get).toHaveBeenCalledWith('tasks/3')
      expect(result).toEqual(task)
    })
  })

  describe('create', () => {
    it('should call POST tasks with json data', async () => {
      const data = { title: 'New task', description: 'Do something' }
      const responseData = { id: 1, ...data }
      mockJson.mockResolvedValue(responseData)

      const result = await taskExchangeService.create(data as any)

      expect(mockApiClient.post).toHaveBeenCalledWith('tasks', { json: data })
      expect(result).toEqual(responseData)
    })
  })

  describe('update', () => {
    it('should call PUT tasks/:id with json data', async () => {
      const data = { title: 'Updated task' }
      const responseData = { id: 2, title: 'Updated task' }
      mockJson.mockResolvedValue(responseData)

      const result = await taskExchangeService.update(2, data as any)

      expect(mockApiClient.put).toHaveBeenCalledWith('tasks/2', { json: data })
      expect(result).toEqual(responseData)
    })
  })

  describe('assign', () => {
    it('should call POST tasks/:id/assign', async () => {
      const task = { id: 4, status: 'assigned' }
      mockJson.mockResolvedValue(task)

      const result = await taskExchangeService.assign(4)

      expect(mockApiClient.post).toHaveBeenCalledWith('tasks/4/assign')
      expect(result).toEqual(task)
    })
  })

  describe('unassign', () => {
    it('should call POST tasks/:id/unassign', async () => {
      const task = { id: 4, status: 'open' }
      mockJson.mockResolvedValue(task)

      const result = await taskExchangeService.unassign(4)

      expect(mockApiClient.post).toHaveBeenCalledWith('tasks/4/unassign')
      expect(result).toEqual(task)
    })
  })

  describe('removeAssignee', () => {
    it('should call DELETE tasks/:taskId/assignees/:memberId', async () => {
      const task = { id: 10, status: 'open' }
      mockJson.mockResolvedValue(task)

      const result = await taskExchangeService.removeAssignee(10, 25)

      expect(mockApiClient.delete).toHaveBeenCalledWith('tasks/10/assignees/25')
      expect(result).toEqual(task)
    })
  })

  describe('markDone', () => {
    it('should call POST tasks/:id/done', async () => {
      const task = { id: 6, status: 'done' }
      mockJson.mockResolvedValue(task)

      const result = await taskExchangeService.markDone(6)

      expect(mockApiClient.post).toHaveBeenCalledWith('tasks/6/done')
      expect(result).toEqual(task)
    })
  })

  describe('approve', () => {
    it('should call POST tasks/:id/approve', async () => {
      const task = { id: 8, status: 'approved' }
      mockJson.mockResolvedValue(task)

      const result = await taskExchangeService.approve(8)

      expect(mockApiClient.post).toHaveBeenCalledWith('tasks/8/approve')
      expect(result).toEqual(task)
    })
  })

  describe('reject', () => {
    it('should call POST tasks/:id/reject', async () => {
      const task = { id: 9, status: 'rejected' }
      mockJson.mockResolvedValue(task)

      const result = await taskExchangeService.reject(9)

      expect(mockApiClient.post).toHaveBeenCalledWith('tasks/9/reject')
      expect(result).toEqual(task)
    })
  })

  describe('remove', () => {
    it('should call DELETE tasks/:id', async () => {
      mockApiClient.delete.mockResolvedValue(undefined)

      await taskExchangeService.remove(12)

      expect(mockApiClient.delete).toHaveBeenCalledWith('tasks/12')
    })
  })
})
