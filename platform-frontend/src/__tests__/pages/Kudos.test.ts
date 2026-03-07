import { describe, expect, it } from 'vitest'

describe('Kudos logic', () => {
  function displayName(firstName: string, lastName: string) {
    return [firstName, lastName].filter(Boolean).join(' ')
  }

  function timeAgo(dateStr: string) {
    const diff = Date.now() - new Date(dateStr).getTime()
    const mins = Math.floor(diff / 60000)
    if (mins < 60) return `${mins} мин. назад`
    const hours = Math.floor(mins / 60)
    if (hours < 24) return `${hours} ч. назад`
    const days = Math.floor(hours / 24)
    return `${days} дн. назад`
  }

  describe('displayName', () => {
    it('joins first and last name', () => {
      expect(displayName('Иван', 'Иванов')).toBe('Иван Иванов')
    })

    it('handles empty last name', () => {
      expect(displayName('Иван', '')).toBe('Иван')
    })

    it('handles both empty', () => {
      expect(displayName('', '')).toBe('')
    })

    it('handles empty first name', () => {
      expect(displayName('', 'Иванов')).toBe('Иванов')
    })
  })

  describe('timeAgo', () => {
    it('returns minutes for recent times', () => {
      const fiveMinAgo = new Date(Date.now() - 5 * 60000).toISOString()
      expect(timeAgo(fiveMinAgo)).toBe('5 мин. назад')
    })

    it('returns hours for times within a day', () => {
      const threeHoursAgo = new Date(Date.now() - 3 * 3600000).toISOString()
      expect(timeAgo(threeHoursAgo)).toBe('3 ч. назад')
    })

    it('returns days for older times', () => {
      const twoDaysAgo = new Date(Date.now() - 2 * 86400000).toISOString()
      expect(timeAgo(twoDaysAgo)).toBe('2 дн. назад')
    })

    it('returns 0 minutes for just now', () => {
      const now = new Date().toISOString()
      expect(timeAgo(now)).toBe('0 мин. назад')
    })
  })
})
