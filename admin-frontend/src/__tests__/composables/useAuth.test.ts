import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock authService
const mockCheckAuth = vi.fn()
const mockLogoutService = vi.fn()
const mockIsAuthenticated = { value: false }
const mockIsLoading = { value: false }

vi.mock('@/services/authService', () => ({
  checkAuth: () => mockCheckAuth(),
  isAuthenticated: mockIsAuthenticated,
  isLoading: mockIsLoading,
  logout: () => mockLogoutService(),
}))

// Mock vue-router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}))

const { useAuth } = await import('@/composables/useAuth')

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockIsAuthenticated.value = false
    mockIsLoading.value = false
  })

  it('calls checkAuth on initialization', () => {
    useAuth()

    expect(mockCheckAuth).toHaveBeenCalled()
  })

  it('exposes isAuthenticated ref', () => {
    mockIsAuthenticated.value = true

    const { isAuthenticated } = useAuth()

    expect(isAuthenticated.value).toBe(true)
  })

  it('exposes isLoading ref', () => {
    mockIsLoading.value = true

    const { isLoading } = useAuth()

    expect(isLoading.value).toBe(true)
  })

  describe('logout', () => {
    it('calls logout service and redirects to login', () => {
      const { logout } = useAuth()

      logout()

      expect(mockLogoutService).toHaveBeenCalled()
      expect(mockPush).toHaveBeenCalledWith('/login')
    })
  })
})
