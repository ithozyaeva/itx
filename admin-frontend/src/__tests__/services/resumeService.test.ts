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
const mockBlob = vi.fn()
const mockApi = {
  get: vi.fn(() => ({ json: mockJson, blob: mockBlob })),
  post: vi.fn(() => ({ json: mockJson })),
  put: vi.fn(() => ({ json: mockJson })),
  delete: vi.fn(),
}
vi.mock('@/lib/api', () => ({ default: mockApi }))

const { resumeService } = await import('@/services/resumeService')

describe('resumeService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('is a BaseService with basePath "resumes"', async () => {
    mockJson.mockResolvedValueOnce({ items: [], total: 0 })

    await resumeService.search()

    expect(mockApi.get).toHaveBeenCalledWith('resumes', expect.any(Object))
  })

  describe('searchWithFilters', () => {
    it('delegates to search with filter params', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await resumeService.searchWithFilters({ workFormat: 'REMOTE', desiredPosition: 'Frontend' })

      expect(mockApi.get).toHaveBeenCalledWith('resumes', {
        searchParams: expect.objectContaining({
          workFormat: 'REMOTE',
          desiredPosition: 'Frontend',
        }),
      })
    })

    it('works without filters', async () => {
      mockJson.mockResolvedValueOnce({ items: [], total: 0 })

      await resumeService.searchWithFilters()

      expect(mockApi.get).toHaveBeenCalledWith('resumes', expect.any(Object))
    })
  })

  describe('downloadArchive', () => {
    it('creates a download link for the archive', async () => {
      const mockBlobValue = new Blob(['zip-data'])
      mockBlob.mockResolvedValueOnce(mockBlobValue)

      const mockCreateObjectURL = vi.fn(() => 'blob:http://test/zip')
      const mockRevokeObjectURL = vi.fn()
      window.URL.createObjectURL = mockCreateObjectURL
      window.URL.revokeObjectURL = mockRevokeObjectURL

      const mockClick = vi.fn()
      const mockRemove = vi.fn()
      const mockElement = { href: '', download: '', click: mockClick, remove: mockRemove } as any
      vi.spyOn(document, 'createElement').mockReturnValueOnce(mockElement)
      vi.spyOn(document.body, 'appendChild').mockImplementationOnce(() => mockElement)

      await resumeService.downloadArchive()

      expect(mockApi.get).toHaveBeenCalledWith('resumes/download', {
        searchParams: {},
      })
      expect(mockElement.download).toBe('resumes.zip')
      expect(mockClick).toHaveBeenCalled()
      expect(mockRemove).toHaveBeenCalled()
      expect(mockRevokeObjectURL).toHaveBeenCalledWith('blob:http://test/zip')
    })

    it('passes filter params to download endpoint', async () => {
      const mockBlobValue = new Blob(['data'])
      mockBlob.mockResolvedValueOnce(mockBlobValue)

      window.URL.createObjectURL = vi.fn(() => 'blob:url')
      window.URL.revokeObjectURL = vi.fn()
      const mockElement = { href: '', download: '', click: vi.fn(), remove: vi.fn() } as any
      vi.spyOn(document, 'createElement').mockReturnValueOnce(mockElement)
      vi.spyOn(document.body, 'appendChild').mockImplementationOnce(() => mockElement)

      await resumeService.downloadArchive({ workFormat: 'OFFICE' })

      expect(mockApi.get).toHaveBeenCalledWith('resumes/download', {
        searchParams: { workFormat: 'OFFICE' },
      })
    })

    it('handles download errors', async () => {
      const error = new Error('Download failed')
      mockBlob.mockRejectedValueOnce(error)

      await resumeService.downloadArchive()

      expect(mockHandleError).toHaveBeenCalledWith(error)
    })
  })
})
