import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockApiClient, mockKyGet } = vi.hoisted(() => {
  const mockJson = vi.fn()
  return {
    mockJson,
    mockApiClient: {
      get: vi.fn(() => ({ ok: true, json: mockJson })),
      post: vi.fn(() => ({ json: mockJson })),
      patch: vi.fn(() => ({ json: mockJson })),
      put: vi.fn(() => ({ json: mockJson })),
      delete: vi.fn(() => ({ json: mockJson })),
    },
    mockKyGet: vi.fn(),
  }
})

vi.mock('@/services/api', () => ({
  apiClient: mockApiClient,
}))

vi.mock('ky', () => ({
  default: {
    get: mockKyGet,
  },
}))

vi.mock('@/composables/useUser', () => ({
  useUser: () => ({ value: null }),
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

import { profileService } from '@/services/profile'
import { handleError } from '@/services/errorService'

describe('profileService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getMe', () => {
    it('should call GET members/me', async () => {
      const userData = { id: 1, firstName: 'John' }
      mockJson.mockResolvedValue(userData)
      mockApiClient.get.mockReturnValue({ ok: true, json: mockJson })

      await profileService.getMe()

      expect(mockApiClient.get).toHaveBeenCalledWith('members/me')
    })

    it('should handle error when request fails', async () => {
      const error = new Error('Network error')
      mockApiClient.get.mockImplementation(() => { throw error })

      await profileService.getMe()

      expect(handleError).toHaveBeenCalledWith(error)
    })

    it('should not parse json if response is not ok', async () => {
      mockApiClient.get.mockReturnValue({ ok: false, json: mockJson })

      await profileService.getMe()

      expect(mockJson).not.toHaveBeenCalled()
    })
  })

  describe('updateMe', () => {
    it('should call PATCH members/me with user data', async () => {
      const newUser = { firstName: 'Jane', lastName: 'Doe', birthday: '1990-01-01', bio: 'Hello', grade: 'Senior', company: 'ACME', tg: '@jane' }
      const responseData = { id: 1, ...newUser }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.patch.mockReturnValue({ json: mockJson })

      const result = await profileService.updateMe(newUser)

      expect(mockApiClient.patch).toHaveBeenCalledWith('members/me', { json: newUser })
      expect(result).toEqual(responseData)
    })
  })

  describe('uploadAvatar', () => {
    it('should call POST members/me/avatar with form data', async () => {
      const formData = new FormData()
      formData.append('avatar', new Blob(['test']), 'avatar.png')
      const responseData = { id: 1, avatar: 'url' }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      const result = await profileService.uploadAvatar(formData)

      expect(mockApiClient.post).toHaveBeenCalledWith('members/me/avatar', { body: formData })
      expect(result).toEqual(responseData)
    })
  })

  describe('updateTags', () => {
    it('should call POST mentors/me/update-prof-tags', async () => {
      const profTags = [{ id: 1, name: 'TypeScript' }]
      const responseData = { id: 1, profTags }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      const result = await profileService.updateTags(profTags as any)

      expect(mockApiClient.post).toHaveBeenCalledWith('mentors/me/update-prof-tags', { json: { profTags } })
      expect(result).toEqual(responseData)
    })
  })

  describe('updateServices', () => {
    it('should call POST mentors/me/update-services', async () => {
      const services = [{ id: 1, name: 'Consulting' }]
      const responseData = { id: 1, services }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      const result = await profileService.updateServices(services as any)

      expect(mockApiClient.post).toHaveBeenCalledWith('mentors/me/update-services', { json: { services } })
      expect(result).toEqual(responseData)
    })
  })

  describe('updateContacts', () => {
    it('should call POST mentors/me/update-contacts', async () => {
      const contacts = [{ type: 'email', value: 'test@test.com' }]
      const responseData = { id: 1, contacts }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      const result = await profileService.updateContacts(contacts as any)

      expect(mockApiClient.post).toHaveBeenCalledWith('mentors/me/update-contacts', { json: { contacts } })
      expect(result).toEqual(responseData)
    })
  })

  describe('updateMentorInfo', () => {
    it('should call POST mentors/me/update-info', async () => {
      const info = { occupation: 'Developer', experience: '5 years' }
      const responseData = { id: 1, ...info }
      mockJson.mockResolvedValue(responseData)
      mockApiClient.post.mockReturnValue({ json: mockJson })

      const result = await profileService.updateMentorInfo(info)

      expect(mockApiClient.post).toHaveBeenCalledWith('mentors/me/update-info', { json: info })
      expect(result).toEqual(responseData)
    })
  })

  describe('getMemberById', () => {
    it('should call GET members/:id', async () => {
      const memberData = { id: 42, firstName: 'Alice' }
      mockJson.mockResolvedValue(memberData)
      mockApiClient.get.mockReturnValue({ json: mockJson })

      const result = await profileService.getMemberById(42)

      expect(mockApiClient.get).toHaveBeenCalledWith('members/42')
      expect(result).toEqual(memberData)
    })
  })

  describe('getAllProfTags', () => {
    it('should call GET /api/profTags using raw ky', async () => {
      const tags = [{ id: 1, name: 'Go' }, { id: 2, name: 'Vue' }]
      const mockTagsJson = vi.fn().mockResolvedValue({ items: tags })
      mockKyGet.mockReturnValue({ json: mockTagsJson })

      const result = await profileService.getAllProfTags()

      expect(mockKyGet).toHaveBeenCalledWith('/api/profTags')
      expect(result).toEqual(tags)
    })
  })
})
