import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
      patch: vi.fn(() => ({ json: mockJson })),
      delete: vi.fn(() => Promise.resolve()),
    },
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

import { resumeService } from '@/services/resume'

describe('resumeService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('upload', () => {
    it('should call POST resumes with FormData containing file only', async () => {
      const file = new File(['content'], 'resume.pdf', { type: 'application/pdf' })
      const responseData = { resume: { id: 1 }, parsed: {} }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      const result = await resumeService.upload({ file })

      expect(mockApiClient.post).toHaveBeenCalledWith('resumes', {
        body: expect.any(FormData),
      })

      const callArgs = mockApiClient.post.mock.calls[0]
      const formData = callArgs[1].body as FormData
      expect(formData.get('file')).toEqual(file)
      expect(formData.get('workExperience')).toBeNull()
      expect(formData.get('desiredPosition')).toBeNull()
      expect(formData.get('workFormat')).toBeNull()

      expect(result).toEqual(responseData)
    })

    it('should include optional fields in FormData when provided', async () => {
      const file = new File(['content'], 'resume.pdf', { type: 'application/pdf' })
      const responseData = { resume: { id: 1 }, parsed: {} }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      await resumeService.upload({
        file,
        workExperience: '5 years',
        desiredPosition: 'Senior Developer',
        workFormat: 'remote' as any,
      })

      const callArgs = mockApiClient.post.mock.calls[0]
      const formData = callArgs[1].body as FormData
      expect(formData.get('file')).toEqual(file)
      expect(formData.get('workExperience')).toBe('5 years')
      expect(formData.get('desiredPosition')).toBe('Senior Developer')
      expect(formData.get('workFormat')).toBe('remote')
    })
  })

  describe('listMine', () => {
    it('should call GET resumes/me', async () => {
      const resumes = [{ id: 1, name: 'resume.pdf' }]
      mockJson.mockResolvedValue(resumes)
      mockApiClient.get.mockReturnValue({ json: mockJson })

      const result = await resumeService.listMine()

      expect(mockApiClient.get).toHaveBeenCalledWith('resumes/me')
      expect(result).toEqual(resumes)
    })
  })

  describe('update', () => {
    it('should call PATCH resumes/:id with payload', async () => {
      const updatedResume = { id: 1, workExperience: '7 years' }
      mockJson.mockResolvedValue(updatedResume)
      mockApiClient.patch.mockReturnValue({ json: mockJson })

      const result = await resumeService.update(1, { workExperience: '7 years' })

      expect(mockApiClient.patch).toHaveBeenCalledWith('resumes/1', {
        json: { workExperience: '7 years' },
      })
      expect(result).toEqual(updatedResume)
    })
  })

  describe('delete', () => {
    it('should call DELETE resumes/:id', async () => {
      mockApiClient.delete.mockResolvedValue(undefined)

      await resumeService.delete(5)

      expect(mockApiClient.delete).toHaveBeenCalledWith('resumes/5')
    })
  })

  describe('download', () => {
    it('should call GET resumes/:id/download and return url', async () => {
      const data = { url: 'https://s3.example.com/resume.pdf' }
      mockJson.mockResolvedValue(data)
      mockApiClient.get.mockReturnValue({ json: mockJson })

      const result = await resumeService.download(3)

      expect(mockApiClient.get).toHaveBeenCalledWith('resumes/3/download')
      expect(result).toEqual(data)
    })
  })
})
