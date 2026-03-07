import type { RouteLocationNormalized } from 'vue-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

const mockIsAuthenticated = ref(false)

vi.mock('@/composables/useAuth', () => ({
  useAuth: () => ({
    isAuthenticated: mockIsAuthenticated,
  }),
}))

const { useAuthGuard } = await import('@/composables/useAuthGuard')

function createRoute(overrides: Partial<RouteLocationNormalized> = {}): RouteLocationNormalized {
  return {
    path: '/',
    fullPath: '/',
    name: undefined,
    hash: '',
    query: {},
    params: {},
    matched: [],
    redirectedFrom: undefined,
    meta: {},
    ...overrides,
  } as RouteLocationNormalized
}

describe('useAuthGuard', () => {
  beforeEach(() => {
    mockIsAuthenticated.value = false
  })

  it('redirects to login when not authenticated and route requires auth', () => {
    const { guardRoute } = useAuthGuard()
    const route = createRoute({ path: '/dashboard', fullPath: '/dashboard', meta: { requiresAuth: true } })

    const result = guardRoute(route)

    expect(result).toEqual({ path: '/login', query: { redirect: '/dashboard' } })
  })

  it('includes redirect query param with the original full path', () => {
    const { guardRoute } = useAuthGuard()
    const route = createRoute({
      path: '/members',
      fullPath: '/members?page=2',
      meta: { requiresAuth: true },
    })

    const result = guardRoute(route)

    expect(result).toEqual({ path: '/login', query: { redirect: '/members?page=2' } })
  })

  it('redirects to dashboard when authenticated and going to login', () => {
    mockIsAuthenticated.value = true
    const { guardRoute } = useAuthGuard()
    const route = createRoute({ path: '/login', fullPath: '/login' })

    const result = guardRoute(route)

    expect(result).toEqual({ path: '/dashboard' })
  })

  it('returns true for normal navigation when authenticated', () => {
    mockIsAuthenticated.value = true
    const { guardRoute } = useAuthGuard()
    const route = createRoute({ path: '/dashboard', fullPath: '/dashboard', meta: { requiresAuth: true } })

    const result = guardRoute(route)

    expect(result).toBe(true)
  })

  it('returns true for public route when not authenticated', () => {
    const { guardRoute } = useAuthGuard()
    const route = createRoute({ path: '/about', fullPath: '/about', meta: {} })

    const result = guardRoute(route)

    expect(result).toBe(true)
  })
})
