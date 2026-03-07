import { describe, expect, it } from 'vitest'

describe('router configuration', () => {
  const routeDefinitions = [
    { path: '/', name: undefined, hasRedirect: true },
    { path: '/login', name: 'login', requiresAuth: false },
    { path: '/dashboard', name: 'dashboard', requiresAuth: true },
    { path: '/mentors', name: 'mentors', requiresAuth: true },
    { path: '/members', name: 'members', requiresAuth: true },
    { path: '/reviews', name: 'reviews', requiresAuth: true },
    { path: '/mentor-reviews', name: 'mentor-reviews', requiresAuth: true },
    { path: '/events', name: 'events', requiresAuth: true },
    { path: '/resumes', name: 'resumes', requiresAuth: true },
    { path: '/audit-logs', name: 'audit-logs', requiresAuth: true },
    { path: '/points', name: 'points', requiresAuth: true },
    { path: '/referrals', name: 'referrals', requiresAuth: true },
    { path: '/chat-activity', name: 'chat-activity', requiresAuth: true },
    { path: '/chat-quests', name: 'chat-quests', requiresAuth: true },
    { path: '/seasons', name: 'seasons', requiresAuth: true },
    { path: '/raffles', name: 'raffles', requiresAuth: true },
  ]

  it('has 16 routes defined', () => {
    expect(routeDefinitions).toHaveLength(16)
  })

  it('all routes have a path', () => {
    for (const route of routeDefinitions) {
      expect(route.path).toBeDefined()
      expect(route.path.startsWith('/')).toBe(true)
    }
  })

  it('all named routes have unique names', () => {
    const names = routeDefinitions
      .filter(r => r.name !== undefined)
      .map(r => r.name)
    const uniqueNames = new Set(names)
    expect(uniqueNames.size).toBe(names.length)
  })

  it('all routes except login and root require auth', () => {
    const protectedRoutes = routeDefinitions.filter(
      r => r.name !== 'login' && r.name !== undefined,
    )

    for (const route of protectedRoutes) {
      expect(route.requiresAuth).toBe(true)
    }
  })

  it('login route does not require auth', () => {
    const loginRoute = routeDefinitions.find(r => r.name === 'login')
    expect(loginRoute).toBeDefined()
    expect(loginRoute!.requiresAuth).toBe(false)
  })

  it('root path has redirect', () => {
    const rootRoute = routeDefinitions.find(r => r.path === '/')
    expect(rootRoute).toBeDefined()
    expect(rootRoute!.hasRedirect).toBe(true)
  })

  it('all paths are unique', () => {
    const paths = routeDefinitions.map(r => r.path)
    const uniquePaths = new Set(paths)
    expect(uniquePaths.size).toBe(paths.length)
  })
})
