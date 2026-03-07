import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockJson = vi.fn()
const mockPublicApi = {
  get: vi.fn(() => ({ json: mockJson })),
}

vi.mock('ky', () => ({
  default: {
    create: vi.fn(() => mockPublicApi),
  },
}))

const { DictionaryService, dictionaryService } = await import('@/services/dictionaryService')

describe('dictionaryService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('exports a DictionaryService class', () => {
    expect(DictionaryService).toBeDefined()
    expect(dictionaryService).toBeInstanceOf(DictionaryService)
  })

  describe('getAllDictionaries', () => {
    it('calls publicApi.get with "dictionaries" and returns parsed JSON', async () => {
      const mockDictionaries = {
        grades: [
          { value: 'junior', label: 'Junior' },
          { value: 'middle', label: 'Middle' },
        ],
        roles: [
          { value: 'developer', label: 'Developer' },
        ],
      }
      mockJson.mockResolvedValueOnce(mockDictionaries)

      const result = await dictionaryService.getAllDictionaries()

      expect(mockPublicApi.get).toHaveBeenCalledWith('dictionaries')
      expect(result).toEqual(mockDictionaries)
    })

    it('propagates errors from the API', async () => {
      mockJson.mockRejectedValueOnce(new Error('Network error'))

      await expect(dictionaryService.getAllDictionaries()).rejects.toThrow('Network error')
    })
  })
})
