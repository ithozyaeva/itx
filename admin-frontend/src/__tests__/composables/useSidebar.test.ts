import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock usePermissions
const mockHasPermission = vi.fn()
vi.mock('@/composables/usePermissions', () => ({
  usePermissions: () => ({
    hasPermission: { value: mockHasPermission },
  }),
}))

const { useSidebar } = await import('@/composables/useSidebar')

describe('useSidebar', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Reset global state
    const { isCollapsed, isMobileOpen } = useSidebar()
    isCollapsed.value = false
    isMobileOpen.value = false
  })

  describe('sidebarItems filtering', () => {
    it('shows all items when user has all permissions', () => {
      mockHasPermission.mockReturnValue(true)

      const { sidebarItems } = useSidebar()

      expect(sidebarItems.value).toHaveLength(16)
    })

    it('shows only items without requiredPermission when user has no permissions', () => {
      mockHasPermission.mockReturnValue(false)

      const { sidebarItems } = useSidebar()

      // Items without requiredPermission: Рефералы, Активность чатов, Задания чатов, Сезоны, Розыгрыши, Мини-игры
      const itemsWithoutPermission = sidebarItems.value
      expect(itemsWithoutPermission.length).toBe(6)
      expect(itemsWithoutPermission.map(i => i.path)).toEqual(
        expect.arrayContaining(['/referrals', '/chat-activity', '/chat-quests', '/seasons', '/raffles', '/minigames']),
      )
    })

    it('selectively filters based on specific permissions', () => {
      mockHasPermission.mockImplementation((perm: string) => {
        return perm === 'can_view_admin_panel' || perm === 'can_view_admin_events'
      })

      const { sidebarItems } = useSidebar()

      const paths = sidebarItems.value.map(i => i.path)
      expect(paths).toContain('/dashboard') // has can_view_admin_panel
      expect(paths).toContain('/events') // has can_view_admin_events
      expect(paths).toContain('/referrals') // no permission required
      expect(paths).not.toContain('/mentors') // requires can_view_admin_mentors
      expect(paths).not.toContain('/members') // requires can_view_admin_members
    })
  })

  describe('sidebar items structure', () => {
    it('all items have title, path, and icon', () => {
      mockHasPermission.mockReturnValue(true)

      const { sidebarItems } = useSidebar()

      for (const item of sidebarItems.value) {
        expect(item.title).toBeTruthy()
        expect(item.path).toMatch(/^\//)
        expect(item.icon).toBeDefined()
      }
    })

    it('has correct paths for all items', () => {
      mockHasPermission.mockReturnValue(true)

      const { sidebarItems } = useSidebar()

      const expectedPaths = [
        '/dashboard', '/mentors', '/members', '/events', '/resumes',
        '/reviews', '/mentor-reviews', '/points', '/chat-activity',
        '/chat-quests', '/seasons', '/raffles', '/minigames',
        '/subscriptions', '/referrals', '/audit-logs',
      ]
      const actualPaths = sidebarItems.value.map(i => i.path)
      expect(actualPaths).toEqual(expectedPaths)
    })
  })

  describe('toggleSidebar', () => {
    it('toggles collapsed state', () => {
      const { isCollapsed, toggleSidebar } = useSidebar()

      expect(isCollapsed.value).toBe(false)
      toggleSidebar()
      expect(isCollapsed.value).toBe(true)
      toggleSidebar()
      expect(isCollapsed.value).toBe(false)
    })
  })

  describe('toggleMobileSidebar', () => {
    it('toggles mobile open state', () => {
      const { isMobileOpen, toggleMobileSidebar } = useSidebar()

      expect(isMobileOpen.value).toBe(false)
      toggleMobileSidebar()
      expect(isMobileOpen.value).toBe(true)
      toggleMobileSidebar()
      expect(isMobileOpen.value).toBe(false)
    })
  })

  describe('closeMobileSidebar', () => {
    it('sets mobile open to false', () => {
      const { isMobileOpen, toggleMobileSidebar, closeMobileSidebar } = useSidebar()

      toggleMobileSidebar() // open
      expect(isMobileOpen.value).toBe(true)

      closeMobileSidebar()
      expect(isMobileOpen.value).toBe(false)
    })

    it('is idempotent when already closed', () => {
      const { isMobileOpen, closeMobileSidebar } = useSidebar()

      closeMobileSidebar()
      expect(isMobileOpen.value).toBe(false)
    })
  })

  describe('shared global state', () => {
    it('isCollapsed is shared across calls', () => {
      const sidebar1 = useSidebar()
      const sidebar2 = useSidebar()

      sidebar1.toggleSidebar()

      expect(sidebar2.isCollapsed.value).toBe(true)
    })

    it('isMobileOpen is shared across calls', () => {
      const sidebar1 = useSidebar()
      const sidebar2 = useSidebar()

      sidebar1.toggleMobileSidebar()

      expect(sidebar2.isMobileOpen.value).toBe(true)
    })
  })
})
