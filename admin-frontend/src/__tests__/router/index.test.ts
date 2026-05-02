import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock useToast (used by authService and errorService)
vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({ toast: vi.fn() }),
}))

// Mock authService - need to control isAuthenticated and checkAuth
let mockIsAuthenticated = false
const mockCheckAuth = vi.fn(() => {
  return mockIsAuthenticated
})

vi.mock('@/services/authService', () => ({
  isAuthenticated: {
    get value() {
      return mockIsAuthenticated
    },
    set value(v: boolean) {
      mockIsAuthenticated = v
    },
  },
  checkAuth: () => mockCheckAuth(),
  logout: vi.fn(),
}))

// Mock errorService
vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

const routerModule = await import('@/router/index')
const router = routerModule.default

describe('router', () => {
  beforeEach(async () => {
    vi.clearAllMocks()
    mockIsAuthenticated = false
    // Reset router to a known state
    await router.push('/login')
  })

  describe('route definitions', () => {
    it('has 21 routes defined', () => {
      // 19 базовых + /daily-tasks + /challenges (геймификация, PR-322)
      expect(router.getRoutes()).toHaveLength(21)
    })

    it('has a login route', () => {
      const loginRoute = router.getRoutes().find(r => r.name === 'login')
      expect(loginRoute).toBeDefined()
      expect(loginRoute!.path).toBe('/login')
    })

    it('has a dashboard route with requiresAuth', () => {
      const route = router.getRoutes().find(r => r.name === 'dashboard')
      expect(route).toBeDefined()
      expect(route!.path).toBe('/dashboard')
      expect(route!.meta.requiresAuth).toBe(true)
    })

    it.each([
      ['mentors', '/mentors'],
      ['members', '/members'],
      ['reviews', '/reviews'],
      ['mentor-reviews', '/mentor-reviews'],
      ['events', '/events'],
      ['resumes', '/resumes'],
      ['audit-logs', '/audit-logs'],
      ['points', '/points'],
      ['referrals', '/referrals'],
      ['chat-activity', '/chat-activity'],
      ['chat-quests', '/chat-quests'],
      ['raffles', '/raffles'],
      ['subscriptions', '/subscriptions'],
      ['feedback', '/feedback'],
      ['moderation', '/moderation'],
    ])('has route "%s" at path "%s" requiring auth', (name, path) => {
      const route = router.getRoutes().find(r => r.name === name)
      expect(route).toBeDefined()
      expect(route!.path).toBe(path)
      expect(route!.meta.requiresAuth).toBe(true)
    })

    it('login route does not require auth', () => {
      const loginRoute = router.getRoutes().find(r => r.name === 'login')
      expect(loginRoute!.meta.requiresAuth).toBeFalsy()
    })

    it('all routes have unique names (excluding root redirect)', () => {
      const names = router.getRoutes()
        .filter(r => r.name !== undefined)
        .map(r => r.name)
      const uniqueNames = new Set(names)
      expect(uniqueNames.size).toBe(names.length)
    })

    it('all routes have unique paths', () => {
      const paths = router.getRoutes().map(r => r.path)
      const uniquePaths = new Set(paths)
      expect(uniquePaths.size).toBe(paths.length)
    })
  })

  describe('auth guard', () => {
    it('redirects unauthenticated users to login with redirect query', async () => {
      mockIsAuthenticated = false

      await router.push('/dashboard')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('login')
      expect(router.currentRoute.value.query.redirect).toBe('/dashboard')
    })

    it('allows authenticated users to access protected routes', async () => {
      mockIsAuthenticated = true

      await router.push('/dashboard')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('dashboard')
    })

    it('redirects authenticated users away from login to dashboard', async () => {
      mockIsAuthenticated = true

      // Navigate to dashboard first so we're not already on /login
      await router.push('/dashboard')
      await router.isReady()

      // Now try /login — should redirect back to dashboard
      await router.push('/login')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('dashboard')
    })

    it('allows unauthenticated users to access login page', async () => {
      mockIsAuthenticated = false

      await router.push('/login')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('login')
    })

    it('calls checkAuth on every navigation', async () => {
      mockCheckAuth.mockClear()
      mockIsAuthenticated = false

      await router.push('/dashboard')
      await router.isReady()

      expect(mockCheckAuth).toHaveBeenCalled()
    })

    it('preserves redirect query param for deep links', async () => {
      mockIsAuthenticated = false

      await router.push('/events')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('login')
      expect(router.currentRoute.value.query.redirect).toBe('/events')
    })
  })
})
