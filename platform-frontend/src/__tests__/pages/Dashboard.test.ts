import { describe, expect, it, vi } from 'vitest'

describe('Dashboard logic', () => {
  function pluralizeDays(n: number): string {
    if (n % 10 === 1 && n % 100 !== 11) return 'день'
    if (n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)) return 'дня'
    return 'дней'
  }

  function questProgress(quest: { currentCount: number, targetCount: number }) {
    return Math.min(100, Math.round((quest.currentCount / quest.targetCount) * 100))
  }

  function isEventLive(event: { date: string }) {
    const now = new Date()
    const eventDate = new Date(event.date)
    const diffMs = now.getTime() - eventDate.getTime()
    return diffMs >= 0 && diffMs < 2 * 60 * 60 * 1000
  }

  function formatQuestDeadline(dateStr: string) {
    const date = new Date(dateStr)
    const now = new Date()
    const diffMs = date.getTime() - now.getTime()
    const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
    if (diffDays <= 0) return 'Истекает'
    if (diffDays === 1) return '1 день'
    if (diffDays <= 7) return `${diffDays} дн.`
    return date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
  }

  describe('pluralizeDays', () => {
    it('returns "день" for 1', () => {
      expect(pluralizeDays(1)).toBe('день')
    })

    it('returns "дня" for 2', () => {
      expect(pluralizeDays(2)).toBe('дня')
    })

    it('returns "дня" for 3 and 4', () => {
      expect(pluralizeDays(3)).toBe('дня')
      expect(pluralizeDays(4)).toBe('дня')
    })

    it('returns "дней" for 5', () => {
      expect(pluralizeDays(5)).toBe('дней')
    })

    it('returns "дней" for 11 (special case)', () => {
      expect(pluralizeDays(11)).toBe('дней')
    })

    it('returns "дней" for 12, 13, 14', () => {
      expect(pluralizeDays(12)).toBe('дней')
      expect(pluralizeDays(13)).toBe('дней')
      expect(pluralizeDays(14)).toBe('дней')
    })

    it('returns "день" for 21', () => {
      expect(pluralizeDays(21)).toBe('день')
    })

    it('returns "дня" for 22', () => {
      expect(pluralizeDays(22)).toBe('дня')
    })

    it('returns "дней" for 25', () => {
      expect(pluralizeDays(25)).toBe('дней')
    })

    it('returns "дней" for 0', () => {
      expect(pluralizeDays(0)).toBe('дней')
    })

    it('returns "дней" for 111', () => {
      expect(pluralizeDays(111)).toBe('дней')
    })

    it('returns "день" for 101', () => {
      expect(pluralizeDays(101)).toBe('день')
    })
  })

  describe('questProgress', () => {
    it('calculates 50%', () => {
      expect(questProgress({ currentCount: 5, targetCount: 10 })).toBe(50)
    })

    it('calculates 100%', () => {
      expect(questProgress({ currentCount: 10, targetCount: 10 })).toBe(100)
    })

    it('caps at 100% when over target', () => {
      expect(questProgress({ currentCount: 15, targetCount: 10 })).toBe(100)
    })

    it('returns 0% for no progress', () => {
      expect(questProgress({ currentCount: 0, targetCount: 10 })).toBe(0)
    })

    it('rounds fractional percentages', () => {
      expect(questProgress({ currentCount: 1, targetCount: 3 })).toBe(33)
      expect(questProgress({ currentCount: 2, targetCount: 3 })).toBe(67)
    })
  })

  describe('isEventLive', () => {
    it('returns true for event that just started', () => {
      const event = { date: new Date().toISOString() }
      expect(isEventLive(event)).toBe(true)
    })

    it('returns true for event started 1 hour ago', () => {
      const event = { date: new Date(Date.now() - 3600000).toISOString() }
      expect(isEventLive(event)).toBe(true)
    })

    it('returns false for event started 3 hours ago', () => {
      const event = { date: new Date(Date.now() - 3 * 3600000).toISOString() }
      expect(isEventLive(event)).toBe(false)
    })

    it('returns false for future event', () => {
      const event = { date: new Date(Date.now() + 3600000).toISOString() }
      expect(isEventLive(event)).toBe(false)
    })

    it('returns false for event exactly 2 hours ago', () => {
      const event = { date: new Date(Date.now() - 2 * 3600000).toISOString() }
      expect(isEventLive(event)).toBe(false)
    })
  })

  describe('formatQuestDeadline', () => {
    it('returns "Истекает" for past dates', () => {
      const past = new Date(Date.now() - 86400000).toISOString()
      expect(formatQuestDeadline(past)).toBe('Истекает')
    })

    it('returns "1 день" for ~1 day from now', () => {
      const tomorrow = new Date(Date.now() + 86400000 * 0.5).toISOString()
      expect(formatQuestDeadline(tomorrow)).toBe('1 день')
    })

    it('returns days for near future (within 7 days)', () => {
      const future = new Date(Date.now() + 86400000 * 4.5).toISOString()
      expect(formatQuestDeadline(future)).toBe('5 дн.')
    })

    it('returns formatted date for more than 7 days', () => {
      const farFuture = new Date(Date.now() + 86400000 * 30).toISOString()
      const result = formatQuestDeadline(farFuture)
      expect(result).toMatch(/\d{2}\.\d{2}/)
    })

    it('returns "Истекает" for current moment', () => {
      const now = new Date(Date.now() - 1000).toISOString()
      expect(formatQuestDeadline(now)).toBe('Истекает')
    })
  })
})
