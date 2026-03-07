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

const { memberService } = await import('@/services/memberService')

describe('memberService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "members"', () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    memberService.search()

    expect(mockApi.get).toHaveBeenCalledWith('members', expect.any(Object))
  })

  it('search returns members registry', async () => {
    const mockMembers = {
      items: [{ id: 1, name: 'Member 1' }],
      total: 1,
    }
    mockJson.mockResolvedValueOnce(mockMembers)

    const result = await memberService.search()

    expect(result).toEqual(mockMembers)
  })

  it('getById fetches a single member', async () => {
    const member = { id: 5, name: 'Member 5' }
    mockJson.mockResolvedValueOnce(member)

    const result = await memberService.getById(5)

    expect(mockApi.get).toHaveBeenCalledWith('members/5')
    expect(result).toEqual(member)
  })

  it('create creates a new member', async () => {
    const newMember = { name: 'New Member' }
    mockJson
      .mockResolvedValueOnce({ id: 1, ...newMember })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await memberService.create(newMember as any)

    expect(mockApi.post).toHaveBeenCalledWith('members', { json: newMember })
    expect(result).toEqual({ id: 1, ...newMember })
  })

  it('update modifies a member', async () => {
    mockJson
      .mockResolvedValueOnce({ id: 2, name: 'Updated' })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await memberService.update(2, { name: 'Updated' } as any)

    expect(mockApi.put).toHaveBeenCalledWith('members/2', {
      json: { name: 'Updated', id: 2 },
    })
    expect(result).toEqual({ id: 2, name: 'Updated' })
  })

  it('delete removes a member', async () => {
    mockApi.delete.mockResolvedValueOnce(undefined)
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    const result = await memberService.delete(3)

    expect(result).toBe(true)
    expect(mockApi.delete).toHaveBeenCalledWith('members/3')
  })
})
