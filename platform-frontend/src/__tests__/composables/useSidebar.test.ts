import { describe, expect, it } from 'vitest'
import { useSidebar } from '@/composables/useSidebar'

describe('useSidebar', () => {
  it('returns isOpen, sidebarGroups, and toggleSidebar', () => {
    const { isOpen, sidebarGroups, toggleSidebar } = useSidebar()

    expect(isOpen).toBeDefined()
    expect(sidebarGroups).toBeDefined()
    expect(toggleSidebar).toBeDefined()
    expect(typeof toggleSidebar).toBe('function')
  })

  it('isOpen starts as false', () => {
    const { isOpen } = useSidebar()

    expect(isOpen.value).toBe(false)
  })

  it('toggleSidebar toggles isOpen from false to true and back', () => {
    const { isOpen, toggleSidebar } = useSidebar()

    expect(isOpen.value).toBe(false)

    toggleSidebar()
    expect(isOpen.value).toBe(true)

    toggleSidebar()
    expect(isOpen.value).toBe(false)
  })

  it('sidebarGroups has 5 groups', () => {
    const { sidebarGroups } = useSidebar()

    expect(sidebarGroups.value).toHaveLength(5)
  })

  it('sidebarGroups contain expected labels', () => {
    const { sidebarGroups } = useSidebar()
    const labels = sidebarGroups.value.map(g => g.label)

    expect(labels).toEqual([undefined, 'Сообщество', 'Знания', 'Активность', 'Бонусы'])
  })

  it.each([
    ['/achievements'],
    ['/leaderboard'],
    ['/kudos'],
    ['/my-stats'],
  ])('does not contain legacy gamification path %s (свернуты в /progress)', (path) => {
    // Регрессионный гард: эти разделы свёрнуты в табы /progress (см. #324, #331).
    // Если кто-то случайно вернёт пункт сюда, sidebar и хаб разойдутся
    // и активные классы перестанут совпадать с реальным URL.
    const { sidebarGroups } = useSidebar()
    const allPaths = sidebarGroups.value.flatMap(g => g.items.map(i => i.path))
    expect(allPaths).not.toContain(path)
  })

  it('contains /progress as the single gamification entry point', () => {
    const { sidebarGroups } = useSidebar()
    const progressItems = sidebarGroups.value.flatMap(g =>
      g.items.filter(i => i.path === '/progress'),
    )
    expect(progressItems).toHaveLength(1)
  })

  it('multiple calls return the same singleton state', () => {
    const first = useSidebar()
    const second = useSidebar()

    expect(first.isOpen).toBe(second.isOpen)
    expect(first.sidebarGroups).toBe(second.sidebarGroups)
    expect(first.toggleSidebar).toBe(second.toggleSidebar)

    first.toggleSidebar()
    expect(second.isOpen.value).toBe(first.isOpen.value)
  })
})
