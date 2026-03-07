import { describe, expect, it, vi } from 'vitest'

const mockUseLocalStorage = vi.fn()

vi.mock('@vueuse/core', () => ({
  useLocalStorage: mockUseLocalStorage,
}))

describe('useToken', () => {
  it('calls useLocalStorage with "tg_token" key and null default', async () => {
    const fakeRef = { value: null }
    mockUseLocalStorage.mockReturnValue(fakeRef)

    const { useToken } = await import('@/composables/useToken')
    const result = useToken()

    expect(mockUseLocalStorage).toHaveBeenCalledWith('tg_token', null)
    expect(result).toBe(fakeRef)
  })
})
