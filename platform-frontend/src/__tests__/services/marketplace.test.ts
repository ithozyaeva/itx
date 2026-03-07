import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
      patch: vi.fn(() => ({ json: mockJson })),
      put: vi.fn(() => ({ json: mockJson })),
      delete: vi.fn(() => Promise.resolve()),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { marketplaceService } from '@/services/marketplace'

describe('marketplaceService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAll', () => {
    it('should call GET marketplace with no params', async () => {
      const data = { items: [], total: 0 }
      mockJson.mockResolvedValue(data)

      const result = await marketplaceService.getAll()

      expect(mockApiClient.get).toHaveBeenCalledWith('marketplace', {
        searchParams: expect.any(URLSearchParams),
      })
      const callArgs = mockApiClient.get.mock.calls[0]
      const searchParams = callArgs[1].searchParams as URLSearchParams
      expect(searchParams.toString()).toBe('')
      expect(result).toEqual(data)
    })

    it('should call GET marketplace with all params', async () => {
      const data = { items: [{ id: 1 }], total: 1 }
      mockJson.mockResolvedValue(data)

      const result = await marketplaceService.getAll({ status: 'active', limit: 10, offset: 20 })

      expect(mockApiClient.get).toHaveBeenCalledWith('marketplace', {
        searchParams: expect.any(URLSearchParams),
      })
      const callArgs = mockApiClient.get.mock.calls[0]
      const searchParams = callArgs[1].searchParams as URLSearchParams
      expect(searchParams.get('status')).toBe('active')
      expect(searchParams.get('limit')).toBe('10')
      expect(searchParams.get('offset')).toBe('20')
      expect(result).toEqual(data)
    })
  })

  describe('getById', () => {
    it('should call GET marketplace/:id', async () => {
      const item = { id: 5, title: 'Laptop' }
      mockJson.mockResolvedValue(item)

      const result = await marketplaceService.getById(5)

      expect(mockApiClient.get).toHaveBeenCalledWith('marketplace/5')
      expect(result).toEqual(item)
    })
  })

  describe('create', () => {
    it('should call POST marketplace with FormData', async () => {
      const data = {
        title: 'Laptop',
        description: 'Good laptop',
        price: '1000',
        city: 'Moscow',
        canShip: true,
        condition: 'used',
        defects: 'none',
        packageContents: 'laptop, charger',
        contactTelegram: '@seller',
      }
      const image = new File(['img'], 'photo.jpg', { type: 'image/jpeg' })
      const responseData = { id: 1, ...data }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      const result = await marketplaceService.create(data as any, image)

      expect(mockApiClient.post).toHaveBeenCalledWith('marketplace', {
        body: expect.any(FormData),
      })

      const callArgs = mockApiClient.post.mock.calls[0]
      const formData = callArgs[1].body as FormData
      expect(formData.get('title')).toBe('Laptop')
      expect(formData.get('description')).toBe('Good laptop')
      expect(formData.get('price')).toBe('1000')
      expect(formData.get('city')).toBe('Moscow')
      expect(formData.get('canShip')).toBe('true')
      expect(formData.get('condition')).toBe('used')
      expect(formData.get('defects')).toBe('none')
      expect(formData.get('packageContents')).toBe('laptop, charger')
      expect(formData.get('contactTelegram')).toBe('@seller')
      expect(formData.get('image')).toEqual(image)
      expect(result).toEqual(responseData)
    })

    it('should not append image if not provided', async () => {
      const data = {
        title: 'Phone',
        description: 'Nice phone',
        price: '500',
        city: 'SPb',
        canShip: false,
        condition: 'new',
        defects: '',
        packageContents: 'phone',
        contactTelegram: '@seller',
      }
      mockJson.mockResolvedValue({ id: 2 })
      mockApiClient.post.mockReturnValue({ json: mockJson })

      await marketplaceService.create(data as any)

      const callArgs = mockApiClient.post.mock.calls[0]
      const formData = callArgs[1].body as FormData
      expect(formData.get('image')).toBeNull()
      expect(formData.get('canShip')).toBe('false')
    })
  })

  describe('update', () => {
    it('should call PATCH marketplace/:id with data', async () => {
      const updateData = { title: 'Updated Laptop' }
      const responseData = { id: 3, title: 'Updated Laptop' }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.patch.mockReturnValue({ json: mockJson })

      const result = await marketplaceService.update(3, updateData)

      expect(mockApiClient.patch).toHaveBeenCalledWith('marketplace/3', { json: updateData })
      expect(result).toEqual(responseData)
    })
  })

  describe('requestPurchase', () => {
    it('should call POST marketplace/:id/request-purchase', async () => {
      const item = { id: 7, status: 'reserved' }
      mockJson.mockResolvedValue(item)

      const result = await marketplaceService.requestPurchase(7)

      expect(mockApiClient.post).toHaveBeenCalledWith('marketplace/7/request-purchase')
      expect(result).toEqual(item)
    })
  })

  describe('cancelPurchase', () => {
    it('should call POST marketplace/:id/cancel-purchase', async () => {
      const item = { id: 7, status: 'active' }
      mockJson.mockResolvedValue(item)

      const result = await marketplaceService.cancelPurchase(7)

      expect(mockApiClient.post).toHaveBeenCalledWith('marketplace/7/cancel-purchase')
      expect(result).toEqual(item)
    })
  })

  describe('markSold', () => {
    it('should call POST marketplace/:id/sold', async () => {
      const item = { id: 7, status: 'sold' }
      mockJson.mockResolvedValue(item)

      const result = await marketplaceService.markSold(7)

      expect(mockApiClient.post).toHaveBeenCalledWith('marketplace/7/sold')
      expect(result).toEqual(item)
    })
  })

  describe('remove', () => {
    it('should call DELETE marketplace/:id', async () => {
      mockApiClient.delete.mockResolvedValue(undefined)

      await marketplaceService.remove(10)

      expect(mockApiClient.delete).toHaveBeenCalledWith('marketplace/10')
    })
  })
})
