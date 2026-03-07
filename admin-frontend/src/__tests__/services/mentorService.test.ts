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

const { mentorService } = await import('@/services/mentorService')

describe('mentorService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "mentors"', () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    mentorService.search()

    expect(mockApi.get).toHaveBeenCalledWith('mentors', expect.any(Object))
  })

  it('search returns mentors registry', async () => {
    const mockMentors = {
      items: [{ id: 1, name: 'Mentor 1' }],
      total: 1,
    }
    mockJson.mockResolvedValueOnce(mockMentors)

    const result = await mentorService.search()

    expect(result).toEqual(mockMentors)
  })

  it('getById fetches a single mentor', async () => {
    const mentor = { id: 5, name: 'Mentor 5' }
    mockJson.mockResolvedValueOnce(mentor)

    const result = await mentorService.getById(5)

    expect(mockApi.get).toHaveBeenCalledWith('mentors/5')
    expect(result).toEqual(mentor)
  })

  it('create creates a new mentor', async () => {
    const newMentor = { name: 'New Mentor' }
    mockJson
      .mockResolvedValueOnce({ id: 1, ...newMentor })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await mentorService.create(newMentor as any)

    expect(mockApi.post).toHaveBeenCalledWith('mentors', { json: newMentor })
    expect(result).toEqual({ id: 1, ...newMentor })
  })

  it('update modifies a mentor', async () => {
    mockJson
      .mockResolvedValueOnce({ id: 2, name: 'Updated' })
      .mockResolvedValueOnce({ items: [], total: 0 })

    const result = await mentorService.update(2, { name: 'Updated' } as any)

    expect(mockApi.put).toHaveBeenCalledWith('mentors/2', {
      json: { name: 'Updated', id: 2 },
    })
    expect(result).toEqual({ id: 2, name: 'Updated' })
  })

  it('delete removes a mentor', async () => {
    mockApi.delete.mockResolvedValueOnce(undefined)
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    const result = await mentorService.delete(3)

    expect(result).toBe(true)
    expect(mockApi.delete).toHaveBeenCalledWith('mentors/3')
  })
})
