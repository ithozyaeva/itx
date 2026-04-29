import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

// Mock all page component imports to avoid loading real Vue SFCs
vi.mock('@/pages/Achievements.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/AutoApplyBot.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Content.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Dashboard.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Events.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Kudos.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Leaderboard.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Marketplace.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/MemberProfile.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/MentorProfile.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Mentors.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/MyPoints.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/MyReviews.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/MyStats.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Quests.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Raffles.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/ReferalLinks.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Resumes.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/TaskExchange.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/User.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/Casino.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/pages/NotificationSettings.vue', () => ({ default: { template: '<div />' } }))

describe('router', () => {
  let router: any

  beforeEach(async () => {
    vi.resetModules()
    setActivePinia(createPinia())
    const mod = await import('@/router/index')
    router = mod.default
  })

  describe('route definitions', () => {
    it('has expected number of routes', () => {
      const routes = router.getRoutes()
      expect(routes.length).toBeGreaterThanOrEqual(20)
    })

    it.each([
      ['/', 'dashboard'],
      ['/me', 'profile'],
      ['/events', 'events'],
      ['/faq', 'faq'],
      ['/mentors', 'mentors'],
      ['/referals', 'referals'],
      ['/resumes', 'resumes'],
      ['/my-reviews', 'myReviews'],
      ['/points', 'myPoints'],
      ['/leaderboard', 'leaderboard'],
      ['/achievements', 'achievements'],
      ['/marketplace', 'marketplace'],
      ['/tasks', 'taskExchange'],
      ['/quests', 'quests'],
      ['/auto-apply', 'autoApplyBot'],
      ['/kudos', 'kudos'],
      ['/raffles', 'raffles'],
      ['/my-stats', 'myStats'],
    ])('route %s has name %s', (path, name) => {
      const route = router.getRoutes().find((r: any) => r.name === name)
      expect(route).toBeDefined()
      expect(route.path).toBe(path)
    })
  })

  describe('parameterized routes', () => {
    it('has memberProfile route with :id param', () => {
      const route = router.getRoutes().find((r: any) => r.name === 'memberProfile')
      expect(route).toBeDefined()
      expect(route.path).toBe('/members/:id')
    })

    it('has mentorProfile route with :id param', () => {
      const route = router.getRoutes().find((r: any) => r.name === 'mentorProfile')
      expect(route).toBeDefined()
      expect(route.path).toBe('/mentors/:id')
    })
  })

  describe('navigation', () => {
    it('allows navigation to dashboard (open without subscription)', async () => {
      await router.push('/')
      await router.isReady()
      expect(router.currentRoute.value.name).toBe('dashboard')
    })

    it('redirects gated route to dashboard when no subscription tier', async () => {
      await router.push('/events')
      await router.isReady()
      expect(router.currentRoute.value.name).toBe('dashboard')
    })

    it('allows navigation to /members/:id without subscription', async () => {
      // Профиль участника — публичная карточка, не должен быть за гейтом
      // (ссылку даёт бот в /whois, юзер не должен улетать на /tariffs).
      await router.push('/members/42')
      await router.isReady()
      expect(router.currentRoute.value.name).toBe('memberProfile')
    })

    it('allows navigation to gated route when user has tier', async () => {
      // Имитируем подписчика — useUser читает tg_user:v2 из localStorage.
      localStorage.setItem('tg_user:v2', JSON.stringify({
        data: {
          id: 1,
          telegramID: 1,
          roles: ['SUBSCRIBER'],
          subscriptionTier: { id: 2, slug: 'foreman', name: 'Бригадир', level: 2 },
        },
        savedAt: Date.now(),
      }))
      vi.resetModules()
      const mod = await import('@/router/index')
      const r = mod.default
      await r.push('/events')
      await r.isReady()
      expect(r.currentRoute.value.name).toBe('events')
      localStorage.removeItem('tg_user:v2')
    })

    it('redirects subscriber away from /tariffs to dashboard', async () => {
      localStorage.setItem('tg_user:v2', JSON.stringify({
        data: {
          id: 1,
          telegramID: 1,
          roles: ['SUBSCRIBER'],
          subscriptionTier: { id: 3, slug: 'master', name: 'Хозяин', level: 3 },
        },
        savedAt: Date.now(),
      }))
      vi.resetModules()
      const mod = await import('@/router/index')
      const r = mod.default
      await r.push('/tariffs')
      await r.isReady()
      expect(r.currentRoute.value.name).toBe('dashboard')
      localStorage.removeItem('tg_user:v2')
    })

    it('allows UNSUBSCRIBER to view /tariffs', async () => {
      // tg_user без subscriptionTier — UNSUBSCRIBER в новой логике.
      localStorage.setItem('tg_user:v2', JSON.stringify({
        data: { id: 1, telegramID: 1, roles: ['UNSUBSCRIBER'] },
        savedAt: Date.now(),
      }))
      vi.resetModules()
      const mod = await import('@/router/index')
      const r = mod.default
      await r.push('/tariffs')
      await r.isReady()
      expect(r.currentRoute.value.name).toBe('tariffs')
      localStorage.removeItem('tg_user:v2')
    })
  })
})
