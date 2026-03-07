import { describe, expect, it } from 'vitest'

describe('Achievements logic', () => {
  const categories = [
    { key: 'all', label: 'Все' },
    { key: 'events', label: 'События' },
    { key: 'points', label: 'Баллы' },
    { key: 'social', label: 'Социальные' },
    { key: 'activity', label: 'Активность' },
  ]

  interface AchievementItem {
    category: string
    name: string
  }

  function filterItems(items: AchievementItem[], category: string) {
    if (category === 'all') return items
    return items.filter(a => a.category === category)
  }

  describe('categories', () => {
    it('has 5 items', () => {
      expect(categories).toHaveLength(5)
    })

    it('first category is "all"', () => {
      expect(categories[0].key).toBe('all')
      expect(categories[0].label).toBe('Все')
    })

    it('contains events category', () => {
      expect(categories.find(c => c.key === 'events')).toBeDefined()
    })

    it('contains all expected keys', () => {
      const keys = categories.map(c => c.key)
      expect(keys).toEqual(['all', 'events', 'points', 'social', 'activity'])
    })
  })

  describe('filterItems', () => {
    const items: AchievementItem[] = [
      { category: 'events', name: 'Первое событие' },
      { category: 'points', name: '100 баллов' },
      { category: 'social', name: 'Первый друг' },
      { category: 'events', name: 'Второе событие' },
      { category: 'activity', name: 'Активист' },
    ]

    it('returns all items when category is "all"', () => {
      expect(filterItems(items, 'all')).toHaveLength(5)
      expect(filterItems(items, 'all')).toBe(items)
    })

    it('filters by events category', () => {
      const result = filterItems(items, 'events')
      expect(result).toHaveLength(2)
      expect(result.every(i => i.category === 'events')).toBe(true)
    })

    it('filters by points category', () => {
      const result = filterItems(items, 'points')
      expect(result).toHaveLength(1)
      expect(result[0].name).toBe('100 баллов')
    })

    it('returns empty array for category with no items', () => {
      expect(filterItems(items, 'nonexistent')).toHaveLength(0)
    })

    it('returns empty array when items is empty', () => {
      expect(filterItems([], 'events')).toHaveLength(0)
    })
  })
})
