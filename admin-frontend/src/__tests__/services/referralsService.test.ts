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

const { referralsService } = await import('@/services/referralsService')

describe('referralsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "admin-referals"', () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    referralsService.search()

    expect(mockApi.get).toHaveBeenCalledWith('admin-referals', expect.any(Object))
  })

  it('search returns referral links registry', async () => {
    const mockReferrals = {
      items: [{ id: 1, code: 'REF1' }],
      total: 1,
    }
    mockJson.mockResolvedValueOnce(mockReferrals)

    const result = await referralsService.search()

    expect(result).toEqual(mockReferrals)
  })

  it('getById fetches a single referral link', async () => {
    const link = { id: 5, code: 'REF5' }
    mockJson.mockResolvedValueOnce(link)

    const result = await referralsService.getById(5)

    expect(mockApi.get).toHaveBeenCalledWith('admin-referals/5')
    expect(result).toEqual(link)
  })

  it('create creates a new referral link', async () => {
    const newLink = { code: 'NEW_REF' }
    mockJson
      .mockResolvedValueOnce({ id: 1, ...newLink })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await referralsService.create(newLink as any)

    expect(mockApi.post).toHaveBeenCalledWith('admin-referals', { json: newLink })
    expect(result).toEqual({ id: 1, ...newLink })
  })

  it('update modifies a referral link', async () => {
    mockJson
      .mockResolvedValueOnce({ id: 2, code: 'Updated' })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await referralsService.update(2, { code: 'Updated' } as any)

    expect(mockApi.put).toHaveBeenCalledWith('admin-referals/2', {
      json: { code: 'Updated', id: 2 },
    })
    expect(result).toEqual({ id: 2, code: 'Updated' })
  })

  it('delete removes a referral link', async () => {
    mockApi.delete.mockResolvedValueOnce(undefined)
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    const result = await referralsService.delete(3)

    expect(result).toBe(true)
    expect(mockApi.delete).toHaveBeenCalledWith('admin-referals/3')
  })
})
