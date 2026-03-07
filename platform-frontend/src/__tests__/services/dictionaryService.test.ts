import { describe, expect, it, vi } from 'vitest'

const { mockJson, mockPublicApi } = vi.hoisted(() => {
  const mockJson = vi.fn()
  const mockPublicApi = {
    get: vi.fn(() => ({ json: mockJson })),
  }
  return {
    mockJson,
    mockPublicApi,
  }
})

vi.mock('ky', () => ({
  default: {
    create: vi.fn(() => mockPublicApi),
  },
}))

import { dictionaryService } from '@/services/dictionaryService'

describe('DictionaryService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAllDictionaries', () => {
    it('should call GET dictionaries on public api', async () => {
      const dictionaries = { grades: [{ value: 'junior', label: 'Junior' }] }
      mockJson.mockResolvedValue(dictionaries)

      const result = await dictionaryService.getAllDictionaries()

      expect(mockPublicApi.get).toHaveBeenCalledWith('dictionaries')
      expect(result).toEqual(dictionaries)
    })
  })
})
