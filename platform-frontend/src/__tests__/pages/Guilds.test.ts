import { describe, expect, it } from 'vitest'

describe('Guilds logic', () => {
  const colorOptions = ['#6366f1', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ef4444', '#06b6d4']

  function displayName(firstName: string, lastName: string) {
    return [firstName, lastName].filter(Boolean).join(' ')
  }

  function isInAnyGuild(guilds: { isMember: boolean }[]) {
    return guilds.some(g => g.isMember)
  }

  describe('colorOptions', () => {
    it('has exactly 8 colors', () => {
      expect(colorOptions).toHaveLength(8)
    })

    it('contains valid hex color strings', () => {
      for (const color of colorOptions) {
        expect(color).toMatch(/^#[0-9a-f]{6}$/)
      }
    })
  })

  describe('displayName', () => {
    it('joins first and last name', () => {
      expect(displayName('John', 'Doe')).toBe('John Doe')
    })

    it('returns only first name when last name is empty', () => {
      expect(displayName('John', '')).toBe('John')
    })

    it('returns only last name when first name is empty', () => {
      expect(displayName('', 'Doe')).toBe('Doe')
    })

    it('returns empty string when both names are empty', () => {
      expect(displayName('', '')).toBe('')
    })
  })

  describe('isInAnyGuild', () => {
    it('returns true when at least one guild has isMember true', () => {
      const guilds = [
        { isMember: false },
        { isMember: true },
        { isMember: false },
      ]
      expect(isInAnyGuild(guilds)).toBe(true)
    })

    it('returns false when no guilds have isMember true', () => {
      const guilds = [
        { isMember: false },
        { isMember: false },
      ]
      expect(isInAnyGuild(guilds)).toBe(false)
    })

    it('returns false for empty guild list', () => {
      expect(isInAnyGuild([])).toBe(false)
    })
  })
})
