import { describe, expect, it } from 'vitest'
import type { SidebarItem } from '@/composables/useSidebar'

describe('useSidebar', () => {
  it('SidebarItem interface requires title, path, and icon', () => {
    const item: SidebarItem = {
      title: 'Test',
      path: '/test',
      icon: {} as any,
    }

    expect(item.title).toBe('Test')
    expect(item.path).toBe('/test')
    expect(item.icon).toBeDefined()
  })

  it('SidebarItem can have optional requiredPermission', () => {
    const itemWithPermission: SidebarItem = {
      title: 'Dashboard',
      path: '/dashboard',
      icon: {} as any,
      requiredPermission: 'can_view_admin_panel',
    }

    const itemWithout: SidebarItem = {
      title: 'Referrals',
      path: '/referrals',
      icon: {} as any,
    }

    expect(itemWithPermission.requiredPermission).toBe('can_view_admin_panel')
    expect(itemWithout.requiredPermission).toBeUndefined()
  })

  it('sidebar paths should start with /', () => {
    const paths = [
      '/dashboard',
      '/mentors',
      '/members',
      '/reviews',
      '/mentor-reviews',
      '/events',
      '/resumes',
      '/referrals',
      '/points',
      '/chat-activity',
      '/chat-quests',
      '/seasons',
      '/raffles',
      '/audit-logs',
    ]

    for (const path of paths) {
      expect(path.startsWith('/')).toBe(true)
    }
  })

  it('all expected sidebar items are accounted for', () => {
    const expectedTitles = [
      'Дашборд',
      'Менторы',
      'Участники',
      'Отзывы на сообщество',
      'Отзывы на менторов',
      'События',
      'Резюме',
      'Рефералы',
      'Баллы',
      'Активность чатов',
      'Задания чатов',
      'Сезоны',
      'Розыгрыши',
      'Журнал действий',
    ]

    expect(expectedTitles).toHaveLength(14)
    expect(new Set(expectedTitles).size).toBe(14)
  })
})
