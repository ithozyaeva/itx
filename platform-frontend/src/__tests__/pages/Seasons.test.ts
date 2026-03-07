import { describe, expect, it } from 'vitest'

describe('Seasons logic', () => {
  function formatDate(d: string) {
    return new Date(d).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short', year: 'numeric' })
  }

  function displayName(firstName: string, lastName: string) {
    return [firstName, lastName].filter(Boolean).join(' ')
  }

  const rankColors = ['text-yellow-500', 'text-zinc-400', 'text-amber-700']

  describe('formatDate', () => {
    it('formats a date in Russian locale', () => {
      const result = formatDate('2025-01-15T00:00:00Z')
      expect(result).toContain('2025')
      expect(result).toContain('15')
    })

    it('formats another date correctly', () => {
      const result = formatDate('2024-12-31T00:00:00Z')
      expect(result).toContain('2024')
      expect(result).toContain('31')
    })
  })

  describe('displayName', () => {
    it('joins first and last name', () => {
      expect(displayName('Анна', 'Петрова')).toBe('Анна Петрова')
    })

    it('handles empty last name', () => {
      expect(displayName('Анна', '')).toBe('Анна')
    })

    it('handles empty first name', () => {
      expect(displayName('', 'Петрова')).toBe('Петрова')
    })

    it('handles both empty', () => {
      expect(displayName('', '')).toBe('')
    })
  })

  describe('rank arrays', () => {
    it('rankColors has 3 items', () => {
      expect(rankColors).toHaveLength(3)
    })

    it('rankColors contains expected values', () => {
      expect(rankColors[0]).toBe('text-yellow-500')
      expect(rankColors[1]).toBe('text-zinc-400')
      expect(rankColors[2]).toBe('text-amber-700')
    })
  })
})
